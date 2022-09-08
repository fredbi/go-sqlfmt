package parser

import "github.com/fredbi/go-sqlfmt/sqlfmt/lexer"

// tokensContext holds the list of tokens being parsed and knows how to provide some context.
type tokensContext struct {
	TokenSource []lexer.Token
	*options
}

func makeTokensContext(tokens []lexer.Token, opts ...Option) tokensContext {
	o := defaultOptions(opts...)

	return tokensContext{
		TokenSource: tokens,
		options:     o,
	}
}

// isAfterComma provides context about the current group relative to a comma.
func (h *tokensContext) isAfterComma(idx int) bool {
	return idx > 0 && h.TokenSource[idx-1].Type == lexer.COMMA
}

// isAfterParenthesis provides context about the current group relative to an opening parenthesis.
func (h *tokensContext) isAfterParenthesis(idx int) bool {
	return idx > 0 && h.TokenSource[idx-1].Type == lexer.STARTPARENTHESIS
}

// isAfterCast provides context about the current group relative to a cast operator '::'
func (h *tokensContext) isAfterCast(idx int) bool {
	return idx > 0 && h.TokenSource[idx-1].Type == lexer.CASTOPERATOR
}
