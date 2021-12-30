package parser

import (
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// Retriever retrieves target SQL clause group from TokenSource.
type Retriever struct {
	TokenSource   []lexer.Token
	result        []group.Reindenter
	indentLevel   int
	endTokenTypes []lexer.TokenType
	endIdx        int

	*options
}

// NewRetriever Creates Retriever that retrieves each target SQL clause.
//
// Each Retriever have endKeywords in order to stop retrieving.
func NewRetriever(tokenSource []lexer.Token, opts ...Option) *Retriever {
	o := defaultOptions(opts...)
	firstTokenType := tokenSource[0].Type

	switch firstTokenType {
	case lexer.SELECT:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfSelect, options: o}
	case lexer.FROM:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfFrom, options: o}
	case lexer.CASE:
		fmt.Printf("DEBUG: NewRetriever CASE, afterComma: %t\n", o.afterComma)

		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfCase, options: o}
	case lexer.JOIN, lexer.INNER, lexer.OUTER, lexer.LEFT, lexer.RIGHT, lexer.NATURAL, lexer.CROSS:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfJoin, options: o}
	case lexer.WHERE:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfWhere, options: o}
	case lexer.ANDGROUP:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfAndGroup, options: o}
	case lexer.ORGROUP:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfOrGroup, options: o}
	case lexer.GROUP:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfGroupBy, options: o}
	case lexer.HAVING:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfHaving, options: o}
	case lexer.ORDER:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfOrderBy, options: o}
	case lexer.LIMIT, lexer.FETCH, lexer.OFFSET:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfLimitClause, options: o}
	case lexer.STARTPARENTHESIS:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfParenthesis, options: o}
	case lexer.UNION, lexer.INTERSECT, lexer.EXCEPT:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfTieClause, options: o}
	case lexer.UPDATE:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfUpdate, options: o}
	case lexer.SET:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfSet, options: o}
	case lexer.RETURNING:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfReturning, options: o}
	case lexer.DELETE:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfDelete, options: o}
	case lexer.INSERT:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfInsert, options: o}
	case lexer.VALUES:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfValues, options: o}
	case lexer.FUNCTION:
		fmt.Printf("DEBUG: NewRetriever FUNCTION, afterComma: %t\n", o.afterComma)

		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfFunction, options: o}
	case lexer.TYPE:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfTypeCast, options: o}
	case lexer.LOCK:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfLock, options: o}
	case lexer.WITH:
		return &Retriever{TokenSource: tokenSource, endTokenTypes: lexer.EndOfWith, options: o}
	default:
		return nil
	}
}

// Retrieve Retrieves a group of SQL clauses.
//
// It returns clause group as slice of Reintenter interface and endIdx for setting offset.
func (r *Retriever) Retrieve() ([]group.Reindenter, int, error) {
	if err := r.appendGroupsToResult(); err != nil {
		return nil, -1, errors.Wrap(err, "appendGroupsToResult failed")
	}

	return r.result, r.endIdx, nil
}

// appendGroupsToResult appends token to result as Reindenter until endTokenType appears.
//
// If a subGroup is found in the target group, subGroup will be appended to result as a Reindenter, calling itself recursively.
//
// It returns an error if it cannot find any endTokenTypes.
func (r *Retriever) appendGroupsToResult() error {
	var (
		idx   int
		token lexer.Token
	)

	for {
		if idx >= len(r.TokenSource) {
			return fmt.Errorf("the retriever couldn't find the endToken for clause started with: %q",
				r.TokenSource[0].Value,
			)
		}

		token = r.TokenSource[idx]

		if r.isEndGroup(token, idx) {
			// TODO(fred): SHOUD ADD END, ) HERE? rather than processing this above
			r.endIdx = idx

			return nil
		}

		if subGroupRetriever := r.getSubGroupRetriever(idx); subGroupRetriever != nil {
			if !containsEndToken(subGroupRetriever.TokenSource, subGroupRetriever.endTokenTypes) {
				return fmt.Errorf("sub group clause started with %q has no end key word",
					subGroupRetriever.TokenSource[0].Value,
				)
			}

			if err := subGroupRetriever.appendGroupsToResult(); err != nil {
				return err
			}

			if err := r.appendSubGroupToResult(subGroupRetriever.result, subGroupRetriever.indentLevel); err != nil {
				return err
			}

			idx = subGroupRetriever.getNextTokenIdx(token.Type, idx)

			continue
		}

		r.result = append(r.result, token)
		idx++
	}
}

