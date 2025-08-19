package room_service

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RoomService struct {
	repo *room_repo.Queries
}

func NewRoomService(repo *room_repo.Queries) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (rs *RoomService) GetRoomsByHotelId(ctx context.Context, hotelId pgtype.UUID) ([]room_repo.Room, error) {

	rooms, err := rs.repo.GetRoomsByHotelId(ctx, hotelId)
	if err != nil {

		zap.S().Errorln("Cannot get Rooms by Hotel Id: ", err)
		return nil, err
	}

	return rooms, nil
}

func (rs *RoomService) GetRoomsByRoomTypeId(ctx context.Context, roomTypeId pgtype.UUID) ([]room_repo.Room, error) {

	rooms, err := rs.repo.GetRoomsByRoomTypeId(ctx, roomTypeId)
	if err != nil {

		zap.S().Errorln("Cannot get Rooms by Hotel Id: ", err)
		return nil, err
	}

	return rooms, nil
}

func (rs *RoomService) GetRoomsById(ctx context.Context, id pgtype.UUID) (*room_repo.Room, error) {

	room, err := rs.repo.GetRoomsById(ctx, id)
	if err != nil {
		if errors.Is(err, common_error.ErrNoRows) {
			zap.S().Info("Phong khong ton tai")
			return nil, common_error.ErrNoRows
		}

		zap.S().Errorln("Cannot get Rooms by Hotel Id")
		return nil, err
	}

	return &room, nil

}

func (rs *RoomService) CreateRoom(ctx context.Context, newRoom *room_repo.CreateRoomParams) error {

	err := rs.repo.CreateRoom(ctx, *newRoom)
	if err != nil {
		zap.S().Errorln("Cannot create Room")
		return err
	}

	return nil
}

func (rs *RoomService) DeleteRoom(ctx context.Context, id pgtype.UUID) error {

	err := rs.repo.DeleteRoomById(ctx, id)
	if err != nil {
		zap.S().Errorln("Cannot delete Room")

		return err
	}

	return nil
}
