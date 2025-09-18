package hotel_service

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type HotelService struct {
	repo        *hotel_repo.Queries
	imageClient image_pb.ImageServiceClient
}

func NewHotelService(repo *hotel_repo.Queries, imageClient image_pb.ImageServiceClient) *HotelService {
	return &HotelService{
		repo:        repo,
		imageClient: imageClient,
	}
}

func (hs *HotelService) GetHotelById(ctx context.Context, id pgtype.UUID) (*hotel_repo.Hotel, error) {

	hotel, err := hs.repo.GetHotelById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common_error.ErrNoRows
		}

		zap.S().Error("Failed to get hotel by ID: ", err)
		return nil, err
	}

	return &hotel, nil
}

func (hs *HotelService) CreateHotel(ctx context.Context, newHotel *hotel_repo.CreateHotelParams) error {
	err := hs.repo.CreateHotel(ctx, hotel_repo.CreateHotelParams{
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
	hotels, err := hs.repo.GetAll(ctx)
	if err != nil {
		zap.S().Error("Failed to get all hotels: ", err)
		return nil, err
	}
	return hotels, nil
}

func (hs *HotelService) GetHotelsByAddress(ctx context.Context, address pgtype.Text, hotelName pgtype.Text) ([]pgtype.UUID, error) {
	hotelIds, err := hs.repo.GetHotelsByAddress(ctx, hotel_repo.GetHotelsByAddressParams{
		Address:   address,
		HotelName: hotelName,
	})
	if err != nil {
		zap.S().Error("Failed to get Hotels By Address: ", err)
		return nil, err
	}

	return hotelIds, nil
}

func (hs *HotelService) FilterHotels(ctx context.Context, roomTypeIds []pgtype.UUID, minPrice pgtype.Int4, maxPrice pgtype.Int4) ([]hotel_repo.FilterHotelsRow, error) {
	result, err := hs.repo.FilterHotels(ctx, hotel_repo.FilterHotelsParams{
		RoomTypeIds: roomTypeIds,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
	})
	zap.S().Infoln("Filter service")
	zap.S().Infoln(result)
	if err != nil {
		zap.S().Errorln("Failed to Filter Hotels: ", err)
		return nil, err
	}

	return result, nil
}

func (hs *HotelService) UpdateHotelById(ctx context.Context, hotelParam *hotel_repo.UpdateHotelByIdParams) error {

	err := hs.repo.UpdateHotelById(ctx, hotel_repo.UpdateHotelByIdParams{
		ID:      hotelParam.ID,
		Name:    hotelParam.Name,
		Address: hotelParam.Address,
	})

	if err != nil {
		return err
	}

	return nil
}

func (hs *HotelService) DeleteHotelById(ctx context.Context, id pgtype.UUID) error {

	err := hs.repo.DeleteHotelById(ctx, id)
	if err != nil {
		zap.S().Errorln("Failed to delete hotel")
		return err
	}

	return nil
}
