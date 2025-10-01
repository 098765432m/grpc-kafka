package image_service

import (
	"context"
	"errors"
	"sync"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/infrastructure/repository/sqlc/image"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type ImageService struct {
	repo *image_repo.Queries
	cld  *cloudinary.Cloudinary
}

func NewImageService(repo *image_repo.Queries, cld *cloudinary.Cloudinary) *ImageService {
	return &ImageService{
		repo: repo,
		cld:  cld,
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

func (is *ImageService) UpdateImageById(ctx context.Context, id pgtype.UUID, publicId string, format string) error {

	oldImage, err := is.repo.GetImageById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Infoln("No record found to update")
			return common_error.ErrNoRows
		}
		zap.S().Errorln("Failed to get image by id: ", err)
		return err
	}

	_, err = is.repo.UpdateImageById(ctx, image_repo.UpdateImageByIdParams{
		ID:       id,
		PublicID: publicId,
		Format:   format,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Infoln("No record found to update")
			return common_error.ErrNoRows
		}
		zap.S().Errorln("Failed to update image by id: ", err)
		return err
	}

	// Delete old image from Cloudinary
	is.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: oldImage.PublicID,
	})

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

func (is *ImageService) GetImagesByUserIds(ctx context.Context, userIds []pgtype.UUID) ([]image_repo.Image, error) {

	images, err := is.repo.GetImagesByUserIds(ctx, userIds)
	if err != nil {
		zap.S().Error("Failed to get Images by User ids: ", err)
		return nil, err
	}

	return images, err
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

func (is *ImageService) GetImagesByRoomTypeIds(ctx context.Context, roomTypeIds []pgtype.UUID) ([]image_repo.Image, error) {

	images, err := is.repo.GetImagesByRoomTypeIds(ctx, roomTypeIds)
	if err != nil {
		zap.S().Errorln("Failed to get Images by Room Type id")
		return nil, err
	}

	return images, nil
}

func (is *ImageService) DeleteImage(ctx context.Context, id pgtype.UUID) error {

	publicId, err := is.repo.DeleteImage(ctx, id)
	if err != nil {
		return err
	}

	// Delete image from Cloudinary
	is.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})

	return nil

}

func (is *ImageService) DeleteImages(ctx context.Context, ids []pgtype.UUID) error {
	publicIds, err := is.repo.DeleteImages(ctx, ids)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	// Delete images from Cloudinary
	// utilize goroutine to speed up the process
	for _, publicId := range publicIds {
		wg.Add(1)
		go func(pid string) {
			defer wg.Done()
			is.cld.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: publicId,
			})
		}(publicId)
	}
	wg.Wait()

	return nil
}

// GET image URL from Cloudinary by public ID
func (is *ImageService) GetImageUrl(ctx context.Context, publicId string) (string, error) {
	img, err := is.cld.Image(publicId)
	if err != nil {
		zap.S().Error("Failed to get image url: ", err)
		return "", err
	}

	url, err := img.String()
	if err != nil {
		zap.S().Error("Failed to convert image to string: ", err)
		return "", err
	}

	return url, nil
}

// GET image URLs from Cloudinary by public IDs
func (is *ImageService) GetImageUrls(ctx context.Context, publicIds []string) ([]string, error) {
	urls := make([]string, 0, len(publicIds))

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errGroup := make([]error, 0, len(publicIds))

	for _, publicId := range publicIds {
		wg.Add(1)
		go func(pid string) {
			defer wg.Done()

			img, err := is.cld.Image(publicId)
			if err != nil {
				mu.Lock()
				errGroup = append(errGroup, err)
				mu.Unlock()
				return
			}

			url, err := img.String()
			if err != nil {
				zap.S().Error("Failed to convert image to string: ", err)
				errGroup = append(errGroup, err)
				return
			}

			urls = append(urls, url)

		}(publicId)
	}

	wg.Wait()
	if len(errGroup) > 0 {
		return nil, common_error.ErrInternalServer
	}

	return urls, nil
}

// GET image URLs
// Return map[publicId]url
func (is *ImageService) GetImageUrlsMap(ctx context.Context, publicIds []string) (map[string]string, error) {
	urlsMap := make(map[string]string)

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errGroup := make([]error, 0, len(publicIds))

	for _, publicId := range publicIds {
		wg.Add(1)
		go func(pid string) {
			defer wg.Done()

			img, err := is.cld.Image(publicId)
			if err != nil {
				mu.Lock()
				errGroup = append(errGroup, err)
				mu.Unlock()
				return
			}

			url, err := img.String()
			if err != nil {
				zap.S().Error("Failed to convert image to string: ", err)
				errGroup = append(errGroup, err)
				return
			}

			urlsMap[pid] = url

		}(publicId)
	}

	wg.Wait()
	if len(errGroup) > 0 {
		return nil, common_error.ErrInternalServer
	}

	return urlsMap, nil
}
