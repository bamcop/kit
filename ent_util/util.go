package ent_util

import (
	"fmt"
	"reflect"
)

// CreateBy 从结构体创建
// CreateBy(client.TOrder.Create(), ent.TOrder{ ID: 1, ResourceFile: "ABC" })
func CreateBy[T1 any, T2 any](ptr T1, obj T2) T1 {
	var (
		rv     = reflect.ValueOf(obj)
		rf     = reflect.TypeOf(ptr)
		fields = reflect.VisibleFields(reflect.TypeOf(obj))
	)

	for i, field := range fields {
		if field.IsExported() {
			if !rv.FieldByName(field.Name).IsZero() {
				fmt.Println(i, field.Name, field.Type)

				fName := "Set" + field.Name
				if _, ok := rf.MethodByName(fName); ok {
					reflect.ValueOf(ptr).MethodByName(fName).
						Call([]reflect.Value{rv.FieldByName(field.Name)})
				}
			}
		}
	}

	return ptr
}
