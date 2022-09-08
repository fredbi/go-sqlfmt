package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/stretchr/testify/require"
)

func TestReindentSelectGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
				lexer.Token{Type: lexer.IDENT, Value: "name"},
				lexer.Token{Type: lexer.COMMA, Value: ","},
				lexer.Token{Type: lexer.IDENT, Value: "age"},
			},
			want: "\nSELECT\n  name\n  , age",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		selectGroup := NewSelect(tt.tokenSource)

		require.NoError(t, selectGroup.Reindent(buf))

		got := buf.String()
		require.Equal(t, tt.want, got)
	}
}

func TestIncrementIndentLevel(t *testing.T) {
	s := NewSelect(nil)
	s.IncrementIndentLevel(1)
	require.Equal(t, 1, s.IndentLevel)
}
