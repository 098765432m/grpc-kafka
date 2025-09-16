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

func (rs *RoomService) GetNumberOfRoomsPerRoomTypeByHotelIds(ctx context.Context, hotelIds []pgtype.UUID) ([]room_repo.GetNumberOfRoomsPerRoomTypeByHotelIdsRow, error) {

	rows, err := rs.repo.GetNumberOfRoomsPerRoomTypeByHotelIds(ctx, hotelIds)
	if err != nil {
		zap.S().Errorln("Failed to get number of Rooms per Room Type by HotelIDs")
		return nil, err
	}

	return rows, nil
}

// Return List of Available rooms for room type id
func (rs *RoomService) GetListOfAvailableRoomsByRoomTypeId(ctx context.Context, roomTypeId pgtype.UUID, numberOfRooms int) ([]pgtype.UUID, error) {

	roomIds, err := rs.repo.GetListOfAvailableRoomsByRoomTypeId(ctx, room_repo.GetListOfAvailableRoomsByRoomTypeIdParams{
		RoomTypeID:    roomTypeId,
		NumberOfRooms: int32(numberOfRooms),
	})
	if err != nil {
		zap.S().Errorln("Failed to get list of available rooms by room type id")
		return nil, err
	}

	if len(roomIds) < numberOfRooms {
		zap.S().Infoln("There is not enough rooms AVAILABLE.")
		return nil, common_error.ErrNoRows
	}

	return roomIds, nil
}

func (rs *RoomService) GetListOfRemainRooms(ctx context.Context, roomTypeId pgtype.UUID, bookedRoomIds []pgtype.UUID, numberOfRooms int) ([]pgtype.UUID, error) {

	zap.L().Info("Request: ", zap.Any("RoomType", roomTypeId), zap.Any("Booked RoomIds", bookedRoomIds), zap.Any("Number of Rooms", numberOfRooms))

	roomIds, err := rs.repo.GetListOfRemainRooms(ctx, room_repo.GetListOfRemainRoomsParams{
		RoomTypeID:    roomTypeId,
		BookedRoomIds: bookedRoomIds,
		NumberOfRooms: int32(numberOfRooms),
	})
	if err != nil {
		zap.S().Errorln("Failed to get list of remain rooms by: ", err)
		return nil, err
	}

	zap.S().Infoln("Remain Rooms: ", roomIds)

	if len(roomIds) < numberOfRooms {
		zap.S().Infoln("There is not enough rooms AVAILABLE.")
		return nil, common_error.ErrNoRows
	}

	return roomIds, nil
}

// Set Status of rooms
func (rs *RoomService) ChangeStatusOfRooms(ctx context.Context, roomIds []pgtype.UUID, room_status room_repo.RoomStatus) error {

	rows, err := rs.repo.ChangeStatusRoomsByIds(ctx, room_repo.ChangeStatusRoomsByIdsParams{
		RoomIds: roomIds,
		Status:  room_status,
	})
	if err != nil {
		zap.S().Error("Failed to change status of rooms", err)
		return err
	}

	// Catch when no rooms found
	if rows == 0 {
		zap.S().Infoln("No Rooms found with given ids")
		return common_error.ErrNoRows
	}

	return nil

}
