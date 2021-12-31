package parser

import (
	"errors"
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
)

// Retriever retrieves target SQL clause group from TokenSource.
type Retriever struct {
	TokenSource []lexer.Token

	result        []group.Reindenter
	indentLevel   int
	endTokenTypes map[lexer.TokenType]struct{}
	endIdx        int

	*options
}

// NewRetriever Creates Retriever that retrieves each target SQL clause.
//
// It returns nil if no token may be used to form a group.
//
// Each Retriever have endKeywords in order to stop retrieving.
func NewRetriever(tokens []lexer.Token, opts ...Option) *Retriever {
	if len(tokens) == 0 {
		panic("invalid source for NewRetriever: must have some tokens")
	}

	o := defaultOptions(opts...)
	endTokenTypes := tokens[0].EndTokenTypes()
	if endTokenTypes == nil {
		return nil
	}

	return &Retriever{
		TokenSource:   tokens,
		endTokenTypes: endTokenTypes,
		result:        make([]group.Reindenter, 0, len(tokens)),
		options:       o,
	}
}

// Retrieve Retrieves a group of SQL clauses.
//
// It returns groups of SQL clauses as a slice of Reintenters.
//
// The returned endIdx indicates the new offset.
func (r *Retriever) Retrieve() ([]group.Reindenter, int, error) {
	if err := r.appendGroupsToResult(); err != nil {
		return nil, -1, fmt.Errorf("appendGroupsToResult failed: %w", err)
	}

	return r.result, r.endIdx, nil
}

// appendGroupsToResult appends token to result as Reindenter until endTokenType appears.
//
// If a subGroup is found in the target group, subGroup will be appended to result as a Reindenter, calling itself recursively.
//
// It returns an error if it cannot find any endTokenTypes.
func (r *Retriever) appendGroupsToResult() error {
	for idx := 0; ; {
		if idx >= len(r.TokenSource) {
			return fmt.Errorf("the retriever couldn't find the endToken for clause started with: %q",
				r.TokenSource[0].Value,
			)
		}
		token := r.TokenSource[idx]

		if r.isEndGroup(token, idx) {
			// TODO(fred): SHOULD ADD END, ) HERE? rather than processing this above
			r.endIdx = idx

			return nil
		}

		subGroupRetriever, offset := r.getSubGroupRetriever(idx)
		if subGroupRetriever == nil {
			r.result = append(r.result, token)
			idx += offset

			continue
		}

		if !subGroupRetriever.containsEndToken() {
			return fmt.Errorf("sub group clause started with %q has no end key word",
				subGroupRetriever.TokenSource[0].Value,
			)
		}

		if err := subGroupRetriever.appendGroupsToResult(); err != nil {
			return err
		}

		group, err := subGroupRetriever.subGroupResult()
		if err != nil {
			return err
		}

		r.result = append(r.result, group)

		idx = subGroupRetriever.getNextTokenIdx(token.Type, idx)
	}
}

// getSubGroupRetriever creates Retriever to retrieve sub group in the target group starting from tokens sliced from idx.
func (r *Retriever) getSubGroupRetriever(idx int) (*Retriever, int) {
	// when idx is equal to 0, target group itself will be Subgroup, which causes an error
	if idx == 0 {
		return nil, 1
	}

	const rangeOfJoinGroupStart = 3
	firstToken := r.TokenSource[0]
	previousToken := r.TokenSource[idx-1]
	token := r.TokenSource[idx]
	nextToken := r.TokenSource[idx+1] // should always work: trailed by EOF token

	afterComma := idx > 0 && r.TokenSource[idx-1].Type == lexer.COMMA
	afterParenthesis := idx > 0 && r.TokenSource[idx-1].Type == lexer.STARTPARENTHESIS
	_, isJoin := lexer.JoinMakers()[token.Type]
	_, isGroup := lexer.GroupMakers()[token.Type]

	switch {
	case r.containIrregularGroupMaker(firstToken, token, nextToken):
		fmt.Printf("DEBUG: irregular group maker in clause starting with %q|... at %q|%q\n",
			firstToken.Value, token.Value, nextToken.Value)

		return nil, 1

		/*
			case nextToken.Type == lexer.OPERATOR && nextToken.Value == "::":
				// TODO: cast_operator_group
		*/

	case token.Type == lexer.STARTPARENTHESIS && previousToken.Type == lexer.FUNCTION:
		return nil, 1

	case token.Type == lexer.FUNCTION && nextToken.Type == lexer.STARTPARENTHESIS:
		fmt.Printf("DEBUG: function group maker in clause starting with %q|... at %q|%q, afterParenthesis: %t\n",
			firstToken.Value,
			token.Value,
			nextToken.Value,
			afterParenthesis,
		)

		return NewRetriever(r.TokenSource[idx:], withOptions(r.options),
			withAfterComma(afterComma),
			withAfterParenthesis(afterParenthesis),
		), 2

	case token.Type == lexer.STARTPARENTHESIS && nextToken.Type == lexer.SELECT:
		subR := NewRetriever(r.TokenSource[idx:], withOptions(r.options),
			withAfterComma(afterComma),
			withAfterParenthesis(afterParenthesis),
		)
		if subR == nil {
			return nil, 1
		}

		subR.indentLevel = r.indentLevel

		// if subquery is found, indentLevel of all tokens until ")" will be incremented
		subR.indentLevel++

		return subR, 1

	case isJoin && idx < rangeOfJoinGroupStart:
		return nil, 1

	case isJoin, isGroup:
		// if group keywords appears in start of join group such as LEFT INNER JOIN, those keywords will be ignored
		// In this case, "INNER" and "JOIN" are group keyword, but should not make subGroup
		subR := NewRetriever(r.TokenSource[idx:], withOptions(r.options),
			withAfterComma(afterComma),
			withAfterParenthesis(afterParenthesis),
		)
		if subR == nil {
			return nil, 1
		}

		subR.indentLevel = r.indentLevel

		return subR, 1

	default:
		return nil, 1
	}
}

