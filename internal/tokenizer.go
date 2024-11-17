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

	// literals
	INTEGER // int
	DECIMAL // dec
	STRING  // str

	// identifiers
	IDENTIFIER          // id
	TYPE_IDENTIFIER     // tid
	NON_WORD_IDENTIFIER // nwid
	BREAK               // BREAK
	CONTINUE            // CONTINUE
	RETURN              // RETURN
	Unexpected          // Unexpected
)

type tokenType uint

type Token struct {
	pos  position
	tt   tokenType
	data string
}

type position struct {
	line int
	col  int
}

type tokenizer struct {
	filname   string
	content   string
	tokens    []Token
	col       int
	line      int
	pos       int
	keepGoing bool
}

func (tt tokenType) String() string {
	descriptions := [25]string{
		`EOF`, `(`, `)`, `{`, `}`, `[`, `]`, `,`, `:`, `;`, `.`, `=`, `==`, `#`,
		`int`, `dec`, `str`, `id`, `tid`, `nwid`, "BREAK", "CONTINUE", "RETURN", "Unexpected",
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

func (t Token) String() string {

	switch t.tt {
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NON_WORD_IDENTIFIER, TYPE_IDENTIFIER:
		return fmt.Sprintf("%s:%s ", t.tt, t.data)
	default:
		return fmt.Sprintf("%v ", t.tt)
	}
}

// Tokenize converts the content into an array of tokens or returns an error if the content is not valid
func Tokenize(path []string, content string) ([]Token, error) {
	t := &tokenizer{path[len(path)-1], content, []Token{}, 0, 1, 0, true}
	tokens, e := t.tokenize()
	return tokens, e
}

func (t *tokenizer) addToken(tt tokenType, data string) {
	col := t.col

	if tt != EOF {
		dataLen := utf8.RuneCountInString(data)
		col = t.col - dataLen + 1
	}
	t.tokens = append(t.tokens, Token{pos(t.line, col), tt, data})
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
	if r == utf8.RuneError {
		t.addToken(Unexpected, fmt.Sprintf("[%s: line:%d: col:%d]: Syntax error: unterminated comment", t.filname, t.line, t.col))
		t.keepGoing = false
		return
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
		return TYPE_IDENTIFIER
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
			return NON_WORD_IDENTIFIER
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
		if r == utf8.RuneError {
			t.addToken(Unexpected, fmt.Sprintf("[%s: line:%d: col:%d]: Syntax error: unterminated string literal", t.filname, t.line, t.col))
			t.keepGoing = false
			return
		}
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
	t.tokens = append(t.tokens, Token{pos, STRING, builder.String()})
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

func (t *tokenizer) tokenize() ([]Token, error) {
	for r := t.nextRune(); t.keepGoing; r = t.nextRune() {
		if r == '\n' {
			t.addCommaIfNeeded()
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
		case '"', '\'':
			t.unReadRune(r)
			t.addStringLiteral()
		default:
			switch {
			case r == '-' && unicode.IsDigit(t.peek()):
				t.unReadRune(r)
				t.addNegativeNumber()
				break
			case unicode.IsDigit(r):
				t.unReadRune(r)
				t.addNumber()
			case r == '/':
				if t.peek() == '/' {
					t.nextRune()
					t.skipComment()
					continue
				} else if t.peek() == '*' {
					t.nextRune()
					t.skipMultilineComment()
					continue
				}
				fallthrough // add `/` as an identifier Token
			case t.isIdentifier(r):
				t.unReadRune(r)
				id := t.readIdentifier()
				t.addToken(lookupIdent(id), id)
			default:
				unexpectedTokenMessage := fmt.Sprintf("[%s: line:%d: col:%d]: Unexpected Token %s", t.filname, t.line, t.col, string(r))
				t.addToken(Unexpected, unexpectedTokenMessage)
				return t.tokens, errors.New(unexpectedTokenMessage)
			}
		}
	}

	if len(t.tokens) > 0 && t.tokens[len(t.tokens)-1].tt == Unexpected {
		return t.tokens, errors.New(t.tokens[len(t.tokens)-1].data)
	} else {
		t.addToken(EOF, "EOF")
		return t.tokens, nil
	}
}

func (t *tokenizer) addCommaIfNeeded() {
	if len(t.tokens) > 0 {
		last := t.tokens[len(t.tokens)-1]
		if last.tt == IDENTIFIER || last.tt == INTEGER || last.tt == DECIMAL || last.tt == STRING || last.tt == NON_WORD_IDENTIFIER || last.tt == TYPE_IDENTIFIER || last.tt == RBRACE || last.tt == RPAREN || last.tt == RBRACKET {
			t.addToken(COMMA, ",")
		}
	}
}
