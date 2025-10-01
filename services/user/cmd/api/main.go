package main

import (
	"context"
	"fmt"
	"net"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	user_service "github.com/098765432m/grpc-kafka/user/internal/application"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/infrastructure/repository/sqlc/user"
	user_handler "github.com/098765432m/grpc-kafka/user/internal/interfaces/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// 1.Config
	utils.Init()

	dsn := viper.GetString("DB_URL")

	// 2. Infras
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := user_repo.New(db)

	// 3. Application
	userService := user_service.NewUserService(repo)

	// 4. Server
	userHandler := user_handler.NewUserGrpcHandler(userService)

	// 5. Start Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", consts.USER_GRPC_PORT))
	if err != nil {
		panic(err)
	}

	grpc := grpc.NewServer()

	user_pb.RegisterUserServiceServer(grpc, userHandler)

	zap.S().Infof("Running grpc server on port %d, %v", consts.USER_GRPC_PORT, err)
	if err := grpc.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
