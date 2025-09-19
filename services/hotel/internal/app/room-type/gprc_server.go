package room_type_app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	room_type_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/room-type"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr int
	conn *pgxpool.Pool
}

func NewGrpcServer(addr int, conn *pgxpool.Pool) *GrpcServer {
	return &GrpcServer{
		addr: addr,
		conn: conn,
	}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// repo
	roomTypeRepo := room_type_repo.New(g.conn)

	// service
	roomTypeService := room_type_service.NewRoomTypeService(roomTypeRepo)

	// Register grpc services
	room_type_pb.RegisterRoomTypeServiceServer(grpcServer, room_type_handler.NewRoomTypeGrpcHandler(roomTypeService))

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
