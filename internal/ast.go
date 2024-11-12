package internal

import "fmt"

type (
	boc struct {
		Name      string
		bocType   *blockType
		blockBody *blockBody
	}

	expression interface {
		value() string
		String() string
	}
	statement interface {
		value() string
		String() string
	}
	blockType struct {
		pos position
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
	ArrayLit struct {
		pos  position
		val  string
		exps []expression
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

func (boc *boc) value() string {
	return boc.Name
}

func (b *boc) String() string {
	return prettyPrint(b, 0)
}

func (bt *blockType) String() string {
	return prettyPrint(bt, 0)
}

func (bb *blockBody) String() string {
	return prettyPrint(bb, 0)
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

func (bl BasicLit) value() string {
	return bl.val
}
func (al ArrayLit) value() string {
	return al.val
}
func (d DictLit) value() string {
	return d.val
}

func (sd ShortDeclaration) value() string {
	return fmt.Sprintf("%s : %s", sd.key.value(), sd.val.value())
}
func (e empty) value() string {
	return "<empty>"
}
