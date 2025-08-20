package api_handler

import (
	"net/http"

	api_dto "github.com/098765432m/grpc-kafka/api-gateway/internal/dto"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HotelHandler struct {
	hotelClient    hotel_pb.HotelServiceClient
	roomTypeClient room_type_pb.RoomTypeServiceClient
	roomClient     room_pb.RoomServiceClient
	userClient     user_pb.UserServiceClient
	imageClient    image_pb.ImageServiceClient
	ratingClient   rating_pb.RatingServiceClient
}

func NewHotelHandler(
	hotelClient hotel_pb.HotelServiceClient,
	roomTypeClient room_type_pb.RoomTypeServiceClient,
	roomClient room_pb.RoomServiceClient,
	userClient user_pb.UserServiceClient,
	imageClient image_pb.ImageServiceClient,
	ratingClient rating_pb.RatingServiceClient,
) *HotelHandler {
	return &HotelHandler{
		hotelClient:    hotelClient,
		roomTypeClient: roomTypeClient,
		roomClient:     roomClient,
		userClient:     userClient,
		imageClient:    imageClient,
		ratingClient:   ratingClient,
	}
}

func (hh *HotelHandler) RegisterRoutes(router *gin.RouterGroup) {
	hotelHandler := router.Group("/hotels")

	hotelHandler.GET("/", hh.GetAll)
	hotelHandler.GET("/:id", hh.GetHotelById)

	hotelHandler.GET("/:id/room-types", hh.GetRoomTypesByHotelId)
	hotelHandler.GET("/:id/rooms", hh.GetRoomsByHotelId)
	hotelHandler.GET("/:id/ratings", hh.GetRatingsByHotelId)

}

func (hh *HotelHandler) GetAll(ctx *gin.Context) {

	hotels, err := hh.hotelClient.GetAllHotels(ctx, &hotel_pb.GetAllHotelsRequest{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to get all Hotels"))
		return
	}

	// Get hotel Ids for images
	hotelIds := make([]string, 0, len(hotels.Hotels)) // Improve By preallocate than just []string
	for _, hotel := range hotels.Hotels {
		hotelIds = append(hotelIds, hotel.Id)
	}

	// Get all images with set of hotel ids
	images, err := hh.imageClient.GetImagesByHotelIds(ctx, &image_pb.GetImagesByHotelIdsRequest{
		HotelIds: hotelIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to get Images By Hotel ids"))
		return
	}

	// Response type
	var responses []api_dto.HotelResponse

	// Data Merging
	for _, hotel := range hotels.Hotels {

		// Merge Hotel
		resp := api_dto.HotelResponse{
			Id:      hotel.Id,
			Name:    hotel.Name,
			Address: hotel.Address,
		}

		// Merge Image into hotel
		for _, img := range images.Images {
			if img.HotelId == hotel.Id { // Append image if match hotelId
				resp.Images = append(resp.Images, api_dto.HotelImage{
					Id:       img.Id,
					PublicId: img.PublicId,
					Format:   img.Format,
				})
			}
		}

		responses = append(responses, resp)
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(responses, "Hotels retrieved successfully"))
}

func (hh *HotelHandler) GetHotelById(ctx *gin.Context) {

	id := ctx.Param("id")
	hotelGrpc, err := hh.hotelClient.GetHotelById(ctx, &hotel_pb.GetHotelByIdRequest{
		Id: id,
	})
	hotel := hotelGrpc.GetHotel()
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Khach san khong ton tai"))
				return
			}
		}

		zap.S().Errorln(err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong khi lay khach san"))
		return
	}

	images, err := hh.imageClient.GetImagesByHotelId(ctx, &image_pb.GetImagesByHotelIdRequest{HotelId: hotel.Id})
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong khi lay hinh anh khach san"))
		return
	}

	resp := api_dto.HotelResponse{
		Id:      hotel.GetId(),
		Name:    hotel.GetName(),
		Address: hotel.GetAddress(),
	}

	for _, img := range images.GetImages() {
		resp.Images = append(resp.Images, api_dto.HotelImage{
			Id:       img.GetId(),
			PublicId: img.GetPublicId(),
			Format:   img.GetFormat(),
		})
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(resp, "Hotel retrieved successfully"))
}

