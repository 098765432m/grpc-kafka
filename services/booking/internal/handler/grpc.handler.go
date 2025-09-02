package booking_handler

import (
	"context"

	booking_service "github.com/098765432m/grpc-kafka/booking/internal/service"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingGrpcHandler struct {
	booking_pb.UnimplementedBookingServiceServer
	service *booking_service.BookingService
}

func NewBookingGrpcHandler(service *booking_service.BookingService) *BookingGrpcHandler {
	return &BookingGrpcHandler{
		service: service,
	}
}

// Return {roomTypeId, Number of occupied rooms in that room type}
func (bg *BookingGrpcHandler) GetNumberOfOccupiedRooms(ctx context.Context, req *booking_pb.GetNumberOfOccupiedRoomsRequest) (*booking_pb.GetNumberOfOccupiedRoomsResponse, error) {
	// Check Are room type Ids valid
	roomTypeIds := make([]pgtype.UUID, 0, len(req.GetRoomTypeIds()))
	for _, roomTypeIdReq := range req.GetRoomTypeIds() {
		var roomTypeId pgtype.UUID
		if err := roomTypeId.Scan(roomTypeIdReq); err != nil {
			zap.S().Info("Invalid Room Id: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Room Id")
		}

		roomTypeIds = append(roomTypeIds, roomTypeId)
	}

	//Check are dates valid
	var checkInDate pgtype.Date
	if err := checkInDate.Scan(req.CheckIn); err != nil {
		zap.S().Info("Invalid date format: ", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid date format")
	}

	if checkInDateValue, err := checkInDate.DateValue(); err == nil {
		zap.L().Info("CheckIn: ", zap.Any("Day", checkInDateValue.Time.Day()), zap.Any("Month", checkInDateValue.Time.Month()), zap.Any("Year", checkInDateValue.Time.Year()))
	}
	var checkOutDate pgtype.Date
	if err := checkOutDate.Scan(req.CheckOut); err != nil {
		zap.S().Info("Invalid date format: ", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid date format")
	}
	if checkOutDateValue, err := checkOutDate.DateValue(); err == nil {
		zap.L().Info("CheckIn: ", zap.Any("Day", checkOutDateValue.Time.Day()), zap.Any("Month", checkOutDateValue.Time.Month()), zap.Any("Year", checkOutDateValue.Time.Year()))
	}
	results, err := bg.service.GetNumberOfOccupiedRooms(ctx, roomTypeIds, checkInDate, checkOutDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	numberOccpiedRoomsResponse := make([]*booking_pb.ResultNumberOfOccupiedRooms, 0, len(results))
	for _, result := range results {
		tempResult := &booking_pb.ResultNumberOfOccupiedRooms{
			RoomTypeId:            result.RoomTypeID.String(),
			NumberOfOccupiedRooms: int32(result.NumberOfOccupiedRooms),
		}
		numberOccpiedRoomsResponse = append(numberOccpiedRoomsResponse, tempResult)
	}

	return &booking_pb.GetNumberOfOccupiedRoomsResponse{
		Results: numberOccpiedRoomsResponse,
	}, nil
}
