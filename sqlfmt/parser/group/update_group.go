package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Update clause.
type Update struct {
	elementReindenter
}

// NewUpdate clause group
func NewUpdate(element []Reindenter, opts ...Option) *Update {
	return &Update{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (u *Update) Reindent(buf *bytes.Buffer) error {
	u.start = 0

	elements, err := processPunctuation(u.Element)
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
			if erw := u.writeWithComma(buf, v, previous, u.IndentLevel); erw != nil {
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