func (hh *HotelHandler) GetRatingsByHotelId(ctx *gin.Context) {
	id := ctx.Param("id")

	ratingGrpcResult, err := hh.ratingClient.GetRatingsByHotelId(ctx, &rating_pb.GetRatingsByHotelIdRequest{
		HotelId: id,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	// Get User Ids for merge user to rating
	var userIds []string
	for _, rating := range ratingGrpcResult.GetRatings() {
		userIds = append(userIds, rating.UserId)
	}

	usersGrpcResult, err := hh.userClient.GetUsersByIds(ctx, &user_pb.GetUsersByIdsRequest{
		Ids: userIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	imagesGrpcResult, err := hh.imageClient.GetImagesByUserIds(ctx, &image_pb.GetImagesByUserIdsRequest{
		UserIds: userIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	// Look up with O(1) -- no need nested loop
	imageMap := make(map[string]*api_dto.RatingImageResponse)
	for _, image := range imagesGrpcResult.GetImages() {
		// only work when image 1-1 relation with user
		imageMap[image.UserId] = &api_dto.RatingImageResponse{
			ImageId:  image.Id,
			PublicId: image.PublicId,
			Format:   image.Format,
			UserId:   image.UserId,
		}
	}

	// Look up with O(1) -- no need nested loop
	userMap := make(map[string]*api_dto.RatingUserResponse)
	for _, user := range usersGrpcResult.GetUsers() {
		userMap[user.Id] = &api_dto.RatingUserResponse{
			UserId:   user.Id,
			Username: user.Username,
			Image: api_dto.RatingImageResponse{
				ImageId:  imageMap[user.Id].ImageId,
				PublicId: imageMap[user.Id].PublicId,
				UserId:   user.Id,
			},
		}
	}

	ratingResult := make([]api_dto.RatingResponse, 0, len(ratingGrpcResult.Ratings))
	for _, rating := range ratingGrpcResult.GetRatings() {

		user := userMap[rating.GetUserId()] // Get User direct from map

		ratingResult = append(ratingResult, api_dto.RatingResponse{
			Id:      rating.Id,
			Rating:  int(rating.Score),
			Comment: rating.GetComment(),
			HotelId: rating.GetHotelId(),
			User: api_dto.RatingUserResponse{
				UserId:   user.UserId,
				Username: user.Username,
				Image:    user.Image,
			},
		})
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(ratingResult, "Lay binh luan thanh cong"))
}

func (hh *HotelHandler) GetRoomTypesByHotelId(ctx *gin.Context) {
	hotelId := ctx.Param("id")

	roomTypesGrpcResult, err := hh.roomTypeClient.GetRoomTypesByHotelId(ctx, &room_type_pb.GetRoomTypesByHotelIdRequest{
		HotelId: hotelId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach loai phong"))
		return
	}

	roomTypeIds := make([]string, 0, len(roomTypesGrpcResult.GetRoomTypes()))
	for _, roomType := range roomTypesGrpcResult.GetRoomTypes() {
		roomTypeIds = append(roomTypeIds, roomType.Id)
	}

	imagesGrpcResult, err := hh.imageClient.GetImagesByRoomTypeIds(ctx, &image_pb.GetImagesByRoomTypeIdsRequest{
		RoomTypeIds: roomTypeIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc hinh anh loai phong"))
		return
	}
	imageMap := make(map[string]*api_dto.RoomTypeImage)
	for _, image := range imagesGrpcResult.GetImages() {
		imageMap[image.RoomTypeId] = &api_dto.RoomTypeImage{
			Id:         image.Id,
			PublicId:   image.PublicId,
			Format:     image.Format,
			RoomTypeId: image.RoomTypeId,
		}
	}

	roomTypesResponse := make([]*api_dto.RoomTypeResponse, 0, len(roomTypeIds))
	for _, roomType := range roomTypesGrpcResult.GetRoomTypes() {
		roomTypeResponse := &api_dto.RoomTypeResponse{
			Id:    roomType.Id,
			Name:  roomType.Name,
			Price: int(roomType.Price),
		}

		if image, exists := imageMap[roomType.Id]; exists {
			roomTypeResponse.Images = append(roomTypeResponse.Images, *image)
		}

		roomTypesResponse = append(roomTypesResponse, roomTypeResponse)
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomTypesResponse, "Lay danh sach loai phong thanh cong"))
}

func (hh *HotelHandler) GetRoomsByHotelId(ctx *gin.Context) {
	hotelId := ctx.Param("id")

	roomsResult, err := hh.roomClient.GetRoomsByHotelId(ctx, &room_pb.GetRoomsByHotelIdRequest{
		HotelId: hotelId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach phong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomsResult.Rooms, "Lay danh sach phong thanh cong"))
}
