package booking_service

import (
	"context"
	"fmt"
	"strings"

	booking_domain "github.com/098765432m/grpc-kafka/booking/internal/domain"
	booking_repo_mapping "github.com/098765432m/grpc-kafka/booking/internal/infrastructure/repository"
	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/infrastructure/repository/sqlc/booking"
	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BookingService struct {
	conn *pgxpool.Pool
	repo *booking_repo.Queries
}

func NewBookingService(conn *pgxpool.Pool, repo *booking_repo.Queries) *BookingService {
	return &BookingService{
		conn: conn,
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

type GetBookingsByUserIdParams struct {
	UserId         pgtype.UUID
	CheckDateStart pgtype.Date
	CheckDateEnd   pgtype.Date
	Size           int
	Offset         int
}

func (bs *BookingService) GetBookingsByUserId(ctx context.Context, params *GetBookingsByUserIdParams) ([]booking_domain.Booking, error) {

	bookings, err := bs.repo.GetBookingsByUserId(ctx, booking_repo.GetBookingsByUserIdParams{
		UserID:         params.UserId,
		CheckDateStart: params.CheckDateStart,
		CheckDateEnd:   params.CheckDateEnd,
		Size:           int32(params.Size),
		NumberOfOffset: int32(params.Offset),
	})
	if err != nil {
		zap.S().Errorln("Failed to Get Bookings by User Id: ", err)
		return nil, err
	}

	result := booking_repo_mapping.FromBookingsRepoToBookingsDomain(bookings)

	return result, nil
}

func (bs *BookingService) CreateBooking(ctx context.Context, bookingParams *booking_repo.CreateBookingParams) error {
	if err := bs.repo.CreateBooking(ctx, *bookingParams); err != nil {
		zap.S().Errorln("Cannot create booking")
		return err
	}

	return nil
}

type NewBooking struct {
	CheckIn    pgtype.Date
	CheckOut   pgtype.Date
	Total      int
	RoomTypeId pgtype.UUID
	UserId     pgtype.UUID
	RoomId     pgtype.UUID
}

func (bs *BookingService) CreateBookings(ctx context.Context, newBookingParams []NewBooking) error {

	stmt := `INSERT INTO bookings (check_in, check_out, total, room_type_id, user_id, room_id) VALUES `

	if len(newBookingParams) == 0 {
		zap.S().Infoln("No New Booking to create")
		return common_error.ErrBadRequest
	}

	args := []any{}
	placeholders := make([]string, 0, len(newBookingParams))

	for index, param := range newBookingParams {

		n := index * 6
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5, n+6))

		args = append(args, param.CheckIn, param.CheckOut, param.Total, param.RoomTypeId, param.UserId, param.RoomId)
	}

	stmt += strings.Join(placeholders, ", ")
	stmt += ";"

	zap.S().Info("Create Booking with statment: ", stmt)

	_, err := bs.conn.Exec(ctx, stmt, args...)
	if err != nil {
		zap.S().Errorln("Cannot create bookings: ", err)
		return err
	}

	return nil
}

func (bs *BookingService) DeleteBookingById(ctx context.Context, id pgtype.UUID) error {

	if err := bs.repo.DeleteBookingById(ctx, id); err != nil {
		zap.S().Errorln("Cannot delete Booking")
		return err
	}

	return nil
}

func (bs *BookingService) DeleteBookingsByIds(ctx context.Context, ids []pgtype.UUID) error {

	if err := bs.repo.DeleteBookingsByIds(ctx, ids); err != nil {
		zap.S().Errorln("Cannot delete Bookings By Ids")
		return err
	}

	return nil
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

// Return number of Occupied rooms for each Room Type By Hotel in a range of time
func (bs *BookingService) GetNumberOfOccupiedRoomsByHotelIds(ctx context.Context, hotelIds []pgtype.UUID, checkIn pgtype.Date, checkOut pgtype.Date) ([]booking_repo.GetNumberOfOccupiedRoomsByHotelIdsRow, error) {
	result, err := bs.repo.GetNumberOfOccupiedRoomsByHotelIds(ctx, booking_repo.GetNumberOfOccupiedRoomsByHotelIdsParams{
		HotelIds: hotelIds,
		CheckIn:  checkIn,
		CheckOut: checkOut,
	})
	if err != nil {
		zap.S().Errorln("Failed to get list of occupied rooms: ", err)
		return nil, err
	}

	return result, nil
}

// Return all rooms that booked or UNAVAILABLE in range of time
func (bs *BookingService) GetUnavailableRoomsByRoomTypeId(ctx context.Context, roomTypeId pgtype.UUID, checkIn pgtype.Date, checkOut pgtype.Date) ([]pgtype.UUID, error) {

	roomIds, err := bs.repo.GetUnavailableRoomsByRoomTypeId(ctx, booking_repo.GetUnavailableRoomsByRoomTypeIdParams{
		RoomTypeID: roomTypeId,
		CheckIn:    checkIn,
		CheckOut:   checkOut,
	})
	if err != nil {
		zap.S().Errorln("Failed to get UNAVAILABLE Rooms by Room Type Id", err)
		return nil, err
	}

	return roomIds, nil
}
