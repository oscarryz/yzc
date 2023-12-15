package internal

import "fmt"

func parse(fileName string, tokens []token) *program {
	p := newParser(fileName, tokens)
	a := p.parse()
	fmt.Printf("%#v", a)
	return a
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

func (p *parser) parse() *program {
	return p.program()
}

// program ::= block_body
func (p *parser) program() *program {
	return &program{p.blockBody()}
}

// block_body::= expression+ | statement*
func (p *parser) blockBody() *blockBody {
	bb := &blockBody{
		[]expression{},
		[]*statement{},
	}
	token := p.token()
	for token.tt != EOF {
		e := p.attemptExpression()
		if e != nil {
			bb.expressions = append(bb.expressions, e)
		} else {
			s := p.attemptStatement()
			bb.statements = append(bb.statements, s)
			if s != nil {
			} else {
				p.syntaxError("BlockBody should contain expressions or statements")
			}
		}
		p.consume()
		token = p.token()
	}
	return bb
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

type blockBody struct {
	expressions []expression
	statements  []*statement
}
type expression interface {
	value() string
}
type statement struct {
	name string
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
	return emptyExpression
}

func (p *parser) statement() {
}

func (p *parser) syntaxError(message string) {
	p.currentToken = len(p.tokens) - 1
	logger.Fatalf("[%s:%s] %s", p.fileName, p.token().pos, message)
}

func (p *parser) attemptExpression() expression {
	token := p.token()
	switch token.tt {
	// literal
	case INTEGER, DECIMAL, STRING:
		return &BasicLit{token.pos, token.tt, token.data}

	}
	return nil
}

func (p *parser) attemptStatement() *statement {
	return nil
}

func (p *parser) consume() {
	//if p.currentToken >= len(p.tokens) {
	//	return token{pos(0, 0), EOF, "EOF"}
	//}
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
