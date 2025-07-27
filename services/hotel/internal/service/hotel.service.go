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

func (hs *HotelService) GetHotel(ctx context.Context, id pgtype.UUID) (any, error) {

	hotel, err := hs.repository.GetHotelById(ctx, id)
	if err != nil {
		zap.S().Error("Failed to get hotel by ID: ", err)
		return nil, err
	}

	return hotel, nil
}

func (hs *HotelService) CreateHotel(ctx context.Context, name string) error {
	err := hs.repository.CreateHotel(ctx, hotel_repo.CreateHotelParams{
		Name: name,
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
