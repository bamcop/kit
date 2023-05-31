package kent_import

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/schema"
	"entgo.io/ent/schema/field"
	"github.com/samber/lo"
)

type fieldBuilder struct {
	obj any
}

func NewFieldBuilder(name string, f any) *fieldBuilder {
	in := []reflect.Value{
		reflect.ValueOf(name),
	}
	ou := reflect.ValueOf(f).Call(in)[0]

	return &fieldBuilder{
		obj: ou.Interface(),
	}
}

func (b *fieldBuilder) Descriptor() *field.Descriptor {
	v := reflect.ValueOf(b.obj).MethodByName("Descriptor").Call(nil)[0]

	return v.Interface().(*field.Descriptor)
}

func (b *fieldBuilder) StructTag(s string) *fieldBuilder {
	in := []reflect.Value{
		reflect.ValueOf(s),
	}

	reflect.ValueOf(b.obj).MethodByName("StructTag").Call(in)

	return b
}

func (b *fieldBuilder) SchemaType(types map[string]string) *fieldBuilder {
	in := []reflect.Value{
		reflect.ValueOf(types),
	}

	f := reflect.ValueOf(b.obj).MethodByName("SchemaType")
	if !f.IsValid() {
		return b
	}

	f.Call(in)

	return b
}

func (b *fieldBuilder) Optional() *fieldBuilder {
	f := reflect.ValueOf(b.obj).MethodByName("Optional")
	if !f.IsValid() {
		return b
	}

	f.Call(nil)

	return b
}

func (b *fieldBuilder) Comment(attrs []schema.Attr) *fieldBuilder {
	attr, ok := lo.Find(attrs, func(item schema.Attr) bool {
		_, ok := item.(*schema.Comment)
		if ok {
			return true
		}
		return false
	})
	if !ok {
		return b
	}

	in := []reflect.Value{
		reflect.ValueOf(attr.(*schema.Comment).Text),
	}

	f := reflect.ValueOf(b.obj).MethodByName("Comment")
	if !f.IsValid() {
		return b
	}

	f.Call(in)

	return b
}

func (b *fieldBuilder) Default(expr schema.Expr) *fieldBuilder {
	if expr == nil {
		return b
	}

	method, ok := reflect.TypeOf(b.obj).MethodByName("Default")
	if !ok {
		return b
	}

	var v any
	t1 := method.Type.In(1).Kind()
	switch t1 {
	case reflect.Int:
		p, _ := strconv.ParseInt(expr.(*schema.Literal).V, 10, 64)
		v = int(p)
	case reflect.Int64:
		p, _ := strconv.ParseInt(expr.(*schema.Literal).V, 10, 64)
		v = p
	case reflect.String:
		p := expr.(*schema.Literal).V
		if strings.HasPrefix(p, "'") && strings.HasSuffix(p, "'") {
			v = strings.Trim(p, "'")
		}
	case reflect.Float64:
		v, _ = strconv.ParseFloat(expr.(*schema.Literal).V, 10)
	case reflect.Bool:
		if expr.(*schema.Literal).V == "0" || expr.(*schema.Literal).V == "1" {
			v = expr.(*schema.Literal).V != "0"
		} else {
			fmt.Println("unexpect value", expr)
		}
	default:
		fmt.Println(1)
	}

	in := []reflect.Value{
		reflect.ValueOf(v),
	}

	reflect.ValueOf(b.obj).MethodByName("Default").Call(in)

	return b
}
