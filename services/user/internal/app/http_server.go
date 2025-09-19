package app

import (
	"fmt"

	common_middleware "github.com/098765432m/grpc-kafka/common/middleware"
	user_handler "github.com/098765432m/grpc-kafka/user/internal/handler"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/repository/user"
	user_service "github.com/098765432m/grpc-kafka/user/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HttpServer struct {
	addr int
	conn *pgxpool.Pool
}

func NewHttpServer(addr int, conn *pgxpool.Pool) *HttpServer {
	return &HttpServer{addr: addr, conn: conn}
}

func (h *HttpServer) Run() {
	router := gin.Default()

	router.Use(common_middleware.CorsMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := router.Group("/api")

	// repo
	repo := user_repo.New(h.conn)

	// service
	service := user_service.NewUserService(repo)

	// handler
	handler := user_handler.NewUserHttpHandler(service)

	handler.RegisterRoutes(api)

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}
