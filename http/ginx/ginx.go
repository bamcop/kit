package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"net/http"
	"reflect"
)

var (
	HandlerBadRequest = defaultHandlerBadRequest
	HandlerError      = defaultHandlerError
	HandlerSuccess    = defaultHandlerSuccess
)

func defaultHandlerBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func defaultHandlerError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"error": err.Error(),
	})
}

func defaultHandlerSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  data,
	})
}

type Validator interface {
	Validate() error
}

var (
	validatorAttrType = reflect.TypeOf((*Validator)(nil)).Elem()
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func Wrap[T any](h func(*gin.Context, T) (any, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ptr := new(T)
		if err := ctx.ShouldBind(ptr); err != nil {
			HandlerBadRequest(ctx, err)
			return
		}

		// 使用 go-playground/validator 验证结构体
		if err := validate.Struct(ptr); err != nil {
			HandlerBadRequest(ctx, err)
			return
		}

		// 使用自定义验证
		v := reflect.ValueOf(*ptr)
		if v.Type().Implements(validatorAttrType) {
			if err := v.Interface().(Validator).Validate(); err != nil {
				HandlerBadRequest(ctx, err)
				return
			}
		}

		resp, err := h(ctx, *ptr)
		if err != nil {
			HandlerError(ctx, err)
			return
		}

		HandlerSuccess(ctx, resp)
	}
}

type IGinContext interface {
	GinContext() *gin.Context
}

func WrapX[T1 IGinContext, T2 any](h func(T1, T2) (any, error), f func(ctx *gin.Context) T1) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ptr := new(T2)
		if err := ctx.ShouldBind(ptr); err != nil {
			HandlerBadRequest(ctx, err)
			return
		}

		// 使用 go-playground/validator 验证结构体
		if err := validate.Struct(ptr); err != nil {
			HandlerBadRequest(ctx, err)
			return
		}

		// 使用自定义验证
		v := reflect.ValueOf(*ptr)
		if v.Type().Implements(validatorAttrType) {
			if err := v.Interface().(Validator).Validate(); err != nil {
				HandlerBadRequest(ctx, err)
				return
			}
		}

		resp, err := h(f(ctx), *ptr)
		if err != nil {
			HandlerError(ctx, err)
			return
		}

		HandlerSuccess(ctx, resp)
	}
}
