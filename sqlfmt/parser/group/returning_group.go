package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Returning clause.
type Returning struct {
	elementReindenter
}

// NewReturning group
func NewReturning(element []Reindenter, opts ...Option) *Returning {
	return &Returning{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (r *Returning) Reindent(buf *bytes.Buffer) error {
	elements, err := r.processPunctuation()
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
			if erw := r.writeWithComma(buf, v, previous, r.IndentLevel); erw != nil {
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
