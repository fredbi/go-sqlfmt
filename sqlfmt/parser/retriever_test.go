package parser

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/stretchr/testify/require"
)

func TestNewRetriever(t *testing.T) {
	testingData := []lexer.Token{
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "name"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.IDENT, Value: "age"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "user"},
		{Type: lexer.EOF, Value: "EOF"},
	}

	r := NewRetriever(testingData)
	want := []lexer.Token{
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "name"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.IDENT, Value: "age"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "user"},
		{Type: lexer.EOF, Value: "EOF"},
	}

	got := r.TokenSource
	require.EqualValues(t, want, got, "initialize retriever failed: want %#v got %#v", want, got)
}

func TestRetrieve(t *testing.T) {
	type want struct {
		stmt    []string
		lastIdx int
	}

	tests := []struct {
		name          string
		source        []lexer.Token
		endTokenTypes []lexer.TokenType
		want          *want
	}{
		{
			name: "normal_test",
			source: []lexer.Token{
				{Type: lexer.SELECT, Value: "SELECT"},
				{Type: lexer.IDENT, Value: "name"},
				{Type: lexer.COMMA, Value: ","},
				{Type: lexer.IDENT, Value: "age"},
				{Type: lexer.FROM, Value: "FROM"},
				{Type: lexer.IDENT, Value: "user"},
				{Type: lexer.EOF, Value: "EOF"},
			},
			endTokenTypes: []lexer.TokenType{lexer.FROM},
			want: &want{
				stmt:    []string{"SELECT", "name", ",", "age"},
				lastIdx: 4,
			},
		},
		{
			name: "normal_test3",
			source: []lexer.Token{
				{Type: lexer.LEFT, Value: "LEFT"},
				{Type: lexer.JOIN, Value: "JOIN"},
				{Type: lexer.IDENT, Value: "xxx"},
				{Type: lexer.ON, Value: "ON"},
				{Type: lexer.IDENT, Value: "xxx"},
				{Type: lexer.IDENT, Value: "="},
				{Type: lexer.IDENT, Value: "xxx"},
				{Type: lexer.WHERE, Value: "WHERE"},
			},
			endTokenTypes: []lexer.TokenType{lexer.WHERE},
			want: &want{
				stmt:    []string{"LEFT", "JOIN", "xxx", "ON", "xxx", "=", "xxx"},
				lastIdx: 7,
			},
		},
		{
			name: "normal_test4",
			source: []lexer.Token{
				{Type: lexer.UPDATE, Value: "UPDATE"},
				{Type: lexer.IDENT, Value: "xxx"},
				{Type: lexer.SET, Value: "SET"},
			},
			endTokenTypes: []lexer.TokenType{lexer.SET},
			want: &want{
				stmt:    []string{"UPDATE", "xxx"},
				lastIdx: 2,
			},
		},
		{
			name: "normal_test5",
			source: []lexer.Token{
				{Type: lexer.INSERT, Value: "INSERT"},
				{Type: lexer.INTO, Value: "INTO"},
				{Type: lexer.IDENT, Value: "xxx"},
				{Type: lexer.VALUES, Value: "VALUES"},
			},
			endTokenTypes: []lexer.TokenType{lexer.VALUES},
			want: &want{
				stmt:    []string{"INSERT", "INTO", "xxx"},
				lastIdx: 3,
			},
		},
		{
			name: "case statement, with comma",
			source: []lexer.Token{
				{Type: lexer.CASE, Value: "SELECT"},
				{Type: lexer.IDENT, Value: "a"},
				{Type: lexer.COMMA, Value: ","},
				{Type: lexer.CASE, Value: "CASE"},
				{Type: lexer.WHEN, Value: "WHEN"},
				{Type: lexer.RESERVEDVALUE, Value: "TRUE"},
				{Type: lexer.THEN, Value: "THEN"},
				{Type: lexer.RESERVEDVALUE, Value: "FALSE"},
				{Type: lexer.END, Value: "END"},
				{Type: lexer.EOF, Value: "EOF"},
			},
			want: &want{
				stmt:    []string{"SELECT", "a", ",", "\n  CASE\n     WHEN TRUE THEN FALSE\n  END"},
				lastIdx: 9,
			},
		},
		{
			name: "case statement, without comma",
			source: []lexer.Token{
				{Type: lexer.CASE, Value: "SELECT"},
				{Type: lexer.CASE, Value: "CASE"},
				{Type: lexer.WHEN, Value: "WHEN"},
				{Type: lexer.RESERVEDVALUE, Value: "TRUE"},
				{Type: lexer.THEN, Value: "THEN"},
				{Type: lexer.RESERVEDVALUE, Value: "FALSE"},
				{Type: lexer.END, Value: "END"},
				{Type: lexer.EOF, Value: "EOF"},
			},
			want: &want{
				stmt:    []string{"SELECT", "\n  CASE\n     WHEN TRUE THEN FALSE\n  END"},
				lastIdx: 7,
			},
		},
	}

	for _, toPin := range tests {
		tt := toPin
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			var (
				gotStmt    []string
				gotLastIdx int
			)
			r := NewRetriever(tt.source)
			require.NotNil(t, r)

			reindenters, gotLastIdx, err := r.Retrieve()
			require.NoError(t, err)

			for _, v := range reindenters {
				if tok, ok := v.(lexer.Token); ok {
					gotStmt = append(gotStmt, tok.Value)
				} else {
					var buf bytes.Buffer
					require.NoError(t, v.Reindent(&buf))
					gotStmt = append(gotStmt, buf.String())
				}
			}

			require.EqualValues(t, tt.want.stmt, gotStmt)
			require.Equal(t, tt.want.lastIdx, gotLastIdx)
		})
	}
}
