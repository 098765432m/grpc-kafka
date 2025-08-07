package app

import (
	"fmt"

	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
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

	hotelRepo := hotel_repo.New(h.conn)
	hotelService := hotel_service.NewHotelService(hotelRepo)
	hotelHandler := hotel_handler.NewHotelHttpHandler(hotelService)

	api := router.Group("/api")
	hotelHandler.RegisterRoutes(api)

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
