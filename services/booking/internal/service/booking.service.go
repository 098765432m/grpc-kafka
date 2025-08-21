package booking_service

import (
	"context"

	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/repository/booking"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type BookingService struct {
	repo *booking_repo.Queries
}

func NewBookingService(repo *booking_repo.Queries) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (bs *BookingService) GetBookingById(ctx context.Context, id pgtype.UUID) (*booking_repo.Booking, error) {

	booking, err := bs.repo.GetBookingById(ctx, id)
	if err != nil {
		zap.S().Errorln("Cannot get Booking by id")

		return nil, err
	}

	return &booking, nil
}

func (bs *BookingService) CreateBooking(ctx context.Context, bookingParams *booking_repo.CreateBookingParams) error {
	if err := bs.repo.CreateBooking(ctx, *bookingParams); err != nil {
		zap.S().Errorln("Cannot create booking")
		return err
	}

	return nil
}

func (bs *BookingService) DeleteBooking(ctx context.Context, id pgtype.UUID) error {

	if err := bs.repo.DeleteBookingById(ctx, id); err != nil {
		zap.S().Errorln("Cannot delete Booking")
		return err
	}

	return nil
}

// Return Room Ids that occupied in a range of time
func (bs *BookingService) GetListOfOccupiedRooms(ctx context.Context, roomIds []pgtype.UUID, checkIn pgtype.Date, checkOut pgtype.Date) ([]pgtype.UUID, error) {
	occupiedRoomIds, err := bs.repo.GetListOfOccupiedRooms(ctx, booking_repo.GetListOfOccupiedRoomsParams{
		RoomIds:  roomIds,
		CheckIn:  checkIn,
		CheckOut: checkOut,
	})
	if err != nil {
		zap.S().Errorln("Failed to get list of occupied rooms: ", err)
		return nil, err
	}

	return occupiedRoomIds, nil
}

// Return number of Occupied rooms for each Room Type in a range of time
func (bs *BookingService) GetNumberOfOccupiedRooms(ctx context.Context, roomTypeIds []pgtype.UUID, checkIn pgtype.Date, checkOut pgtype.Date) ([]booking_repo.GetNumberOfOccupiedRoomsRow, error) {
	result, err := bs.repo.GetNumberOfOccupiedRooms(ctx, booking_repo.GetNumberOfOccupiedRoomsParams{
		RoomTypeIds: roomTypeIds,
		CheckIn:     checkIn,
		CheckOut:    checkOut,
	})
	if err != nil {
		zap.S().Errorln("Failed to get list of occupied rooms: ", err)
		return nil, err
	}

	return result, nil
}
