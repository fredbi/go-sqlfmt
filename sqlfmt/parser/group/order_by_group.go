package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// OrderBy clause.
type OrderBy struct {
	elementReindenter
}

func NewOrderBy(element []Reindenter, opts ...Option) *OrderBy {
	return &OrderBy{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (o *OrderBy) Reindent(buf *bytes.Buffer) error {
	o.start = 0

	elements, err := o.processPunctuation()
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
			if erw := o.writeWithComma(buf, v, previous, o.IndentLevel); erw != nil {
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
