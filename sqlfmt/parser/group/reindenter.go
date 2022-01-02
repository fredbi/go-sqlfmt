package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

const (
	NewLine          = "\n"
	WhiteSpace       = " "
	DoubleWhiteSpace = "  "
)

var (
	_ Reindenter = &elementReindenter{}
	_ Reindenter = &AndGroup{}
	_ Reindenter = &Case{}
	_ Reindenter = &Delete{}
	_ Reindenter = &From{}
	_ Reindenter = &Function{}
	_ Reindenter = &GroupBy{}
	_ Reindenter = &Having{}
	_ Reindenter = &Insert{}
	_ Reindenter = &Join{}
	_ Reindenter = &LimitClause{}
	_ Reindenter = &Lock{}
	_ Reindenter = &OrderBy{}
	_ Reindenter = &OrGroup{}
	_ Reindenter = &Parenthesis{}
	_ Reindenter = &Returning{}
	_ Reindenter = &Select{}
	_ Reindenter = &Set{}
	_ Reindenter = &Subquery{}
	_ Reindenter = &TieClause{}
	_ Reindenter = &TypeCast{}
	_ Reindenter = &Update{}
	_ Reindenter = &Values{}
	_ Reindenter = &Where{}
	_ Reindenter = &With{}
)

type (
	// Reindenter interface
	// specific values of Reindenter would be clause group or token.
	Reindenter interface {
		Reindent(*bytes.Buffer) error
		IncrementIndentLevel(int)
		GetStart() int
	}

	baseReindenter struct {
		start int
		*options
	}

	elementReindenter struct {
		Element []Reindenter
		baseReindenter
	}
)

func (g baseReindenter) GetStart() int {
	return g.start
}

// IncrementIndentLevel increments by its specified indent level.
func (g *baseReindenter) IncrementIndentLevel(lev int) {
	g.IndentLevel += lev
}

// writeComma writes a comma token with different indentation styles.
func (g *baseReindenter) writeComma(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch g.commaStyle {
	case CommaStyleRight:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s%s",
			token.FormattedValue(),
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			WhiteSpace,
		))
	default:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s%s",
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			DoubleWhiteSpace,
			token.FormattedValue(),
		))
	}
}

func (g *baseReindenter) write(buf *bytes.Buffer, token lexer.Token, indent int) error {
	switch {
	case token.IsNeedNewLineBefore():
		buf.WriteString(
			fmt.Sprintf(
				"%s%s%s",
				NewLine,
				strings.Repeat(DoubleWhiteSpace, indent),
				token.FormattedValue(),
			))
	case token.Type == lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.DO:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s",
			NewLine,
			token.FormattedValue(),
			WhiteSpace,
		))
	case token.Type == lexer.WITH:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			NewLine,
			token.FormattedValue(),
		))
	case token.Type == lexer.CASTOPERATOR, token.Type == lexer.WS:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.TYPE && g.hasCastBefore:
		buf.WriteString(token.FormattedValue())
	default:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			WhiteSpace,
			token.FormattedValue(),
		))
	}

	return nil
}

func (g *baseReindenter) writeWithComma(buf *bytes.Buffer, token lexer.Token, indent int) error {
	columnCount := g.start
	defer func() {
		g.start = columnCount
	}()

	switch {
	case token.IsNeedNewLineBefore():
		buf.WriteString(fmt.Sprintf(
			"%s%s%s",
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			token.FormattedValue()),
		)
	case token.Type == lexer.BY:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			WhiteSpace,
			token.FormattedValue()),
		)
	case token.Type == lexer.COMMA:
		g.writeComma(buf, token, indent)
	case token.Type == lexer.CASTOPERATOR, token.Type == lexer.WS:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.TYPE && g.hasCastBefore:
		buf.WriteString(token.FormattedValue())
	default:
		return fmt.Errorf("can not reindent %q", token.FormattedValue())
	}

	return nil
}

func newElementReindenter(element []Reindenter, opts ...Option) elementReindenter {
	o := defaultOptions(opts...)

	return elementReindenter{
		Element: element,
		baseReindenter: baseReindenter{
			options: o,
		},
	}
}

// Reindent reindents its elements.
func (e *elementReindenter) Reindent(buf *bytes.Buffer) error {
	elements, err := e.processPunctuation()
	if err != nil {
		return err
	}

	return e.elementsTokenApply(elements, buf, e.write)
}

func (e *elementReindenter) processPunctuation() ([]Reindenter, error) {
	elements, err := processPunctuation(e.Element)
	if err != nil {
		return nil, err
	}

	return elements, nil
}

func (e *elementReindenter) elementsTokenApply(
	elements []Reindenter,
	buf *bytes.Buffer,
	tokenWriter func(*bytes.Buffer, lexer.Token, int) error,
) error {
	for _, el := range elements {
		switch token := el.(type) {
		case lexer.Token:
			if err := tokenWriter(buf, token, e.IndentLevel); err != nil {
				return err
			}
		default:
			if err := el.Reindent(buf); err != nil {
				return err
			}
		}
	}

	return nil
}
