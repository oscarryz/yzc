package internal

import "testing"

var l = 1

func tok(tt tokenType) token {
	l++
	return token{
		pos:  position{l, 1},
		tt:   tt,
		data: "1",
	}
}
func Test_parser_exploreExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []token
		want   bool
	}{
		{
			"Int literal",
			[]token{tok(INTEGER), tok(EOF)},
			true,
		},
		{
			"Closing brace",
			[]token{tok(RBRACE), tok(EOF)},
			false,
		},
		{
			"Int literal",
			[]token{tok(INTEGER), tok(RBRACE), tok(EOF)},
			true,
		},
		{
			"Parenthesis expression",
			[]token{tok(LPAREN), tok(INTEGER), tok(RPAREN), tok(EOF)},
			true,
		},
		//{
		//	"Block Body",
		//	[]token{tok(LBRACE), tok(INTEGER), tok(RBRACE), tok(EOF)},
		//	true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				tokens: tt.tokens,
			}
			if got := p.exploreExpression(); got != tt.want {
				t.Errorf("exploreExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
