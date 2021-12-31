package lexer

import (
	"bytes"
	"fmt"
)

// Token is a token struct.
type Token struct {
	Type  TokenType
	Value string
	*options
}

// MakeToken builds an immutable token.
func MakeToken(ttype TokenType, value string, opts ...Option) Token {
	o := defaultOptions(opts...)

	return Token{
		Type:    ttype,
		Value:   value,
		options: o,
	}
}

// Reindent is a placeholder for implementing Reindenter interface.
func (t Token) Reindent(buf *bytes.Buffer) error { return nil }

// GetStart is a placeholder for implementing Reindenter interface.
func (t Token) GetStart() int { return 0 }

// IncrementIndentLevel is a placeholder implementing Reindenter interface.
func (t Token) IncrementIndentLevel(lev int) {}

func (t Token) formatKeyword() string {
	if t.options == nil {
		return t.Value
	}
	in := t.Value

	switch t.Type {
	case STRING, IDENT:
		// no op
	case RESERVEDVALUE:
		switch t.Value {
		case "NAN":
			in = "NaN"
		case "INFINITY":
			in = "Infinity"
		case "-INFINITY":
			in = "-Infinity"
		}
	case FUNCTION:
		recased, ok := casedFunctions.Get([]byte(t.Value))
		if ok {
			in = recased.(string)
		}
	}

	if t.recaser != nil {
		in = t.recaser(in)
	}

	if t.colorizer != nil {
		in = t.colorizer(t.Type, in)
	}

	return in
}

func (t Token) formatPunctuation() string {
	if t.Type == SEMICOLON {
		return fmt.Sprintf("%s%s", NewLine, t.Value)
	}

	return t.Value
}

// FormattedValue returns the token with some formatting options.
func (t Token) FormattedValue() string {
	switch t.Type {
	case EOF,
		WS,
		NEWLINE,
		COMMA,
		SEMICOLON,
		STARTPARENTHESIS,
		ENDPARENTHESIS,
		STARTBRACKET,
		ENDBRACKET,
		STARTBRACE,
		ENDBRACE:
		// ANDGROUP,
		// ORGROUP:
		return t.formatPunctuation()
	default:
		return t.formatKeyword()
	}
}

// IsNeedNewLineBefore returns true if token needs new line before written in buffer.
func (t Token) IsNeedNewLineBefore() bool {
	var ttypes = []TokenType{
		SELECT, UPDATE, INSERT, DELETE,
		ANDGROUP,
		FROM, GROUP, ORGROUP,
		ORDER, HAVING, LIMIT, OFFSET, FETCH, RETURNING,
		SET, UNION, INTERSECT, EXCEPT, VALUES,
		WHERE, ON, USING, UNION, EXCEPT, INTERSECT,
	}
	for _, v := range ttypes {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsKeyWordInSelect returns true if token is a keyword in select group.
func (t Token) IsKeyWordInSelect() bool {
	return t.Type == SELECT ||
		t.Type == EXISTS ||
		t.Type == DISTINCT ||
		t.Type == DISTINCTROW ||
		t.Type == INTO ||
		t.Type == AS ||
		t.Type == GROUP ||
		t.Type == ORDER ||
		t.Type == BY ||
		t.Type == ON ||
		t.Type == RETURNING ||
		t.Type == SET ||
		t.Type == UPDATE
}
