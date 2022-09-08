package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// GroupBy clause.
//
// nolint:revive
type GroupBy struct {
	elementReindenter
}

// NewGroupBy group;
func NewGroupBy(element []Reindenter, opts ...Option) *GroupBy {
	return &GroupBy{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (g *GroupBy) Reindent(buf *bytes.Buffer) error {
	g.start = 0

	elements, err := g.processPunctuation()
	if err != nil {
		return err
	}

	reindenters := separate(elements)
	for i, el := range reindenters {
		var previous Reindenter
		if i > 0 {
			previous = reindenters[i-1]
		}
		switch v := el.(type) {
		case lexer.Token:
			if erw := g.writeWithComma(buf, v, previous, g.IndentLevel); erw != nil {
				return erw
			}
		default:
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}
