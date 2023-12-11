package internal

import "bytes"

type tokenizer struct {
	tokens []token
}

func (t *tokenizer) add(tok token) {
	t.tokens = append(t.tokens, tok)
}

type tokenType uint
type token struct {
	position int
	line     int
	tt       tokenType
	data     string
}

func newToken(pos, line int, tt tokenType, data string) token {
	return token{pos, line, tt, data}
}

const (
	LBRACE tokenType = iota
	RBRACE
)

func tokenize(e error, content []byte) []token {
	t := &tokenizer{}
	tokens, e := t.tokenize(content)
	return tokens
}

func (t *tokenizer) tokenize(content []byte) ([]token, error) {
	buffer := bytes.NewBuffer(content)
	i := 0
	l := 0
	r, _, err := buffer.ReadRune()
	if err != nil {
		return nil, err
	}
	switch r {
	case '{':
		t.add(newToken(i, l, LBRACE, "{"))
	case '}':
		t.add(newToken(i, l, RBRACE, "}"))
	default:

	}
	return t.tokens, nil
}
