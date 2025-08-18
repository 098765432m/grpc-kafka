package booking_app

import (
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type BookingGrpcServer struct {
	addr int
	conn *pgx.Conn
}

func NewBookingGrpcServer(addr int, conn *pgx.Conn) *BookingGrpcServer {
	return &BookingGrpcServer{
		addr: addr,
		conn: conn,
	}
}

func (bg *BookingGrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", bg.addr))
	if err != nil {
		zap.S().Fatalf("Failed to run grpc server on port %d: %v", bg.addr, err)
	}

	grpcServer := grpc.NewServer()

	// Register grpc server
	// hotel_pb.RegisterHotelServiceServer(grpcServer, hotel_handler.NewHotelGrpcHandler(g.hotelService))

	log.Printf("Running grpc server on port: %d\n", bg.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
