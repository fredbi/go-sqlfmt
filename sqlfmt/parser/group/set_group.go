package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Set clause.
type Set struct {
	elementReindenter
}

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

	for _, el := range separate(elements) {
		switch v := el.(type) {
		case lexer.Token:
			if erw := s.writeWithComma(buf, v, s.IndentLevel); erw != nil {
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
