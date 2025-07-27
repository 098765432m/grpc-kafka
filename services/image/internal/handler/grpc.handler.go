package image_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/repository/image"
	image_service "github.com/098765432m/grpc-kafka/image/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
)

type ImageGrpcHandler struct {
	image_pb.UnimplementedImageServiceServer
	service *image_service.ImageService
}

func NewImageGrpcHandler(service *image_service.ImageService) *ImageGrpcHandler {
	return &ImageGrpcHandler{
		service: service,
	}
}

func (ig *ImageGrpcHandler) GetImage(ctx context.Context, req *image_pb.GetImageRequest) (*image_pb.GetImageResponse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		return nil, err
	}

	image, err := ig.service.GetImage(ctx, id)
	if err != nil {
		return nil, err
	}

	var hotelId string
	if image.HotelID.Valid {
		hotelId = image.HotelID.String()
	} else {
		hotelId = ""
	}

	return &image_pb.GetImageResponse{
		Image: &image_pb.Image{
			Id:       image.ID.String(),
			PublicId: image.PublicID,
			Format:   image.Format,
			HotelId:  hotelId,
		},
	}, nil
}

func (ig *ImageGrpcHandler) GetHotelImages(ctx context.Context, req *image_pb.GetHotelImagesRequest) (*image_pb.GetHotelImagesResponse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	images, err := ig.service.GetHotelImages(ctx, hotelId)
	if err != nil {
		return nil, err
	}

	var hotelImages []*image_pb.HotelImage
	for _, img := range images {
		hotelImage := &image_pb.HotelImage{
			Id:       img.ID.String(),
			PublicId: img.PublicID,
			Format:   img.Format,
			HotelId:  img.HotelID.String(),
		}

		hotelImages = append(hotelImages, hotelImage)
	}

	return &image_pb.GetHotelImagesResponse{
		Images: hotelImages,
	}, nil
}

func (ig *ImageGrpcHandler) UploadImage(ctx context.Context, req *image_pb.UploadImageRequest) (*image_pb.UploadImageResponse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	err := ig.service.UpdaloadImage(ctx, image_repo.UploadImageParams{
		PublicID: req.PublicId,
		Format:   req.Format,
		HotelID:  hotelId,
	})

	if err != nil {
		return nil, err
	}

	return &image_pb.UploadImageResponse{}, nil

}

func (ig *ImageGrpcHandler) DeleteImage(ctx context.Context, req *image_pb.DeleteImageRequest) (*image_pb.DeleteImageResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		return nil, err
	}

	err := ig.service.DeleteImage(ctx, id)
	if err != nil {
		return nil, err
	}

	return &image_pb.DeleteImageResponse{}, nil
}

func (ig *ImageGrpcHandler) DeleteImages(ctx context.Context, req *image_pb.DeleteImagesRequest) (*image_pb.DeleteImagesResponse, error) {

	// Convert string to pgtype UUID
	var ids []pgtype.UUID
	for _, req_id := range req.Ids {
		var uuid pgtype.UUID
		if err := uuid.Scan(req_id); err != nil {
			return nil, err
		}

		if uuid.Valid {
			ids = append(ids, uuid)
		}
	}

	err := ig.service.DeleteImages(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &image_pb.DeleteImagesResponse{}, nil
}
