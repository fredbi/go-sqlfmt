package group

import (
	"bytes"
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Lock clause.
type Lock struct {
	elementReindenter
}

func NewLock(element []Reindenter, opts ...Option) *Lock {
	return &Lock{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindent its elements.
func (l *Lock) Reindent(buf *bytes.Buffer) error {
	elements, err := l.processPunctuation()
	if err != nil {
		return err
	}

	return l.elementsTokenApply(elements, buf, l.writeLock)
}

func (l *Lock) writeLock(buf *bytes.Buffer, token lexer.Token, _ Reindenter, _ int) error {
	switch token.Type {
	case lexer.LOCK, lexer.IN:
		buf.WriteString(fmt.Sprintf("%s%s", NewLine, token.FormattedValue()))
	case lexer.CASTOPERATOR:
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
