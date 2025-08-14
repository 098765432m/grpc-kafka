package app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/hotel"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr         int
	hotelService *hotel_service.HotelService
}

func NewGrpcServer(addr int,
	hotelService *hotel_service.HotelService) *GrpcServer {
	return &GrpcServer{addr: addr, hotelService: hotelService}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// Register grpc services
	hotel_pb.RegisterHotelServiceServer(grpcServer, hotel_handler.NewHotelGrpcHandler(g.hotelService))

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
