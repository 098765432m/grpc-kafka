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

func (bg *BookingGrpcHandler) CreateBookings(ctx context.Context, req *booking_pb.NewBookingParam) (*booking_pb.Empty, error) {

	zap.L().Info("Create Bookings", zap.Any("Req", req.NewBookings))

	newBookings := []booking_service.NewBooking{}
	for _, param := range req.NewBookings {
		zap.L().Info("New Booking: ", zap.Any("CheckIn", param.CheckIn), zap.Any("CheckOut", param.CheckOut), zap.Any("Total", param.Total), zap.Any("RoomTypeId", param.RoomTypeId), zap.Any("UserId", param.UserId), zap.Any("RoomId", param.RoomId))

		// Chech are dates valid
		var checkIn pgtype.Date
		if err := checkIn.Scan(param.CheckIn); err != nil {
			zap.S().Info("Invalid Check In date format on create bookings: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Check In date format")
		}

		var checkOut pgtype.Date
		if err := checkOut.Scan(param.CheckOut); err != nil {
			zap.S().Info("Invalid Check Out date format on create bookings: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Check Out date format")
		}

		// check In < Check Out
		if checkIn.Time.After(checkOut.Time) {
			zap.S().Info("Check In date must be before Check Out date")
			return nil, status.Error(codes.InvalidArgument, "Check In date must be before Check Out date")
		}

		// Check are UUIDs valid
		var roomTypeId pgtype.UUID
		if err := roomTypeId.Scan(param.RoomTypeId); err != nil {
			zap.S().Info("Invalid Room Type UUID on create bookings: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Room Type Id")
		}

		var userId pgtype.UUID
		if err := userId.Scan(param.UserId); err != nil {
			zap.S().Info("Invalid User UUID on create bookings: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid User UUID")
		}
		var roomId pgtype.UUID
		if err := roomId.Scan(param.RoomId); err != nil {
			zap.S().Info("Invalid Room UUID on create bookings: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Room Id")
		}

		newBookings = append(newBookings, booking_service.NewBooking{
			CheckIn:    checkIn,
			CheckOut:   checkOut,
			Total:      int(param.Total),
			RoomTypeId: roomTypeId,
			UserId:     userId,
			RoomId:     roomId,
		})
	}

	err := bg.service.CreateBookings(ctx, newBookings)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong dat phong duoc")
	}

	return &booking_pb.Empty{}, nil
}

func (bg *BookingGrpcHandler) DeleteBookingsByIds(ctx context.Context, req *booking_pb.DeleteBookingRequest) (*booking_pb.Empty, error) {

	ids := make([]pgtype.UUID, 0, len(req.GetBookingIds()))

	for _, idReq := range req.GetBookingIds() {
		var id pgtype.UUID
		if err := id.Scan(idReq); err != nil {
			zap.S().Info("Invalid Booking Id: ", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Booking Id")
		}
		ids = append(ids, id)
	}

	err := bg.service.DeleteBookingsByIds(ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong xoa duoc booking")
	}

	return &booking_pb.Empty{}, nil
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
