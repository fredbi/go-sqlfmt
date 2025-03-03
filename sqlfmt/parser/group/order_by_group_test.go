package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentOrderByGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.ORDER, Value: "ORDER"},
				lexer.Token{Type: lexer.BY, Value: "BY"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
			},
			want: "\nORDER BY\n  xxxxxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		orderByGroup := NewOrderBy(tt.tokenSource)

		if err := orderByGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
