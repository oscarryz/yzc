package internal

import (
	"fmt"
)

type (
	expression interface {
		String() string
		stringValue() string
		dataType() Type
	}
	statement interface {
		value() string
		String() string
	}

	Boc struct {
		expressions []expression
		statements  []statement
	}

	BasicLit struct {
		pos       position
		tt        tokenType
		val       string
		basicType Type
	}
	ArrayLit struct {
		pos         position
		expressions []expression
		arrayType   *ArrayType
	}
	// DictLit represents a dictionary literal [k1:v1 k2:v2] or [String]Int for empty dictionary
	DictLit struct {
		pos      position
		dictType *DictType
		keys     []expression
		values   []expression
	}

	ShortDeclaration struct {
		pos      position
		variable *Variable
		value    expression
	}

	KeyValue struct {
		pos position
		key expression
		val expression
	}

	ParenthesisExp struct {
		lparen      position
		expressions []expression
		rparen      position
	}
	Variable struct {
		pos     position
		name    string // empty name means only return type is expressed, single uppercase name generic
		varType Type
	}
)

func (k *KeyValue) String() string {
	return prettyPrint(k, 0)
}

func (k *KeyValue) stringValue() string {
	return fmt.Sprintf("%s : %s", k.key.stringValue(), k.val.stringValue())
}

func (k *KeyValue) dataType() Type {
	return k.val.dataType()
}

func (v *Variable) String() string {
	return prettyPrint(v, 0)
}

func (v *Variable) stringValue() string {
	return v.name
}

func (v *Variable) dataType() Type {
	return v.varType
}

func (boc *Boc) String() string {
	return prettyPrint(boc, 0)
}

func (boc *Boc) stringValue() string {
	return "{}"
}

func (boc *Boc) dataType() Type {
	return newBocType()
}

func (bl *BasicLit) String() string {
	return prettyPrint(bl, 0)
}

func (bl *BasicLit) stringValue() string {
	return bl.val
}

func (bl *BasicLit) dataType() Type {
	return bl.basicType
}

func (al *ArrayLit) String() string {
	return prettyPrint(al, 0)
}

func (al *ArrayLit) stringValue() string {
	return al.String()
}

func (al *ArrayLit) dataType() Type {
	return al.arrayType
}

func (d *DictLit) String() string {
	return prettyPrint(d, 0)
}

func (d *DictLit) stringValue() string {
	return d.String()
}

func (d *DictLit) dataType() Type {
	return d.dictType
}

func (sd *ShortDeclaration) String() string {
	return prettyPrint(sd, 0)
}
func (sd *ShortDeclaration) stringValue() string {
	return fmt.Sprintf("%s : %s", sd.variable, sd.value.stringValue())
}

func (sd *ShortDeclaration) dataType() Type {
	return sd.value.dataType()
}
