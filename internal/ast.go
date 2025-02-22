package internal

import (
	"fmt"
)

type (
	expression interface {
		String() string
		value() string
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
		pos       position
		exps      []expression
		arrayType *ArrayType
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
		val      expression
	}

	KeyValue struct {
		pos position
		key expression
		val expression
	}

	ParenthesisExp struct {
		lparen position
		exps   []expression
		rparen position
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

func (k *KeyValue) value() string {
	return fmt.Sprintf("%s : %s", k.key.value(), k.val.value())
}

func (k *KeyValue) dataType() Type {
	return k.val.dataType()
}

func (v *Variable) String() string {
	return prettyPrint(v, 0)
}

func (v *Variable) value() string {
	return v.name // TODO: what's the value of a variable?
}

func (v *Variable) dataType() Type {
	return v.varType
}

func (boc *Boc) String() string {
	return prettyPrint(boc, 0)
}

func (boc *Boc) value() string {
	return "{}"
}

func (boc *Boc) dataType() Type {
	return newBocType()
}

func (bl *BasicLit) String() string {
	return prettyPrint(bl, 0)
}

func (bl *BasicLit) value() string {
	return bl.val
}

func (bl *BasicLit) dataType() Type {
	return bl.basicType
}

func (al *ArrayLit) String() string {
	return prettyPrint(al, 0)
}

func (al *ArrayLit) value() string {
	return al.String()
}

func (al *ArrayLit) dataType() Type {
	return al.arrayType
}

func (d *DictLit) String() string {
	return prettyPrint(d, 0)
}

func (d *DictLit) value() string {
	return d.String()
}

func (d *DictLit) dataType() Type {
	return d.dictType
}

func (sd *ShortDeclaration) String() string {
	return prettyPrint(sd, 0)
}
func (sd *ShortDeclaration) value() string {
	return fmt.Sprintf("%s : %s", sd.variable, sd.val.value())
}

func (sd *ShortDeclaration) dataType() Type {
	return sd.val.dataType()
}
