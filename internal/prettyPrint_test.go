package internal

import (
	"strings"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "Boc",
			input: &Boc{
				expressions: []expression{
					&BasicLit{
						pos: pos(1, 1),
						tt:  STRING,
						val: "Hello",
					},
				},
				statements: []statement{},
			},
			expected: `Boc (
    BasicLit (
        pos: line: 1 col: 1
        tt: str
        val: Hello
    )
)
`,
		},

		{
			name: "BasicLit",
			input: &BasicLit{
				pos: pos(1, 1),
				tt:  STRING,
				val: "Hello",
			},
			expected: `BasicLit (
    pos: line: 1 col: 1
    tt: str
    val: Hello
)
`,
		},
		{
			name: "ArrayLit",
			input: &ArrayLit{
				pos: pos(1, 1),
				arrayType: &ArrayType{
					elemType: &IntType{},
				},
				exps: []expression{
					&BasicLit{
						pos: pos(2, 1),
						tt:  INTEGER,
						val: "1",
					},
				},
			},
			expected: `ArrayLit(
    pos: line: 1 col: 1
    arrayType: 
		ArrayType(
			elemType: IntType
		)
    exps: [
        BasicLit(
            pos: line: 2 col: 1
            tt: int
            val: 1
        )
    ]
)
`,
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
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(s, " ", ""),
			"\n", ""),
		"\t", "")
}
