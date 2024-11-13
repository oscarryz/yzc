package internal

import (
	"fmt"
	"strings"
)

type parser struct {
	fileName     string
	tokens       []token
	currentToken int
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
func (p *parser) expect(t tokenType) error {
	if p.nextToken().tt != t {
		return p.syntaxError(fmt.Sprintf("expected %s", t))
	}
	return nil
}

func (p *parser) peek() token {
	return p.tokenPlus(1)
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

		switch p.token().tt {
		case COMMA:
			p.consume()
			continue
		case RBRACE:
			p.consume() // consume the RBRACE
			fallthrough
		case EOF:
			return bb, nil
		default:
			return nil, p.syntaxError("expected \",\", NEWLINE or RBRACE. Got \"" + p.token().data + "\"")

		}
	}
	//return bb, nil

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
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NONWORDIDENTIFIER:
		return p.parseLiteralOrShortDeclaration(token)
	case LBRACE:
		return p.parseBlockLiteral()
	case RBRACE:
		return &empty{}, nil
	case LBRACKET:
		return p.parseArrayOrDictionaryLiteral()
	case EOF:
		return &empty{}, nil
	default:
		return nil, nil
	}
}

func (p *parser) parseLiteralOrShortDeclaration(token token) (expression, error) {
	if p.peek().tt == COLON {
		p.consume()
		p.consume()
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &ShortDeclaration{token.pos, &BasicLit{token.pos, token.tt, token.data}, val}, nil
	}
	p.consume()
	return &BasicLit{token.pos, token.tt, token.data}, nil
}

func (p *parser) parseBlockLiteral() (expression, error) {
	p.consume()
	bb, err := p.blockBody()
	if err != nil {
		return nil, err
	}
	return &Boc{"", nil, bb}, nil
}

func (p *parser) parseArrayOrDictionaryLiteral() (expression, error) {
	ap := p.token().pos
	p.consume()

	// s/p.token()/peek
	if p.token().tt == RBRACKET {
		p.consume()
		err := p.expect(TYPEIDENTIFIER)
		if err != nil {
			return nil, err
		}
		return p.parseTypedArrayLiteral(ap)
	} else if p.token().tt == TYPEIDENTIFIER {
		return p.parseEmptyDictionaryLiteral(ap)
	} else {
		return p.parseNonEmptyArrayOrDictionaryLiteral(ap)
	}
}

func (p *parser) parseTypedArrayLiteral(ap position) (expression, error) {
	//p.consume()
	p.rewind(1)
	ct := p.token()
	p.consume()
	return &ArrayLit{ap, &BasicLit{ct.pos, ct.tt, ct.data}, []expression{}}, nil
}

func (p *parser) parseEmptyDictionaryLiteral(ap position) (expression, error) {
	//p.consume()
	keyType := p.nextToken().data
	if err := p.expect(RBRACKET); err != nil {
		return nil, err
	}
	if err := p.expect(TYPEIDENTIFIER); err != nil {
		return nil, err
	}
	p.rewind(1)
	valType := p.nextToken().data
	return &DictLit{ap, "[]", keyType, valType, []expression{}, []expression{}}, nil
}

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

		if p.token().tt == RBRACKET {
			p.consume()
			if insideDict {
				return dl, nil
			}

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
	}
}
func (p *parser) statement() (statement, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *parser) syntaxError(message string) error {
	//p.currentToken = len(p.tokens) - 1
	return fmt.Errorf("[%s %s] %s", p.fileName, p.token().pos, message)
}
