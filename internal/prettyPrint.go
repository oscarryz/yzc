package internal

import (
	"fmt"
	"strings"
)

func prettyPrint(v interface{}, indent int) string {
	var sb strings.Builder

	switch v := v.(type) {
	case *boc:
		sb.WriteString(indentStr(indent) + "boc {\n")
		sb.WriteString(indentStr(indent+2) + "Name: " + v.Name + "\n")
		if v.bocType != nil {
			sb.WriteString(prettyPrint(v.bocType, indent+2))
		}
		if v.blockBody != nil {
			sb.WriteString(prettyPrint(v.blockBody, indent+2))
		}
		sb.WriteString(indentStr(indent) + "}\n")
	case *blockType:
		sb.WriteString(indentStr(indent) + "blockType {\n")
		sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "tt: " + v.tt.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "val: " + v.val + "\n")
		sb.WriteString(indentStr(indent) + "}\n")
	case *blockBody:
		sb.WriteString(indentStr(indent) + "blockBody {\n")
		for _, exp := range v.expressions {
			sb.WriteString(prettyPrint(exp, indent+2))
		}
		for _, stmt := range v.statements {
			sb.WriteString(prettyPrint(stmt, indent+2))
		}
		sb.WriteString(indentStr(indent) + "}\n")
	case *BasicLit:
		sb.WriteString(indentStr(indent) + "BasicLit {\n")
		sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "tt: " + v.tt.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "val: " + v.val + "\n")
		sb.WriteString(indentStr(indent) + "}\n")
	case *ArrayLit:
		sb.WriteString(indentStr(indent) + "ArrayLit {\n")
		sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "tt: " + v.tt.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "val: " + v.val + "\n")
		sb.WriteString(indentStr(indent+2) + "exps: [\n")
		for _, exp := range v.exps {
			sb.WriteString(prettyPrint(exp, indent+4))
		}
		sb.WriteString(indentStr(indent+2) + "]\n")
		sb.WriteString(indentStr(indent) + "}\n")
	case *empty:
		sb.WriteString(indentStr(indent) + "<empty>\n")
	// Add more cases for other types as needed
	default:
		sb.WriteString(indentStr(indent) + fmt.Sprintf("%+v\n", v))
	}

	return sb.String()
}

func indentStr(indent int) string {
	return strings.Repeat("  ", indent)
}
