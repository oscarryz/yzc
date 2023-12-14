package internal

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	EOF tokenType = iota
	// delimiters
	LBRACE
	RBRACE
	COLON
	PERIOD
	SEMICOLON
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	EQL
	COMMA
	NEWLINE

	// literals
	NUMBER
	DECIMAL
	STRING

	// identifiers
	IDENTIFIER
	TYPEIDENTIFIER
	PUNCTIDENTIFIER
	BREAK
	CONTINUE
	RETURN
	ILLEGAL
)

func (tt tokenType) String() string {
	descriptions := [23]string{
		`EOF`,
		`{`,
		`}`,
		`:`,
		`.`,
		`;`,
		`(`,
		`)`,
		`[`,
		`]`,
		`=`,
		`,`,
		"NEWLINE",
		`num`,
		`dec`,
		`str`,
		`id`,
		`tid`,
		`punctid`,
		"BREAK",
		"CONTINUE",
		"RETURN",
		"ILLEGAL",
	}
	return descriptions[tt]
}

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

func (p position) String() string {
	return fmt.Sprintf("line: %d col: %d", p.line, p.col)
}

func pos(line, col int) position {
	return position{line, col}
}

func (t token) String() string {
	switch t.tt {
	case NUMBER, DECIMAL, STRING, IDENTIFIER, TYPEIDENTIFIER:
		return fmt.Sprintf("%s:%s ", t.tt, t.data)
	default:
		return fmt.Sprintf("%#s ", t.tt)
	}

}

type tokenizer struct {
	content   string
	tokens    []token
	col       int
	line      int
	pos       int
	keepGoing bool
}

func tokenize(fileName string, content string) ([]token, error) {
	t := &tokenizer{content, []token{}, 0, 1, 0, true}
	tokens, e := t.tokenize()
	printTokens(tokens)
	return tokens, e
}
func printTokens(tokens []token) {
	ll := 1
	fmt.Printf("Tokens: \n%d: ", ll)
	for _, t := range tokens {
		if ll != t.pos.line {
			ll = t.pos.line
			fmt.Println()
			fmt.Printf("%d: ", ll)
		}
		fmt.Printf("%v", t)
	}
	fmt.Println()
}

func (t *tokenizer) addToken(tt tokenType, data string) {
	t.tokens = append(t.tokens, token{pos(t.line, t.col), tt, data})
}
func (t *tokenizer) nextRune() rune {
	r, w := utf8.DecodeRuneInString(t.content[t.pos:])
	if r == utf8.RuneError {
		// we're done tokenizing
		t.keepGoing = false
	}
	t.pos += w
	t.col++
	return r
}
func (t *tokenizer) unReadRune(r rune) {
	t.pos -= utf8.RuneLen(r)
	t.col--
}
func (t *tokenizer) peek() rune {
	r := t.nextRune()
	t.unReadRune(r)
	return r
}

func (t *tokenizer) skipComment() {
	for r := t.nextRune(); r != '\n'; /*&& r != utf8.RuneError*/ r = t.nextRune() {
	}
}
func (t *tokenizer) skipMultilineComment() {
	r := t.nextRune()
	for {
		if r == '*' && t.peek() == '/' {
			t.nextRune()
			return
		}
		r = t.nextRune()
	}
}
func lookupIdent(identifier string) tokenType {
	runes := []rune(identifier)
	if unicode.IsUpper(runes[0]) {
		return TYPEIDENTIFIER
	}
	switch identifier {
	case "break":
		return BREAK
	case "continue":
		return CONTINUE
	case "return":
		return RETURN
	default:
		// all punct
		allPunctuation := true
		for _, r := range runes {
			allPunctuation = allPunctuation && !unicode.IsLetter(r)
		}
		if allPunctuation {
			return PUNCTIDENTIFIER
		}
	}
	return IDENTIFIER

}

func (t *tokenizer) addStringLiteral() {

	// todo: scape literals
	opening := t.nextRune()
	r := t.nextRune()
	builder := strings.Builder{}
	for r != opening {
		builder.WriteRune(r)
		if r == '\n' {
			t.line++
			t.col = 0
		}
		r = t.nextRune()
	}
	t.addToken(STRING, builder.String())
}

func (t *tokenizer) isIdentifier(r rune) bool {
	return !unicode.IsSpace(r) &&
		unicode.IsPrint(r) &&
		!unicode.IsDigit(r) &&
		!strings.ContainsRune("{}[]().,:;\"'`", r)
}

func (t *tokenizer) readIdentifier() string {
	builder := strings.Builder{}
	r := t.nextRune()
	for (unicode.IsPrint(r) ||
		unicode.IsDigit(r)) &&
		!strings.ContainsRune("{}[]().,:;\"'`", r) &&
		!unicode.IsSpace(r) &&
		r != utf8.RuneError {
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
	for /*r != utf8.RuneError &&*/ unicode.IsDigit(r) || r == '.' {
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
	tt := NUMBER
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
		case '=':
			t.addToken(EQL, "=")
		case ',':
			t.addToken(COMMA, ",")
		case '/':
			if t.peek() == '/' {
				t.nextRune()
				t.skipComment()
			} else if t.peek() == '*' {
				t.nextRune()
				t.skipMultilineComment()
			}
		case '-':
			if unicode.IsDigit(t.peek()) {
				t.unReadRune(r)
				t.addNegativeNumber()
			}
		case '"', '`', '\'':
			t.unReadRune(r)
			t.addStringLiteral()
		default:
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
