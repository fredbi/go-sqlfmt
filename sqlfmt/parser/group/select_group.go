package group

import (
	"bytes"
	"fmt"
	"io"
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

func (s *Select) writeIndent(w io.Writer, element Reindenter) {
	switch v := element.(type) {
	case lexer.Token:
		fmt.Fprintf(
			w,
			"%s%s",
			strings.Repeat(WhiteSpace, s.indentSize*s.indentLevel),
			token.FormattedValue(),
		)
	default:
		element.Reindent(w)
	}
}

// Reindent reindents the elements in a SELECT statement.
//
// {SELECT ...} {FROM}
// {SELECT ...} {UNION|UNION ALL|INTERSECT|EXCEPT}
// {SELECT ...} {EOF}
// (SELECT ...}) -> see Parenthesis group
//
// Individual fields are preceded by a line feed.
// Individual fields are indented.
func (s *Select) Reindent(buf *bytes.Buffer) error {
	if len(s.Element) == 0 {
		return nil
	}

	s.writeIndent(buf, s.Element[0])
	linefeed(buf)

	if len(s.Element) == 1 {
		return nil
	}

	s.start = 0

	s.writeIndent(buf, s.Element[1])
	if len(s.Element) == 2 {
		return nil
	}

	for i, element := range s.Element[2:] {
		previous := s.Element[i+1]

		if !isComma(previous) {
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
