package group

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// separate elements by comma and the reserved keywords in a select clause.
func separate(rs []Reindenter) []Reindenter {
	var (
		skipRange, count int
		result, buf      []Reindenter
	)

	flushBuf := func() {
		if len(buf) > 0 {
			result = append(result, buf...)
			buf = buf[:0]
		}
	}

	for _, r := range rs {
		switch token := r.(type) {
		case lexer.Token:
			switch {
			case skipRange > 0:
				skipRange--

			case token.IsKeywordInSelect():
				flushBuf()
				result = append(result, token)

			case token.Type == lexer.COMMA:
				flushBuf()
				result = append(result, token)
				count = 0

			default:
				if count == 0 {
					buf = append(buf, token)
				} else {
					// NULLIFY THIS FOR THE MOMENT: CAN'T EASILY DEMONSTRATE WHERE THIS COMES IN PLAY
					// buf = append(buf, lexer.MakeToken(lexer.WS, "#" /*WhiteSpace*/))
					buf = append(buf, token)
				}

				count++
			}

		default:
			flushBuf()
			result = append(result, r)
		}
	}

	flushBuf()

	return result
}

// process bracket, singlequote and brace.
//
// TODO: more elegant.
//
// TODO(fred): neutered out - for the moment, can't easily demonstrate what this is useful for.
func processPunctuation(rs []Reindenter) ([]Reindenter, error) {
	/*
		var (
			result    []Reindenter
			skipRange int
		)

		for i, v := range rs {
			token, ok := v.(lexer.Token)
			if !ok {
				result = append(result, v)

				continue
			}

			// simple token

			fmt.Printf("DEBUG: processPunctuation: token [%d]:%q\n", i, token.Value)
			switch {
			case skipRange > 0:
				skipRange--
			case token.Type == lexer.STARTBRACE || token.Type == lexer.STARTBRACKET:
					surrounding, sr, err := extractSurroundingArea(rs[i:])
					if err != nil {
						return nil, err
					}
						result = append(result, lexer.Token{
							Type:  lexer.SURROUNDING,
							Value: surrounding,
						})
					skipRange += sr
			default:
				result = append(result, token)
			}
		}

		return result, nil
	*/
	return rs, nil
}

// returns surrounding area including punctuation such as {xxx, xxx}.
func extractSurroundingArea(rs []Reindenter) (string, int, error) {
	var (
		countOfStart int
		countOfEnd   int
		result       string
		skipRange    int
	)

	for i, r := range rs {
		if token, ok := r.(lexer.Token); ok {
			switch {
			case token.Type == lexer.COMMA || token.Type == lexer.STARTBRACKET || token.Type == lexer.STARTBRACE || token.Type == lexer.ENDBRACKET || token.Type == lexer.ENDBRACE:
				result += fmt.Sprint(token.FormattedValue())
				// for next token of StartToken
			case i == 1:
				result += fmt.Sprint(token.FormattedValue())
			default:
				result += fmt.Sprint(WhiteSpace + token.FormattedValue())
			}

			if token.Type == lexer.STARTBRACKET || token.Type == lexer.STARTBRACE || token.Type == lexer.STARTPARENTHESIS {
				countOfStart++
			}
			if token.Type == lexer.ENDBRACKET || token.Type == lexer.ENDBRACE || token.Type == lexer.ENDPARENTHESIS {
				countOfEnd++
			}
			if countOfStart == countOfEnd {
				break
			}
			skipRange++
		} else {
			// TODO: should support group type in surrounding area?
			// I have not encountered any groups in surrounding area so far
			return "", -1, errors.New("group type is not supposed be here")
		}
	}

	return result, skipRange, nil
}
