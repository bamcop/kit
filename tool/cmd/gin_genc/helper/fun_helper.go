package helper

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
)

func TargetShape(f *ast.FuncDecl) bool {
	if f.Type == nil || f.Type.Params == nil || f.Type.Results == nil {
		return false
	}
	if !(len(f.Type.Params.List) == 2 && len(f.Type.Results.List) == 2) {
		return false
	}

	// 检查第一个返回值是 `interface{}` 或者 `any`
	r1 := f.Type.Results.List[0]
	if _, ok := r1.Type.(*ast.InterfaceType); !ok {
		v, ok := r1.Type.(*ast.Ident)
		if !ok {
			return false
		}
		if v.Name != "any" {
			return false
		}
	}

	// 检查第二个返回值是 `error`
	r2 := f.Type.Results.List[1]
	if v, ok := r2.Type.(*ast.Ident); !ok {
		return false
	} else {
		if v.Name != "error" {
			return false
		}
	}

	// 检查第一个参数的类型是某种 `context`
	p1 := f.Type.Params.List[0]
	if !strings.HasSuffix(AstString(p1.Type), "Context") {
		return false
	}

	// 检查第二个参数的类型是匿名的 `struct`
	p2 := f.Type.Params.List[1]
	if _, ok := p2.Type.(*ast.StructType); ok {
		return true
	}

	// 也支持命名类型的 `struct`
	ident, ok := p2.Type.(*ast.Ident)
	if !ok {
		return false
	}
	spec, ok := ident.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return false
	}
	if _, ok := spec.Type.(*ast.StructType); !ok {
		return false
	}
	return true
}

func IsRawGinContext(f *ast.FuncDecl, checker *types.Checker) bool {
	tp := checker.TypeOf(f.Type.Params.List[0].Type)
	return tp.String() == "*github.com/gin-gonic/gin.Context"
}

func TargetFunc(f *ast.FuncDecl, checker *types.Checker) ([]*ast.StructType, error) {
	// 检查函数体
	// 此处基于假设: 最后一个 `stmt` 是 `ast.ReturnStmt`, 而且不是死代码, 而且是函数正确调用的返回值(错误已经提前返回了)
	stmt := f.Body.List[len(f.Body.List)-1]
	if _, ok := stmt.(*ast.ReturnStmt); !ok {
		return nil, errors.New("not ast.ReturnStmt")
	}

	results := stmt.(*ast.ReturnStmt).Results
	// 第一个返回值应该是结构体字面量
	value, ok := results[0].(*ast.CompositeLit)
	if !ok {
		return nil, errors.New("not ast.CompositeLit")
	}

	expr, ok := value.Type.(*ast.SelectorExpr)
	if !ok {
		ident, ok := value.Type.(*ast.Ident)
		if !ok {
			return nil, errors.New("not ast.Ident")
		}
		if ident.Obj == nil {
			return nil, errors.New("ident.Obj nil")
		}
		spec, ok := ident.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			return nil, errors.New("not ast.TypeSpec")
		}
		if _, ok := spec.Type.(*ast.StructType); !ok {
			return nil, errors.New("not ast.StructType")
		}

		if _, ok := f.Type.Params.List[1].Type.(*ast.StructType); ok {
			return []*ast.StructType{
				f.Type.Params.List[1].Type.(*ast.StructType),
				spec.Type.(*ast.StructType),
			}, nil
		} else {
			return []*ast.StructType{
				f.Type.Params.List[1].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType),
				spec.Type.(*ast.StructType),
			}, nil
		}
	}

	ident, ok := expr.X.(*ast.Ident)
	if !ok {
		return nil, errors.New("SelectorExpr.X not ast.Ident")
	}
	if !(ident.Name == "gin" && expr.Sel.Name == "H") {
		return nil, errors.New("results[0] not gin.H")
	}

	resp := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

	for _, elt := range value.Elts {
		kv := elt.(*ast.KeyValueExpr)

		var (
			v    = checker.TypeOf(kv.Value)
			strs = strings.Split(v.String(), "/")
		)
		expr, err := parser.ParseExpr(strs[len(strs)-1])
		if err != nil {
			panic(err)
		}

		var (
			key  = strings.Trim(kv.Key.(*ast.BasicLit).Value, `"`)
			name = strcase.ToCamel(key)
			tag  = &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", key),
			}
		)

		resp.Fields.List = append(resp.Fields.List, &ast.Field{
			Doc: nil,
			Names: []*ast.Ident{
				{
					Name: name,
				},
			},
			Type:    expr,
			Tag:     tag,
			Comment: nil,
		})
	}

	if _, ok := f.Type.Params.List[1].Type.(*ast.StructType); ok {
		return []*ast.StructType{
			f.Type.Params.List[1].Type.(*ast.StructType),
			resp,
		}, nil
	} else {
		return []*ast.StructType{
			f.Type.Params.List[1].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType),
			resp,
		}, nil
	}
}
