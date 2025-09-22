package room_type_service

import (
	"context"

	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RoomTypeService struct {
	repo *room_type_repo.Queries
}

func NewRoomTypeService(repo *room_type_repo.Queries) *RoomTypeService {
	return &RoomTypeService{
		repo: repo,
	}
}

func (rts *RoomTypeService) GetRoomTypeById(ctx context.Context, id pgtype.UUID) (*room_type_repo.RoomType, error) {
	roomType, err := rts.repo.GetRoomTypeById(ctx, id)
	if err != nil {
		zap.S().Errorln("Cannot get Room Type By id")
		return nil, err
	}

	return &roomType, nil
}

func (rts *RoomTypeService) GetRoomTypesByHotelId(ctx context.Context, hotelId pgtype.UUID) ([]room_type_repo.GetRoomTypesByHotelIdRow, error) {
	roomTypes, err := rts.repo.GetRoomTypesByHotelId(ctx, hotelId)
	if err != nil {
		zap.S().Errorln("Cannot get Room Type By id")
		return nil, err
	}

	return roomTypes, nil
}

func (rts *RoomTypeService) CreateRoomType(ctx context.Context, newRoomType *room_type_repo.CreateRoomTypeParams) error {

	err := rts.repo.CreateRoomType(ctx, *newRoomType)
	if err != nil {
		zap.S().Errorln("Failed to create Room Type: ", err)
		return err
	}

	return nil
}

func (rts *RoomTypeService) DeleteRoomTypeById(ctx context.Context, id pgtype.UUID) error {

	err := rts.repo.DeleteRoomTypeById(ctx, id)
	if err != nil {
		zap.S().Errorln("Delete Room by id")
		return err
	}

	return nil
}
