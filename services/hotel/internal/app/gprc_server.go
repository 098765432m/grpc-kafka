package app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/hotels"
	"github.com/098765432m/grpc-kafka/hotel/internal/handler"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr int
}

func NewGrpcServer(addr int) *GrpcServer {
	return &GrpcServer{addr: addr}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// Register our grpc services
	hotels.RegisterHotelServiceServer(grpcServer, handler.NewHotelGrpcHandler())

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
