package main

import (
	"context"
	"fmt"
	"net"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/application/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/application/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/application/room-type"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/hotel"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/room"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/room-type"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/interfaces/hotel"
	room_handler "github.com/098765432m/grpc-kafka/hotel/internal/interfaces/room"
	room_type_handler "github.com/098765432m/grpc-kafka/hotel/internal/interfaces/room-type"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// 1. Config
	utils.Init()

	dsn := viper.GetString("DB_URL")

	// 2. Infras
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		zap.S().Fatalln("Failed to connect to database: ", err)
	}
	defer db.Close()

	// dial grpc
	imageConn := utils.NewGrpcClient(fmt.Sprintf(":%d", consts.IMAGE_GRPC_PORT))
	imageClient := image_pb.NewImageServiceClient(imageConn)

	hotelRepo := hotel_repo.New(db)
	roomTypeRepo := room_type_repo.New(db)
	roomRepo := room_repo.New(db)

	// 3. Application
	hotelService := hotel_service.NewHotelService(hotelRepo, imageClient)
	roomTypeService := room_type_service.NewRoomTypeService(roomTypeRepo)
	roomService := room_service.NewRoomService(roomRepo)

	// 4. Server
	hotelHandler := hotel_handler.NewHotelGrpcHandler(hotelService, roomTypeService, roomService)
	roomTypeHandler := room_type_handler.NewRoomTypeGrpcHandler(roomTypeService)
	roomHandler := room_handler.NewRoomGrpcHandler(roomService)

	// 5. Start Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", consts.HOTEL_GRPC_PORT))
	if err != nil {
		zap.S().Fatalln("Failed to Start server on port: ", consts.HOTEL_GRPC_PORT, err)

	}

	grpc := grpc.NewServer()

	hotel_pb.RegisterHotelServiceServer(grpc, hotelHandler)
	room_type_pb.RegisterRoomTypeServiceServer(grpc, roomTypeHandler)
	room_pb.RegisterRoomServiceServer(grpc, roomHandler)

	zap.S().Infof("Running grpc server on port %d, %v", consts.HOTEL_GRPC_PORT, err)
	if err := grpc.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
