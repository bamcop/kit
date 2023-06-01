package skia_gen

import (
	"time"
)

type BarRequest struct {
	ID int64 `json:"id"`
}
type BarResponse struct {
	Path string `json:"path"`
}
type FooRequest struct {
	ID int64 `json:"id"`
}
type FooResponse struct {
	Path string `json:"path"`
}
type HelloRequest struct {
	Id int64 `json:"id"`
}
type HelloResponse struct {
	Now time.Time `json:"now"`
}
