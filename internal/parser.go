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
	// literal
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NONWORDIDENTIFIER:
		//// if : then is short declaration
		if p.peek().tt == COLON {
			p.consume()
			p.consume()
			val, e := p.expression()
			if e != nil {
				return nil, e
			}
			return &ShortDeclaration{token.pos, &BasicLit{token.pos, token.tt, token.data}, val}, nil
		}
		p.consume()
		return &BasicLit{token.pos, token.tt, token.data}, nil
	// block literal
	case LBRACE:
		p.consume()
		bb, e := p.blockBody() // will consume the RBRACE if found
		if e != nil {
			return nil, e
		}
		return &Boc{"", nil, bb}, nil
	// Array or Dictionary literal
	case RBRACE:
		return &empty{}, nil

		// [ for array or dictionary
	case LBRACKET:
		ap := p.token().pos
		// closing bracket means empty array
		if p.peek().tt == RBRACKET {
			// eg [] Int
			p.consume()
			p.consume()
			// we want type identifier e.g. []Int
			if e := p.expect(TYPEIDENTIFIER); e != nil {
				return nil, e
			}
			return &ArrayLit{ap, "[]", []expression{}}, nil
		} else if p.peek().tt == TYPEIDENTIFIER {
			// Empty dictionary literal: [String]Int
			p.consume()
			keyType := p.nextToken().data
			if e := p.expect(RBRACKET); e != nil {
				return nil, e
			}
			if e := p.expect(TYPEIDENTIFIER); e != nil {
				return nil, e
			}
			p.rewind(1)
			valType := p.nextToken().data
			return &DictLit{ap, "[]", keyType, valType, []expression{}, []expression{}}, nil
		} else {
			// Could be array literal or dictionary literal
			// eg [1 2 3]
			p.consume() // consume [
			var exps []expression
			dl := &DictLit{ap, "[]", "", "", []expression{}, []expression{}}
			insideDict := false
			for {
				// first element or first key
				expr, e := p.expression()

				// if the expr is a short declaration, then it's a dictionary
				if sd, ok := expr.(*ShortDeclaration); ok {
					insideDict = true
					key := sd.key
					val := sd.val
					dl.keys = append(dl.keys, key)
					dl.values = append(dl.values, val)
				} else if expr != nil {
					// we are inside an array
					exps = append(exps, expr)
				} else {
					return nil, e
				}

				//if t.tt == COLON {
				//	// we are inside a dictionary
				//	insideDict = true
				//
				//	p.consume()
				//	key := expr
				//	val, e := p.expression()
				//	if e != nil {
				//		return nil, e
				//	}
				//	if val == nil {
				//		return nil, p.syntaxError("expected value")
				//	}
				//	dl.keys = append(dl.keys, key)
				//	dl.values = append(dl.values, val)
				//} else {
				// we are inside an array
				//exps = append(exps, expr)
				//}

				if p.token().tt == RBRACKET && insideDict {
					p.consume()
					return dl, nil
				} else if p.token().tt == RBRACKET {
					p.consume()
					return &ArrayLit{ap, "[]", exps}, nil
				}
			}
		}
	case EOF:
		return &empty{}, nil
	default:
		return nil, nil
	}
}

func (p *parser) statement() (statement, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *parser) syntaxError(message string) error {
	//p.currentToken = len(p.tokens) - 1
	return fmt.Errorf("[%s %s] %s", p.fileName, p.token().pos, message)
}
