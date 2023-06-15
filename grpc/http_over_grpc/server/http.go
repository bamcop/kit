package main

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/bamcop/kit/grpc/http_over_grpc"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func StartHTTPServer() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		conn, err := LocalServer.conn(DefaultClientIdentifier)
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
			return
		}

		req := &pb.HTTPRequest{
			Method:  c.Request.Method,
			Url:     c.Request.URL.String(),
			Headers: nil,
			Body:    nil,
		}

		resp, err := conn.Handle(context.Background(), req)
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		slog.Info("resp", slog.Any("code", resp.Code), slog.Any("body", string(resp.Body)))

		for _, header := range resp.Headers {
			header := header
			c.Header(header.Key, header.Values[0])
		}
		c.Data(int(resp.Code), c.Writer.Header().Get("Content-Type"), resp.Body)
	})

	r.Run(fmt.Sprintf(":%d", HTTPPort))
}
