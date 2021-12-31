package parser

import (
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
)

type Parser struct {
	offset int
	tokens []lexer.Token
	*options
}

// New SQL parser.
func New(opts ...Option) *Parser {
	o := defaultOptions(opts...)

	return &Parser{
		options: o,
	}
}

// ParseTokens parses Tokens, creating slice of Reindenter's.
//
// Each Reindenter is a group of SQL clauses such as SelectGroup, FromGroup ...etc.
func (p *Parser) Parse(tokens []lexer.Token) ([]group.Reindenter, error) {
	if err := isStartSupportedClause(tokens[0]); err != nil {
		return nil, err
	}
	p.tokens = tokens

	return p.parseTokens()
}

func (p *Parser) parseTokens() ([]group.Reindenter, error) {
	if len(p.tokens) == 0 || p.offset >= len(p.tokens) || p.tokens[p.offset].Type == lexer.EOF {
		return nil, nil
	}

	result := make([]group.Reindenter, 0, len(p.tokens[p.offset:]))

	r := NewRetriever(p.tokens[p.offset:],
		withOptions(p.options.CloneWithOptions(
			withAfterComma(p.isAfterComma()),
			withAfterParenthesis(p.isAfterParenthesis()),
		)),
	)
	if r == nil {
		return nil, nil
	}

	// extra groups from tokens starting at this offset
	elements, endIdx, err := r.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("parseTokens failed: %w", err)
	}

	group, err := r.createGroup(elements)
	if err != nil {
		return nil, fmt.Errorf("parseTokens failed to create group: %w", err)
	}

	if group != nil {
		result = append(result, group)
	}

	p.offset += endIdx

	// parse any remaining tokens
	next, err := p.parseTokens()
	if err != nil {
		return nil, err
	}
	result = append(result, next...)

	return result, nil
}

// isAfterComma provides context about the current group relative to a comma.
func (p *Parser) isAfterComma() bool {
	return p.offset > 0 && p.tokens[p.offset-1].Type == lexer.COMMA
}

// isAfterParenthesis provides context about the current group relative to an opening parenthesis.
func (p *Parser) isAfterParenthesis() bool {
	return p.offset > 0 && p.tokens[p.offset-1].Type == lexer.STARTPARENTHESIS
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
