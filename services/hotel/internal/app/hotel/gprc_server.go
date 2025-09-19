package hotel_app

import (
	"fmt"
	"log"
	"net"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/hotel"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	addr int
	conn *pgxpool.Pool
}

func NewGrpcServer(addr int,
	conn *pgxpool.Pool) *GrpcServer {
	return &GrpcServer{addr: addr, conn: conn}
}

func (g *GrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.addr))
	if err != nil {
		log.Fatalf("failed to run grpc server on port %d: %v", g.addr, err)
	}

	grpcServer := grpc.NewServer()

	// repo
	hotelRepo := hotel_repo.New(g.conn)
	roomTypeRepo := room_type_repo.New(g.conn)
	roomRepo := room_repo.New(g.conn)

	// dial grpc
	imageConn := utils.NewGrpcClient(fmt.Sprintf(":%d", consts.IMAGE_GRPC_PORT))
	imageClient := image_pb.NewImageServiceClient(imageConn)

	// service
	hotelService := hotel_service.NewHotelService(hotelRepo, imageClient)
	roomTypeService := room_type_service.NewRoomTypeService(roomTypeRepo)
	roomService := room_service.NewRoomService(roomRepo)

	// Register grpc services
	hotel_pb.RegisterHotelServiceServer(grpcServer, hotel_handler.NewHotelGrpcHandler(hotelService, roomTypeService, roomService))

	log.Printf("Running grpc server on port: %d\n", g.addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to start grpc server: ", err)
	}
}
