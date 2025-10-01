package main

import (
	"context"
	"fmt"
	"net"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	image_service "github.com/098765432m/grpc-kafka/image/internal/application"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/infrastructure/repository/sqlc/image"
	image_handler "github.com/098765432m/grpc-kafka/image/internal/interfaces"
	"github.com/cloudinary/cloudinary-go/v2"
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
	cloudinaryUrl := viper.GetString("CLOUDINARY_URL")

	// 2. Infra
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := image_repo.New(db)
	cloudinaryClient, _ := cloudinary.NewFromURL(cloudinaryUrl)

	// 3. Application
	service := image_service.NewImageService(repo, cloudinaryClient)

	// 4. Server
	handler := image_handler.NewImageGrpcHandler(service)

	// 5. Start Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", consts.IMAGE_GRPC_PORT))
	if err != nil {
		panic(err)
	}

	grpc := grpc.NewServer()

	image_pb.RegisterImageServiceServer(grpc, handler)

	zap.S().Infof("Running grpc server on port %d, %v", consts.IMAGE_GRPC_PORT, err)
	if err := grpc.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
