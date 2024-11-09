package internal

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
	ArrayLit struct {
		pos  position
		tt   tokenType
		val  string
		exps []expression
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
func (e *empty) String() string {
	return prettyPrint(e, 0)
}
