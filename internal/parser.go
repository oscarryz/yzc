package internal

import "fmt"

func parse(fileName string, tokens []token) (*program, error) {
	p := newParser(fileName, tokens)
	return p.parse()
}

func newParser(fileName string, tokens []token) *parser {
	return &parser{
		fileName,
		tokens,
		0,
		nil,
	}
}

type parser struct {
	fileName     string
	tokens       []token
	currentToken int
	prog         *program
}

func (p *parser) token() token {
	return p.tokenPlus(0)
}
func (p *parser) nextToken() token {
	if p.currentToken >= len(p.tokens) {
		return token{pos(0, 0), EOF, "EOF"}
	}
	t := p.token()
	p.currentToken++
	return t
}

func (p *parser) tokenPlus(i int) token {
	return p.tokens[p.currentToken+i]
}

func (p *parser) parse() (*program, error) {
	return p.program()
}

// program ::= block_body
func (p *parser) program() (*program, error) {
	body, err := p.blockBody()
	return &program{body}, err
}

// block_body::= expression+ | statement*
func (p *parser) blockBody() (*blockBody, error) {
	bb := &blockBody{
		[]expression{},
		[]statement{},
	}
	token := p.token()
	for token.tt != EOF {
		if p.exploreExpression() {
			bb.expressions = append(bb.expressions, p.expression())
		} else if p.exploreStatement() {
			p.statement()
		} else {
			return nil, p.syntaxError("BlockBody should contain expressions or statements")
		}
		p.consume()
		token = p.token()
	}
	return bb, nil
}

type program struct {
	blockBody *blockBody
}

func (a *program) Bytes() []byte {
	return []byte(`package main
func main() {
    print("Hello world (from parser)")
}`)
}

func (p *program) String() string {
	return fmt.Sprintf("blockBody: %#v", p.blockBody)
}

type blockBody struct {
	expressions []expression
	statements  []statement
}

func (bb *blockBody) String() string {
	return fmt.Sprintf("expressions: %#v statements: %#v", bb.expressions, bb.statements)
}

type expression interface {
	value() string
}
type statement interface {
	value() string
}

func (bl BasicLit) value() string {
	return bl.val
}

type BasicLit struct {
	pos position
	tt  tokenType
	val string
}

type ParenthesisExp struct {
	lparen position
	exps   []expression
	rparen position
}

type empty struct{}

func (e empty) value() string {
	return "<empty>"
}

var emptyExpression = empty{}

func (p *parser) expression() expression {

	token := p.token()
	switch token.tt {
	// literal
	case INTEGER, DECIMAL, STRING:
		return &BasicLit{token.pos, token.tt, token.data}

	}
	return emptyExpression
}

func (p *parser) statement() statement {
	return nil
}

func (p *parser) syntaxError(message string) error {
	p.currentToken = len(p.tokens) - 1
	return fmt.Errorf("[%s %s] %s", p.fileName, p.token().pos, message)
}

func (p *parser) exploreExpression() bool {
	token := p.token()
	switch token.tt {
	case INTEGER, DECIMAL, STRING:
		return true

	}
	return false

}

func (p *parser) exploreStatement() bool {
	return false
}

func (p *parser) consume() {
	p.currentToken++
}

func variableDefinition() {

}

func memberAccess() {

}

func blockInvocation() {

}

func variableDeclaration() {

}
