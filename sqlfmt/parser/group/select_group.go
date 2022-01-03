package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Select clause.
type Select struct {
	elementReindenter
}

// NewSelect group
func NewSelect(element []Reindenter, opts ...Option) *Select {
	return &Select{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindens its elements.
func (s *Select) Reindent(buf *bytes.Buffer) error {
	s.start = 0

	elements, err := s.processPunctuation()
	if err != nil {
		return err
	}

	indenters := separate(elements)
	for i, element := range indenters {
		var previous Reindenter
		if i > 0 {
			previous = indenters[i-1]
		}

		switch v := element.(type) {
		case lexer.Token:
			s.writeSelect(buf, v, previous, &s.start, s.IndentLevel)
		case *Case:
			if tok, ok := elements[i-1].(lexer.Token); ok && tok.Type == lexer.COMMA {
				v.hasCommaBefore = true
			}

			if eri := v.Reindent(buf); eri != nil {
				return eri
			}

			// Case group in Select clause must be in column area
			s.start++
		case *Parenthesis:
			v.InColumnArea = true
			v.ColumnCount = s.start
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
			s.start++
		case *Subquery:
			if token, ok := elements[i-1].(lexer.Token); ok {
				if token.Type == lexer.EXISTS {
					if eri := v.Reindent(buf); eri != nil {
						return eri
					}

					continue
				}
			}
			v.InColumnArea = true
			v.ColumnCount = s.start
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
		case *Function:
			v.InColumnArea = true
			v.ColumnCount = s.start
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
			s.start++
		case Reindenter:
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
			s.start++
		default:
			return fmt.Errorf("can not reindent %#v", v)
		}
	}

	return nil
}

func (s *Select) writeSelect(buf *bytes.Buffer, token lexer.Token, previous Reindenter, start *int, indent int) {
	columnCount := *start
	defer func() {
		*start = columnCount
	}()

	switch token.Type {
	case lexer.SELECT, lexer.INTO:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s",
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			token.FormattedValue(),
		))
	case lexer.AS, lexer.DISTINCT, lexer.DISTINCTROW, lexer.GROUP, lexer.ON:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			WhiteSpace,
			token.FormattedValue(),
		))
	case lexer.EXISTS:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			WhiteSpace,
			token.FormattedValue(),
		))
		columnCount++
	case lexer.CASTOPERATOR, lexer.WS:
		buf.WriteString(token.FormattedValue())
	case lexer.COMMA:
		s.writeComma(buf, token, indent)
	case lexer.TYPE:
		if s.hasCastBefore || isCastOperator(previous) {
			buf.WriteString(token.FormattedValue())

			break
		}

		_ = s.write(buf, token, previous, indent)
	default:
		_ = s.write(buf, token, previous, indent)
	}
}
