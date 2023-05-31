package helper

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

func AstString(node interface{}) string {
	var b bytes.Buffer
	if err := printer.Fprint(&b, token.NewFileSet(), node); err != nil {
		panic(err)
	}
	return b.String()
}

func PrintGenDecl(name string, nodes []ast.GenDecl, imports []*ast.ImportSpec) {
	var decls []ast.Decl
	for _, node := range nodes {
		node := node
		decls = append(decls, &node)
	}

	f := &ast.File{
		Name:    ast.NewIdent(name),
		Imports: imports,
		Decls:   decls,
	}

	fset := token.NewFileSet()
	err := printer.Fprint(os.Stdout, fset, f)
	if err != nil {
		panic(err)
	}
}

func MakeASTFile(pkgName string, nodes []ast.GenDecl, imports []*ast.ImportSpec, paths []string) string {
	var decls []ast.Decl
	if len(imports) > 0 {
		ImportDecl := MakeASTImportDecl(imports, paths)
		decls = append(decls, &ImportDecl)
	}

	for _, node := range nodes {
		node := node
		decls = append(decls, &node)
	}

	f := &ast.File{
		Name:    ast.NewIdent(pkgName),
		Imports: imports,
		Decls:   decls,
	}

	var b bytes.Buffer
	err := printer.Fprint(&b, token.NewFileSet(), f)
	if err != nil {
		panic(err)
	}

	str, err := GoFmtFile(b.String())
	if err != nil {
		panic(err)
	}

	return str
}

func GoFmtFile(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	ast.SortImports(fset, file)

	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := printer.Fprint(buffer, fset, file); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func MakeASTGenDecl(name string, Type *ast.StructType, pkgPath string) ast.GenDecl {
	for i, _ := range Type.Fields.List {
		field := Type.Fields.List[i]
		if v, ok := field.Type.(*ast.Ident); ok {
			if v.Obj != nil {
				field.Type = &ast.SelectorExpr{
					X:   ast.NewIdent(filepath.Base(pkgPath)),
					Sel: ast.NewIdent(v.Name),
				}
			}
		}
	}

	node := ast.GenDecl{
		Doc:    nil,
		TokPos: 0,
		Tok:    token.TYPE,
		Lparen: 0,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Doc:     nil,
				Name:    ast.NewIdent(name),
				Assign:  0,
				Type:    Type,
				Comment: nil,
			},
		},
		Rparen: 0,
	}

	return node
}

func MakeASTImportDecl(items []*ast.ImportSpec, paths []string) ast.GenDecl {
	var specs []ast.Spec

	for _, item := range paths {
		specs = append(specs, &ast.ImportSpec{
			Doc:  nil,
			Name: nil,
			Path: &ast.BasicLit{
				ValuePos: 0,
				Kind:     token.STRING,
				Value:    strconv.Quote(item),
			},
			Comment: nil,
			EndPos:  0,
		})
	}

	for _, item := range items {
		item := item

		var name *ast.Ident
		if item.Name != nil {
			name = ast.NewIdent(item.Name.Name)
		}
		specs = append(specs, &ast.ImportSpec{
			Doc:  nil,
			Name: name,
			Path: &ast.BasicLit{
				ValuePos: 0,
				Kind:     item.Path.Kind,
				Value:    item.Path.Value,
			},
			Comment: nil,
			EndPos:  0,
		})
	}

	node := ast.GenDecl{
		Doc:    nil,
		TokPos: 0,
		Tok:    token.IMPORT,
		Lparen: 0,
		Specs:  specs,
		Rparen: 0,
	}

	return node
}
