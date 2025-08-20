package main

import (
	"fmt"
	"strconv"

	api_handler "github.com/098765432m/grpc-kafka/api-gateway/internal/handler"
	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	common_middleware "github.com/098765432m/grpc-kafka/common/middleware"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Init logger
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	router := gin.Default()

	// CORS config
	router.Use(common_middleware.CorsMiddleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	api := router.Group("/api")

	hotelConn := utils.NewGrpcClient(strconv.Itoa(consts.HOTEL_GRPC_PORT))
	defer hotelConn.Close()

	roomTypeConn := utils.NewGrpcClient(strconv.Itoa(consts.ROOM_TYPE_GRPC_PORT))
	defer hotelConn.Close()

	roomConn := utils.NewGrpcClient(strconv.Itoa(consts.ROOM_GRPC_PORT))
	defer hotelConn.Close()

	userConn := utils.NewGrpcClient(strconv.Itoa(consts.USER_GRPC_PORT))
	defer userConn.Close()

	imageConn := utils.NewGrpcClient(strconv.Itoa(consts.IMAGE_GRPC_PORT))
	defer imageConn.Close()

	ratingConn := utils.NewGrpcClient(strconv.Itoa(consts.RATING_GRPC_PORT))
	defer ratingConn.Close()

	hotelClient := hotel_pb.NewHotelServiceClient(hotelConn)
	roomTypeClient := room_type_pb.NewRoomTypeServiceClient(roomTypeConn)
	roomClient := room_pb.NewRoomServiceClient(roomConn)
	userClient := user_pb.NewUserServiceClient(userConn)
	imageClient := image_pb.NewImageServiceClient(imageConn)
	ratingClient := rating_pb.NewRatingServiceClient(ratingConn)

	userHandler := api_handler.NewUserHandler(userClient, imageClient)
	userHandler.RegisterRoutes(api)

	hotelHandler := api_handler.NewHotelHandler(hotelClient, roomTypeClient, roomClient, userClient, imageClient, ratingClient)
	hotelHandler.RegisterRoutes(api)

	roomTypeHandler := api_handler.NewRoomTypeGrpchandler(roomTypeClient, roomClient)
	roomTypeHandler.RegisterRoutes(api)

	roomHandler := api_handler.NewRoomGrpchandler(roomClient)
	roomHandler.RegisterRoutes(api)

	imageHandler := api_handler.NewImageHandler(imageClient)
	imageHandler.RegisterRoutes(api)

	zap.S().Infoln("Running api-gateway on port ", consts.API_GATEWAY_PORT)

	if err := router.Run(fmt.Sprintf(":%d", consts.API_GATEWAY_PORT)); err != nil {
		zap.S().Fatalf("Failed to start HTTP server: %v\n", err)

	}
}
