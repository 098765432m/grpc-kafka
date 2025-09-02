package booking_app

import (
	"fmt"
	"log"
	"net"

	booking_handler "github.com/098765432m/grpc-kafka/booking/internal/handler"
	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/repository/booking"
	booking_service "github.com/098765432m/grpc-kafka/booking/internal/service"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
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

	// repo
	repo := booking_repo.New(bg.conn)

	// service
	service := booking_service.NewBookingService(bg.conn, repo)

	// handler
	handler := booking_handler.NewBookingGrpcHandler(service)

	// Register grpc server
	booking_pb.RegisterBookingServiceServer(grpcServer, handler)

	log.Printf("Running grpc server on port: %d\n", bg.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
