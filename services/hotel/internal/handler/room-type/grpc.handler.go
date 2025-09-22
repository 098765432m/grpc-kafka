package room_type_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomTypeGrpcHandler struct {
	room_type_pb.UnimplementedRoomTypeServiceServer
	service *room_type_service.RoomTypeService
}

func NewRoomTypeGrpcHandler(service *room_type_service.RoomTypeService) *RoomTypeGrpcHandler {
	return &RoomTypeGrpcHandler{
		service: service,
	}
}

func (rtg *RoomTypeGrpcHandler) GetRoomTypeById(ctx context.Context, req *room_type_pb.GetRoomTypeByIdRequest) (*room_type_pb.GetRoomTypeByIdResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		return nil, err
	}

	roomType, err := rtg.service.GetRoomTypeById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &room_type_pb.GetRoomTypeByIdResponse{
		RoomType: &room_type_pb.RoomType{
			Id:      roomType.ID.String(),
			Name:    roomType.Name,
			Price:   uint32(roomType.Price),
			HotelId: roomType.HotelID.String(),
		},
	}, nil
}

func (rtg *RoomTypeGrpcHandler) GetRoomTypesByHotelId(ctx context.Context, req *room_type_pb.GetRoomTypesByHotelIdRequest) (*room_type_pb.GetRoomTypesByHotelIdResponse, error) {

	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	roomTypes, err := rtg.service.GetRoomTypesByHotelId(ctx, hotelId)
	if err != nil {
		return nil, err
	}

	var grpcRoomTypes []*room_type_pb.GetRoomTypesByHotelIdRow
	for _, roomType := range roomTypes {
		grpcRoomType := &room_type_pb.GetRoomTypesByHotelIdRow{
			Id:            roomType.ID.String(),
			Name:          roomType.Name,
			Price:         uint32(roomType.Price),
			HotelId:       roomType.HotelID.String(),
			NumberOfRooms: uint32(roomType.NumberOfRooms),
		}

		grpcRoomTypes = append(grpcRoomTypes, grpcRoomType)
	}

	return &room_type_pb.GetRoomTypesByHotelIdResponse{
		RoomTypes: grpcRoomTypes,
	}, nil
}

func (rtg *RoomTypeGrpcHandler) CreateRoomType(ctx context.Context, req *room_type_pb.CreateRoomTypeRequest) (*room_type_pb.CreateRoomTypeResponse, error) {

	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		zap.S().Infoln("Invalid Hotel UUID: ", err)
		return nil, status.Error(codes.InvalidArgument, "Hotel UUID khong hop le")
	}

	err := rtg.service.CreateRoomType(ctx, &room_type_repo.CreateRoomTypeParams{
		Name:    req.Name,
		Price:   req.Price,
		HotelID: hotelId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Khong tao duoc loai phong")
	}

	return &room_type_pb.CreateRoomTypeResponse{}, nil
}

// Delete Room By Id
func (rtg *RoomTypeGrpcHandler) DeleteRoomTypeById(ctx context.Context, req *room_type_pb.DeleteRoomTypeByIdRequest) (*room_type_pb.DeleteRoomTypeByIdResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Info("Invalid Room Type UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi Room Type UUID")
	}

	err := rtg.service.DeleteRoomTypeById(ctx, id)
	if err != nil {
		zap.S().Info("Cannot delete Room by id: ", err)
		return nil, status.Error(codes.Internal, "Khong the xoa Room bang id")
	}

	return &room_type_pb.DeleteRoomTypeByIdResponse{}, nil
}
