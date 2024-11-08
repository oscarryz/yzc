package internal

type (
	boc struct {
		Name      string
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
