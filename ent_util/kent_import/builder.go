package kent_import

import (
	"reflect"

	"ariga.io/atlas/sql/schema"
	"entgo.io/ent/schema/field"
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

func (b *fieldBuilder) Default(expr schema.Expr) *fieldBuilder {
	if expr == nil {
		return b
	}

	in := []reflect.Value{
		reflect.ValueOf(expr),
	}

	reflect.ValueOf(b.obj).MethodByName("Default").Call(in)

	return b
}
