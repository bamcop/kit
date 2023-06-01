package skia_gen

import (
	"net/http"

	"github.com/bamcop/kit/http/ginx"
	"github.com/gin-gonic/gin"

	"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv"
	"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv/model"
)

func NewRouter(engine *gin.Engine, prefix string) {
	engine.Use(cors())
	r := engine.Group(prefix)

	r.POST("bar", ginx.WrapX(model.Global{}.Bar, srv.NewAppContext))
	r.POST("foo", ginx.WrapX(model.Foo, srv.NewAppContext))
	r.POST("hello", ginx.Wrap(model.Hello))

}

// cors 需要注意的是，如果要发送Cookie，Access-Control-Allow-Origin就不能设为星号，必须指定明确的、与请求网页一致的域名
// https://www.ruanyifeng.com/blog/2016/04/cors.html
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		Origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Header.Get("Access-Control-Request-Headers") != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", c.Request.Header.Get("Access-Control-Request-Headers"))
		}
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, DELETE, POST, PUT")
			c.Writer.Header().Set("Allow", "OPTIONS, GET, DELETE, POST, PUT")
			c.Writer.Header().Set("Cache-Control", "max-age=604800")
			c.Writer.Header().Set("Content-Length", "0")
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
