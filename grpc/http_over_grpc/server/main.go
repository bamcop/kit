package main

import (
	"fmt"

	"github.com/bamcop/kit/preset"
)

const (
	HTTPPort = 8341
	GrpcPort = 9158
)

func main() {
	preset.SetDefaultSlog("./http_over_grpc.log")

	go StartGrpc(fmt.Sprintf(":%d", GrpcPort))

	StartHTTPServer()
}
