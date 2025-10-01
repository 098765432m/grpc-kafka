package image_handler

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	image_service "github.com/098765432m/grpc-kafka/image/internal/application"
	image_repo "github.com/098765432m/grpc-kafka/image/internal/infrastructure/repository/sqlc/image"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Reuse to get image urls map
func (ig *ImageGrpcHandler) getImageUrlsMap(ctx context.Context, images []image_repo.Image) (map[string]string, error) {

	imagePublicIds := make([]string, 0, len(images))
	for _, image := range images {
		imagePublicIds = append(imagePublicIds, image.PublicID)
	}

	// Get public IDs Map
	imageUrlsMap, err := ig.service.GetImageUrlsMap(ctx, imagePublicIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	return imageUrlsMap, nil
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

	url, err := ig.service.GetImageUrl(ctx, image.PublicID)
	if err != nil {
		zap.S().Errorln("Failed to get image url: ", err)
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	return &image_pb.GetImageByIdResponse{
		Image: &image_pb.Image{
			Id:         image.ID.String(),
			PublicId:   image.PublicID,
			Format:     image.Format,
			Url:        url,
			HotelId:    image.HotelID.String(),
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

	imageUrlsMap, err := ig.getImageUrlsMap(ctx, images)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay image urls map duoc")
	}

	var hotelImages []*image_pb.HotelImage
	for _, img := range images {
		hotelImage := &image_pb.HotelImage{
			Id:       img.ID.String(),
			Url:      imageUrlsMap[img.PublicID],
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

	url, err := ig.service.GetImageUrl(ctx, image.PublicID)
	if err != nil {
		zap.S().Errorln("Failed to get image url: ", err)
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	return &image_pb.GetImageByUserIdResponse{
		Image: &image_pb.UserImage{
			Id:       image.ID.String(),
			Url:      url,
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

	imageUrlsMap, err := ig.getImageUrlsMap(ctx, images)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay image urls map duoc")
	}

	var roomTypeImages []*image_pb.RoomTypeImage
	for _, img := range images {
		hotelImage := &image_pb.RoomTypeImage{
			Id:         img.ID.String(),
			Url:        imageUrlsMap[img.PublicID],
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

func (ig *ImageGrpcHandler) GetImagesByRoomTypeIds(ctx context.Context, req *image_pb.GetImagesByRoomTypeIdsRequest) (*image_pb.GetImagesByRoomTypeIdsResponse, error) {
	var roomTypeIds []pgtype.UUID
	for _, roomTypeIdReq := range req.RoomTypeIds {
		var roomTypeId pgtype.UUID
		if err := roomTypeId.Scan(roomTypeIdReq); err != nil {
			zap.S().Info("Invalid Room Type UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi UUID loai phong")
		}

		roomTypeIds = append(roomTypeIds, roomTypeId)
	}

	images, err := ig.service.GetImagesByRoomTypeIds(ctx, roomTypeIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi lay hinh anh loai phong")
	}

	imageUrlsMap, err := ig.getImageUrlsMap(ctx, images)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay image urls map duoc")
	}

	var roomTypeImages []*image_pb.RoomTypeImage
	for _, img := range images {
		roomTypeImage := &image_pb.RoomTypeImage{
			Id:         img.ID.String(),
			Url:        imageUrlsMap[img.PublicID],
			PublicId:   img.PublicID,
			Format:     img.Format,
			RoomTypeId: img.RoomTypeID.String(),
		}

		roomTypeImages = append(roomTypeImages, roomTypeImage)
	}

	return &image_pb.GetImagesByRoomTypeIdsResponse{
		Images: roomTypeImages,
	}, nil
}

func (ig *ImageGrpcHandler) GetImagesByHotelIds(ctx context.Context, req *image_pb.GetImagesByHotelIdsRequest) (*image_pb.GetImagesByHotelIdsResponse, error) {
	var hotelIds []pgtype.UUID

	// ChecK UUID
	for _, hotelId := range req.HotelIds {
		var tempHotelId pgtype.UUID
		if err := tempHotelId.Scan(hotelId); err != nil {
			zap.S().Info("Invalid Hotel UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi UUID khach san")
		}

		hotelIds = append(hotelIds, tempHotelId)
	}

	images, err := ig.service.GetImagesByHotelIds(ctx, hotelIds)
	if err != nil {
		zap.S().Errorln("Failed to get images by hotel ids: ", err)
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	imageUrlsMap, err := ig.getImageUrlsMap(ctx, images)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay image urls map duoc")
	}

	var hotelImages []*image_pb.HotelImage
	for _, img := range images {
		hotelImage := &image_pb.HotelImage{
			Id:       img.ID.String(),
			Url:      imageUrlsMap[img.PublicID],
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

func (ig *ImageGrpcHandler) GetImagesByUserIds(ctx context.Context, req *image_pb.GetImagesByUserIdsRequest) (*image_pb.GetImagesByUserIdsResponse, error) {
	userIds := make([]pgtype.UUID, 0, len(req.GetUserIds()))
	for _, userIdReq := range req.GetUserIds() {
		var id pgtype.UUID
		if err := id.Scan(userIdReq); err != nil {
			zap.S().Info("Invalid User UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi user UUID")
		}

		userIds = append(userIds, id)
	}

	images, err := ig.service.GetImagesByUserIds(ctx, userIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "Khong the tim images bang user ids")
	}

	imageUrlsMap, err := ig.getImageUrlsMap(ctx, images)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the lay image urls map duoc")
	}

	imagesGrpcResult := make([]*image_pb.UserImage, 0, len(images))
	for _, image := range images {
		imageGrpc := &image_pb.UserImage{
			Id:       image.ID.String(),
			Url:      imageUrlsMap[image.PublicID],
			PublicId: image.PublicID,
			Format:   image.Format,
			UserId:   image.UserID.String(),
		}

		imagesGrpcResult = append(imagesGrpcResult, imageGrpc)
	}
	return &image_pb.GetImagesByUserIdsResponse{
		Images: imagesGrpcResult,
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
