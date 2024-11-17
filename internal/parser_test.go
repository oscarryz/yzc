package internal

import (
	"github.com/go-test/deep"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		path         []string
		tokens       []Token
		want         *Boc
		wantErr      bool
		errorMessage string
	}{
		{
			name: "Empty file",
			path: []string{"empty.yz"},
			tokens: []Token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "empty.yz",
						},
						val: &Boc{
							expressions: []expression{},
							statements:  []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name: "Nested directory",
			path: []string{"parent", "simple.yz"},
			tokens: []Token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "parent",
						},
						val: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos: pos(0, 0),
									key: &BasicLit{
										pos: pos(0, 0),
										tt:  IDENTIFIER,
										val: "simple.yz",
									},
									val: &Boc{
										expressions: []expression{},
										statements:  []statement{},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
			//want: &ShortDeclaration{
			//				pos(1, 1),
			//				&BasicLit{
			//					pos(1, 1),
			//					IDENTIFIER,
			//					"parent",
			//				},
			//				&Boc{
			//					expressions: []expression{},
			//					statements:  []statement{},
			//				},
			//},
			//want: &Boc{
			//
			//	expressions: []expression{
			//		&Boc{
			//			expressions: []expression{},
			//			statements:  []statement{},
			//		},
			//	},
			//	statements: []statement{},
			//},
		},
		{
			name: "Literal expressions",
			path: []string{"literals.yz"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "literals.yz",
						},
						val: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									INTEGER,
									"1",
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name: "Literal expressions string",
			path: []string{"string_literal.yz"},
			tokens: []Token{
				{pos(1, 1), STRING, "Hello world"},
				{pos(1, 12), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "string_literal.yz",
						},
						val: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									STRING,
									"Hello world",
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name: "Block literal",
			path: []string{"block_literal.yz"},
			tokens: []Token{
				{pos(1, 1), LBRACE, "{"},
				{pos(1, 2), RBRACE, "}"},
				{pos(1, 3), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "block_literal.yz",
						},
						val: &Boc{
							expressions: []expression{
								&Boc{
									expressions: []expression{},
									statements:  []statement{},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name: "Two literals",
			path: []string{"two_literals.yz"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), COMMA, ","},
				{pos(1, 3), STRING, "Hello world"},
				{pos(1, 14), EOF, "EOF"},
			}, want: &Boc{

				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "two_literals.yz",
						},
						val: &Boc{
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
				statements: []statement{},
			},
		},
		{
			name: "Invalid expression expression",
			path: []string{"invalid_expression.yz"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), INTEGER, "2"},
				{pos(1, 3), EOF, "EOF"},
			}, wantErr: true,
			errorMessage: "[line: 1 col: 2] expected \",\" or \"}\". Got \"2\"",
		},
		{
			name: "Two literals with new line",
			path: []string{"two_literals_newline.yz"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), COMMA, "\n"},
				{pos(2, 1), STRING, "Hello world"},
				{pos(2, 12), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "two_literals_newline.yz",
						},
						val: &Boc{
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
				statements: []statement{},
			},
		},
		{
			name: "Array literal []Int",
			path: []string{"array_literal.yz"},
			tokens: []Token{
				{pos(1, 1), LBRACKET, "["},
				{pos(1, 2), RBRACKET, "]"},
				{pos(1, 3), TYPE_IDENTIFIER, "Int"},
				{pos(1, 6), EOF, "EOF"},
			},
			want: &Boc{
				statements: []statement{},
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "array_literal.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									&BasicLit{
										pos(1, 3),
										TYPE_IDENTIFIER, "Int",
									},
									[]expression{},
								},
							},
							statements: []statement{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.path, tt.tokens)
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
		parents      []string
		source       string
		want         *Boc
		wantErr      bool
		errorMessage string
	}{
		{
			name:    "Two literals",
			parents: []string{"two_literals.yz"},
			source: `[] Int
"Hello"`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "two_literals.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									&BasicLit{
										pos(1, 4),
										TYPE_IDENTIFIER,
										"Int",
									},
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
				statements: []statement{},
			},
		},
		{
			name:    "Array literal [1, 2, 3] is a [Int]",
			parents: []string{"array_literal.yz"},
			source:  `[1, 2, 3]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "array_literal.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									&BasicLit{
										pos(1, 2),
										INTEGER,
										"1",
									},
									[]expression{
										&BasicLit{
											pos(1, 2),
											INTEGER,
											"1",
										},
										&BasicLit{
											pos(1, 5),
											INTEGER,
											"2",
										},
										&BasicLit{
											pos(1, 8),
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
				statements: []statement{},
			},
		},
		{
			name:    "Array of arrays [[1 2] []Int] is an [][Int]",
			parents: []string{"array_of_arrays.yz"},
			source:  `[[1, 2], []Int]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "array_of_arrays.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									&ArrayLit{
										pos(1, 2),
										&BasicLit{
											pos(1, 3),
											INTEGER, "1",
										},
										[]expression{}},
									[]expression{
										&ArrayLit{
											pos(1, 2),
											&BasicLit{
												pos(1, 3),
												INTEGER,
												"1",
											},
											[]expression{
												&BasicLit{
													pos(1, 3),
													INTEGER,
													"1",
												},
												&BasicLit{
													pos(1, 6),
													INTEGER,
													"2",
												},
											},
										},
										&ArrayLit{
											pos(1, 10),
											&BasicLit{
												pos(1, 12),
												TYPE_IDENTIFIER,
												"Int",
											},
											[]expression{},
										},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Array of blocks",
			parents: []string{"array_of_blocks.yz"},
			source:  `[{1}, {2}]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "array_of_blocks.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									&Boc{
										expressions: []expression{},
										statements:  []statement{},
									},
									[]expression{
										&Boc{
											expressions: []expression{
												&BasicLit{
													pos(1, 3),
													INTEGER,
													"1",
												},
											},
											statements: []statement{},
										},
										&Boc{
											expressions: []expression{
												&BasicLit{
													pos(1, 8),
													INTEGER,
													"2",
												},
											},
											statements: []statement{},
										},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Empty dictionary literal [String]Int",
			parents: []string{"empty_dictionary_literal.yz"},
			source:  `[String]Int`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "empty_dictionary_literal.yz",
						},
						val: &Boc{
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
				statements: []statement{},
			},
		},
		{
			name:    "Dictionary literal [k1:v1 k2:v2]",
			parents: []string{"dictionary_literal.yz"},
			source:  `[k1:v1, k2:v2]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "dictionary_literal.yz",
						},
						val: &Boc{
							expressions: []expression{
								&DictLit{
									pos(1, 1),
									"[]",
									"",
									"",
									[]expression{
										&BasicLit{
											pos(1, 2),
											IDENTIFIER,
											"k1",
										},
										&BasicLit{
											pos(1, 9),
											IDENTIFIER,
											"k2",
										},
									},
									[]expression{
										&BasicLit{
											pos(1, 5),
											IDENTIFIER,
											"v1",
										},
										&BasicLit{
											pos(1, 12),
											IDENTIFIER,
											"v2",
										},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Dictionary literal of type [String][String]",
			parents: []string{"dictionary_literal_type.yz"},
			source: `[
    "name": ["Yz"]
    "type system": ["static", "strong", "structural"]
]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "dictionary_literal_type.yz",
						},
						val: &Boc{
							expressions: []expression{
								&DictLit{
									pos(1, 1),
									"[]",
									"",
									"",
									[]expression{
										&BasicLit{
											pos(2, 5),
											STRING,
											"name",
										},
										&BasicLit{
											pos(3, 5),
											STRING,
											"type system",
										},
									},
									[]expression{
										&ArrayLit{
											pos(2, 13),
											&BasicLit{
												pos(2, 14),
												STRING,
												"Yz",
											},
											[]expression{
												&BasicLit{
													pos(2, 14),
													STRING,
													"Yz",
												},
											},
										},
										&ArrayLit{
											pos(3, 20),
											&BasicLit{
												pos(3, 21),
												STRING,
												"static",
											},
											[]expression{
												&BasicLit{
													pos(3, 21),
													STRING,
													"static",
												},
												&BasicLit{
													pos(3, 31),
													STRING,
													"strong",
												},
												&BasicLit{
													pos(3, 41),
													STRING,
													"structural",
												},
											},
										},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Short declaration",
			parents: []string{"short_declaration.yz"},
			source:  `a : 1`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "short_declaration.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&BasicLit{
										pos(1, 1),
										IDENTIFIER,
										"a",
									},
									&BasicLit{
										pos(1, 5),
										INTEGER,
										"1",
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Short declaration with block and array",
			parents: []string{"short_declaration_block_array.yz"},
			source: `language: {
		   name: "Yz"
		   features: ["static", "strong", "structural"]
		}`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "short_declaration_block_array.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&BasicLit{
										pos: pos(1, 1),
										tt:  IDENTIFIER,
										val: "language",
									},
									&Boc{
										expressions: []expression{
											&ShortDeclaration{
												pos(2, 6),
												&BasicLit{
													pos(2, 6),
													IDENTIFIER,
													"name",
												},
												&BasicLit{
													pos(2, 12),
													STRING,
													"Yz",
												},
											},
											&ShortDeclaration{
												pos(3, 6),
												&BasicLit{
													pos: pos(3, 6),
													tt:  IDENTIFIER,
													val: "features",
												}, &ArrayLit{
													pos(3, 16),
													&BasicLit{
														pos(3, 17),
														STRING,
														"static",
													},
													[]expression{
														&BasicLit{
															pos(3, 17),
															STRING,
															"static",
														},
														&BasicLit{
															pos(3, 27),
															STRING,
															"strong",
														},
														&BasicLit{
															pos(3, 37),
															STRING,
															"structural",
														},
													},
												},
											},
										},
										statements: []statement{},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Closing bracket",
			parents: []string{"closing.yz"},
			source: `dictionary: [
		"ready" : false
]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "closing.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&BasicLit{
										pos(1, 1),
										IDENTIFIER,
										"dictionary",
									},
									&DictLit{
										pos(1, 13),
										"[]",
										"",
										"",
										[]expression{
											&BasicLit{
												pos(2, 3),
												STRING,
												"ready",
											},
										},
										[]expression{
											&BasicLit{
												pos(2, 13),
												IDENTIFIER,
												"false",
											},
										},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
		{
			name:    "Short declaration with literal, array and dictionary",
			parents: []string{"short_declaration_literal_array_dictionary.yz"},
			source: `main: {
        msg: "Hello"
        array: [1, 2, 3 ]
        dictionary: [
                "ready" :false
        ]
}`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						key: &BasicLit{
							pos: pos(0, 0),
							tt:  IDENTIFIER,
							val: "short_declaration_literal_array_dictionary.yz",
						},
						val: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&BasicLit{
										pos(1, 1),
										IDENTIFIER,
										"main",
									},
									&Boc{
										expressions: []expression{
											&ShortDeclaration{
												pos(2, 9),
												&BasicLit{
													pos(2, 9),
													IDENTIFIER,
													"msg",
												},
												&BasicLit{
													pos(2, 14),
													STRING,
													"Hello",
												},
											},
											&ShortDeclaration{
												pos(3, 9),
												&BasicLit{
													pos(3, 9),
													IDENTIFIER,
													"array",
												},
												&ArrayLit{
													pos(3, 16),
													&BasicLit{
														pos(3, 17),
														INTEGER,
														"1",
													},
													[]expression{
														&BasicLit{
															pos(3, 17),
															INTEGER,
															"1",
														},
														&BasicLit{
															pos(3, 20),
															INTEGER,
															"2",
														},
														&BasicLit{
															pos(3, 23),
															INTEGER,
															"3",
														},
													},
												},
											},
											&ShortDeclaration{
												pos(4, 9),
												&BasicLit{
													pos(4, 9),
													IDENTIFIER,
													"dictionary",
												},
												&DictLit{
													pos(4, 21),
													"[]",
													"",
													"",
													[]expression{
														&BasicLit{
															pos(5, 17),
															STRING,
															"ready",
														},
													},
													[]expression{
														&BasicLit{
															pos(5, 26),
															IDENTIFIER,
															"false",
														},
													},
												},
											},
										},
										statements: []statement{},
									},
								},
							},
							statements: []statement{},
						},
					},
				},
				statements: []statement{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.parents, tt.source)
			if err != nil {
				t.Errorf("Tokenize() error = \"%v\"", err)
				return
			}
			got, err := Parse(tt.parents, tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s\nParse() error = \"%v\", wantErr %v", tt.source, err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMessage {
				t.Errorf("Parse() error = \"%v\", wantErr %v", err, tt.errorMessage)
				return
			}

			deep.CompareUnexportedFields = true
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("diff = %v", diff)
			}
		})
	}
}
