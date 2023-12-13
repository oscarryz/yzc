package internal

import (
	"fmt"
)

type program struct {
	expressions []expression
	statements  []statement
}
type expression interface {
	value() string
}
type statement struct {
	name string
}
type ast struct {
}

func (a *ast) Bytes() []byte {
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

func parse(fileName string, tokens []token) (*ast, error) {
	p := newParser(fileName, tokens)
	a, e := p.parse()
	if e != nil {
		return nil, e
	}
	return a, nil
}

func newParser(fileName string, tokens []token) *parser {
	return &parser{
		fileName,
		tokens,
		0,
		&program{
			[]expression{},
			[]statement{}}}
}

func (p *parser) parse() (*ast, error) {
	p.program()
	return &ast{}, nil
}

// PROGRAM ::= EXPRESSION+ | STATEMENT*
func (p *parser) program() {

	token := p.nextToken()
	for token.tt != EOF {
		switch token.tt {
		// an expression
		case LPAREN, NUMBER, STRING, DECIMAL:
			p.addExpression(p.expression(token))
		// might be a statement
		case RETURN, CONTINUE, BREAK, TYPEIDENTIFIER:
			p.statement(token)
		}
		// if neither we need another token to know
		if token.tt == IDENTIFIER {
			nt := p.nextToken()
			switch nt.tt {
			case EOF:
				break
			case LPAREN, PERIOD:
				p.expression(token) // blockinvocation, member access
			case TYPEIDENTIFIER, COLON:
				p.statement(token, nt)
			}
		}
		token = p.nextToken()
	}
	fmt.Println("Program: {")
	for _, stmts := range p.prog.statements {
		fmt.Printf("stmt: %s\n", stmts)
	}
	for _, expr := range p.prog.expressions {
		fmt.Printf("expr: %s\n", expr.value())
	}
	fmt.Println("}")

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

func (p *parser) expression(t token) expression {
	switch t.tt {
	case NUMBER, DECIMAL, STRING:
		return &BasicLit{t.pos, t.tt, t.data}
	case LPAREN:
		e := p.expression(p.nextToken())
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

func (p *parser) addExpression(e expression) {
	fmt.Printf("adding expr: %#v\n", e)
	p.prog.expressions = append(p.prog.expressions, e)
}
func (p *parser) statement(tokens ...token) {
	fmt.Printf("stmt: %v\n", tokens)
}

func variableDefinition() {

}

func memberAccess() {

}

func blockInvocation() {

}

func variableDeclaration() {

}

func (p *parser) nextToken() token {
	if p.currentToken >= len(p.tokens) {
		return token{pos(0, 0), EOF, "EOF"}
	}
	t := p.tokens[p.currentToken]
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
