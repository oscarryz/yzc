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
		pos position
		tt  tokenType
		val string
	}
	ArrayLit struct {
		pos       position
		exps      []expression
		arrayType Type
	}
	// DictLit represents a dictionary literal [k1:v1 k2:v2] or [String]Int for empty dictionary
	DictLit struct {
		pos     position
		val     string
		keyType string
		valType string
		keys    []expression
		values  []expression
	}

	ShortDeclaration struct {
		pos position
		key expression
		val expression
	}

	ParenthesisExp struct {
		lparen position
		exps   []expression
		rparen position
	}
)

func (boc *Boc) String() string {
	return prettyPrint(boc, 0)
}

func (boc *Boc) value() string {
	return "{}"
}

func (boc *Boc) dataType() Type {
	return new(BocType)
}

func (bl *BasicLit) String() string {
	return prettyPrint(bl, 0)
}

func (bl *BasicLit) value() string {
	return bl.val
}

func (bl *BasicLit) dataType() Type {
	switch bl.tt {
	case INTEGER:
		return new(IntType)
	case DECIMAL:
		return new(DecimalType)
	case STRING:
		return new(StringType)
	default:
		return nil
	}
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
	return d.val
}

func (d *DictLit) dataType() Type {
	return new(DictType)
}

func (sd *ShortDeclaration) String() string {
	return prettyPrint(sd, 0)
}
func (sd *ShortDeclaration) value() string {
	return fmt.Sprintf("%s : %s", sd.key.value(), sd.val.value())
}

func (sd *ShortDeclaration) dataType() Type {
	return sd.val.dataType()
}
