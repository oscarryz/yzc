package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	EOF tokenType = iota
	// punctuation
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	PERIOD    // .
	ASSIGN    // =
	EQUALS    // ==
	HASH      // #
	NEWLINE   // \n

	// literals
	INTEGER // int
	DECIMAL // dec
	STRING  // str

	// identifiers
	IDENTIFIER        // id
	TYPEIDENTIFIER    // tid
	NONWORDIDENTIFIER // nwid
	BREAK             // BREAK
	CONTINUE          // CONTINUE
	RETURN            // RETURN
	ILLEGAL           // ILLEGAL
)

type tokenType uint

type token struct {
	pos  position
	tt   tokenType
	data string
}

type position struct {
	line int
	col  int
}

type tokenizer struct {
	content   string
	tokens    []token
	col       int
	line      int
	pos       int
	keepGoing bool
}

func (tt tokenType) String() string {
	descriptions := [25]string{
		`EOF`, `(`, `)`, `{`, `}`, `[`, `]`, `,`, `:`, `;`, `.`, `=`, `==`, `#`, "NEWLINE",
		`int`, `dec`, `str`, `id`, `tid`, `nwid`, "BREAK", "CONTINUE", "RETURN", "ILLEGAL",
	}
	vot := int(tt)
	if vot > len(descriptions) {
		return strconv.Itoa(vot)
	} else {
		return descriptions[tt]
	}
}

func (p position) String() string {
	return fmt.Sprintf("line: %d col: %d", p.line, p.col)
}

func pos(line, col int) position {
	return position{line, col}
}

func (t token) String() string {

	switch t.tt {
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NONWORDIDENTIFIER, TYPEIDENTIFIER:
		return fmt.Sprintf("%s:%s ", t.tt, t.data)
	default:
		return fmt.Sprintf("%v ", t.tt)
	}
}

// Tokenize converts the content into an array of tokens or returns an error if the content is not valid
func Tokenize(fileName string, content string) ([]token, error) {
	t := &tokenizer{content, []token{}, 0, 1, 0, true}
	tokens, e := t.tokenize()
	printTokens(tokens)
	return tokens, e
}

func printTokens(tokens []token) {
	ll := 1
	var builder strings.Builder
	builder.WriteString("Tokens:\n")
	builder.WriteString(fmt.Sprintf("%d: ", ll))
	for _, t := range tokens {
		if ll != t.pos.line {
			ll = t.pos.line
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("%d: ", ll))
		}
		builder.WriteString(fmt.Sprintf("%v", t))
	}
	builder.WriteString("\n")
	logger.Println(builder.String())
}

func (t *tokenizer) addToken(tt tokenType, data string) {
	col := t.col

	if tt != EOF {
		dataLen := utf8.RuneCountInString(data)
		col = t.col - dataLen + 1
	}
	t.tokens = append(t.tokens, token{pos(t.line, col), tt, data})
}

func (t *tokenizer) nextRune() rune {
	r, w := utf8.DecodeRuneInString(t.content[t.pos:])
	if r == utf8.RuneError {
		t.keepGoing = false
	}
	t.pos += w
	t.col++
	return r
}

func (t *tokenizer) unReadRune(r rune) {
	if r == utf8.RuneError {
		t.keepGoing = false
	} else {
		t.pos -= utf8.RuneLen(r)
	}
	t.col--
}

func (t *tokenizer) peek() rune {
	r := t.nextRune()
	t.unReadRune(r)
	return r
}

func (t *tokenizer) skipComment() {
	// read until new line or EOF
	for {
		r := t.nextRune()
		if r == '\n' {
			t.addToken(NEWLINE, "\n")
			t.line++
			t.col = 0
			return
		}
		if r == utf8.RuneError {
			t.keepGoing = false
			return

		}
	}
}

func (t *tokenizer) skipMultilineComment() {
	r := t.nextRune()
	if r == '\n' {
		t.line++
		t.col = 0
	}

	for {
		if r == '*' && t.peek() == '/' {
			t.nextRune()
			if r == '\n' {
				t.line++
				t.col = 0
			}
			return
		}
		r = t.nextRune()
		if r == '\n' {
			t.line++
			t.col = 0
		}
	}
}

