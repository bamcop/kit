package skia_gen

import (
	"path/filepath"

	"github.com/bamcop/kit/http/ginx"
	"github.com/gin-gonic/gin"

	"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv"
	"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv/model"
)

func AddRouteGroupDefault(engine *gin.Engine, prefix string, handlers ...gin.HandlerFunc) {
	r := engine.Group(filepath.Join(prefix, "default"), handlers...)

	r.POST("bar", ginx.WrapX(model.Global{}.Bar, srv.NewAppContext))
	r.POST("hello", ginx.Wrap(model.Hello))
}

func AddRouteGroupG1(engine *gin.Engine, prefix string, handlers ...gin.HandlerFunc) {
	r := engine.Group(filepath.Join(prefix, "g_1"), handlers...)

	r.POST("foo", ginx.WrapX(model.Foo, srv.NewAppContext))
}
