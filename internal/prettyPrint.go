package internal

import (
	"fmt"
	"strings"
)

func prettyPrint(v interface{}, indent int) string {
	var sb strings.Builder

	switch v := v.(type) {
	// ASTs
	case *Boc:
		sb.WriteString(indentStr(indent) + "Boc(\n")
		if v.expressions == nil {
			sb.WriteString(prettyPrint("exprs: nil", indent+2))
		}
		if v.statements == nil {
			sb.WriteString(prettyPrint("stmts: nil", indent+2))
		}
		for _, exp := range v.expressions {
			sb.WriteString(prettyPrint(exp, indent+2))
		}
		for _, stmt := range v.statements {
			sb.WriteString(prettyPrint(stmt, indent+2))
		}
		sb.WriteString(indentStr(indent) + ")\n")
	case *BasicLit:
		sb.WriteString(indentStr(indent) + "BasicLit(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "tt: " + v.tt.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "val: " + v.val + "\n")
		sb.WriteString(indentStr(indent+2) + "basicType: " + prettyPrint(v.basicType, 0) + "\n")
		sb.WriteString(indentStr(indent) + ")\n")
	case *ArrayLit:
		sb.WriteString(indentStr(indent) + "ArrayLit(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "arrayType: " + prettyPrint(v.arrayType, 0))
		sb.WriteString(indentStr(indent+2) + "exps: [\n")
		for _, exp := range v.exps {
			sb.WriteString(prettyPrint(exp, indent+4))
		}
		sb.WriteString(indentStr(indent+2) + "]\n")
		sb.WriteString(indentStr(indent) + ")\n")
	case *DictLit:
		sb.WriteString(indentStr(indent) + "DictLit(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "dictType:\n" + prettyPrint(v.dictType, indent+4))
		sb.WriteString(indentStr(indent+2) + "keys: [\n")
		for _, exp := range v.keys {
			sb.WriteString(prettyPrint(exp, indent+4))

		}
		sb.WriteString(indentStr(indent+2) + "]\n")
		sb.WriteString(indentStr(indent+2) + "values: [\n")
		for _, exp := range v.values {
			sb.WriteString(prettyPrint(exp, indent+4))
		}
		sb.WriteString(indentStr(indent+2) + "]\n")
		sb.WriteString(indentStr(indent) + ")\n")
	case *ShortDeclaration:
		sb.WriteString(indentStr(indent) + "ShortDeclaration(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(prettyPrint(v.variable, indent+2))
		sb.WriteString(prettyPrint(v.val, indent+2))
		sb.WriteString(indentStr(indent) + ")\n")
	case *KeyValue:
		sb.WriteString(indentStr(indent) + "KeyValue(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(prettyPrint(v.key, indent+2))
		sb.WriteString(prettyPrint(v.val, indent+2))
		sb.WriteString(indentStr(indent) + ")\n")
	case *Variable:
		sb.WriteString(indentStr(indent) + "Var(\n")
		//sb.WriteString(indentStr(indent+2) + "pos: " + v.pos.String() + "\n")
		sb.WriteString(indentStr(indent+2) + "name: " + v.name + "\n")
		sb.WriteString(indentStr(indent+2) + "varType: " + prettyPrint(v.varType, 0) + "\n")
		sb.WriteString(indentStr(indent) + ")\n")
		// Types
	case *IntType:
		sb.WriteString(indentStr(indent) + "IntType")
	case *DecimalType:
		sb.WriteString(indentStr(indent) + "DecimalType")
	case *StringType:
		sb.WriteString(indentStr(indent) + "StringType")
	case *ArrayType:
		sb.WriteString(indentStr(indent) + "ArrayType(")
		sb.WriteString(prettyPrint(v.elemType, 0))
		sb.WriteString(")\n")
	case *DictType:
		sb.WriteString(indentStr(indent) + "DictType(\n")
		sb.WriteString(indentStr(indent+2) + "key:\n" + prettyPrint(v.keyType, indent+4))
		sb.WriteString(indentStr(indent+2) + "val:\n" + prettyPrint(v.valType, indent+4))
		sb.WriteString(indentStr(indent) + ")")
	case *BocType:
		sb.WriteString(indentStr(indent) + "BocType")
	case *TBD:
		sb.WriteString(indentStr(indent) + "TBD\n")
	// Add more cases for other types as needed
	case nil:
		sb.WriteString(indentStr(indent) + "nil\n")
	default:
		sb.WriteString(indentStr(indent) + fmt.Sprintf("%+v\n", v))
	}

	return sb.String()
}

func indentStr(indent int) string {
	return strings.Repeat("  ", indent)
}
