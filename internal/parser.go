package internal

type ast struct{}

func (a *ast) Bytes() []byte {
	return []byte(`package main
func main() {
    print("Hello world (from parser)")
}`)
}

type parser struct {
}

func parse(tokens []token) (*ast, error) {
	p := &parser{}
	a, e := p.parse(tokens)
	if e != nil {
		return nil, e
	}
	return a, nil
}
func (p *parser) parse(tokens []token) (*ast, error) {
	return &ast{}, nil

}
