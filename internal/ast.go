package internal

import "fmt"

type (
	Boc struct {
		Name        string
		expressions []expression
		statements  []statement
	}

	expression interface {
		value() string
		String() string
	}
	statement interface {
		value() string
		String() string
	}

	BasicLit struct {
		pos position
		tt  tokenType
		val string
	}
	ArrayLit struct {
		pos       position
		arrayType expression
		exps      []expression
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

	empty struct{}
)

func (boc *Boc) value() string {
	return boc.Name
}

func (boc *Boc) String() string {
	return prettyPrint(boc, 0)
}

func (bl *BasicLit) String() string {
	return prettyPrint(bl, 0)
}

func (al *ArrayLit) String() string {
	return prettyPrint(al, 0)
}

func (d *DictLit) String() string {
	return prettyPrint(d, 0)
}

func (sd *ShortDeclaration) String() string {
	return prettyPrint(sd, 0)
}
func (e *empty) String() string {
	return prettyPrint(e, 0)
}

func (bl *BasicLit) value() string {
	return bl.val
}
func (al *ArrayLit) value() string {
	return al.arrayType.value()
}
func (d *DictLit) value() string {
	return d.val
}

func (sd *ShortDeclaration) value() string {
	return fmt.Sprintf("%s : %s", sd.key.value(), sd.val.value())
}
func (e *empty) value() string {
	return "<empty>"
}
