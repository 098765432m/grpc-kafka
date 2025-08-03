package rating_app

import (
	"fmt"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	rating_handler "github.com/098765432m/grpc-kafka/rating/internal/handler"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr int
	conn *pgx.Conn
}

func NewGrpcServer(addr int, conn *pgx.Conn) *GrpcServer {
	return &GrpcServer{addr: addr, conn: conn}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		zap.S().Fatalf("Failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	//repo
	repo := rating_repo.New(g.conn)
	//service
	ratingService := rating_service.NewRatingService(repo)

	// Register grpc services
	rating_pb.RegisterRatingServiceServer(grpcServer, rating_handler.NewRatingGrpcHandler(ratingService))

	zap.S().Infof("Running grpc server on port %d, %v", g.addr, err)
	if err := grpcServer.Serve(lis); err != nil {
		zap.S().Fatal("Failed to start grpc server: ", err)
	}
}
