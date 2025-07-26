package image_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
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
