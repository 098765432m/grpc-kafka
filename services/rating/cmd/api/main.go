package main

import (
	"context"
	"fmt"
	"net"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/application"
	rating_redis "github.com/098765432m/grpc-kafka/rating/internal/infrastructure/redis"
	rating_handler "github.com/098765432m/grpc-kafka/rating/internal/interfaces/grpc"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// 1. Config
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	err := gotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	viper.AutomaticEnv()

	dsn := viper.GetString("DB_URL")
	redisUrl := viper.GetString("REDIS_URL")

	// 2. Infra
	// Postgres
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := rating_repo.New(db)

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	ratingRedis := rating_redis.NewRedisRatingCache(rdb)

	// 3. Application
	ratingService := rating_service.NewRatingService(repo, ratingRedis)

	// 4. Server
	ratingHandler := rating_handler.NewRatingGrpcHandler(ratingService)

	// 5. Start Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", consts.RATING_GRPC_PORT))
	if err != nil {
		panic(err)
	}

	grpc := grpc.NewServer()

	rating_pb.RegisterRatingServiceServer(grpc, ratingHandler)

	zap.S().Infof("Running grpc server on port %d, %v", consts.RATING_GRPC_PORT, err)
	if err := grpc.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
