package internal

import (
	"fmt"
	"strings"
)

type parser struct {
	fileName     string
	tokens       []token
	currentToken int
	prog         *boc
}

func Parse(fileName string, tokens []token) (*boc, error) {
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

// token returns the current token without advancing the parser.
func (p *parser) token() token {
	return p.tokenPlus(0)
}

// nextToken returns the current token and advances the parser.
func (p *parser) nextToken() token {
	if p.currentToken >= len(p.tokens) {
		return token{pos(0, 0), EOF, "EOF"}
	}
	t := p.token()
	p.currentToken++
	return t
}

// rewind rewinds the parser by n tokens.
func (p *parser) rewind(n int) {
	p.currentToken -= n
}

// tokenPlus returns the token i tokens ahead of the current token.
func (p *parser) tokenPlus(i int) token {
	return p.tokens[p.currentToken+i]
}

// consume advances the parser by one token.
func (p *parser) consume() {
	p.currentToken++
}

// expect returns true if the next token is of type t.
func (p *parser) expect(t tokenType) bool {
	return p.nextToken().tt == t
}

// parse parses the input file and returns the boc.
func (p *parser) parse() (*boc, error) {
	// splits the file Name into directories and file Name without extension
	parts := strings.Split(p.fileName, "/")
	fileNameWithoutExtension := strings.Split(parts[len(parts)-1], ".")[0]

	leaf, e := p.boc()
	if e != nil {
		return nil, e
	}
	leaf.Name = fileNameWithoutExtension
	// Creates the parent bocs
	// for a/b/c.yz will creates
	// boc{ Name: "a", bocType: nil, blockBody:
	//		boc: { Name: "b", bocType: nil, blockBody:
	//			boc: { Name: "c", bocType: nil, blockBody: nil } } }
	for i := len(parts) - 2; i >= 0; i-- {
		leaf = &boc{
			Name:    parts[i],
			bocType: nil,
			blockBody: &blockBody{
				[]expression{leaf},
				[]statement{},
			},
		}
	}

	return leaf, nil
}

// boc ::= block_body
func (p *parser) boc() (*boc, error) {
	bb, e := p.blockBody()
	if e != nil {
		return nil, e
	}
	return &boc{"", nil, bb}, nil
}

// block_body ::= (expression | statement) ("," (expression | statement))* | ""
func (p *parser) blockBody() (*blockBody, error) {
	bb := &blockBody{
		[]expression{},
		[]statement{},
	}
	// Checks if there is an expression or a statement
	// if there's an expression adds it to the expressions slice
	// if there's a statement adds it to the statements slice
	// if there's a comma, it continues to parse the next expression or statement
	// Checks if there is an expression or a statement
	for {
		expr, e := p.expression()
		if e == nil {
			bb.expressions = append(bb.expressions, expr)
		} else {
			stmt, e := p.statement()
			if e == nil {
				bb.statements = append(bb.statements, stmt)
			} else {
				return nil, e
			}
		}

		if p.token().tt != COMMA {
			break
		}
		p.consume() // consume the comma
	}
	return bb, nil



}

//   expression ::= block_invocation
//    | method_invocation
//    | parenthesized_expressions
//    | type_instantiation
//    | array_access
//    | dictionary_access
//    | member_access
//    | literal
//    | variable
//    | assignment
//    | variable_short_definition

func (p *parser) expression() (expression, error) {

	token := p.token()
	switch token.tt {
	// literal
	case INTEGER, DECIMAL, STRING:
		p.consume()
		return &BasicLit{token.pos, token.tt, token.data}, nil
	case LBRACE:
		p.consume()
		bb, e := p.blockBody()
		if e != nil {
			return nil, e
		}
		if !p.expect(RBRACE) {
			return nil, p.syntaxError("expected }")
		}
		return &boc{"", nil, bb}, nil
	case RBRACE:
		return &empty{}, nil
	case EOF:
		return &empty{}, nil
	}
	return nil, nil // p.syntaxError("expected expression")
}

func (p *parser) statement() (statement, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *parser) syntaxError(message string) error {
	p.currentToken = len(p.tokens) - 1
	return fmt.Errorf("[%s %s] %s", p.fileName, p.token().pos, message)
}

func (p *boc) String() string {
	return p.value()
}

func (bb *blockBody) String() string {
	return fmt.Sprintf("expressions: %#v statements: %#v", bb.expressions, bb.statements)
}

func (bl BasicLit) value() string {
	return bl.val
}

func (e empty) value() string {
	return "<empty>"
}
