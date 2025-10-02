package main

import (
	"context"
	"fmt"
	"net"

	booking_service "github.com/098765432m/grpc-kafka/booking/internal/application"
	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/infrastructure/repository/sqlc/booking"
	booking_handler "github.com/098765432m/grpc-kafka/booking/internal/interfaces"
	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
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

	repo := booking_repo.New(db)

	// 3. Application
	service := booking_service.NewBookingService(db, repo)

	// 4. Server
	handler := booking_handler.NewBookingGrpcHandler(service)

	// 5. Start Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", consts.BOOKING_GRPC_PORT))
	if err != nil {
		zap.S().Fatalln("Failed to Start server on port: ", consts.BOOKING_GRPC_PORT, err)
	}

	grpc := grpc.NewServer()

	booking_pb.RegisterBookingServiceServer(grpc, handler)

	zap.S().Infof("Running grpc server on port %d, %v", consts.BOOKING_GRPC_PORT, err)
	if err := grpc.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
