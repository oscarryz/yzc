package internal

import (
	"fmt"
	"testing"
)

func TestTokenizer_Tokenize(t *testing.T) {
	tests := []struct {
		name    string
		path    []string
		content string
		want    []Token
	}{

		{
			"All the tokens",
			[]string{"test.yz"},
			`( ) { } [ ] , : ; . = # => when
1 1.0 "Hello, World!" 'Hello, name!' a Point + break continue return /* This is a block comment */ // This is a line comment`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: LPAREN, data: "("},
				{pos: position{line: 1, col: 3}, tt: RPAREN, data: ")"},
				{pos: position{line: 1, col: 5}, tt: LBRACE, data: "{"},
				{pos: position{line: 1, col: 7}, tt: RBRACE, data: "}"},
				{pos: position{line: 1, col: 9}, tt: LBRACKET, data: "["},
				{pos: position{line: 1, col: 11}, tt: RBRACKET, data: "]"},
				{pos: position{line: 1, col: 13}, tt: COMMA, data: ","},
				{pos: position{line: 1, col: 15}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 17}, tt: SEMICOLON, data: ";"},
				{pos: position{line: 1, col: 19}, tt: PERIOD, data: "."},
				{pos: position{line: 1, col: 21}, tt: ASSIGN, data: "="},
				{pos: position{line: 1, col: 23}, tt: HASH, data: "#"},
				{pos: position{line: 1, col: 25}, tt: THEN_ARROW, data: "=>"},
				{pos: position{line: 1, col: 28}, tt: WHEN, data: "when"},
				{pos: position{line: 2, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 2, col: 3}, tt: DECIMAL, data: "1.0"},
				{pos: position{line: 2, col: 7}, tt: STRING, data: "Hello, World!"},
				{pos: position{line: 2, col: 23}, tt: STRING, data: "Hello, name!"},
				{pos: position{line: 2, col: 38}, tt: IDENTIFIER, data: "a"},
				{pos: position{line: 2, col: 40}, tt: TYPE_IDENTIFIER, data: "Point"},
				{pos: position{line: 2, col: 46}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 2, col: 48}, tt: BREAK, data: "break"},
				{pos: position{line: 2, col: 54}, tt: CONTINUE, data: "continue"},
				{pos: position{line: 2, col: 63}, tt: RETURN, data: "return"},
				{pos: position{line: 2, col: 126}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Line comment + EOF",
			[]string{"test.yz"},
			`// This is a comment`,
			[]Token{
				{pos: position{line: 1, col: 22}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Line comment",
			[]string{"test.yz"},
			`// This is a comment
1 + 2`,
			[]Token{
				{pos: position{line: 2, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 2, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 2, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 2, col: 6}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Code and line comment",
			[]string{"test.yz"},
			`1 + 2 // This is a comment
a: 3`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 2, col: 1}, tt: IDENTIFIER, data: "a"},
				{pos: position{line: 2, col: 2}, tt: COLON, data: ":"},
				{pos: position{line: 2, col: 4}, tt: INTEGER, data: "3"},
				{pos: position{line: 2, col: 5}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Block comment",
			[]string{"test.yz"},
			`a:1 /* This is a block comment
that spans multiple lines
*/
b:2`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "a"},
				{pos: position{line: 1, col: 2}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 3}, tt: INTEGER, data: "1"},
				{pos: position{line: 3, col: 3}, tt: COMMA, data: ","},
				{pos: position{line: 4, col: 1}, tt: IDENTIFIER, data: "b"},
				{pos: position{line: 4, col: 2}, tt: COLON, data: ":"},
				{pos: position{line: 4, col: 3}, tt: INTEGER, data: "2"},
				{pos: position{line: 4, col: 4}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Block comment followed by line comment",
			[]string{"test.yz"},
			`a:1 /* This is a block comment */ // This is a line comment`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "a"},
				{pos: position{line: 1, col: 2}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 3}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 61}, tt: EOF, data: "EOF"},
			},
		},
		{
			"String literals",
			[]string{"test.yz"},
			`"Hello, World!"`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "Hello, World!"},
				{pos: position{line: 1, col: 16}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Multiline string literals",
			[]string{"test.yz"},
			`"Hi,
World!"`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "Hi,\nWorld!"},
				{pos: position{line: 2, col: 8}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Multiline string literals with escape",
			[]string{"test.yz"},
			`"Hi, \"
World!"`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "Hi, \"\nWorld!"},
				{pos: position{line: 2, col: 8}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Multiline string literals preserving indentation",
			[]string{"test.yz"},
			// There are trailing spaces in the first line of the string
			`"Hi,   
	World!   .
	."`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "Hi,   \n\tWorld!   .\n\t."},
				{pos: position{line: 3, col: 4}, tt: EOF, data: "EOF"},
			},
		},
		{
			"String literals with various escapes",
			[]string{"test.yz"},
			`"The quick\nbrown fox\tjumps over\\the lazy dog\"She said, 'Hello'\\xUnknown escape"`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "The quick\nbrown fox\tjumps over\\the lazy dog\"She said, 'Hello'\\xUnknown escape"},
				{pos: position{line: 1, col: 85}, tt: EOF, data: "EOF"},
			},
		},
		{
			"String interpolation",
			[]string{"test.yz"},
			"'Hello, `name`!'",
			[]Token{
				{pos: position{line: 1, col: 1}, tt: STRING, data: "Hello, `name`!"},
				{pos: position{line: 1, col: 17}, tt: EOF, data: "EOF"},
			},
		},
		{
			"String literals",
			[]string{"test.yz"},
			`["One" 'Two' "'Three'" '"Four"']`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: LBRACKET, data: "["},
				{pos: position{line: 1, col: 2}, tt: STRING, data: "One"},
				{pos: position{line: 1, col: 8}, tt: STRING, data: "Two"},
				{pos: position{line: 1, col: 14}, tt: STRING, data: "'Three'"},
				{pos: position{line: 1, col: 24}, tt: STRING, data: "\"Four\""},
				{pos: position{line: 1, col: 32}, tt: RBRACKET, data: "]"},
				{pos: position{line: 1, col: 33}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Integers",
			[]string{"test.yz"},
			`1 9876324`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: INTEGER, data: "9876324"},
				{pos: position{line: 1, col: 10}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Decimal literals",
			[]string{"test.yz"},
			`1.0`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: DECIMAL, data: "1.0"},
				{pos: position{line: 1, col: 4}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Negative numbers",
			[]string{"test.yz"},
			`minusThree: -1 - -2.0
plusOne: -1 + -2.0
`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "minusThree"},
				{pos: position{line: 1, col: 11}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 13}, tt: INTEGER, data: "-1"},
				{pos: position{line: 1, col: 16}, tt: NON_WORD_IDENTIFIER, data: "-"},
				{pos: position{line: 1, col: 18}, tt: DECIMAL, data: "-2.0"},
				{pos: position{line: 1, col: 22}, tt: COMMA, data: ","},
				{pos: position{line: 2, col: 1}, tt: IDENTIFIER, data: "plusOne"},
				{pos: position{line: 2, col: 8}, tt: COLON, data: ":"},
				{pos: position{line: 2, col: 10}, tt: INTEGER, data: "-1"},
				{pos: position{line: 2, col: 13}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 2, col: 15}, tt: DECIMAL, data: "-2.0"},
				{pos: position{line: 2, col: 19}, tt: COMMA, data: ","},
				{pos: position{line: 3, col: 1}, tt: EOF, data: "EOF"},
			},
		},

		{
			"Equals sign examples",
			[]string{"test.yz"},
			`== => =< =a =`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: EQUALS, data: "=="},
				{pos: position{line: 1, col: 4}, tt: THEN_ARROW, data: "=>"},
				{pos: position{line: 1, col: 7}, tt: NON_WORD_IDENTIFIER, data: "=<"},
				{pos: position{line: 1, col: 10}, tt: IDENTIFIER, data: "=a"},
				{pos: position{line: 1, col: 13}, tt: ASSIGN, data: "="},
				{pos: position{line: 1, col: 14}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Non ascii characters",
			[]string{"test.yz"},
			`message: 👋🌍`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "message"},
				{pos: position{line: 1, col: 8}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 10}, tt: NON_WORD_IDENTIFIER, data: "👋🌍"},
				{pos: position{line: 1, col: 12}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Non word identifiers examples",
			[]string{"test.yz"},
			`+ - * / % ~ < > ! & | ^ += -= /= ~= != <= >= && || ++ -- >>= <<= >> << |> <- -> `,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "-"},
				{pos: position{line: 1, col: 5}, tt: NON_WORD_IDENTIFIER, data: "*"},
				{pos: position{line: 1, col: 7}, tt: NON_WORD_IDENTIFIER, data: "/"},
				{pos: position{line: 1, col: 9}, tt: NON_WORD_IDENTIFIER, data: "%"},
				{pos: position{line: 1, col: 11}, tt: NON_WORD_IDENTIFIER, data: "~"},
				{pos: position{line: 1, col: 13}, tt: NON_WORD_IDENTIFIER, data: "<"},
				{pos: position{line: 1, col: 15}, tt: NON_WORD_IDENTIFIER, data: ">"},
				{pos: position{line: 1, col: 17}, tt: NON_WORD_IDENTIFIER, data: "!"},
				{pos: position{line: 1, col: 19}, tt: NON_WORD_IDENTIFIER, data: "&"},
				{pos: position{line: 1, col: 21}, tt: NON_WORD_IDENTIFIER, data: "|"},
				{pos: position{line: 1, col: 23}, tt: NON_WORD_IDENTIFIER, data: "^"},
				{pos: position{line: 1, col: 25}, tt: NON_WORD_IDENTIFIER, data: "+="},
				{pos: position{line: 1, col: 28}, tt: NON_WORD_IDENTIFIER, data: "-="},
				{pos: position{line: 1, col: 31}, tt: NON_WORD_IDENTIFIER, data: "/="},
				{pos: position{line: 1, col: 34}, tt: NON_WORD_IDENTIFIER, data: "~="},
				{pos: position{line: 1, col: 37}, tt: NON_WORD_IDENTIFIER, data: "!="},
				{pos: position{line: 1, col: 40}, tt: NON_WORD_IDENTIFIER, data: "<="},
				{pos: position{line: 1, col: 43}, tt: NON_WORD_IDENTIFIER, data: ">="},
				{pos: position{line: 1, col: 46}, tt: NON_WORD_IDENTIFIER, data: "&&"},
				{pos: position{line: 1, col: 49}, tt: NON_WORD_IDENTIFIER, data: "||"},
				{pos: position{line: 1, col: 52}, tt: NON_WORD_IDENTIFIER, data: "++"},
				{pos: position{line: 1, col: 55}, tt: NON_WORD_IDENTIFIER, data: "--"},
				{pos: position{line: 1, col: 58}, tt: NON_WORD_IDENTIFIER, data: ">>="},
				{pos: position{line: 1, col: 62}, tt: NON_WORD_IDENTIFIER, data: "<<="},
				{pos: position{line: 1, col: 66}, tt: NON_WORD_IDENTIFIER, data: ">>"},
				{pos: position{line: 1, col: 69}, tt: NON_WORD_IDENTIFIER, data: "<<"},
				{pos: position{line: 1, col: 72}, tt: NON_WORD_IDENTIFIER, data: "|>"},
				{pos: position{line: 1, col: 75}, tt: NON_WORD_IDENTIFIER, data: "<-"},
				{pos: position{line: 1, col: 78}, tt: NON_WORD_IDENTIFIER, data: "->"},
				{pos: position{line: 1, col: 81}, tt: EOF, data: "EOF"},
			},
		},

		{
			"Printable UTF-8 characters (additional symbols and emojis)",
			[]string{"test.yz"},
			`©®™✓✔✕✖✗✘✙✚✛✜✢✣✤✥✦✧✨⭐✩✪✫✬✭✮✯✰✱✲✳✴✵✶✷✸✹✺✻✼✽✾✿❀❁❂❃❄❅❆❇❈❉❊❋❌❍❎❏❐❑❒❖❗❘❙❚❛❜❝❞❡❢❣❤❥❦❧`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: NON_WORD_IDENTIFIER, data: "©®™✓✔✕✖✗✘✙✚✛✜✢✣✤✥✦✧✨⭐✩✪✫✬✭✮✯✰✱✲✳✴✵✶✷✸✹✺✻✼✽✾✿❀❁❂❃❄❅❆❇❈❉❊❋❌❍❎❏❐❑❒❖❗❘❙❚❛❜❝❞❡❢❣❤❥❦❧"},
				{pos: position{line: 1, col: 80}, tt: EOF, data: "EOF"},
			},
		},
		{
			name:    "Common mathematical symbols in Unicode and UTF-8",
			path:    []string{"test.yz"},
			content: `∑ ∏ ∫ ∞ √ ∇ ≈ ≠ ≤ ≥ ± ∂ ∃ ∀ ∈ ∉ ∋ ∅ ∧ ∨ ∩ ∪ ⊂ ⊃ ⊆ ⊇ ⊕ ⊗ ⊥`,
			want: []Token{
				{pos: position{line: 1, col: 1}, tt: NON_WORD_IDENTIFIER, data: "∑"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "∏"},
				{pos: position{line: 1, col: 5}, tt: NON_WORD_IDENTIFIER, data: "∫"},
				{pos: position{line: 1, col: 7}, tt: NON_WORD_IDENTIFIER, data: "∞"},
				{pos: position{line: 1, col: 9}, tt: NON_WORD_IDENTIFIER, data: "√"},
				{pos: position{line: 1, col: 11}, tt: NON_WORD_IDENTIFIER, data: "∇"},
				{pos: position{line: 1, col: 13}, tt: NON_WORD_IDENTIFIER, data: "≈"},
				{pos: position{line: 1, col: 15}, tt: NON_WORD_IDENTIFIER, data: "≠"},
				{pos: position{line: 1, col: 17}, tt: NON_WORD_IDENTIFIER, data: "≤"},
				{pos: position{line: 1, col: 19}, tt: NON_WORD_IDENTIFIER, data: "≥"},
				{pos: position{line: 1, col: 21}, tt: NON_WORD_IDENTIFIER, data: "±"},
				{pos: position{line: 1, col: 23}, tt: NON_WORD_IDENTIFIER, data: "∂"},
				{pos: position{line: 1, col: 25}, tt: NON_WORD_IDENTIFIER, data: "∃"},
				{pos: position{line: 1, col: 27}, tt: NON_WORD_IDENTIFIER, data: "∀"},
				{pos: position{line: 1, col: 29}, tt: NON_WORD_IDENTIFIER, data: "∈"},
				{pos: position{line: 1, col: 31}, tt: NON_WORD_IDENTIFIER, data: "∉"},
				{pos: position{line: 1, col: 33}, tt: NON_WORD_IDENTIFIER, data: "∋"},
				{pos: position{line: 1, col: 35}, tt: NON_WORD_IDENTIFIER, data: "∅"},
				{pos: position{line: 1, col: 37}, tt: NON_WORD_IDENTIFIER, data: "∧"},
				{pos: position{line: 1, col: 39}, tt: NON_WORD_IDENTIFIER, data: "∨"},
				{pos: position{line: 1, col: 41}, tt: NON_WORD_IDENTIFIER, data: "∩"},
				{pos: position{line: 1, col: 43}, tt: NON_WORD_IDENTIFIER, data: "∪"},
				{pos: position{line: 1, col: 45}, tt: NON_WORD_IDENTIFIER, data: "⊂"},
				{pos: position{line: 1, col: 47}, tt: NON_WORD_IDENTIFIER, data: "⊃"},
				{pos: position{line: 1, col: 49}, tt: NON_WORD_IDENTIFIER, data: "⊆"},
				{pos: position{line: 1, col: 51}, tt: NON_WORD_IDENTIFIER, data: "⊇"},
				{pos: position{line: 1, col: 53}, tt: NON_WORD_IDENTIFIER, data: "⊕"},
				{pos: position{line: 1, col: 55}, tt: NON_WORD_IDENTIFIER, data: "⊗"},
				{pos: position{line: 1, col: 57}, tt: NON_WORD_IDENTIFIER, data: "⊥"},
				{pos: position{line: 1, col: 58}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Identifiers with various scripts and their values",
			[]string{"test.yz"},
			`¡hola!: "latin"
привет: "cyrillic"
变量: "chinese"
변수: "korean"
変数: "japanese"
über: "german"
नमस्ते: "hindi"
สวัสดี: "thai"
ሰላም: "amharic"
γειά: "greek"
æøå: "nordic"
`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "¡hola!"},
				{pos: position{line: 1, col: 7}, tt: COLON, data: ":"},
				{pos: position{line: 1, col: 9}, tt: STRING, data: "latin"},
				{pos: position{line: 1, col: 16}, tt: COMMA, data: ","},
				{pos: position{line: 2, col: 1}, tt: IDENTIFIER, data: "привет"},
				{pos: position{line: 2, col: 7}, tt: COLON, data: ":"},
				{pos: position{line: 2, col: 9}, tt: STRING, data: "cyrillic"},
				{pos: position{line: 2, col: 19}, tt: COMMA, data: ","},
				{pos: position{line: 3, col: 1}, tt: IDENTIFIER, data: "变量"},
				{pos: position{line: 3, col: 3}, tt: COLON, data: ":"},
				{pos: position{line: 3, col: 5}, tt: STRING, data: "chinese"},
				{pos: position{line: 3, col: 14}, tt: COMMA, data: ","},
				{pos: position{line: 4, col: 1}, tt: IDENTIFIER, data: "변수"},
				{pos: position{line: 4, col: 3}, tt: COLON, data: ":"},
				{pos: position{line: 4, col: 5}, tt: STRING, data: "korean"},
				{pos: position{line: 4, col: 13}, tt: COMMA, data: ","},
				{pos: position{line: 5, col: 1}, tt: IDENTIFIER, data: "変数"},
				{pos: position{line: 5, col: 3}, tt: COLON, data: ":"},
				{pos: position{line: 5, col: 5}, tt: STRING, data: "japanese"},
				{pos: position{line: 5, col: 15}, tt: COMMA, data: ","},
				{pos: position{line: 6, col: 1}, tt: IDENTIFIER, data: "über"},
				{pos: position{line: 6, col: 5}, tt: COLON, data: ":"},
				{pos: position{line: 6, col: 7}, tt: STRING, data: "german"},
				{pos: position{line: 6, col: 15}, tt: COMMA, data: ","},
				{pos: position{line: 7, col: 1}, tt: IDENTIFIER, data: "नमस्ते"},
				{pos: position{line: 7, col: 7}, tt: COLON, data: ":"},
				{pos: position{line: 7, col: 9}, tt: STRING, data: "hindi"},
				{pos: position{line: 7, col: 16}, tt: COMMA, data: ","},
				{pos: position{line: 8, col: 1}, tt: IDENTIFIER, data: "สวัสดี"},
				{pos: position{line: 8, col: 7}, tt: COLON, data: ":"},
				{pos: position{line: 8, col: 9}, tt: STRING, data: "thai"},
				{pos: position{line: 8, col: 15}, tt: COMMA, data: ","},
				{pos: position{line: 9, col: 1}, tt: IDENTIFIER, data: "ሰላም"},
				{pos: position{line: 9, col: 4}, tt: COLON, data: ":"},
				{pos: position{line: 9, col: 6}, tt: STRING, data: "amharic"},
				{pos: position{line: 9, col: 15}, tt: COMMA, data: ","},
				{pos: position{line: 10, col: 1}, tt: IDENTIFIER, data: "γειά"},
				{pos: position{line: 10, col: 5}, tt: COLON, data: ":"},
				{pos: position{line: 10, col: 7}, tt: STRING, data: "greek"},
				{pos: position{line: 10, col: 14}, tt: COMMA, data: ","},
				{pos: position{line: 11, col: 1}, tt: IDENTIFIER, data: "æøå"},
				{pos: position{line: 11, col: 4}, tt: COLON, data: ":"},
				{pos: position{line: 11, col: 6}, tt: STRING, data: "nordic"},
				{pos: position{line: 11, col: 14}, tt: COMMA, data: ","},
				{pos: position{line: 12, col: 1}, tt: EOF, data: "EOF"},
			},
		},

		{
			"Simple",
			[]string{"test.yz"},
			`1 + 2`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Simple with spaces",
			[]string{"test.yz"},
			`1 + 2`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Simple with newline",
			[]string{"test.yz"},
			`1 + 2
`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: COMMA, data: ","},
				{pos: position{line: 2, col: 1}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Simple with newline and spaces",
			[]string{"test.yz"},
			`1 + 2
`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NON_WORD_IDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: COMMA, data: ","},
				{pos: position{line: 2, col: 1}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Hash and parenthesis",
			[]string{"test.yz"},
			`fn #(Int)`,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "fn"},
				{pos: position{line: 1, col: 4}, tt: HASH, data: "#"},
				{pos: position{line: 1, col: 5}, tt: LPAREN, data: "("},
				{pos: position{line: 1, col: 6}, tt: TYPE_IDENTIFIER, data: "Int"},
				{pos: position{line: 1, col: 9}, tt: RPAREN, data: ")"},
				{pos: position{line: 1, col: 10}, tt: EOF, data: "EOF"},
			},
		},

		{
			"Empty file",
			[]string{"test.yz"},
			``,
			[]Token{
				{pos: position{line: 1, col: 1}, tt: EOF, data: "EOF"},
			},
		},
		{
			"Empty file with newline",
			[]string{"test.yz"},
			`  
  `,
			[]Token{
				{pos: position{line: 2, col: 3}, tt: EOF, data: "EOF"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := Tokenize(tt.path, tt.content)
			if e != nil {
				t.Errorf("Tokenize() error = %v", e)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Tokenize() = \ngot =  %v\n"+
						"       (pos:%v), \nwant = %v\n"+
						"       (pos:%v)", got, got[i].pos, tt.want, tt.want[i].pos)
					return
				}
			}
		})
	}
}

func TestTokenizer_SyntaxErr(t *testing.T) {
	tests := []struct {
		name    string
		path    []string
		content string
		want    error
	}{
		{
			"Unclosed string",
			[]string{"test.yz"},
			`"`,
			fmt.Errorf("[test.yz: line:1: col:2]: Syntax error: unterminated string literal"),
		},
		{
			"Unclosed multiline comment",
			[]string{"test.yz"},
			`/*`,
			fmt.Errorf("[test.yz: line:1: col:3]: Syntax error: unterminated comment"),
		},
		{
			"Backtick strings",
			[]string{"test.yz"},
			"`hola`",
			fmt.Errorf("[test.yz: line:1: col:1]: Unexpected Token `"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := Tokenize(tt.path, tt.content)

			if got == nil {
				t.Errorf("Tokenize() error = %v, want %v", got, tt.want)
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("Tokenize() error = %v, want %v", got, tt.want)
			}
		})
	}
}
