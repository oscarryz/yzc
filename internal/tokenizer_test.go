package internal

// Tests the tokenizer.Tokenize function

import "testing"

func TestTokenizer_Tokenize(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		content  string
		want     []token
	}{
		{
			"Simple",
			"test.yz",
			`1 + 2`,
			[]token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NONWORDIDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: EOF, data: "EOF" },
			},
		},
		{
			"Simple with spaces",
			"test.yzc",
			`1 + 2`,
			[]token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NONWORDIDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 1, col: 6}, tt: EOF, data: "EOF" },
			},
		},
		{
			"Simple with newline",
			"test.yzc",
			`1 + 2
`,
			[]token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NONWORDIDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 2, col: 1}, tt: EOF, data: "EOF" },
			},
		},
		{
			"Simple with newline and spaces",
			"test.yzc",
			`1 + 2
`,
			[]token{
				{pos: position{line: 1, col: 1}, tt: INTEGER, data: "1"},
				{pos: position{line: 1, col: 3}, tt: NONWORDIDENTIFIER, data: "+"},
				{pos: position{line: 1, col: 5}, tt: INTEGER, data: "2"},
				{pos: position{line: 2, col: 1}, tt: EOF, data: "EOF" },
			},
		},
		{
			"Hash and parenthesis",
			"test.yzc",
			`fn #(Int)`,
			[]token{
				{pos: position{line: 1, col: 1}, tt: IDENTIFIER, data: "fn"},
				{pos: position{line: 1, col: 4}, tt: HASH, data: "#"},
				{pos: position{line: 1, col: 5}, tt: LPAREN, data: "("},
				{pos: position{line: 1, col: 6}, tt: TYPEIDENTIFIER, data: "Int"},
				{pos: position{line: 1, col: 9}, tt: RPAREN, data: ")"},
				{pos: position{line: 1, col: 10}, tt: EOF, data: "EOF" },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := Tokenize(tt.fileName, tt.content)
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
					t.Errorf("Tokenize() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}
