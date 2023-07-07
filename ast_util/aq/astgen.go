package aq

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

var (
	Gen gen
)

type gen struct{}

func (g gen) MustNewExpr(src string) dst.Node {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		panic(err)
	}

	node, err := decorator.Decorate(token.NewFileSet(), expr)
	if err != nil {
		panic(err)
	}

	return node
}

func (g gen) MustNewFunc(src string) dst.Node {
	src = fmt.Sprintf("package main\n\n%s\n", src)

	root, err := decorator.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	if len(root.Decls) == 0 {
		panic("len(root.Decls) == 0")
	}

	v, ok := root.Decls[0].(*dst.FuncDecl)
	if !ok {
		panic("root.Decls[0] not ast.FuncDecl")
	}

	return v
}
