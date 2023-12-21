package internal

import (
	"reflect"
	"testing"
)

func tok(tt tokenType) token {
	return token{pos(1, 1), tt, "1"}
}
func Test_parser_program(t *testing.T) {
	EOF := tok(EOF)
	tests := []struct {
		name    string
		tokens  []token
		want    *program
		wantErr bool
	}{
		// TODO: Add test cases.

		{"Empty file",
			[]token{EOF},
			&program{&blockBody{
				[]expression{},
				[]statement{},
			},
			}, false,
		},
		{"Test literals",
			[]token{tok(INTEGER), EOF},
			&program{&blockBody{
				[]expression{
					&BasicLit{pos(1, 1), INTEGER, "1"},
				},
				[]statement{},
			}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser("a.yz", tt.tokens)
			got, err := p.program()
			if (err != nil) != tt.wantErr {
				t.Errorf("program() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("program() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_parser_exploreExpression(t *testing.T) {

	type fields struct {
		tokens []token
	}
	tests := []struct {
		name   string
		tokens []token
		want   bool
	}{
		{
			"Integer",
			[]token{tok(INTEGER), tok(EOF)},
			true,
		}, {
			"Decimal",
			[]token{tok(DECIMAL), tok(EOF)},
			true,
		}, {
			"String",
			[]token{tok(STRING), tok(EOF)},
			true,
		}, {
			"Identifier",
			[]token{tok(IDENTIFIER), tok(EOF)},
			true,
		}, {
			"Statement",
			[]token{tok(BREAK), tok(EOF)},
			false,
		}, {
			"Member access",
			[]token{
				tok(IDENTIFIER),
				tok(PERIOD),
				tok(IDENTIFIER),
				tok(EOF)},
			true,
		}, {
			"Dictionary access",
			[]token{
				tok(IDENTIFIER),
				tok(LBRACE),
				tok(INTEGER),
				tok(COLON),
				tok(STRING),
				tok(RBRACE),
				tok(EOF)},
			true,
		},
		{
			"Array access",
			[]token{
				tok(IDENTIFIER),
				tok(LBRACE),
				tok(INTEGER),
				tok(RBRACE),
				tok(EOF)},
			true,
		},
		//{
		//	"Literal_non_word_invocation",
		//	[]token{
		//		tok(STRING),
		//		tok(NONWORDIDENTIFIER),
		//		tok(STRING),
		//		tok(EOF)},
		//	true,
		//},
		{
			"Parenthesized_expressions",
			[]token{
				tok(LPAREN),
				tok(STRING),
				tok(RPAREN),
				tok(EOF)},
			true,
		},
		{
			"Method_invocation",
			[]token{
				tok(IDENTIFIER),
				tok(PERIOD),
				tok(IDENTIFIER),
				tok(LPAREN),
				tok(RPAREN),
				tok(EOF)},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser("a.yz", tt.tokens)
			if got := p.exploreExpression(); got != tt.want {
				t.Errorf("exploreExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
