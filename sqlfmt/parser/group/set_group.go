package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Set clause.
type Set struct {
	elementReindenter
}

// NewSet group in update clause
func NewSet(element []Reindenter, opts ...Option) *Set {
	return &Set{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (s *Set) Reindent(buf *bytes.Buffer) error {
	s.start = 0

	elements, err := s.processPunctuation()
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
			if erw := s.writeWithComma(buf, v, previous, s.IndentLevel); erw != nil {
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
