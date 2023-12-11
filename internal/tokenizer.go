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
	BREAK
	CONTINUE
	RETURN
	ILLEGAL
)

func (tt tokenType) String() string {
	descriptions := [21]string{
		"EOF",
		"LBRACE",
		"RBRACE",
		"COLON",
		"PERIOD",
		"SEMICOLON",
		"LPAREN",
		"RPAREN",
		"LBRACKET",
		"RBRACKET",
		"EQL",
		"COMMA",
		"NEWLINE",
		"NUMBER",
		"DECIMAL",
		"STRING",
		"IDENTIFIER",
		"BREAK",
		"CONTINUE",
		"RETURN",
		"ILLEGAL",
	}
	return descriptions[tt]
}

type tokenType uint

type token struct {
	line int
	col  int
	tt   tokenType
	data string
}

func (t token) String() string {
	return fmt.Sprintf("token: l:%d c:%d %s:%s", t.line, t.col, t.tt, t.data)

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
	// todo: create a first token with the file name.
	//  e.g. a.yz -> token: tt:identifier data:'a'
	//  maybe somewhere else
	fmt.Printf("Tokenizing: %s\n", fileName)
	t := &tokenizer{content, []token{}, 1, 1, 0, true}
	tokens, e := t.tokenize()
	return tokens, e
}

func (t *tokenizer) addToken(tt tokenType, data string) {
	t.tokens = append(t.tokens, token{t.line, t.col, tt, data})
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
	switch identifier {
	case "break":
		return BREAK
	case "continue":
		return CONTINUE
	case "return":
		return RETURN
	default:
		return IDENTIFIER
	}
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
	return unicode.IsPrint(r) &&
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
	for t.keepGoing {
		r := t.nextRune()

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
		case '\n':
			t.addToken(NEWLINE, "\n")
			t.line++
			t.col = 1
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
	return t.tokens, nil
}