// check tokens contain endTokenType.
func containsEndToken(tokenSource []lexer.Token, endTokenTypes []lexer.TokenType) bool {
	for _, tok := range tokenSource {
		for _, endttype := range endTokenTypes {
			if tok.Type == endttype {
				return true
			}
		}
	}

	return false
}

// isEndGroup determines if token is the end token.
func (r *Retriever) isEndGroup(token lexer.Token, idx int) bool {
	for _, endTokenType := range r.endTokenTypes {
		// ignore endTokens when first token type is equal to endTokenType
		// because first token type might be a endTokenType.
		// For example "AND","OR"
		// isRangeOfJoinStart ignores if endTokenType appears in start of Join clause such as LEFT OUTER JOIN, INNER JOIN etc ...
		if idx == 0 || r.isRangeOfJoinStart(idx) {
			return false
		}

		if token.Type == endTokenType || token.Type == lexer.EOF {
			return true
		}
	}

	return false
}

// getSubGroupRetriever creates Retriever to retrieve sub group in the target group starting from tokens sliced from idx.
func (r *Retriever) getSubGroupRetriever(idx int) *Retriever {
	// when idx is equal to 0, target group itself will be Subgroup, which causes an error
	if idx == 0 {
		return nil
	}

	firstToken := r.TokenSource[0]
	token := r.TokenSource[idx]
	nextToken := r.TokenSource[idx+1] // should always work: trailed by EOF token

	if r.containIrregularGroupMaker(firstToken, token, nextToken) {
		fmt.Printf("DEBUG: irregular group maker in clause starting with %q|... at %q|%q\n", firstToken.Value, token.Value, nextToken.Value)

		return nil
	}

	afterComma := idx > 0 && r.TokenSource[idx-1].Type == lexer.COMMA

	if token.Type == lexer.STARTPARENTHESIS && nextToken.Type == lexer.SELECT {
		subR := NewRetriever(r.TokenSource[idx:], withOptions(r.options), withAfterComma(afterComma))
		if subR == nil {
			return nil
		}

		subR.indentLevel = r.indentLevel

		// if subquery is found, indentLevel of all tokens until ")" will be incremented
		subR.indentLevel++

		return subR
	}

	if token.IsJoinStart() {
		// if group keywords appears in start of join group such as LEFT INNER JOIN, those keywords will be ignored
		// In this case, "INNER" and "JOIN" are group keyword, but should not make subGroup
		rangeOfJoinGroupStart := 3
		if idx < rangeOfJoinGroupStart {
			return nil
		}
		subR := NewRetriever(r.TokenSource[idx:], withOptions(r.options), withAfterComma(afterComma))
		if subR == nil {
			return nil
		}

		subR.indentLevel = r.indentLevel

		return subR
	}

	for _, v := range lexer.TokenTypesOfGroupMaker {
		if token.Type == v {
			subR := NewRetriever(r.TokenSource[idx:], withOptions(r.options), withAfterComma(afterComma))
			if subR == nil {
				return nil
			}

			subR.indentLevel = r.indentLevel

			return subR
		}
	}

	return nil
}

// TODO: this should be cleaned up - there is nothing irregular in these constructs
//
// func (r *Retriever) containIrregularGroupMaker(ttype lexer.TokenType, idx int) bool {
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

	if firstToken.Type == lexer.FUNCTION && (ttype == lexer.STARTPARENTHESIS || ttype == lexer.FROM) {
		return true
	}

	if ttype == lexer.TYPE && !(nextToken.Type == lexer.STARTPARENTHESIS) {
		return true
	}
	/* TODO
	 */

	return false
}

// if group key words to make join group such as "LEFT" or "OUTER" appear within idx is in range of join group, any keyword must be ignored not be made into a sub group.
func (r *Retriever) isRangeOfJoinStart(idx int) bool {
	firstTokenType := r.TokenSource[0].Type
	for _, v := range lexer.TokenTypesOfJoinMaker {
		joinStartRange := 3
		if v == firstTokenType && idx < joinStartRange {
			return true
		}
	}

	return false
}

