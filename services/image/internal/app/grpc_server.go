package app

import (
	"fmt"
	"log"
	"net"

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
	// hotel_pb.RegisterHotelServiceServer(grpcServer)

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
