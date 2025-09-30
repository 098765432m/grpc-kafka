package rating_app

import (
	"fmt"
	"net"
	"time"

	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	rating_handler "github.com/098765432m/grpc-kafka/rating/internal/handler"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_redis "github.com/098765432m/grpc-kafka/rating/internal/repository/redis"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr      int
	redisPort int
	conn      *pgxpool.Pool
}

func NewGrpcServer(addr int, redisPort int, conn *pgxpool.Pool) *GrpcServer {
	return &GrpcServer{addr: addr, redisPort: redisPort, conn: conn}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		zap.S().Fatalf("Failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("redis://localhost:%d", g.redisPort),
	})

	//repo
	repo := rating_repo.New(g.conn)
	ratingRedis := rating_redis.NewRedisRatingCache(redisClient, 30*time.Second)

	//service
	ratingService := rating_service.NewRatingService(repo, ratingRedis)

	// Register grpc services
	rating_pb.RegisterRatingServiceServer(grpcServer, rating_handler.NewRatingGrpcHandler(ratingService))

	zap.S().Infof("Running grpc server on port %d, %v", g.addr, err)
	if err := grpcServer.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
