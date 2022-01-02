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

	for i, element := range separate(elements) {
		switch v := element.(type) {
		case lexer.Token:
			if erw := s.writeSelect(buf, v, &s.start, s.IndentLevel); erw != nil {
				return fmt.Errorf("writeSelect failed: %w", erw)
			}
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

func (s *Select) writeSelect(buf *bytes.Buffer, token lexer.Token, start *int, indent int) error {
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
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case lexer.EXISTS:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
		columnCount++
	case lexer.CASTOPERATOR, lexer.WS:
		buf.WriteString(token.FormattedValue())
	case lexer.COMMA:
		s.writeComma(buf, token, indent)
	case lexer.TYPE:
		if s.hasCastBefore {
			buf.WriteString(token.FormattedValue())

			break
		}

		s.write(buf, token, indent)
	default:
		s.write(buf, token, indent)
	}

	return nil
}