// appendSubGroupToResult makes Reindenter from subGroup result and append it to result.
func (r *Retriever) appendSubGroupToResult(result []group.Reindenter, lev int) error {
	subGroup, err := r.createGroup(result)
	if err != nil {
		return err
	}
	if subGroup == nil {
		return fmt.Errorf("can not make sub group result :%#v", result)
	}

	subGroup.IncrementIndentLevel(lev)
	r.result = append(r.result, subGroup)

	return nil
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
		return nil, errors.New("expected tokens to be passed to createGroup")
	}

	switch firstToken.Type {
	case lexer.SELECT:
		return group.NewSelect(tokenSource, r.ToGroupOptions()...), nil
	case lexer.FROM:
		return group.NewFrom(tokenSource, r.ToGroupOptions()...), nil
	case lexer.JOIN, lexer.INNER, lexer.OUTER, lexer.LEFT, lexer.RIGHT, lexer.NATURAL, lexer.CROSS, lexer.LATERAL:
		return group.NewJoin(tokenSource, r.ToGroupOptions()...), nil
	case lexer.WHERE:
		return group.NewWhere(tokenSource, r.ToGroupOptions()...), nil
	case lexer.ANDGROUP:
		return group.NewAndGroup(tokenSource, r.ToGroupOptions()...), nil
	case lexer.ORGROUP:
		return group.NewOrGroup(tokenSource, r.ToGroupOptions()...), nil
	case lexer.GROUP:
		return group.NewGroupBy(tokenSource, r.ToGroupOptions()...), nil
	case lexer.ORDER:
		return group.NewOrderBy(tokenSource, r.ToGroupOptions()...), nil
	case lexer.HAVING:
		return group.NewHaving(tokenSource, r.ToGroupOptions()...), nil
	case lexer.LIMIT, lexer.OFFSET, lexer.FETCH:
		return group.NewLimitClause(tokenSource, r.ToGroupOptions()...), nil
	case lexer.UNION, lexer.INTERSECT, lexer.EXCEPT:
		return group.NewTieClause(tokenSource, r.ToGroupOptions()...), nil
	case lexer.UPDATE:
		return group.NewUpdate(tokenSource, r.ToGroupOptions()...), nil
	case lexer.SET:
		return group.NewSet(tokenSource, r.ToGroupOptions()...), nil
	case lexer.RETURNING:
		return group.NewReturning(tokenSource, r.ToGroupOptions()...), nil
	case lexer.LOCK:
		return group.NewLock(tokenSource, r.ToGroupOptions()...), nil
	case lexer.INSERT:
		return group.NewInsert(tokenSource, r.ToGroupOptions()...), nil
	case lexer.VALUES:
		return group.NewValues(tokenSource, r.ToGroupOptions()...), nil
	case lexer.DELETE:
		return group.NewDelete(tokenSource, r.ToGroupOptions()...), nil
	case lexer.WITH:
		return group.NewWith(tokenSource, r.ToGroupOptions()...), nil
	case lexer.CASE:
		// endKeyWord of CASE group("END") has to be included in the group, so it is appended to result
		endToken := lexer.MakeToken(lexer.END, "END", lexer.WithOptionsFrom(firstToken))
		tokenSource = append(tokenSource, endToken)

		return group.NewCase(tokenSource, r.ToGroupOptions()...), nil
	case lexer.STARTPARENTHESIS:
		// endKeyWord of subQuery group (")") has to be included in the group, so it is appended to result
		endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(firstToken))
		tokenSource = append(tokenSource, endToken)

		if _, isSubQuery := tokenSource[1].(*group.Select); isSubQuery {
			return group.NewSubquery(tokenSource, r.ToGroupOptions()...), nil
		}

		return group.NewParenthesis(tokenSource, r.ToGroupOptions()...), nil
	case lexer.FUNCTION:
		fmt.Printf("DEBUG: createGroup: add ) then NewFunction group\n")
		endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(firstToken))
		tokenSource = append(tokenSource, endToken)

		return group.NewFunction(tokenSource, r.ToGroupOptions()...), nil
	case lexer.TYPE:
		endToken := lexer.MakeToken(lexer.ENDPARENTHESIS, ")", lexer.WithOptionsFrom(firstToken))
		tokenSource = append(tokenSource, endToken)

		return group.NewTypeCast(tokenSource, r.ToGroupOptions()...), nil
	}

	return nil, nil
}
