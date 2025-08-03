package rating_app

import (
	rating_handler "github.com/098765432m/grpc-kafka/rating/internal/handler"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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
			"message": "pong",
		})
	})

	ratingRepo := rating_repo.New(h.conn)
	ratingService := rating_service.NewRatingService(ratingRepo)
	ratingHandler := rating_handler.NewHotelHttpHandler(ratingService)

	api := router.Group("/api")
	ratingHandler.RegisterRoutes(api)

	zap.S().Infof("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(); err != nil {
		zap.S().Errorf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
