package booking_app

import (
	"fmt"

	booking_handler "github.com/098765432m/grpc-kafka/booking/internal/handler"
	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/repository/booking"
	booking_service "github.com/098765432m/grpc-kafka/booking/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type HttpServer struct {
	addr int
	conn *pgx.Conn
}

func NewHttpServer(addr int, conn *pgx.Conn) *HttpServer {
	return &HttpServer{
		addr: addr,
		conn: conn,
	}
}

func (h *HttpServer) Run() (*gin.Engine, error) {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	// repo
	repo := booking_repo.New(h.conn)

	// service
	service := booking_service.NewBookingService(repo)

	// handler
	handler := booking_handler.NewBookingHttpHandler(service)

	api := router.Group("/api")
	handler.RegisterRoutes(api)

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
