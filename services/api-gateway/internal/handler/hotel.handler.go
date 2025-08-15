package api_handler

import (
	"errors"
	"net/http"

	api_dto "github.com/098765432m/grpc-kafka/api-gateway/internal/dto"
	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type HotelHandler struct {
	hotelClient hotel_pb.HotelServiceClient
	imageClient image_pb.ImageServiceClient
}

func NewHotelHandler(
	hotelClient hotel_pb.HotelServiceClient,
	imageClient image_pb.ImageServiceClient,
) *HotelHandler {
	return &HotelHandler{
		hotelClient: hotelClient,
		imageClient: imageClient,
	}
}

func (hh *HotelHandler) RegisterRoutes(router *gin.RouterGroup) {
	hotelHandler := router.Group("/hotels")

	hotelHandler.GET("/", hh.GetAll)
}

func (hh *HotelHandler) GetAll(ctx *gin.Context) {

	hotels, err := hh.hotelClient.GetAllHotels(ctx, &hotel_pb.GetAllHotelsRequest{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get all Hotels"))
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
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get Images By Hotel ids"))
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
				resp.Images = append(resp.Images, api_dto.Image{
					Id:       img.Id,
					PublicId: img.PublicId,
					Format:   img.Format,
				})
			}
		}

		responses = append(responses, resp)
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(responses, "Hotels retrieved successfully"))
}

func (hh *HotelHandler) GetHotelById(ctx *gin.Context) {

	var id *pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Id Khach san khong hop le"))
		return
	}
	hotelGrpc, err := hh.hotelClient.GetHotelById(ctx, &hotel_pb.GetHotelByIdRequest{
		Id: id.String(),
	})
	hotel := hotelGrpc.GetHotel()
	if err != nil {

		switch {
		case errors.Is(err, common_error.ErrNoRows):
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Khach san khong ton tai"))
			return

		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Khong the "))
			return
		}
	}

	images, err := hh.imageClient.GetImagesByHotelId(ctx, &image_pb.GetImagesByHotelIdRequest{HotelId: hotel.Id})
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Khong the "))
		return
	}

	resp := api_dto.HotelResponse{
		Id:      hotel.GetId(),
		Name:    hotel.GetName(),
		Address: hotel.GetAddress(),
	}

	for _, img := range images.GetImages() {
		resp.Images = append(resp.Images, api_dto.Image{
			Id:       img.GetId(),
			PublicId: img.GetPublicId(),
			Format:   img.GetFormat(),
		})
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(resp, "Hotel retrieved successfully"))
}
