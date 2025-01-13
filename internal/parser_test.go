package internal

import (
	"fmt"
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
			path: []string{"empty"},
			tokens: []Token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "empty",
							varType: newBocType(),
						},
						value: &Boc{
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
			path: []string{"parent", "simple"},
			tokens: []Token{
				{pos(0, 0), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "parent",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos: pos(0, 0),
									variable: &Variable{
										pos:     pos(0, 0),
										name:    "simple",
										varType: newBocType(),
									},
									value: &Boc{
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
		},
		{
			name: "Literal expressions",
			path: []string{"literals"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "literals",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									INTEGER,
									"1",
									&IntType{},
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
			path: []string{"string_literal"},
			tokens: []Token{
				{pos(1, 1), STRING, "Hello world"},
				{pos(1, 12), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "string_literal",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									STRING,
									"Hello world",
									&StringType{},
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
			path: []string{"block_literal"},
			tokens: []Token{
				{pos(1, 1), LBRACE, "{"},
				{pos(1, 2), RBRACE, "}"},
				{pos(1, 3), EOF, "EOF"},
			},
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "block_literal",
							varType: newBocType(),
						},
						value: &Boc{
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
			path: []string{"two_literals"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), COMMA, ","},
				{pos(1, 3), STRING, "Hello world"},
				{pos(1, 14), EOF, "EOF"},
			}, want: &Boc{

				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "two_literals",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									INTEGER,
									"1",
									&IntType{},
								},
								&BasicLit{
									pos(1, 3),
									STRING,
									"Hello world",
									&StringType{},
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
			path: []string{"invalid_expression"},
			tokens: []Token{
				{pos(1, 1), INTEGER, "1"},
				{pos(1, 2), INTEGER, "2"},
				{pos(1, 3), EOF, "EOF"},
			}, wantErr: true,
			errorMessage: "[line: 1 col: 2] expected \",\" or \"}\". Got \"2\"",
		},
		{
			name: "Two literals with new line",
			path: []string{"two_literals_newline"},
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
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "two_literals_newline",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&BasicLit{
									pos(1, 1),
									INTEGER,
									"1",
									&IntType{},
								},
								&BasicLit{
									pos(2, 1),
									STRING,
									"Hello world",
									&StringType{},
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
			path: []string{"array_literal"},
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
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_literal",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{},
									&ArrayType{elemType: &IntType{}},
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
			parents: []string{"two_literals"},
			source: `[] Int
"Hello"`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "two_literals",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{},
									&ArrayType{elemType: &IntType{}},
								},
								&BasicLit{
									pos(2, 1),
									STRING,
									"Hello",
									&StringType{},
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
			parents: []string{"array_literal"},
			source:  `[1, 2, 3]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_literal",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{
										&BasicLit{
											pos(1, 2),
											INTEGER,
											"1",
											&IntType{},
										},
										&BasicLit{
											pos(1, 5),
											INTEGER,
											"2",
											&IntType{},
										},
										&BasicLit{
											pos(1, 8),
											INTEGER,
											"3",
											&IntType{},
										},
									},
									&ArrayType{elemType: &IntType{}},
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
			name:    "Array of arrays [[1 2]] is an [][Int]",
			parents: []string{"array_of_arrays_2"},
			source:  `[[1, 2],[1, 2]]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_of_arrays_2",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{
										&ArrayLit{
											pos(1, 2),
											[]expression{
												&BasicLit{
													pos(1, 3),
													INTEGER,
													"1",
													&IntType{},
												},
												&BasicLit{
													pos(1, 6),
													INTEGER,
													"2",
													&IntType{},
												},
											},
											&ArrayType{elemType: &IntType{}}},
										&ArrayLit{
											pos(1, 9),
											[]expression{
												&BasicLit{
													pos(1, 10),
													INTEGER,
													"1",
													&IntType{},
												},
												&BasicLit{
													pos(1, 13),
													INTEGER,
													"2",
													&IntType{},
												},
											},
											&ArrayType{elemType: &IntType{}}},
									},
									&ArrayType{elemType: &IntType{}},
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
			parents: []string{"array_of_arrays"},
			source:  `[[1, 2], []Int]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_of_arrays",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{
										&ArrayLit{
											pos(1, 2),
											[]expression{
												&BasicLit{
													pos(1, 3),
													INTEGER,
													"1",
													&IntType{},
												},
												&BasicLit{
													pos(1, 6),
													INTEGER,
													"2",
													&IntType{},
												},
											},
											&ArrayType{elemType: &IntType{}}},
										&ArrayLit{
											pos(1, 10),
											[]expression{},
											&ArrayType{elemType: &IntType{}}},
									},
									&ArrayType{elemType: &IntType{}},
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
			parents: []string{"array_of_blocks"},
			source:  `[{1}, {2}]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_of_blocks",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{
										&Boc{
											expressions: []expression{
												&BasicLit{
													pos(1, 3),
													INTEGER,
													"1",
													&IntType{},
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
													&IntType{},
												},
											},
											statements: []statement{},
										},
									},
									&ArrayType{elemType: newBocType()},
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
			parents: []string{"empty_dictionary_literal"},
			source:  `[String]Int`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "empty_dictionary_literal",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&DictLit{
									pos(1, 1),
									&DictType{
										keyType: &StringType{},
										valType: &IntType{},
									},
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
			parents: []string{"dictionary_literal"},
			source:  `[k1:v1, k2:v2]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "dictionary_literal",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&DictLit{
									pos(1, 1),
									&DictType{
										keyType: newTBD(),
										valType: newTBD(),
									},
									[]expression{
										&Variable{
											pos(1, 2),
											"k1",
											newTBD(),
										},
										&Variable{
											pos(1, 9),
											"k2",
											newTBD(),
										},
									},
									[]expression{
										&Variable{
											pos(1, 5),
											"v1",
											newTBD(),
										},
										&Variable{
											pos(1, 12),
											"v2",
											newTBD(),
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
			parents: []string{"dictionary_literal_type"},
			source: `[
    "name": ["Yz"]
    "type system": ["static", "strong", "structural"]
]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "dictionary_literal_type",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&DictLit{
									pos(1, 1),
									&DictType{
										keyType: &StringType{},
										valType: &ArrayType{elemType: &StringType{}},
									},
									[]expression{
										&BasicLit{
											pos(2, 5),
											STRING,
											"name",
											&StringType{},
										},
										&BasicLit{
											pos(3, 5),
											STRING,
											"type system",
											&StringType{},
										},
									},
									[]expression{
										&ArrayLit{
											pos(2, 13),

											[]expression{
												&BasicLit{
													pos(2, 14),
													STRING,
													"Yz",
													&StringType{},
												},
											},
											&ArrayType{elemType: &StringType{}},
										},
										&ArrayLit{
											pos(3, 20),
											[]expression{
												&BasicLit{
													pos(3, 21),
													STRING,
													"static",
													&StringType{},
												},
												&BasicLit{
													pos(3, 31),
													STRING,
													"strong",
													&StringType{},
												},
												&BasicLit{
													pos(3, 41),
													STRING,
													"structural",
													&StringType{},
												},
											},
											&ArrayType{elemType: &StringType{}},
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
			parents: []string{"short_declaration"},
			source:  `a : 1`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "short_declaration",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&Variable{
										pos(1, 1),
										"a",
										&IntType{},
									},
									&BasicLit{
										pos(1, 5),
										INTEGER,
										"1",
										&IntType{},
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
			parents: []string{"short_declaration_block_array"},
			source: `language: {
		   name: "Yz"
		   features: ["static", "strong", "structural"]
		}`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "short_declaration_block_array",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&Variable{
										pos:     pos(1, 1),
										name:    "language",
										varType: newBocType(),
									},
									&Boc{
										expressions: []expression{
											&ShortDeclaration{
												pos(2, 6),
												&Variable{
													pos(2, 6),
													"name",
													&StringType{},
												},
												&BasicLit{
													pos(2, 12),
													STRING,
													"Yz",
													&StringType{},
												},
											},
											&ShortDeclaration{
												pos(3, 6),
												&Variable{
													pos:     pos(3, 6),
													name:    "features",
													varType: &ArrayType{elemType: &StringType{}},
												}, &ArrayLit{
													pos(3, 16),
													[]expression{
														&BasicLit{
															pos(3, 17),
															STRING,
															"static",
															&StringType{},
														},
														&BasicLit{
															pos(3, 27),
															STRING,
															"strong",
															&StringType{},
														},
														&BasicLit{
															pos(3, 37),
															STRING,
															"structural",
															&StringType{},
														},
													},
													&ArrayType{elemType: &StringType{}},
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
			parents: []string{"closing"},
			source: `dictionary: [
		"ready" : false
]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "closing",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&Variable{
										pos(1, 1),
										"dictionary",
										&DictType{
											keyType: &StringType{},
											valType: newTBD(),
										},
									},
									&DictLit{
										pos(1, 13),
										&DictType{
											keyType: &StringType{},
											valType: newTBD(),
										},
										[]expression{
											&BasicLit{
												pos(2, 3),
												STRING,
												"ready",
												&StringType{},
											},
										},
										[]expression{
											&Variable{
												pos(2, 13),
												"false",
												newTBD(),
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
			parents: []string{"short_declaration_literal_array_dictionary"},
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
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "short_declaration_literal_array_dictionary",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&Variable{
										pos(1, 1),
										"main",
										newBocType(),
									},
									&Boc{
										expressions: []expression{
											&ShortDeclaration{
												pos(2, 9),
												&Variable{
													pos(2, 9),
													"msg",
													&StringType{},
												},
												&BasicLit{
													pos(2, 14),
													STRING,
													"Hello",
													&StringType{},
												},
											},
											&ShortDeclaration{
												pos(3, 9),
												&Variable{
													pos(3, 9),
													"array",
													&ArrayType{elemType: &IntType{}},
												},
												&ArrayLit{
													pos(3, 16),
													[]expression{
														&BasicLit{
															pos(3, 17),
															INTEGER,
															"1",
															&IntType{},
														},
														&BasicLit{
															pos(3, 20),
															INTEGER,
															"2",
															&IntType{},
														},
														&BasicLit{
															pos(3, 23),
															INTEGER,
															"3",
															&IntType{},
														},
													},
													&ArrayType{elemType: &IntType{}},
												},
											},
											&ShortDeclaration{
												pos(4, 9),
												&Variable{
													pos(4, 9),
													"dictionary",
													&DictType{
														keyType: &StringType{},
														valType: newTBD(),
													},
												},
												&DictLit{
													pos(4, 21),
													&DictType{
														keyType: &StringType{},
														valType: newTBD(),
													},
													[]expression{
														&BasicLit{
															pos(5, 17),
															STRING,
															"ready",
															&StringType{},
														},
													},
													[]expression{
														&Variable{
															pos(5, 26),
															"false",
															newTBD(),
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
			name:    "Array of dictionaries",
			parents: []string{"array_of_dictionaries"},
			source: `[
	[
		"ready" : false
	],
	[
		"done" : true
	]
]`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "array_of_dictionaries",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ArrayLit{
									pos(1, 1),
									[]expression{
										&DictLit{
											pos(2, 2),
											&DictType{
												keyType: &StringType{},
												valType: newTBD(),
											},
											[]expression{
												&BasicLit{
													pos(3, 3),
													STRING,
													"ready",
													&StringType{},
												},
											},
											[]expression{
												&Variable{
													pos(3, 13),
													"false",
													newTBD(),
												},
											},
										},
										&DictLit{
											pos(5, 2),
											&DictType{
												keyType: &StringType{},
												valType: newTBD(),
											},
											[]expression{
												&BasicLit{
													pos(6, 3),
													STRING,
													"done",
													&StringType{},
												},
											}, []expression{
												&Variable{
													pos(6, 12),
													"true",
													newTBD(),
												},
											},
										},
									},
									&ArrayType{elemType: &DictType{keyType: &StringType{}, valType: newTBD()}},
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
			name:    "Variable",
			parents: []string{"variable"},
			source:  `a: 1`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "variable",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&ShortDeclaration{
									pos(1, 1),
									&Variable{
										pos(1, 1),
										"a",
										&IntType{},
									},
									&BasicLit{
										pos(1, 4),
										INTEGER,
										"1",
										&IntType{},
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
			name:    "Variable",
			parents: []string{"variable"},
			source:  `a`,
			want: &Boc{
				expressions: []expression{
					&ShortDeclaration{
						pos: pos(0, 0),
						variable: &Variable{
							pos:     pos(0, 0),
							name:    "variable",
							varType: newBocType(),
						},
						value: &Boc{
							expressions: []expression{
								&Variable{
									pos:     pos(1, 1),
									name:    "a",
									varType: newTBD(),
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
				t.Logf("Want AST: %s", tt.want)
				fmt.Printf("%s", diff)
				t.Errorf("diff = %v", diff)
			}
		})
	}
}
