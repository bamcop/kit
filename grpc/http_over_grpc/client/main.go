package main

import (
	"fmt"
)

// 默认服务器地址 47.97.37.242
const (
	HTTPLocalPort  = 50877
	RemoteGrpcAddr = "47.97.37.242:9158"
)

func main() {
	go StartLocalHTTPServer(fmt.Sprintf(":%d", HTTPLocalPort))

	StartLocalGrpc(RemoteGrpcAddr)
}
