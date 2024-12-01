package internal

import (
	"fmt"
	"strings"
)

type parser struct {
	tokens       []Token
	currentIndex int
	Token
	prog *Boc
}

func Parse(parents []string, tokens []Token) (*Boc, error) {
	p := newParser(tokens)
	debugCurrentTokenAtPosition(tokens, 0)
	leaf, e := p.boc()
	if e != nil {
		return nil, e
	}
	// Creates the intermediate parents.

	for i := len(parents) - 1; i >= 0; i-- {
		leaf = &Boc{
			expressions: []expression{
				&ShortDeclaration{
					pos: pos(0, 0),
					key: &BasicLit{
						pos: pos(0, 0),
						tt:  IDENTIFIER,
						val: parents[i],
					},
					val: leaf,
				},
			},
			statements: []statement{},
		}
	}
	return leaf, nil
}

func newParser(tokens []Token) *parser {
	return &parser{
		tokens,
		0,
		tokens[0],
		nil,
	}
}

// consume advances the parser by one Token.
func (p *parser) consume() {
	p.currentIndex++
	if p.currentIndex >= len(p.tokens) {
		p.Token = Token{p.pos, EOF, "EOF"}
	} else {
		p.Token = p.tokens[p.currentIndex]
	}
}

// expect returns true if the next Token is of type t.
func (p *parser) expect(t tokenType) error {
	if p.tt != t {
		return p.syntaxError(fmt.Sprintf("expected %s", t))
	}
	return nil
}

// parse parses the input file and returns the Boc.

//// Boc ::= block_body
//func (p *parser) boc() (*Boc, error) {
//	bb, e := p.boc()
//	if e != nil {
//		return nil, e
//	}
//	return &Boc{"", bb}, nil
//}

// block_body ::= (expression | statement) ((","|"\n") (expression | statement))* | ""
func (p *parser) boc() (*Boc, error) {
	bb := &Boc{
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
			if e == nil && stmt != nil {
				bb.statements = append(bb.statements, stmt)
			} else if e != nil {
				return nil, e
			}
		}

		switch p.tt {
		case COMMA:
			p.consume()
			continue
		case RBRACE:
			p.consume() // consume the RBRACE
			fallthrough
		case EOF:
			return bb, nil
		default:
			return nil, p.syntaxError("expected \",\" or \"}\". Got \"" + p.data + "\"")

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
	token := p.tt
	switch token {
	case INTEGER, DECIMAL, STRING, IDENTIFIER, NON_WORD_IDENTIFIER:
		return p.parseLiteralOrShortDeclaration()
	case LBRACE:
		p.consume() // consume the {
		return p.parseBlockLiteral()
	case RBRACE:
		return nil, nil
	case LBRACKET:
		ap := p.pos
		p.consume()
		return p.parseArrayOrDictionaryLiteral(ap)
	case EOF:
		return nil, nil
	default:
		return nil, nil
	}
}

// a, 1, "hello", 1.0
// a: 1, b: "hello", c: 1.0
func (p *parser) parseLiteralOrShortDeclaration() (expression, error) {
	token := p.tt
	ctp := p.pos
	ctd := p.data

	p.consume() // consume the  literal that brought us here
	bl := &BasicLit{ctp, token, ctd}
	if p.tt == COLON {
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
	return p.boc()
}

func (p *parser) parseArrayOrDictionaryLiteral(ap position) (expression, error) {
	if p.tt == RBRACKET {
		p.consume()
		return p.parseTypedArrayLiteral(ap)
	} else if p.tt == TYPE_IDENTIFIER {
		return p.parseEmptyDictionaryLiteral(ap)
	} else {
		return p.parseNonEmptyArrayOrDictionaryLiteral(ap)
	}
}

// e.g [] Int, current position is at the TYPE_IDENTIFIER
func (p *parser) parseTypedArrayLiteral(ap position) (expression, error) {
	err := p.expect(TYPE_IDENTIFIER)
	if err != nil {
		return nil, err
	}
	ct := p.tt
	//ctp := p.pos
	ctd := p.data
	p.consume()
	at := new(ArrayType)
	// TODO: need to handle arrays and dictionaries e.g. [][Int], [][String:Int]
	switch ct {
	case INTEGER:
		at.elemType = new(IntType)
	case DECIMAL:
		at.elemType = new(DecimalType)
	case STRING:
		at.elemType = new(StringType)
	case TYPE_IDENTIFIER:
		switch ctd {
		case "Int":
			at.elemType = new(IntType)
		case "Decimal":
			at.elemType = new(DecimalType)
		case "String":
			at.elemType = new(StringType)
		}
	default:
		at.elemType = new(TBD)
	}

	return &ArrayLit{ap, []expression{}, at}, nil
}

// [ String ] Int
func (p *parser) parseEmptyDictionaryLiteral(ap position) (expression, error) {
	keyType := p.data
	p.consume()
	if err := p.expect(RBRACKET); err != nil {
		return nil, err
	}
	p.consume() // consume the RBRACKET
	if err := p.expect(TYPE_IDENTIFIER); err != nil {
		return nil, err
	}
	valType := p.data
	p.consume()
	return &DictLit{ap, "[]", keyType, valType, []expression{}, []expression{}}, nil
}

// [ (expression (, )?)+ ]
// [ (expression : expression (, )?)+ ]
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

		ct := p.tt
		p.consume()

		if ct == COMMA {
			ct = p.tt
			if ct == RBRACKET {
				p.consume() // consume the RBRACKET
				if insideDict {
					return dl, nil
				}
				return createArrayLiteral(ap, exps)
			}
			continue
		}

		if ct == RBRACKET {
			if insideDict {
				return dl, nil
			}
			return createArrayLiteral(ap, exps)
		}

		return nil, p.syntaxError("expected \",\" or \"]\". Got " + p.data)
	}
}

func createArrayLiteral(ap position, exps []expression) (expression, error) {
	switch exps[0].(type) {
	case *ArrayLit:
		al, _ := exps[0].(*ArrayLit)
		et := al.arrayType.(*ArrayType).elemType
		at := new(ArrayType)
		at.elemType = et
		return &ArrayLit{ap, exps, at}, nil
	case *Boc:
		// I think I need the type here
		//bt := &Boc{[]expression{}, []statement{}}
		bt := new(BocType)
		return &ArrayLit{ap, exps, bt}, nil
	default:
		at := new(ArrayType)
		at.elemType = exps[0].dataType()
		return &ArrayLit{ap, exps, at}, nil

	}
}
func (p *parser) statement() (statement, error) {
	return nil, nil // fmt.Errorf("not implemented")
}

func (p *parser) syntaxError(message string) error {
	//p.currentIndex = len(p.tokens) - 1
	return fmt.Errorf("[%s] %s", p.pos, message)
}

func debugCurrentTokenAtPosition(tokens []Token, index int) string {
	ll := 1
	var builder strings.Builder
	builder.WriteString("Tokens:\n")
	builder.WriteString(fmt.Sprintf("%d: ", ll))
	for i, t := range tokens {
		if ll != t.pos.line {
			ll = t.pos.line
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("%d: ", ll))
		}
		if i == index {
			builder.WriteString(fmt.Sprintf(" ->  `%v`   ", t.data))
		} else {
			builder.WriteString(fmt.Sprintf(" %v ", t.data))
		}
	}
	if index >= len(tokens) {
		builder.WriteString(" <- [ beyond EOF ] ")
	}
	builder.WriteString("\n")
	return builder.String()
}
