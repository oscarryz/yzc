package internal

import (
	"fmt"
	"strings"
)

type parser struct {
	fileName     string
	tokens       []token
	currentIndex int
	prog         *Boc
}

func Parse(fileName string, tokens []token) (*Boc, error) {
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
func (p *parser) currentToken() tokenType {
	return p.tokens[p.currentIndex].tt
}
func (p *parser) currentTokenData() string {
	return p.tokens[p.currentIndex].data
}
func (p *parser) currentTokenPos() position {
	return p.tokens[p.currentIndex].pos
}

// Returns the next token without advancing the parser.
func (p *parser) peek() token {
	return p.tokens[p.currentIndex+1]
}

// consume advances the parser by one token.
func (p *parser) consume() {
	p.currentIndex++
}

// expect returns true if the next token is of type t.
func (p *parser) expect(t tokenType) error {
	if p.currentToken() != t {
		return p.syntaxError(fmt.Sprintf("expected %s", t))
	}
	return nil
}

// parse parses the input file and returns the Boc.
func (p *parser) parse() (*Boc, error) {
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
	// Boc{ Name: "a", bocType: nil, blockBody:
	//		Boc: { Name: "b", bocType: nil, blockBody:
	//			Boc: { Name: "c", bocType: nil, blockBody: nil } } }
	for i := len(parts) - 2; i >= 0; i-- {
		leaf = &Boc{
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

// Boc ::= block_body
func (p *parser) boc() (*Boc, error) {
	bb, e := p.blockBody()
	if e != nil {
		return nil, e
	}
	return &Boc{"", nil, bb}, nil
}

// block_body ::= (expression | statement) ((","|"\n") (expression | statement))* | ""
func (p *parser) blockBody() (*blockBody, error) {
	bb := &blockBody{
		[]expression{},
		[]statement{},
	}
	// Checks if there is an expression or a statement
	// if there's an expression adds it to the expressions slice
	// if there's a statement adds it to the statements slice
	// if there's a comma or newline, it continues to parse the next expression or statement
	// Checks if there is an expression or a statement
	for {
		expr, e := p.expression()
		if e == nil && expr != nil {
			bb.expressions = append(bb.expressions, expr)
		} else if e != nil {
			return nil, e
		} else {
			stmt, e := p.statement()
			if e == nil {
				bb.statements = append(bb.statements, stmt)
			} else {
				return nil, e
			}
		}

		// p.token
		switch p.currentToken() {
		case COMMA:
			p.consume()
			continue
		case RBRACE:
			p.consume() // consume the RBRACE
			fallthrough
		case EOF:
			return bb, nil
		default:
			return nil, p.syntaxError("expected \",\" or \"}\". Got \"" + p.currentTokenData() + "\"")

		}
	}
	//return bb, nil

}

// expression ::= block_invocation
//
//	| method_invocation
//	| parenthesized_expressions
//	| type_instantiation
//	| array_access
//	| dictionary_access
//	| member_access
//	| literal
//	| variable
//	| assignment
//	| variable_short_definition
func (p *parser) expression() (expression, error) {
	token := p.currentToken()
	switch token {
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NONWORDIDENTIFIER:
		return p.parseLiteralOrShortDeclaration()
	case LBRACE:
		p.consume() // consume the {
		return p.parseBlockLiteral()
	case RBRACE:
		return &empty{}, nil
	case LBRACKET:
		ap := p.currentTokenPos()
		p.consume()
		return p.parseArrayOrDictionaryLiteral(ap)
	case EOF:
		return &empty{}, nil
	default:
		return nil, nil
	}
}

// a, 1, "hello", 1.0
// a: 1, b: "hello", c: 1.0
func (p *parser) parseLiteralOrShortDeclaration() (expression, error) {
	token := p.currentToken()
	ctp := p.currentTokenPos()
	ctd := p.currentTokenData()

	p.consume() // consume the  literal that brought us here
	bl := &BasicLit{ctp, token, ctd}
	if p.currentToken() == COLON {
		p.consume() // consume the COLON
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &ShortDeclaration{ctp, bl, val}, nil
	}
	return bl, nil
}

func (p *parser) parseBlockLiteral() (expression, error) {
	bb, err := p.blockBody()
	return &Boc{"", nil, bb}, err
}

func (p *parser) parseArrayOrDictionaryLiteral(ap position) (expression, error) {
	if p.currentToken() == RBRACKET {
		p.consume()
		return p.parseTypedArrayLiteral(ap)
	} else if p.currentToken() == TYPEIDENTIFIER {
		return p.parseEmptyDictionaryLiteral(ap)
	} else {
		return p.parseNonEmptyArrayOrDictionaryLiteral(ap)
	}
}

// e.g [] Int, current position is at the TYPEIDENTIFIER
func (p *parser) parseTypedArrayLiteral(ap position) (expression, error) {
	err := p.expect(TYPEIDENTIFIER)
	if err != nil {
		return nil, err
	}
	ct := p.currentToken()
	ctp := p.currentTokenPos()
	ctd := p.currentTokenData()
	p.consume()
	return &ArrayLit{ap, &BasicLit{ctp, ct, ctd}, []expression{}}, nil
}

// [ String ] Int
func (p *parser) parseEmptyDictionaryLiteral(ap position) (expression, error) {
	keyType := p.currentTokenData()
	p.consume()
	if err := p.expect(RBRACKET); err != nil {
		return nil, err
	}
	p.consume() // consume the RBRACKET
	if err := p.expect(TYPEIDENTIFIER); err != nil {
		return nil, err
	}
	valType := p.currentTokenData()
	p.consume()
	return &DictLit{ap, "[]", keyType, valType, []expression{}, []expression{}}, nil
}

// [ expression ("," expression)* ]
// [ key ":" value ("," key ":" value)* ]
func (p *parser) parseNonEmptyArrayOrDictionaryLiteral(ap position) (expression, error) {
	var exps []expression
	dl := &DictLit{ap, "[]", "", "", []expression{}, []expression{}}
	insideDict := false

	for {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if sd, ok := expr.(*ShortDeclaration); ok {
			insideDict = true
			dl.keys = append(dl.keys, sd.key)
			dl.values = append(dl.values, sd.val)
		} else if expr != nil {
			exps = append(exps, expr)
		} else {
			return nil, err
		}

		ct := p.currentToken()
		p.consume()

		if ct == COMMA {
			ct = p.currentToken()
			if ct != RBRACKET {
				continue
			}
		}

		if ct == RBRACKET {
			if insideDict {
				p.consume()
				return dl, nil
			}
			return createArrayLiteral(ap, exps)
		}

		return nil, p.syntaxError("expected \",\" or \"]\". Got " + p.currentTokenData())
	}
}

func createArrayLiteral(ap position, exps []expression) (expression, error) {
	switch exps[0].(type) {
	case *ArrayLit:
		al, _ := exps[0].(*ArrayLit)
		at := &ArrayLit{al.pos, al.arrayType, []expression{}}
		return &ArrayLit{ap, at, exps}, nil
	case *Boc:
		bl, _ := exps[0].(*Boc)
		bt := &Boc{bl.Name, nil, &blockBody{[]expression{}, []statement{}}}
		return &ArrayLit{ap, bt, exps}, nil
	default:
		return &ArrayLit{ap, exps[0], exps}, nil

	}
}
func (p *parser) statement() (statement, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *parser) syntaxError(message string) error {
	//p.currentIndex = len(p.tokens) - 1
	return fmt.Errorf("[%s %s] %s", p.fileName, p.currentTokenPos(), message)
}
