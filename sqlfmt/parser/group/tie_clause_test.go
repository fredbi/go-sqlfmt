package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentUnionGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case1",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.UNION, Value: "UNION"},
				lexer.Token{Type: lexer.ALL, Value: "ALL"},
			},
			want: "\nUNION ALL",
		},
		{
			name: "normal case2",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.INTERSECT, Value: "INTERSECT"},
			},
			want: "\nINTERSECT",
		},
		{
			name: "normal case3",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.EXCEPT, Value: "EXCEPT"},
			},
			want: "\nEXCEPT",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		unionGroup := NewTieClause(tt.tokenSource)

		if err := unionGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
