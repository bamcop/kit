package dstfmt

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func H(filename string, hAppend func(string, *bytes.Buffer), imports []string) {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	f, err := decorator.Parse(b)
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		decl := decl
		if fn, ok := decl.(*dst.FuncDecl); ok {
			fn.Decs.Start.Append("\n")

			if stmt, ok := fn.Body.List[0].(*dst.ReturnStmt); ok {
				if cl, ok := stmt.Results[0].(*dst.CompositeLit); ok {
					for i, elt := range cl.Elts {
						elt := elt
						if call, ok := elt.(*dst.CallExpr); ok {
							// ent.Field Builder 换行
							//if sel, ok := call.Fun.(*dst.SelectorExpr); ok {
							//	sel.Sel.Decs.Start.Append("\n")
							//}

							if i == 0 {
								call.Decs.Start.Append("\n")
							}
							call.Decs.End.Append("\n")
						}
					}
				}
			}
		}
	}

	for _, decl := range f.Decls {
		decl := decl

		// 添加注解
		if fn, ok := decl.(*dst.FuncDecl); ok && fn.Name.Name == "Annotations" && len(fn.Body.List) > 0 {
			if result, ok := fn.Body.List[0].(*dst.ReturnStmt); ok && len(result.Results) > 0 {
				if expr, ok := result.Results[0].(*dst.CompositeLit); ok {
					expr.Elts = append(expr.Elts, &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "entsql"},
							Sel: &dst.Ident{Name: "WithComments"},
						},
						Args: []dst.Expr{
							&dst.Ident{Name: "true"},
						},
					})
				}
			}
		}

		// 链式调用换行
		if fn, ok := decl.(*dst.FuncDecl); ok && fn.Name.Name == "Fields" {
			FmtFuncFields(fn)
		}
	}

	var buff bytes.Buffer
	if err := decorator.Fprint(&buff, f); err != nil {
		panic(err)
	}

	if hAppend != nil {
		hAppend(filename, &buff)
	}

	// 添加 import
	content := AddImports(buff.String(), imports)

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		panic(err)
	}
}

func FmtFuncFields(f *dst.FuncDecl) {
	if len(f.Body.List) == 0 {
		return
	}

	stmt, ok := f.Body.List[0].(*dst.ReturnStmt)
	if !ok {
		return
	}

	if len(stmt.Results) == 0 {
		return
	}

	comp, ok := stmt.Results[0].(*dst.CompositeLit)
	if !ok {
		return
	}

	for _, elt := range comp.Elts {
		elt := elt
		if expr, ok := elt.(*dst.CallExpr); ok {
			FmtCallExpr(expr)
		}
	}
}

func FmtCallExpr(expr *dst.CallExpr) {
	if fun, ok := expr.Fun.(*dst.SelectorExpr); ok {
		x, ok := fun.X.(*dst.CallExpr)
		if ok {
			fun.Sel.Decs.NodeDecs.Before = dst.NewLine
			FmtCallExpr(x)
		}
	}
}

func AddImports(content string, imports []string) []byte {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "nil", content, 0)
	if err != nil {
		panic(err)
	}

	importDecl := f.Decls[0].(*ast.GenDecl)
	var header bytes.Buffer
	if err := printer.Fprint(&header, fset, importDecl); err != nil {
		panic(err)
	}

	for _, s := range imports {
		if strings.Contains(header.String(), s) {
			continue
		}

		newImportSpec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s\"", s),
			},
		}
		importDecl.Specs = append(importDecl.Specs, newImportSpec)
	}

	// Sort the imports.
	ast.SortImports(fset, f)

	var buff bytes.Buffer
	if err := printer.Fprint(&buff, fset, f); err != nil {
		panic(err)
	}

	return buff.Bytes()
}

func FmtDir(dirname string, hAppend func(string, *bytes.Buffer), addImports []string) {
	entities, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, entity := range entities {
		entity := entity

		if !entity.IsDir() && strings.HasSuffix(entity.Name(), ".go") {
			H(filepath.Join(dirname, entity.Name()), hAppend, addImports)
		}
	}
}
