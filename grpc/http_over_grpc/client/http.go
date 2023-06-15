package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartLocalHTTPServer(addr string) {
	r := gin.Default()

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": 314,
		})
	})

	r.Run(addr)
}
