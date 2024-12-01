package internal

type Kind int

const (
	INT Kind = iota
	DEC
	STR
	ARRAY
	DICT
	BOC
)

type (
	Type interface {
		Kind() Kind
		String() string
	}
	IntType struct {
		Type
	}
	DecimalType struct {
		Type
	}
	StringType struct {
		Type
	}
	ArrayType struct {
		elemType Type
		Type
	}
	DictType struct {
		keyType Type
		valType Type
		Type
	}
	BocType struct {
		variables []*Field
		Type
	}
	Field struct {
		name     string // nil means just type, single uppercase letter means generic
		dataType Type
	}
	TBD struct {
		Type
	}
)
