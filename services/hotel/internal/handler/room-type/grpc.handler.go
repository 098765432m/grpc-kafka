package room_type_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/jackc/pgx/v5/pgtype"
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

func (rtg *RoomTypeGrpcHandler) GetHotelById(ctx context.Context, req *room_type_pb.GetRoomTypeByIdRequest) (*room_type_pb.GetRoomTypeByIdResponse, error) {

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
			Price:   roomType.Price,
			HotelId: roomType.HotelID.String(),
		},
	}, nil
}

func (rtg *RoomTypeGrpcHandler) GetHotelByHotelId(ctx context.Context, req *room_type_pb.GetRoomTypesByHotelIdRequest) (*room_type_pb.GetRoomTypesByHotelIdResponse, error) {

	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	roomTypes, err := rtg.service.GetRoomTypesByHotelId(ctx, hotelId)
	if err != nil {
		return nil, err
	}

	var grpcRoomTypes []*room_type_pb.RoomType
	for _, roomType := range roomTypes {
		grpcRoomType := &room_type_pb.RoomType{
			Id:      roomType.ID.String(),
			Name:    roomType.Name,
			Price:   roomType.Price,
			HotelId: roomType.HotelID.String(),
		}

		grpcRoomTypes = append(grpcRoomTypes, grpcRoomType)
	}

	return &room_type_pb.GetRoomTypesByHotelIdResponse{
		RoomTypes: grpcRoomTypes,
	}, nil
}
