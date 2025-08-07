package hotel_service

import (
	"context"

	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type HotelService struct {
	repository *hotel_repo.Queries
}

func NewHotelService(repo *hotel_repo.Queries) *HotelService {
	return &HotelService{
		repository: repo,
	}
}

func (hs *HotelService) GetHotel(ctx context.Context, id pgtype.UUID) (*hotel_repo.Hotel, error) {

	hotel, err := hs.repository.GetHotelById(ctx, id)
	if err != nil {
		zap.S().Error("Failed to get hotel by ID: ", err)
		return nil, err
	}

	return &hotel, nil
}

func (hs *HotelService) CreateHotel(ctx context.Context, newHotel *hotel_repo.CreateHotelParams) error {
	err := hs.repository.CreateHotel(ctx, hotel_repo.CreateHotelParams{
		Name:    newHotel.Name,
		Address: newHotel.Address,
	})
	if err != nil {
		zap.S().Error("Failed to create hotel: ", err)
		return err
	}

	return nil
}

func (hs *HotelService) GetAll(ctx context.Context) ([]hotel_repo.Hotel, error) {
	hotels, err := hs.repository.GetAll(ctx)
	if err != nil {
		zap.S().Error("Failed to get all hotels: ", err)
		return nil, err
	}

	return hotels, nil
}

func (hs *HotelService) UpdateHotel(ctx context.Context, hotelParam *hotel_repo.UpdateHotelByIdParams) error {

	err := hs.repository.UpdateHotelById(ctx, hotel_repo.UpdateHotelByIdParams{
		ID:      hotelParam.ID,
		Name:    hotelParam.Name,
		Address: hotelParam.Address,
	})

	if err != nil {
		return err
	}

	return nil
}

func (hs *HotelService) DeleteHotel(ctx context.Context, id pgtype.UUID) error {

	err := hs.repository.DeleteHotelById(ctx, id)
	if err != nil {
		zap.S().Errorln("Failed to delete hotel")
		return err
	}

	return nil
}
