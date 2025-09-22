package room_handler

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomGrpcHandler struct {
	room_pb.UnimplementedRoomServiceServer
	service *room_service.RoomService
}

func NewRoomGrpcHandler(service *room_service.RoomService) *RoomGrpcHandler {
	return &RoomGrpcHandler{
		service: service,
	}
}

func (rg *RoomGrpcHandler) GetRoomById(ctx context.Context, req *room_pb.GetRoomByIdRequest) (*room_pb.GetRoomByIdResponse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Info("Invalid Room UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi Room UUID")
	}

	room, err := rg.service.GetRoomsById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			return nil, status.Error(codes.NotFound, "Phong khong ton tai")
		}

		return nil, status.Error(codes.Internal, "Khong the lay Room bang id")
	}

	return &room_pb.GetRoomByIdResponse{
		Room: &room_pb.Room{
			Id:      room.ID.String(),
			Name:    room.Name,
			Status:  string(room.Status.RoomStatus),
			HotelId: room.HotelID.String(),
		},
	}, nil
}

func (rg *RoomGrpcHandler) GetRoomsByRoomTypeId(ctx context.Context, req *room_pb.GetRoomsByRoomTypeIdRequest) (*room_pb.GetRoomsByRoomTypeIdResponse, error) {
	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.GetRoomTypeId()); err != nil {

		zap.S().Info("Invalid Room Type UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi Room type UUID")
	}

	rooms, err := rg.service.GetRoomsByRoomTypeId(ctx, roomTypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong lay Rooms bang RoomType id")
	}

	roomsGrpcResult := make([]*room_pb.Room, 0, len(rooms))
	for _, room := range rooms {
		roomGrpc := &room_pb.Room{
			Id:         room.ID.String(),
			Name:       room.Name,
			Status:     string(room.Status.RoomStatus),
			RoomTypeId: room.RoomTypeID.String(),
			HotelId:    room.HotelID.String(),
		}

		roomsGrpcResult = append(roomsGrpcResult, roomGrpc)
	}

	return &room_pb.GetRoomsByRoomTypeIdResponse{
		Rooms: roomsGrpcResult,
	}, nil

}

func (rg *RoomGrpcHandler) GetRoomsByHotelId(ctx context.Context, req *room_pb.GetRoomsByHotelIdRequest) (*room_pb.GetRoomsByHotelIdResponse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.GetHotelId()); err != nil {

		zap.S().Info("Invalid Hotel UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi Hotel UUID")
	}

	rooms, err := rg.service.GetRoomsByHotelId(ctx, hotelId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong lay Rooms bang Hotel id")
	}

	roomsGrpcResult := make([]*room_pb.Room, 0, len(rooms))
	for _, room := range rooms {
		roomGrpc := &room_pb.Room{
			Id:         room.ID.String(),
			Name:       room.Name,
			Status:     string(room.Status.RoomStatus),
			RoomTypeId: room.RoomTypeID.String(),
			HotelId:    room.HotelID.String(),
		}

		roomsGrpcResult = append(roomsGrpcResult, roomGrpc)
	}

	return &room_pb.GetRoomsByHotelIdResponse{
		Rooms: roomsGrpcResult,
	}, nil
}

// Get Number of rooms per RoomType by Hotel Ids
// -> return [room_type_id, number_of_rooms]
func (rg *RoomGrpcHandler) GetNumberOfRoomsPerRoomTypeByHotelIds(ctx context.Context, req *room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsRequest) (*room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsResponse, error) {

	// convert to UUIDs
	hotelIds, err := utils.ToPgUuidArray(req.GetHotelIds())
	if err != nil {
		zap.S().Infoln("Invalid Hotel UUID")
		return nil, status.Error(codes.InvalidArgument, "Hotel Id khong hop le")
	}

	rows, err := rg.service.GetNumberOfRoomsPerRoomTypeByHotelIds(ctx, hotelIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong lay duoc danh sach so luong phong")
	}

	results := make([]*room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsRow, 0, len(rows))
	for _, row := range rows {
		result := &room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsRow{
			RoomTypeId:    row.RoomTypeID.String(),
			NumberOfRooms: int32(row.TotalRooms),
		}

		results = append(results, result)
	}

	return &room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsResponse{
		Results: results,
	}, nil
}

