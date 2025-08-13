package app

import (
	"fmt"

	common_middleware "github.com/098765432m/grpc-kafka/common/middleware"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/hotel"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
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

	router.Use(common_middleware.CorsMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	hotelRepo := hotel_repo.New(h.conn)
	hotelService := hotel_service.NewHotelService(hotelRepo)

	roomTypeRepo := room_type_repo.New(h.conn)
	roomTypeService := room_type_service.NewRoomTypeService(roomTypeRepo)

	roomRepo := room_repo.New(h.conn)
	roomService := room_service.NewRoomService(roomRepo)

	hotelHandler := hotel_handler.NewHotelHttpHandler(hotelService, roomTypeService, roomService)

	api := router.Group("/api")
	hotelHandler.RegisterRoutes(api)

	fmt.Printf("Running HTTP server on port %d\n", h.addr)
	if err := router.Run(fmt.Sprintf(":%d", h.addr)); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
		return nil, err
	}

	return router, nil
}
