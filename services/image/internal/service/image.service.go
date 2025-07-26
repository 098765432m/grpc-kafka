package image_service

import (
	"context"

	image_repo "github.com/098765432m/grpc-kafka/image/internal/repository/image"
	"github.com/jackc/pgx/v5/pgtype"
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

func (is *ImageService) UpdaloadImage(ctx context.Context, newImage image_repo.UploadImageParams) error {
	err := is.repo.UploadImage(ctx, newImage)
	if err != nil {
		return err
	}

	return nil
}

func (is *ImageService) GetHotelImages(ctx context.Context, hotelId pgtype.UUID) ([]image_repo.Image, error) {
	images, err := is.repo.GetHotelImages(ctx, hotelId)
	if err != nil {
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