func (rg *RoomGrpcHandler) GetListOfAvailableRoomsByRoomTypeId(ctx context.Context, req *room_pb.GetListOfAvailableRoomsByRoomTypeIdRequest) (*room_pb.GetListOfAvailableRoomsByRoomTypeIdResponse, error) {

	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.GetRoomTypeId()); err != nil {
		zap.S().Info("Invalid Room Type UUID")
		return nil, status.Error(codes.InvalidArgument, "Khong the lay danh sach phong")
	}

	roomIds, err := rg.service.GetListOfAvailableRoomsByRoomTypeId(ctx, roomTypeId, int(req.GetNumberOfRooms()))
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			return nil, status.Error(codes.NotFound, "Khong du phong trong")
		}

		return nil, status.Error(codes.Internal, "Khong the lay danh sach phong")
	}

	zap.S().Infoln("room ids: ", roomIds)

	roomIdsStr := make([]string, 0, len(roomIds))
	for _, roomId := range roomIds {
		roomIdsStr = append(roomIdsStr, roomId.String())
	}
	return &room_pb.GetListOfAvailableRoomsByRoomTypeIdResponse{
		RoomIds: roomIdsStr,
	}, nil
}

// Get List Of Remain Rooms (NOT_BOOKED ROOM)
func (rs *RoomGrpcHandler) GetListOfRemainRooms(ctx context.Context, req *room_pb.GetListOfRemainRoomsRequest) (*room_pb.GetListOfRemainRoomsResponse, error) {

	// Convert to pgtype.UUID
	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.GetRoomTypeId()); err != nil {
		zap.S().Info("Invalid UUID format: ", err)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	bookedRoomIds, err := utils.ToPgUuidArray(req.GetBookedRoomIds())
	if err != nil {
		zap.S().Infoln(err)
		return nil, status.Error(codes.InvalidArgument, "Loi UUID khong hop le")
	}

	roomIds, err := rs.service.GetListOfRemainRooms(ctx, roomTypeId, bookedRoomIds, int(req.NumberOfRooms))
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay danh sach phong trong")
	}

	// To UUIDs string
	roomIdsStr := utils.ToPgUuidString(roomIds)

	return &room_pb.GetListOfRemainRoomsResponse{
		RoomIds: roomIdsStr,
	}, nil
}

// Create Room
func (rg *RoomGrpcHandler) CreateRoom(ctx context.Context, req *room_pb.CreateRoomRequest) (*room_pb.CreateRoomResponse, error) {

	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.RoomTypeId); err != nil {
		zap.S().Info("Invalid UUID format: ", err)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		zap.S().Info("Invalid UUID format: ", err)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	var roomStatus room_repo.NullRoomStatus
	if err := roomStatus.Scan(req.GetStatus()); err != nil {
		zap.S().Info("Invalid Room Status: ", err)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	err := rg.service.CreateRoom(ctx, &room_repo.CreateRoomParams{
		Name:       req.Name,
		Status:     string(roomStatus.RoomStatus),
		RoomTypeID: roomTypeId,
		HotelID:    hotelId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong tao duoc phong")
	}

	return &room_pb.CreateRoomResponse{}, nil
}

// set status = MAINTAINED to rooms
func (rg *RoomGrpcHandler) SetMaintainedStatusToRooms(ctx context.Context, req *room_pb.SetOccupiedStatusToRoomsRequest) (*room_pb.SetOccupiedStatusToRoomsResponse, error) {

	// Convert []string to []pgtype.UUID
	roomIds, err := utils.ToPgUuidArray(req.GetRoomIds())
	if err != nil {
		zap.S().Infoln("Invalid Rooms UUID")
		return nil, status.Error(codes.InvalidArgument, "UUIDs khong hop le")
	}

	err = rg.service.ChangeStatusOfRooms(ctx, roomIds, room_repo.RoomStatusMAINTAINED)
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			return nil, status.Error(codes.NotFound, "Khong tim thay phong")
		}
		return nil, status.Error(codes.Internal, "Khong the doi trang thai phong")
	}

	return nil, nil
}
