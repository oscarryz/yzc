package internal

import "fmt"

type (
	boc struct {
		name      string
		bocType   *blockType
		blockBody *blockBody
	}

	expression interface {
		value() string
	}
	statement interface {
		value() string
	}
	blockType struct {
		pos position
		tt  tokenType
		val string
	}
	blockBody struct {
		expressions []expression
		statements  []statement
	}

	BasicLit struct {
		pos position
		tt  tokenType
		val string
	}

	ParenthesisExp struct {
		lparen position
		exps   []expression
		rparen position
	}

	empty struct{}
)

func (b *boc) Print(indent string) {
	fmt.Printf("%s%s\n", indent, b.name)
	if b.bocType != nil {
		fmt.Printf("%s  bocType: %s\n", indent, b.bocType.val)
	}
	if b.blockBody != nil {
		fmt.Printf("%s  blockBody:\n", indent)
		for _, expr := range b.blockBody.expressions {
			fmt.Printf("%s    expression: %s\n", indent, expr.value())
		}
		for _, stmt := range b.blockBody.statements {
			fmt.Printf("%s    statement: %s\n", indent, stmt.value())
		}
	}
}

func (boc *boc) value() string {
	return boc.name
}
