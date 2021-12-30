package parser

import (
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// TODO: calling each Retrieve function is not smart, so should be refactored
// TODO(fred): I assume we could start with a retriever at the top level...

// ParseTokens parses Tokens, creating slice of Reindenter's.
//
// Each Reindenter is group of SQL clauses such as SelectGroup, FromGroup ...etc.
func ParseTokens(tokens []lexer.Token, opts ...Option) ([]group.Reindenter, error) {
	if err := isStartSupportedClause(tokens[0]); err != nil {
		return nil, err
	}

	var (
		result []group.Reindenter
	)

	for offset := 0; tokens[offset].Type != lexer.EOF; {
		afterComma := offset > 0 && tokens[offset-1].Type == lexer.COMMA

		r := NewRetriever(tokens[offset:], append(opts, withAfterComma(afterComma))...)
		elements, endIdx, err := r.Retrieve()
		if err != nil {
			return nil, errors.Wrap(err, "ParseTokens failed")
		}

		group, err := r.createGroup(elements)
		if err != nil {
			return nil, err
		}

		result = append(result, group)

		offset += endIdx
	}

	return result, nil
}

// isStartSupportedClause picks valid SQL statement starters.
//
// NOTE: unsupported at this moment:
//   * DDL statements
//   * EXECUTE
//   * EXPLAIN
//   * PREPARE
//.
func isStartSupportedClause(token lexer.Token) error {
	ttype := token.Type
	if ttype == lexer.SELECT ||
		ttype == lexer.UPDATE ||
		ttype == lexer.DELETE ||
		ttype == lexer.INSERT ||
		ttype == lexer.LOCK ||
		ttype == lexer.WITH {
		return nil
	}

	return fmt.Errorf("can not parse: not a valid start of sql statement: %q", token.Value)
}
