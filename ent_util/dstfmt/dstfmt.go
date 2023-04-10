package dstfmt

import (
	"bytes"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"os"
	"path/filepath"
	"strings"
)

func H(filename string) {
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

	var buff bytes.Buffer
	if err := decorator.Fprint(&buff, f); err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, buff.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func FmtDir(dirname string) {
	entities, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, entity := range entities {
		entity := entity

		if !entity.IsDir() && strings.HasSuffix(entity.Name(), ".go") {
			H(filepath.Join(dirname, entity.Name()))
		}
	}
}