// TODO: this should be cleaned up - there is nothing irregular in these constructs

func (r *Retriever) containIrregularGroupMaker(firstToken, token, nextToken lexer.Token) bool {
	ttype := token.Type

	// in order not to make ORDER BY subGroup in Function group
	// this is a solution of window function
	if firstToken.Type == lexer.FUNCTION && ttype == lexer.ORDER {
		return true
	}

	// in order to ignore "(" in TypeCast group
	if firstToken.Type == lexer.TYPE && ttype == lexer.STARTPARENTHESIS {
		return true
	}

	// in order to ignore ORDER BY in window function
	if firstToken.Type == lexer.STARTPARENTHESIS && ttype == lexer.ORDER {
		return true
	}

	if firstToken.Type == lexer.FUNCTION && ttype == lexer.FROM {
		return true
	}

	if ttype == lexer.TYPE && !(nextToken.Type == lexer.STARTPARENTHESIS) {
		return true
	}

	return false
}

// if group key words to make join group such as "LEFT" or "OUTER" appear within idx is in range of join group, any keyword must be ignored not be made into a sub group.
func (r *Retriever) isRangeOfJoinStart(idx int) bool {
	const joinStartRange = 3
	if idx >= joinStartRange {
		return false
	}

	_, ok := lexer.JoinMakers()[r.TokenSource[0].Type]

	return ok
}

// subGroupResult makes a Reindenter from a subgroup.
func (r *Retriever) subGroupResult() (group.Reindenter, error) {
	subGroup, err := r.createGroup(r.result)
	if err != nil {
		return nil, err
	}

	if subGroup == nil {
		return nil, fmt.Errorf("can not make sub group from result :%#v", r.result)
	}

	subGroup.IncrementIndentLevel(r.indentLevel)

	return subGroup, nil
}

// getNextTokenIdx prepares idx for next token value.
func (r *Retriever) getNextTokenIdx(ttype lexer.TokenType, idx int) int {
	// if subGroup is PARENTHESIS group or CASE group, endIdx will be index of "END" or ")".
	//
	// In this case, next token must start after those end keyword, so it adds 1 to idx
	//
	// TODO: should not that be the same with [...]?

	switch ttype {
	case lexer.STARTPARENTHESIS, lexer.CASE, lexer.FUNCTION, lexer.TYPE:
		idx += r.endIdx + 1
	default:
		idx += r.endIdx
	}

	return idx
}

// createGroup creates each clause group from slice of tokens, returning it as Reindenter interface.
func (r *Retriever) createGroup(tokenSource []group.Reindenter) (group.Reindenter, error) {
	if len(tokenSource) == 0 {
		return nil, errors.New("empty token source passed to createGroup")
	}

	firstToken, ok := tokenSource[0].(lexer.Token)
	if !ok {
		return nil, errors.New("expected token to be passed to createGroup")
	}

	group := r.getGroupFromTokens(firstToken, tokenSource)

	return group, nil
}

func (r *Retriever) isEndToken(token lexer.Token) bool {
	if token.Type == lexer.EOF {
		return true
	}

	if _, ok := r.endTokenTypes[token.Type]; ok {
		return true
	}

	return false
}

// check tokens contain endTokenType.
func (r *Retriever) containsEndToken() bool {
	for _, token := range r.TokenSource {
		if r.isEndToken(token) {
			return true
		}
	}

	return false
}

// isEndGroup determines if token is the end token.
func (r *Retriever) isEndGroup(token lexer.Token, idx int) bool {
	// ignore endTokens when first token type is equal to endTokenType
	// because first token type might be a endTokenType.
	// For example "AND","OR"
	// isRangeOfJoinStart ignores if endTokenType appears in start of Join clause such as LEFT OUTER JOIN, INNER JOIN etc ...
	if idx == 0 || r.isRangeOfJoinStart(idx) {
		return false
	}

	return r.isEndToken(token)
}
