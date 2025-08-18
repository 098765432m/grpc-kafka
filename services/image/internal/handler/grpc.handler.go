package image_handler

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
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

func (ig *ImageGrpcHandler) GetImageById(ctx context.Context, req *image_pb.GetImageByIdRequest) (*image_pb.GetImageByIdResponse, error) {
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

	return &image_pb.GetImageByIdResponse{
		Image: &image_pb.Image{
			Id:         image.ID.String(),
			PublicId:   image.PublicID,
			Format:     image.Format,
			HotelId:    hotelId,
			UserId:     image.UserID.String(),
			RoomTypeId: image.RoomTypeID.String(),
		},
	}, nil
}

func (ig *ImageGrpcHandler) GetImagesByHotelId(ctx context.Context, req *image_pb.GetImagesByHotelIdRequest) (*image_pb.GetImagesByHotelIdResponse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	images, err := ig.service.GetImagesByHotelId(ctx, hotelId)
	if err != nil {
		if errors.Is(err, common_error.ErrNoRows) {
			return nil, nil
		}
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

	return &image_pb.GetImagesByHotelIdResponse{
		Images: hotelImages,
	}, nil
}

func (ig *ImageGrpcHandler) GetImageByUserId(ctx context.Context, req *image_pb.GetImageByUserIdRequest) (*image_pb.GetImageByUserIdResponse, error) {
	var userId pgtype.UUID
	if err := userId.Scan(req.UserId); err != nil {
		return nil, err
	}

	image, err := ig.service.GetImageByUserId(ctx, userId)
	if err != nil {
		switch {

		case errors.Is(err, common_error.ErrNoRows):
			return nil, nil
		default:

			return nil, err
		}
	}

	return &image_pb.GetImageByUserIdResponse{
		Image: &image_pb.UserImage{
			Id:       image.ID.String(),
			PublicId: image.PublicID,
			Format:   image.Format,
			UserId:   image.UserID.String(),
		},
	}, nil
}

func (ig *ImageGrpcHandler) GetImagesByRoomTypeId(ctx context.Context, req *image_pb.GetImagesByRoomTypeIdRequest) (*image_pb.GetImagesByRoomTypeIdResponse, error) {
	var roomTypeId pgtype.UUID
	if err := roomTypeId.Scan(req.RoomTypeId); err != nil {
		return nil, err
	}

	images, err := ig.service.GetImagesByRoomTypeId(ctx, roomTypeId)
	if err != nil {
		return nil, err
	}

	var roomTypeImages []*image_pb.RoomTypeImage
	for _, img := range images {
		hotelImage := &image_pb.RoomTypeImage{
			Id:         img.ID.String(),
			PublicId:   img.PublicID,
			Format:     img.Format,
			RoomTypeId: img.RoomTypeID.String(),
		}

		roomTypeImages = append(roomTypeImages, hotelImage)
	}

	return &image_pb.GetImagesByRoomTypeIdResponse{
		Images: roomTypeImages,
	}, nil
}

func (ig *ImageGrpcHandler) GetImagesByHotelIds(ctx context.Context, req *image_pb.GetImagesByHotelIdsRequest) (*image_pb.GetImagesByHotelIdsResponse, error) {
	var hotelIds []pgtype.UUID

	// ChecK UUID
	for _, hotelId := range req.HotelIds {
		var tempHotelId pgtype.UUID
		if err := tempHotelId.Scan(hotelId); err != nil {
			return nil, err
		}

		hotelIds = append(hotelIds, tempHotelId)
	}

	images, err := ig.service.GetImagesByHotelIds(ctx, hotelIds)
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

	return &image_pb.GetImagesByHotelIdsResponse{
		Images: hotelImages,
	}, nil
}

func (ig *ImageGrpcHandler) UploadImage(ctx context.Context, req *image_pb.UploadImageRequest) (*image_pb.UploadImageResponse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.HotelId); err != nil {
		return nil, err
	}

	err := ig.service.UploadHotelImage(ctx, image_repo.UploadHotelImageParams{
		PublicID: req.PublicId,
		Format:   req.Format,
		HotelID:  hotelId,
	})

	if err != nil {
		return nil, err
	}

	return &image_pb.UploadImageResponse{}, nil

}

func (ig *ImageGrpcHandler) DeleteImageById(ctx context.Context, req *image_pb.DeleteImageByIdRequest) (*image_pb.DeleteImageByIdResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		return nil, err
	}

	err := ig.service.DeleteImage(ctx, id)
	if err != nil {
		return nil, err
	}

	return &image_pb.DeleteImageByIdResponse{}, nil
}

func (ig *ImageGrpcHandler) DeleteImagesByIds(ctx context.Context, req *image_pb.DeleteImagesByIdsRequest) (*image_pb.DeleteImagesByIdsResponse, error) {

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

	return &image_pb.DeleteImagesByIdsResponse{}, nil
}
