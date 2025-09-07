package booking_service

import (
	"context"
	"fmt"
	"strings"

	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/repository/booking"
	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type BookingService struct {
	conn *pgx.Conn
	repo *booking_repo.Queries
}

func NewBookingService(conn *pgx.Conn, repo *booking_repo.Queries) *BookingService {
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
	placeholders := make([]string, len(newBookingParams))

	for index, param := range newBookingParams {

		n := index * 6
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5, n+6))

		args = append(args, param.CheckIn.Time.String(), param.CheckOut.Time.String(), param.Total, param.RoomTypeId.String(), param.UserId.String(), param.RoomId.String())
	}

	stmt += strings.Join(placeholders, ", ")
	stmt += ";"

	zap.S().Info("Create Booking with statment: ", stmt)

	_, err := bs.conn.Exec(ctx, stmt, args...)
	if err != nil {
		zap.S().Errorln("Cannot create bookings")
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
