package internal

import (
	"fmt"
)

type program struct {
	blockBody *blockBody
}
type blockBody struct {
	expressions []expression
	statements  []statement
}
type expression interface {
	value() string
}
type statement struct {
	name string
}

func (a *program) Bytes() []byte {
	return []byte(`package main
func main() {
    print("Hello world (from parser)")
}`)
}

type parser struct {
	fileName     string
	tokens       []token
	currentToken int
	prog         *program
}

func parse(fileName string, tokens []token) *program {
	p := newParser(fileName, tokens)
	a := p.parse()
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
		[]statement{},
	}
	token := p.nextToken()
	for token.tt != EOF {
		switch token.tt {
		// an expression
		case LPAREN, NUMBER, STRING, DECIMAL:
			bb.expressions = append(bb.expressions, p.expression())
		// might be a statement
		case RETURN, CONTINUE, BREAK, TYPEIDENTIFIER:
			p.statement()
		}
		// if neither we need another token to know
		if token.tt == IDENTIFIER {
			nt := p.nextToken()
			switch nt.tt {
			case EOF:
				break
			case LPAREN, PERIOD, RBRACE:
				p.expression() // blockinvocation, member access, array access
			case TYPEIDENTIFIER, COLON:
				p.statement()
			}
		}
		token = p.nextToken()
	}
	fmt.Println("Program: {")
	for _, stmts := range bb.statements {
		fmt.Printf("stmt: %s\n", stmts)
	}
	for _, expr := range bb.expressions {
		fmt.Printf("expr: %s\n", expr.value())
	}
	fmt.Println("}")
	return bb
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
	exp    expression
	rparen position
}

func (pe *ParenthesisExp) value() string {
	return fmt.Sprintf(`( %s )`, pe.exp.value())
}

type empty struct{}

func (e empty) value() string {
	return "<empty>"
}

var emptyExpression = empty{}

func (p *parser) expression() expression {
	t := p.token()
	switch t.tt {
	case NUMBER, DECIMAL, STRING:
		return &BasicLit{t.pos, t.tt, t.data}
	case LPAREN:
		e := p.expression()
		rparen := p.nextToken()
		if rparen.tt != RPAREN {
			logger.Fatalf(`Missing ")" at %v`, rparen.pos)
		}
		pe := &ParenthesisExp{
			t.pos,
			e,
			rparen.pos,
		}
		return pe
	}

	return emptyExpression
}

func (p *parser) statement() {
}

func variableDefinition() {

}

func memberAccess() {

}

func blockInvocation() {

}

func variableDeclaration() {

}

func (p *parser) token() token {
	return p.tokens[p.currentToken]
}
func (p *parser) nextToken() token {
	if p.currentToken >= len(p.tokens) {
		return token{pos(0, 0), EOF, "EOF"}
	}
	t := p.token()
	p.currentToken++
	return t
}

/*
src/garden/vegetables/Asparagus.yz
	garden/
	garden.yz
	vegetables.yz
	garden/
		vegetables/
			Asparagus.yz
*/
