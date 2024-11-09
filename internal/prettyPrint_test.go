package internal

import (
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "boc",
			input: &boc{
				Name: "test_boc",
				blockBody: &blockBody{
					expressions: []expression{
						&BasicLit{
							pos: pos(1, 1),
							tt:  STRING,
							val: "Hello",
						},
					},
					statements: []statement{},
				},
			},
			expected: `boc {
    Name: test_boc
    blockBody {
        BasicLit {
            pos: line: 1 col: 1
            tt: str
            val: Hello
        }
    }
}
`,
		},
		{
			name: "blockType",
			input: &blockType{
				pos: pos(1, 1),
				tt:  TYPEIDENTIFIER,
				val: "Int",
			},
			expected: `blockType {
    pos: line: 1 col: 1
    tt: tid
    val: Int
}
`,
		},
		{
			name: "blockBody",
			input: &blockBody{
				expressions: []expression{
					&BasicLit{
						pos: pos(1, 1),
						tt:  INTEGER,
						val: "1",
					},
				},
				statements: []statement{},
			},
			expected: `blockBody {
    BasicLit {
        pos: line: 1 col: 1
        tt: int
        val: 1
    }
}
`,
		},
		{
			name: "BasicLit",
			input: &BasicLit{
				pos: pos(1, 1),
				tt:  STRING,
				val: "Hello",
			},
			expected: `BasicLit {
    pos: line: 1 col: 1
    tt: str
    val: Hello
}
`,
		},
		{
			name: "ArrayLit",
			input: &ArrayLit{
				pos: pos(1, 1),
				tt:  LBRACKET,
				val: "[]",
				exps: []expression{
					&BasicLit{
						pos: pos(2, 1),
						tt:  INTEGER,
						val: "1",
					},
				},
			},
			expected: `ArrayLit {
    pos: line: 1 col: 1
    tt: [
    val: []
    exps: [
        BasicLit {
            pos: line: 2 col: 1
            tt: int
            val: 1
        }
    ]
}
`,
		},
		{
			name:     "\tempty struct{}\n",
			input:    &empty{},
			expected: "<empty>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prettyPrint(tt.input, 0)
			if removeSpaces(result) != removeSpaces(tt.expected) {
				t.Errorf("prettyPrint() got:\n%v, want:\n%v", result, tt.expected)
			}
		})
	}
}

func removeSpaces(s string) string {
	return s // strings.ReplaceAll(s, " ", "")
}