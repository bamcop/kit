package model

import (
	"time"

	"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv"
	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context, req struct {
	Id int64 `json:"id"`
}) (any, error) {
	return gin.H{
		"now": time.Now(),
	}, nil
}

func Foo(ctx *srv.Context, req struct {
	ID int64 `json:"id"`
}) (interface{}, error) {
	return gin.H{
		"path": "1",
	}, nil
}

type Global struct{}

func (Global) Bar(app *srv.Context, req struct {
	ID int64 `json:"id"`
}) (any, error) {
	return gin.H{
		"path": "1",
	}, nil
}
