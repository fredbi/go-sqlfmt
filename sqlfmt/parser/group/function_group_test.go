//nolint:dupl
package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentFunctionGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.FUNCTION, Value: "SUM"},
				lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
				lexer.Token{Type: lexer.IDENT, Value: "xxx"},
				lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
			},
			want: " SUM(xxx)",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		functionGroup := NewFunction(tt.tokenSource)

		if err := functionGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
