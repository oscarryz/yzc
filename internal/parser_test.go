package internal

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		tokens   []token
		want     *boc
		wantErr  bool
	}{
		{
			name:     "Empty file",
			fileName: "empty.yz",
			tokens: []token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &boc{
				Name:    "empty",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{&empty{}},
					statements:  []statement{},
				},
			},
		},
		{
			name:     "Nested directory",
			fileName: "parent/simple.yz",
			tokens: []token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &boc{
				Name:    "parent",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&boc{
							Name: "simple",
							blockBody: &blockBody{
								expressions: []expression{&empty{}},
								statements:  []statement{},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Literal expressions",
			fileName: "literals.yz",
			tokens: []token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), EOF, "EOF"},
			},
			want: &boc{
				Name:    "literals",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&BasicLit{
							pos(1, 1),
							INTEGER,
							"1",
						}},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Literal expressions string",
			fileName: "string_literal.yz",
			tokens: []token{
				{pos(1, 1), STRING, "Hello world"},
				{pos(1, 12), EOF, "EOF"},
			},
			want: &boc{
				Name:    "string_literal",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&BasicLit{
							pos(1, 1),
							STRING,
							"Hello world",
						}},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Block literal",
			fileName: "block_literal.yz",
			tokens: []token{
				{pos(1, 1), LBRACE, "{"},
				{pos(1, 2), RBRACE, "}"},
				{pos(1, 3), EOF, "EOF"},
			},
			want: &boc{
				Name:    "block_literal",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&boc{
							Name: "",
							blockBody: &blockBody{
								expressions: []expression{&empty{}},
								statements:  []statement{},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Two literals",
			fileName: "two_literals.yz",
			tokens: []token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), COMMA, ","},
				{pos(1, 3), STRING, "Hello world"},
				{pos(1, 14), EOF, "EOF"},
			}, want: &boc{
				Name:    "two_literals",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&BasicLit{
							pos(1, 1),
							INTEGER,
							"1",
						},
						&BasicLit{
							pos(1, 3),
							STRING,
							"Hello world",
						},
					},
					statements: []statement{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.fileName, tt.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = \"%v\", wantErr %v", err, tt.wantErr)
				return
			}
			// compare go recursively
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v\n want %v", got, tt.want)
			}
		})
	}
}
