package main

import (
	"fmt"
	"strconv"

	api_handler "github.com/098765432m/grpc-kafka/api-gateway/internal/handler"
	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
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

	userConn := utils.NewGrpcClient(strconv.Itoa(consts.USER_GRPC_PORT))
	defer userConn.Close()

	imageConn := utils.NewGrpcClient(strconv.Itoa(consts.IMAGE_GRPC_PORT))
	defer imageConn.Close()

	hotelClient := hotel_pb.NewHotelServiceClient(hotelConn)
	userClient := user_pb.NewUserServiceClient(userConn)
	imageClient := image_pb.NewImageServiceClient(imageConn)

	userHandler := api_handler.NewUserHandler(userClient, imageClient)
	userHandler.RegisterRoutes(api)

	hotelHandler := api_handler.NewHotelHandler(hotelClient, imageClient)
	hotelHandler.RegisterRoutes(api)

	imageHandler := api_handler.NewImageHandler(imageClient)
	imageHandler.RegisterRoutes(api)

	zap.S().Infoln("Running api-gateway on port ", consts.API_GATEWAY_PORT)

	if err := router.Run(fmt.Sprintf(":%d", consts.API_GATEWAY_PORT)); err != nil {
		zap.S().Fatalf("Failed to start HTTP server: %v\n", err)

	}
}
