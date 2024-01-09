package internal

import (
	"fmt"
)

type parser struct {
	fileName     string
	tokens       []token
	currentToken int
	prog         *program
}

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
func (p *parser) rewind(n int) {
	p.currentToken -= n
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
	tok := p.nextToken()
	for tok.tt != EOF {
		p.exploreExpression()
	}

	return bb, nil
}
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

func (p *parser) exploreExpression() ruleType {
	return p.explore(expressionTrie())
}

// TODO: repetition
// TODO: return the trie type
func (p *parser) explore(trie *Trie) ruleType {
	node := trie
	t := p.nextToken()
	for t.tt != EOF {
		//if node.isComplex {
		//	return p.explore(expressionTrie())
		//}
		var nt *Trie
		var ok bool
		if nt, ok = find(node.children, t.tt); !ok {
			// we didn't find it... are we at the end of the trie?
			if len(node.children) == 1 && node.children[0].isLeaf {
				return node.children[len(node.children)-1].tt
				//return node.tt // we are, so yes!
			}
			// we're not at the end, probably we have complex children?
			var keepGoing bool
			if ch, ok := filterComplex(node.children); ok {
				for _, ct := range ch {
					switch ct.tt {
					case BODY:
						tt := p.exploreBody()
						if tt != ILLEGAL {
							keepGoing = true
						} else {
							keepGoing = false
						}
						if keepGoing {
							nt = ct
						}
					case EXPR:
						// return the token first
						p.rewind(1)
						tt := p.exploreExpression()
						if tt != ILLEGAL {
							keepGoing = true
						} else {
							keepGoing = false
						}
						// TODO: have to rewind the whole expression
						//  temporarily we rewind only 1 because we know it was
						//  a 1 items expr
						p.rewind(1)
						if keepGoing {
							nt = ct
						}
					case STMT:
						tt := p.exploreStatement()
						if tt != ILLEGAL {
							keepGoing = true
						} else {
							keepGoing = false
						}
						if keepGoing {
							nt = ct
						}
					}
				}
			} else {
				keepGoing = false
			}
			// TODO: backtrack and try another branch
			if !keepGoing {
				return ILLEGAL
			}
		}
		node = nt
		t = p.nextToken()
	}
	// if we were at the end of the leaf...
	if len(node.children) == 1 && node.children[0].isLeaf {
		return node.children[len(node.children)-1].tt
	} else {
		return ILLEGAL
	}
}

func (p *parser) exploreStatement() ruleType {
	return ILLEGAL
}

func (p *parser) exploreBody() ruleType {
	return ILLEGAL
}

func (p *parser) consume() {
	p.currentToken++
}

/*
Some utility functions below for debugging
*/
func (a *program) Bytes() []byte {
	return []byte(`package main
func main() {
    print("Hello world (from parser)")
}`)
}

func (p *program) String() string {
	return fmt.Sprintf("blockBody: %#v", p.blockBody)
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

var emptyExpression = empty{}
