package api_handler

import (
	"net/http"
	"strconv"

	api_dto "github.com/098765432m/grpc-kafka/api-gateway/internal/dto"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
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
	bookingClient  booking_pb.BookingServiceClient
}

type HotelHandlerImpl struct {
	HotelClient    hotel_pb.HotelServiceClient
	RoomTypeClient room_type_pb.RoomTypeServiceClient
	RoomClient     room_pb.RoomServiceClient
	UserClient     user_pb.UserServiceClient
	ImageClient    image_pb.ImageServiceClient
	RatingClient   rating_pb.RatingServiceClient
	BookingClient  booking_pb.BookingServiceClient
}

func NewHotelHandler(hotelHandlerImpl *HotelHandlerImpl) *HotelHandler {
	return &HotelHandler{
		hotelClient:    hotelHandlerImpl.HotelClient,
		roomTypeClient: hotelHandlerImpl.RoomTypeClient,
		roomClient:     hotelHandlerImpl.RoomClient,
		userClient:     hotelHandlerImpl.UserClient,
		imageClient:    hotelHandlerImpl.ImageClient,
		ratingClient:   hotelHandlerImpl.RatingClient,
		bookingClient:  hotelHandlerImpl.BookingClient,
	}
}

func (hh *HotelHandler) RegisterRoutes(router *gin.RouterGroup) {
	hotelHandler := router.Group("/hotels")

	hotelHandler.GET("/", hh.GetAll)

	hotelHandler.GET("/:id", hh.GetHotelById)
	hotelHandler.GET("/:id/room-types", hh.GetRoomTypesByHotelId)
	hotelHandler.GET("/:id/rooms", hh.GetRoomsByHotelId)
	hotelHandler.GET("/:id/ratings", hh.GetRatingsByHotelId)
	hotelHandler.GET("/:id/available-room-types", hh.GetAvailableRoomTypes)

	hotelHandler.GET("/filter", hh.FilterHotels)
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
					Url:      img.GetUrl(),
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

		zap.S().Errorln("Loi khong tim duoc khach san", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong khi lay khach san"))
		return
	}

	images, err := hh.imageClient.GetImagesByHotelId(ctx, &image_pb.GetImagesByHotelIdRequest{HotelId: hotel.Id})
	if err != nil {
		zap.S().Info("Loi ko lay hinh anh khi tim khach san theo ID ", err)

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
			Url:      img.GetUrl(),
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
			Url:        image.GetUrl(),
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
			Price: uint(roomType.Price),
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

func (hh *HotelHandler) GetAvailableRoomTypes(ctx *gin.Context) {
	hotelId := ctx.Param("id")

	checkIn := ctx.Query("check_in")
	checkOut := ctx.Query("check_out")

	zap.L().Info("Get Dates: ", zap.Any("check_in", checkIn), zap.Any("check_out", checkOut))

	roomTypesGrpcResult, err := hh.roomTypeClient.GetRoomTypesByHotelId(ctx, &room_type_pb.GetRoomTypesByHotelIdRequest{
		HotelId: hotelId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach loai phong bang khach san"))
		return
	}

	roomTypeIds := make([]string, 0, len(roomTypesGrpcResult.GetRoomTypes()))
	for _, roomType := range roomTypesGrpcResult.RoomTypes {
		roomTypeIds = append(roomTypeIds, roomType.Id)
	}

	imagesGrpcResult, err := hh.imageClient.GetImagesByRoomTypeIds(ctx, &image_pb.GetImagesByRoomTypeIdsRequest{RoomTypeIds: roomTypeIds})
	if err != nil {
		zap.S().Info("Failed to get Images by Room Type ids: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach hinh anh bang loai phong"))
		return
	}

	imageMap := make(map[string][]api_dto.RoomTypeImage)
	for _, image := range imagesGrpcResult.GetImages() {
		imageMap[image.RoomTypeId] = append(imageMap[image.RoomTypeId], api_dto.RoomTypeImage{
			Id:         image.Id,
			Url:        image.GetUrl(),
			PublicId:   image.PublicId,
			Format:     image.Format,
			RoomTypeId: image.RoomTypeId,
		})
	}

	NumberOfOccupiedRoomsResult, err := hh.bookingClient.GetNumberOfOccupiedRooms(ctx, &booking_pb.GetNumberOfOccupiedRoomsRequest{
		RoomTypeIds: roomTypeIds,
		CheckIn:     checkIn,
		CheckOut:    checkOut,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc so phong da duoc book"))
		return
	}

	numberOccupiedMap := make(map[string]uint)
	for _, result := range NumberOfOccupiedRoomsResult.Results {
		numberOccupiedMap[result.RoomTypeId] = uint(result.NumberOfOccupiedRooms)
	}

	roomTypes := make([]*api_dto.GetNumberOfAvailableRoomsDtoResponse, 0, len(roomTypeIds))
	for _, roomType := range roomTypesGrpcResult.GetRoomTypes() {

		tempRoomType := &api_dto.GetNumberOfAvailableRoomsDtoResponse{
			Id:                     roomType.Id,
			Name:                   roomType.Name,
			Price:                  uint(roomType.Price),
			HotelId:                roomType.HotelId,
			Images:                 imageMap[roomType.Id],
			NumberOfAvailableRooms: uint8(roomType.NumberOfRooms) - uint8(numberOccupiedMap[roomType.Id]),
		}

		roomTypes = append(roomTypes, tempRoomType)
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomTypes, "Thanh cong"))
}

/*
Goal: Tra ve Hotel co phong AVAILABLE trong khung thoi gian CheckIn va CheckOut
Address, va Ten Hotel, trong khung Price

# Tra ve thong tin khach san va gia phong dau tien

1. Tim kiem cac khach san trong khu vuc
-> Return hotelIds (LIMIT 20) /^

2. Tra ve so luong phong per roomType bang HotelIds
-> Return [roomTypeId, So luong phong] /^

3. Su dung cac hotel ids do tra ve cac phong da duoc booking trong khoang thoi gian check in check out
-> Return [roomTypeId, count (so luong phong da dat)] /^

4. O Application, xem neu <so phong da dat> < <so phong> -> AVAILABLE

5. Su dung roomTypeIds de xac dinh khach san con phong nao con trong
-> Return [hotelId, minPrice] (tra ve gia min cua nhung phong trong)
*/

func (hh *HotelHandler) FilterHotels(ctx *gin.Context) {
	// Get Request Params
	hotelName := ctx.Query("hotel_name")
	address := ctx.Query("address")
	checkIn := ctx.Query("check_in")
	checkOut := ctx.Query("check_out")
	minPrice := ctx.Query("min_price")
	maxPrice := ctx.Query("max_price")
	minPriceInt, err := strconv.Atoi(ctx.Query("min_price"))
	if err != nil {
		if minPrice == "" {
			minPriceInt = -1
		} else {
			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Gia khong hop le"))
			return
		}
	}

	maxPriceInt, err := strconv.Atoi(ctx.Query("max_price"))
	if err != nil {
		if maxPrice == "" {
			maxPriceInt = -1
		} else {
			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Gia khong hop le"))
			return
		}
	}

	zap.L().Info("Check Request of Filter Hotels", zap.Any("hotelName", hotelName), zap.Any("check In", checkIn), zap.Any("Check Out", checkOut), zap.Any("Min Price", minPrice), zap.Any("Max Price", maxPrice))

	// Get hotels by address
	resultHotelsByAddress, err := hh.hotelClient.GetHotelsByAddress(ctx, &hotel_pb.GetHotelsByAddressRequest{
		HotelName: hotelName,
		Address:   address,
	})
	if err != nil {
		zap.S().Infoln("Failed to get Hotels by Address: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	zap.L().Info("Filter Hotel By Address ", zap.Any("hotel", resultHotelsByAddress.HotelIds))
	// GetNumber Of rooms to determine which rooms will be AVAILABLE
	resultNumberOfRooms, err := hh.roomClient.GetNumberOfRoomsPerRoomTypeByHotelIds(ctx, &room_pb.GetNumberOfRoomsPerRoomTypeByHotelIdsRequest{
		HotelIds: resultHotelsByAddress.GetHotelIds(),
	})
	if err != nil {
		zap.S().Infoln("Failed to get Number of Rooms Per RoomType By HotelIds: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	allRoomTypeIds := make([]string, 0, len(resultNumberOfRooms.Results))

	// numberOfRoomsInRoomTypesMap := make(map[string]int)
	for _, result := range resultNumberOfRooms.Results {
		allRoomTypeIds = append(allRoomTypeIds, result.GetRoomTypeId())
		// numberOfRoomsInRoomTypesMap[result.RoomTypeId] = int(result.NumberOfRooms)
	}

	zap.S().Infoln("All Room TypeIds")
	zap.S().Infoln(allRoomTypeIds)

	// Get Number of rooms Occupied in booked time
	resultNumberOfOccupiedRooms, err := hh.bookingClient.GetNumberOfOccupiedRooms(ctx, &booking_pb.GetNumberOfOccupiedRoomsRequest{
		RoomTypeIds: allRoomTypeIds,
		CheckIn:     checkIn,
		CheckOut:    checkOut,
	})
	if err != nil {
		zap.S().Infoln("Failed to get Number of Occupied Rooms: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	numberOfOccupiedRoomsInRoomTypesMap := make(map[string]int)
	for _, result := range resultNumberOfOccupiedRooms.Results {
		numberOfOccupiedRoomsInRoomTypesMap[result.RoomTypeId] = int(result.NumberOfOccupiedRooms)
	}

	zap.L().Info("Nunber of Occupied Rooms", zap.Any("result", resultNumberOfOccupiedRooms))

	// Determine which roomtype is available
	// If occupied < total_rooms is AVAILABLE
	availableRoomTypeIds := make([]string, 0, len(allRoomTypeIds))

	for _, result := range resultNumberOfRooms.Results {
		if int32(numberOfOccupiedRoomsInRoomTypesMap[result.RoomTypeId]) < result.NumberOfRooms {
			availableRoomTypeIds = append(availableRoomTypeIds, result.RoomTypeId)
		}
	}

	zap.S().Infoln("Available Room Type Ids")
	zap.S().Infoln(availableRoomTypeIds)

	zap.L().Info("Check Before filter: ", zap.Any("RoomType IDs", availableRoomTypeIds), zap.Any("Min Price", minPrice), zap.Any("Max Price", maxPrice))

	// Get Hotel and Min Price through AVAILABLE Room Type
	hotelRows, err := hh.hotelClient.FilterHotels(ctx, &hotel_pb.FilterHotelsRequest{
		RoomTypeIds: availableRoomTypeIds,
		MinPrice:    int32(minPriceInt),
		MaxPrice:    int32(maxPriceInt),
	})
	if err != nil {
		zap.S().Infoln("Failed to FilterHotels: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	response := hotelRows.FilterHotelRows

	zap.S().Infoln(response)

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(response, "Thanh cong"))
}
