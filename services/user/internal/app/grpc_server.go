package app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	user_handler "github.com/098765432m/grpc-kafka/user/internal/handler"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/repository/user"
	user_service "github.com/098765432m/grpc-kafka/user/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr int
	conn *pgxpool.Pool
}

func NewGrpcServer(addr int, conn *pgxpool.Pool) *GrpcServer {
	return &GrpcServer{addr: addr, conn: conn}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// repo
	repo := user_repo.New(g.conn)

	// service
	service := user_service.NewUserService(repo)

	// handler
	handler := user_handler.NewUserGrpcHandler(service)

	user_pb.RegisterUserServiceServer(grpcServer, handler)

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
