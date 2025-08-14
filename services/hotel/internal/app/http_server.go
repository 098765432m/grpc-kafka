package app

import (
	"fmt"

	common_middleware "github.com/098765432m/grpc-kafka/common/middleware"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/hotel"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	addr int
	// conn *pgx.Conn
	hotelHandler *hotel_handler.HotelHttpHandler
}

func NewHttpServer(addr int, hotelHandler *hotel_handler.HotelHttpHandler) *HttpServer {
	return &HttpServer{addr: addr, hotelHandler: hotelHandler}
}

func (h *HttpServer) Run() (*gin.Engine, error) {
	router := gin.Default()

	router.Use(common_middleware.CorsMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")
	h.hotelHandler.RegisterRoutes(api)

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
