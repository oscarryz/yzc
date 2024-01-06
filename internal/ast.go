package internal

type (
	program struct {
		blockBody *blockBody
	}
	expression interface {
		value() string
	}
	statement interface {
		value() string
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
