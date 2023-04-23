package kent_import

import (
	"reflect"

	"entgo.io/ent/schema/field"
)

type fieldBuilder[T any] struct {
	obj T
}

func NewFieldBuilder[T any](name string, f func(str string) T) *fieldBuilder[T] {
	return &fieldBuilder[T]{
		obj: f(name),
	}
}

func (b *fieldBuilder[T]) Descriptor() *field.Descriptor {
	v := reflect.ValueOf(b.obj).MethodByName("Descriptor").Call(nil)[0]

	return v.Interface().(*field.Descriptor)
}

func (b *fieldBuilder[T]) StructTag(s string) *fieldBuilder[T] {
	in := []reflect.Value{
		reflect.ValueOf(s),
	}

	reflect.ValueOf(b.obj).MethodByName("StructTag").Call(in)

	return b
}

func (b *fieldBuilder[T]) SchemaType(types map[string]string) *fieldBuilder[T] {
	in := []reflect.Value{
		reflect.ValueOf(types),
	}

	reflect.ValueOf(b.obj).MethodByName("SchemaType").Call(in)

	return b
}

func (b *fieldBuilder[T]) Default(v any) *fieldBuilder[T] {
	in := []reflect.Value{
		reflect.ValueOf(v),
	}

	reflect.ValueOf(b.obj).MethodByName("Default").Call(in)

	return b
}
