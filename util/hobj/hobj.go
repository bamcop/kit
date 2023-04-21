package hobj

import (
	"github.com/jinzhu/copier"
)

func Copy[T any](from any) (T, error) {
	var to = new(T)

	if err := copier.Copy(to, from); err != nil {
		return *to, err
	}
	return *to, nil
}
