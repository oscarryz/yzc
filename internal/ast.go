package internal

import (
	"fmt"
	"strings"
)

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

func (b *boc) Print(indent string) string {
	// string buffer
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s%s\n", indent, b.name))
	if b.bocType != nil {
		builder.WriteString(fmt.Sprintf("%s  bocType: %s\n", indent, b.bocType.val))
	}
	if b.blockBody != nil {
		builder.WriteString(fmt.Sprintf("%s  blockBody:\n", indent))
		for _, expr := range b.blockBody.expressions {
			builder.WriteString(fmt.Sprintf("%s    expression: %s\n", indent, expr.value()))
		}
		for _, stmt := range b.blockBody.statements {

			builder.WriteString(fmt.Sprintf("%s    statement: %s\n", indent, stmt.value()))
		}
	}
	return builder.String()
}

func (boc *boc) value() string {
	return boc.Print("")
}
