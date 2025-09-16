package booking_handler

import (
	"context"
	"time"

	booking_service "github.com/098765432m/grpc-kafka/booking/internal/service"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
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

func (bg *BookingGrpcHandler) DeleteBookingById(ctx context.Context, req *booking_pb.DeleteBookingByIdRequest) (*booking_pb.Empty, error) {

	var id pgtype.UUID
	if err := id.Scan(req.GetBookingId()); err != nil {
		zap.S().Infoln("Invalud Booking UUID")
		return nil, status.Error(codes.InvalidArgument, "Booking UUID khong hop le")
	}

	err := bg.service.DeleteBookingById(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong xoa duoc booking")
	}

	return &booking_pb.Empty{}, nil
}

func (bg *BookingGrpcHandler) DeleteBookingsByIds(ctx context.Context, req *booking_pb.DeleteBookingByIdsRequest) (*booking_pb.Empty, error) {

	ids, err := utils.ToPgUuidArray(req.GetBookingIds())
	if err != nil {
		zap.S().Infoln("Invalid Booking UUID")
		return nil, status.Error(codes.InvalidArgument, "Booking UUID khong hop le")
	}

	err = bg.service.DeleteBookingsByIds(ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong xoa duoc booking")
	}

	return &booking_pb.Empty{}, nil
}

// Return {roomTypeId, Number of occupied rooms in that room type}
func (bg *BookingGrpcHandler) GetNumberOfOccupiedRooms(ctx context.Context, req *booking_pb.GetNumberOfOccupiedRoomsRequest) (*booking_pb.GetNumberOfOccupiedRoomsResponse, error) {
	// Check Are room type Ids valid
	roomTypeIds, err := utils.ToPgUuidArray(req.GetRoomTypeIds())
	if err != nil {
		zap.S().Infoln("Invalid Room Type UUID")
		return nil, status.Error(codes.InvalidArgument, "Room Type ID khong hop le")
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

// Return {hotelId, roomTypeId, Number of occupied rooms in that room type}
func (bg *BookingGrpcHandler) GetNumberOfOccupiedRoomsByHotelIds(ctx context.Context, req *booking_pb.GetNumberOfOccupiedRoomsByHotelIdsRequest) (*booking_pb.GetNumberOfOccupiedRoomsByHotelIdsResponse, error) {
	// Check Are Hotel Ids valid
	hotelIds, err := utils.ToPgUuidArray(req.GetHotelIds())
	if err != nil {
		zap.S().Infoln("Invalid Hotel UUID")
		return nil, status.Error(codes.InvalidArgument, "Hotel ID khong hop le")
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
	results, err := bg.service.GetNumberOfOccupiedRoomsByHotelIds(ctx, hotelIds, checkInDate, checkOutDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the tra ve danh sach phong booked")
	}

	numberOccpiedRoomsHotelIdsResponse := make([]*booking_pb.ResultNumberOfOccupiedRoomsByHotelIds, 0, len(results))
	for _, result := range results {
		tempResult := &booking_pb.ResultNumberOfOccupiedRoomsByHotelIds{
			RoomTypeId:            result.RoomTypeID.String(),
			NumberOfOccupiedRooms: int32(result.NumberOfOccupiedRooms),
		}
		numberOccpiedRoomsHotelIdsResponse = append(numberOccpiedRoomsHotelIdsResponse, tempResult)
	}

	return &booking_pb.GetNumberOfOccupiedRoomsByHotelIdsResponse{
		Results: numberOccpiedRoomsHotelIdsResponse,
	}, nil
}

func (bg *BookingGrpcHandler) GetUnavailableRoomsByRoomTypeId(ctx context.Context, req *booking_pb.GetUnavailableRoomsByRoomTypeIdRequest) (*booking_pb.GetUnavailableRoomsByRoomTypeIdResponse, error) {

	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.RoomTypeId); err != nil {
		zap.S().Info("Invalid Room UUID: ", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid Room Id")
	}

	//Check dates valid
	var checkInDate pgtype.Date
	if err := checkInDate.Scan(req.CheckIn); err != nil {
		zap.S().Info("Invalid date format: ", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid date format")
	}

	var checkOutDate pgtype.Date
	if err := checkOutDate.Scan(req.CheckOut); err != nil {
		zap.S().Info("Invalid date format: ", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid date format")
	}

	if checkInDate.Time.After(checkOutDate.Time) {
		zap.S().Info("Invalid date - Check In is After Check out: ")
		return nil, status.Error(codes.InvalidArgument, "Invalid check in check out")
	}

	if time.Now().After(checkInDate.Time) {
		zap.S().Info("Invalid date - Check In is oudated: ")
		return nil, status.Error(codes.InvalidArgument, "Invalid check in check out")
	}

	roomIds, err := bg.service.GetUnavailableRoomsByRoomTypeId(ctx, roomTypeId, checkInDate, checkOutDate)
	if err != nil {
		zap.S().Errorln("Failed to Get UNAVAILABLE Rooms by Room Type Id: ", err)
		return nil, status.Error(codes.Internal, "Loi khong the danh sach phong da duoc dat truoc")
	}

	roomIdsStr := make([]string, 0, len(roomIds))
	for _, roomId := range roomIds {
		roomIdsStr = append(roomIdsStr, roomId.String())
	}

	return &booking_pb.GetUnavailableRoomsByRoomTypeIdResponse{
		RoomIds: roomIdsStr,
	}, nil

}
