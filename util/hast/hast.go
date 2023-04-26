package hast

import (
	"go/ast"
)

var (
	Field = field{}
)

type (
	field struct{}
)

func (f field) IsAny(obj *ast.Field) bool {
	panic("not implemented")
}

func (f field) IsAnonymousStruct(obj *ast.Field) bool {
	panic("not implemented")
}
