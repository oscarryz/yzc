package internal

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		tokens       []token
		want         *boc
		wantErr      bool
		errorMessage string
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
		{
			name:     "Invalid expression expression",
			fileName: "invalid_expression.yz",
			tokens: []token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), INTEGER, "2"},
				{pos(1, 3), EOF, "EOF"},
			}, wantErr: true,
			errorMessage: "[invalid_expression.yz line: 1 col: 2] expected ,",
		},
		{
			name:     "Two literals with new line",
			fileName: "two_literals_newline.yz",
			tokens: []token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), NEWLINE, "\n"},
				{pos(2, 1), STRING, "Hello world"},
				{pos(2, 12), EOF, "EOF"},
			}, want: &boc{
				Name:    "two_literals_newline",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&BasicLit{
							pos(1, 1),
							INTEGER,
							"1",
						},
						&BasicLit{
							pos(2, 1),
							STRING,
							"Hello world",
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Array literal",
			fileName: "array_literal.yz",
			tokens: []token{
				{pos(1, 1), LBRACKET, "["},
				{pos(1, 2), RBRACKET, "]"},
				{pos(1, 3), TYPEIDENTIFIER, "Int"},
				{pos(1, 6), EOF, "EOF"},
			}, want: &boc{
				Name:    "array_literal",
				bocType: nil,
				blockBody: &blockBody{
					statements: []statement{},
					expressions: []expression{
						&ArrayLit{
							pos(1, 1),
							"[]",
							[]expression{},
						},
					},
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
			if tt.wantErr && err.Error() != tt.errorMessage {
				t.Errorf("Parse() error = \"%v\", wantErr %v", err, tt.errorMessage)
				return
			}
			// compare go recursively
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v\n want %v", got, tt.want)
			}
		})
	}
}

// test tokenizer + parse
func TestParse_TokenizeAndParse(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		source       string
		want         *boc
		wantErr      bool
		errorMessage string
	}{
		{
			name:     "Two literals",
			fileName: "two_literals.yz",
			source: `[] Int
"Hello"`,
			want: &boc{
				Name:    "two_literals",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&ArrayLit{
							pos(1, 1),
							"[]",
							[]expression{},
						},
						&BasicLit{
							pos(2, 1),
							STRING,
							"Hello",
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Array literal [1 2 3]",
			fileName: "array_literal.yz",
			source:   `[1 2 3]`,
			want: &boc{
				Name:    "array_literal",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&ArrayLit{
							pos(1, 1),
							"[]",
							[]expression{
								&BasicLit{
									pos(1, 2),
									INTEGER,
									"1",
								},
								&BasicLit{
									pos(1, 4),
									INTEGER,
									"2",
								},
								&BasicLit{
									pos(1, 6),
									INTEGER,
									"3",
								},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Array of arrays	[[1 2] []Int]",
			fileName: "array_of_arrays.yz",
			source:   `[[1 2] []Int]`,
			want: &boc{
				Name:    "array_of_arrays",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&ArrayLit{
							pos(1, 1),
							"[]",
							[]expression{
								&ArrayLit{
									pos(1, 2),
									"[]",
									[]expression{
										&BasicLit{
											pos(1, 3),
											INTEGER,
											"1",
										},
										&BasicLit{
											pos(1, 5),
											INTEGER,
											"2",
										},
									},
								},
								&ArrayLit{
									pos(1, 8),
									"[]",
									[]expression{},
								},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Array of blocks",
			fileName: "array_of_blocks.yz",
			source:   `[{1} {2}]`,
			want: &boc{
				Name:    "array_of_blocks",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&ArrayLit{
							pos(1, 1),
							"[]",
							[]expression{
								&boc{
									Name: "",
									blockBody: &blockBody{
										expressions: []expression{
											&BasicLit{

												pos(1, 3),
												INTEGER,
												"1",
											},
										},
										statements: []statement{},
									},
								},
								&boc{
									Name: "",
									blockBody: &blockBody{
										expressions: []expression{

											&BasicLit{
												pos(1, 7),
												INTEGER,

												"2",
											},
										},
										statements: []statement{},
									},
								},
							},
						},
					},
					statements: []statement{},
				},
			},
		},
		{
			name:     "Empty dictionary literal [String]Int",
			fileName: "empty_dictionary_literal.yz",
			source:   `[String]Int`,
			want: &boc{
				Name:    "empty_dictionary_literal",
				bocType: nil,
				blockBody: &blockBody{
					expressions: []expression{
						&DictLit{
							pos(1, 1),
							"[]",
							"String",
							"Int",
							[]expression{},
							[]expression{},
						},
					},
					statements: []statement{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.fileName, tt.source)
			if err != nil {
				t.Errorf("Tokenize() error = \"%v\"", err)
				return
			}
			got, err := Parse(tt.fileName, tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = \"%v\", wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMessage {
				t.Errorf("Parse() error = \"%v\", wantErr %v", err, tt.errorMessage)
				return
			}
			// compare go recursively
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v\n want %v", got, tt.want)
			}
		})
	}
}
