package internal

import "testing"

var l = 1

func tok(tt tokenType) token {
	l++
	return token{
		pos:  position{l, 1},
		tt:   tt,
		data: "1i",
	}
}
func Test_parser_exploreExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []token
		want   ruleType
	}{
		{
			"Int literal",
			[]token{tok(INTEGER), tok(EOF)},
			INTEGER_EXPR,
		},
		{
			"Closing brace",
			[]token{tok(RBRACE), tok(EOF)},
			ILLEGAL,
		},
		{
			"Input:\" 1 }\" will return Integer Expr as it stops i the first `1` and shouldn't process the `}`",
			[]token{tok(INTEGER), tok(RBRACE), tok(EOF)},
			INTEGER_EXPR,
		},
		{
			"Parenthesis expression",
			[]token{tok(LPAREN), tok(INTEGER), tok(RPAREN), tok(EOF), tok(RBRACE)},
			PARENTHESIS_EXPR,
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