func lookupIdent(identifier string) tokenType {
	runes := []rune(identifier)
	allUpper := true
	for _, r := range runes {
		allUpper = allUpper && unicode.IsUpper(r)
	}
	if allUpper {
		return IDENTIFIER
	}
	if unicode.IsUpper(runes[0]) {
		return TYPEIDENTIFIER
	}
	switch identifier {
	case "=":
		return ASSIGN
	case "==":
		return EQUALS
	case "break":
		return BREAK
	case "continue":
		return CONTINUE
	case "return":
		return RETURN
	default:
		allNonLetter := true
		for _, r := range runes {
			allNonLetter = allNonLetter && !unicode.IsLetter(r)
		}
		if allNonLetter {
			return NONWORDIDENTIFIER
		}
	}
	return IDENTIFIER
}

func (t *tokenizer) addStringLiteral() {
	opening := t.nextRune()
	pos := position{t.line, t.col}
	r := t.nextRune()
	builder := strings.Builder{}
	for r != opening {
		if r == '\\' {
			// Handle escape sequences
			next := t.nextRune()
			switch next {
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			case '\\':
				builder.WriteRune('\\')
			case '"':
				builder.WriteRune('"')
			case '\'':
				builder.WriteRune('\'')
			default:
				builder.WriteRune('\\')
				builder.WriteRune(next)
			}
		} else {
			builder.WriteRune(r)
		}
		if r == '\n' {
			t.line++
			t.col = 0
		}
		r = t.nextRune()
	}
	t.tokens = append(t.tokens, token{pos, STRING, builder.String()})
}

func (t *tokenizer) isIdentifier(r rune) bool {
	return !unicode.IsSpace(r) && unicode.IsPrint(r) && !unicode.IsDigit(r) && !strings.ContainsRune("{}[]#().,:;\"'`", r)
}

func (t *tokenizer) readIdentifier() string {
	builder := strings.Builder{}
	r := t.nextRune()
	for (unicode.IsPrint(r) || unicode.IsDigit(r)) && !strings.ContainsRune("{}[]#().,:;\"'`", r) && !unicode.IsSpace(r) && r != utf8.RuneError {
		builder.WriteRune(r)
		r = t.nextRune()
	}
	t.unReadRune(r)
	return builder.String()
}

func (t *tokenizer) readNumber(positive bool) {
	r := t.nextRune()
	builder := strings.Builder{}
	if !positive {
		builder.WriteRune('-')
	}
	seenPoint := false
	for unicode.IsDigit(r) || r == '.' {
		if r == '.' && seenPoint {
			break
		}
		if r == '.' {
			seenPoint = true
		}
		builder.WriteRune(r)
		r = t.nextRune()
	}
	t.unReadRune(r)
	tt := INTEGER
	if seenPoint {
		tt = DECIMAL
	}
	t.addToken(tt, builder.String())
}

func (t *tokenizer) addNumber() {
	t.readNumber(true)
}

func (t *tokenizer) addNegativeNumber() {
	t.nextRune()
	t.readNumber(false)
}

func (t *tokenizer) tokenize() ([]token, error) {
	for r := t.nextRune(); t.keepGoing; r = t.nextRune() {
		if r == '\n' {
			t.addToken(NEWLINE, "\n")
			t.line++
			t.col = 0
		}
		if unicode.IsSpace(r) {
			continue
		}
		switch r {
		case '{':
			t.addToken(LBRACE, "{")
		case '}':
			t.addToken(RBRACE, "}")
		case ':':
			t.addToken(COLON, ":")
		case ';':
			t.addToken(SEMICOLON, ";")
		case '.':
			t.addToken(PERIOD, ".")
		case '(':
			t.addToken(LPAREN, "(")
		case ')':
			t.addToken(RPAREN, ")")
		case '[':
			t.addToken(LBRACKET, "[")
		case ']':
			t.addToken(RBRACKET, "]")
		case ',':
			t.addToken(COMMA, ",")
		case '#':
			t.addToken(HASH, "#")
		case '"', '`', '\'':
			t.unReadRune(r)
			t.addStringLiteral()
		default:
			if r == '/' {
				if t.peek() == '/' {
					t.nextRune()
					t.skipComment()
					continue
				} else if t.peek() == '*' {
					t.nextRune()
					t.skipMultilineComment()
					continue
				}
			}
			if r == '-' && unicode.IsDigit(t.peek()) {
				t.unReadRune(r)
				t.addNegativeNumber()
				break
			}
			if unicode.IsDigit(r) {
				t.unReadRune(r)
				t.addNumber()
			} else if t.isIdentifier(r) {
				t.unReadRune(r)
				id := t.readIdentifier()
				t.addToken(lookupIdent(id), id)
			} else {
				t.addToken(ILLEGAL, string(r))
				return t.tokens, errors.New(fmt.Sprintf("Illegal token at: %d %d", t.line, t.col))
			}
		}
	}
	t.addToken(EOF, "EOF")
	return t.tokens, nil
}
