package internal

type (
	boc struct {
		name 	string
		bocType 	*blockType
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

func (boc *boc) value() string {
	return boc.name
}