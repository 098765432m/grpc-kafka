package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type HttpServer struct {
	addr int
	conn *pgx.Conn
}

func NewHttpServer(addr int, conn *pgx.Conn) *HttpServer {
	return &HttpServer{addr: addr, conn: conn}
}

func (h *HttpServer) Run() (*gin.Engine, error) {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
