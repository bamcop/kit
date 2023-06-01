package srv

import (
	"github.com/gin-gonic/gin"
)

type (
	Context struct {
		Context *gin.Context
	}
)

func (app *Context) GinContext() *gin.Context {
	return app.Context
}

func NewAppContext(ctx *gin.Context) *Context {
	return &Context{
		Context: ctx,
	}
}
