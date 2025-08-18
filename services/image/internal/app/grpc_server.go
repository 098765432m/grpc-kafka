package image_app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	image_handler "github.com/098765432m/grpc-kafka/image/internal/handler"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/repository/image"
	image_service "github.com/098765432m/grpc-kafka/image/internal/service"
	"github.com/jackc/pgx/v5"
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
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// Register our grpc services
	repo := image_repo.New(g.conn)
	service := image_service.NewImageService(repo)

	image_pb.RegisterImageServiceServer(grpcServer, image_handler.NewImageGrpcHandler(service))

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
