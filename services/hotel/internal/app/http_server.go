package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	addr int
}

func NewHttpServer(addr int) *HttpServer {
	return &HttpServer{addr: addr}
}

func (h *HttpServer) Run() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}
