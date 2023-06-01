package main

import (
	"github.com/bamcop/kit/preset"
	"github.com/bamcop/kit/tool/cmd/gin_genc"
)

func main() {
	preset.SetDefaultSlog("")

	gin_genc.NewApp(
		"/Users/bamcop/StartKit/Golang/kit/tool/cmd/gin_genc/example/srv",
		"/Users/bamcop/StartKit/Golang/kit/tool/cmd/gin_genc/example/web",
		"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv",
		gin_genc.WithCtxProvider("srv.NewAppContext"),
		gin_genc.WithExtraImport([]string{"github.com/bamcop/kit/tool/cmd/gin_genc/example/srv"}),
	).Start()
}
