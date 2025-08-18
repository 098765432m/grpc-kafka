package image_service

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/repository/image"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type ImageService struct {
	repo *image_repo.Queries
}

func NewImageService(repo *image_repo.Queries) *ImageService {
	return &ImageService{
		repo: repo,
	}
}

func (is *ImageService) GetImage(ctx context.Context, id pgtype.UUID) (*image_repo.Image, error) {
	image, err := is.repo.GetImageById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (is *ImageService) UploadHotelImage(ctx context.Context, newImage image_repo.UploadHotelImageParams) error {
	err := is.repo.UploadHotelImage(ctx, newImage)
	if err != nil {
		return err
	}

	return nil
}

func (is *ImageService) UploadUserImage(ctx context.Context, newImage image_repo.UploadUserImageParams) error {
	err := is.repo.UploadUserImage(ctx, newImage)
	if err != nil {
		return err
	}

	return nil
}

func (is *ImageService) UploadRoomTypeImage(ctx context.Context, newImage image_repo.UploadRoomTypeImageParams) error {
	err := is.repo.UploadRoomTypeImage(ctx, newImage)
	if err != nil {
		return err
	}

	return nil
}

func (is *ImageService) GetImagesByHotelId(ctx context.Context, hotelId pgtype.UUID) ([]image_repo.Image, error) {
	images, err := is.repo.GetImagesByHotelId(ctx, hotelId)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (is *ImageService) GetImagesByHotelIds(ctx context.Context, hotelIds []pgtype.UUID) ([]image_repo.Image, error) {
	images, err := is.repo.GetImagesByHotelIds(ctx, hotelIds)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (is *ImageService) GetImageByUserId(ctx context.Context, userId pgtype.UUID) (*image_repo.Image, error) {

	image, err := is.repo.GetImageByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Info("No Record return ")
			return nil, common_error.ErrNoRows
		}

		zap.S().Errorln("Failed to get Image by User id: ", err)
		return nil, err
	}

	return &image, nil
}

func (is *ImageService) GetImagesByRoomTypeId(ctx context.Context, roomTypeId pgtype.UUID) ([]image_repo.Image, error) {

	images, err := is.repo.GetImagesByRoomTypeId(ctx, roomTypeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Errorln("Failed to get Image by User id: ", err)
			return nil, common_error.ErrNoRows
		}
		zap.S().Errorln("Failed to get Images by Room Type id")
		return nil, err
	}

	return images, nil
}

func (is *ImageService) DeleteImage(ctx context.Context, id pgtype.UUID) error {
	err := is.repo.DeleteImage(ctx, id)
	if err != nil {
		return err
	}

	return nil

}

func (is *ImageService) DeleteImages(ctx context.Context, ids []pgtype.UUID) error {
	err := is.repo.DeleteImages(ctx, ids)
	if err != nil {
		return err
	}

	return nil
}
