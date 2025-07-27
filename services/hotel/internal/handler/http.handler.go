package hotel_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type HotelHttpHandler struct {
	service *hotel_service.HotelService
}

func NewHotelHttpHandler(service *hotel_service.HotelService) *HotelHttpHandler {
	return &HotelHttpHandler{
		service: service,
	}
}

func (hh *HotelHttpHandler) RegisterRoutes(handler *gin.RouterGroup) {
	hotels := handler.Group("/hotels")

	hotels.GET("/:id", hh.GetHotel)
	hotels.GET("/", hh.GetAll)
	hotels.POST("/", hh.CreateHotel)
}

func (hh *HotelHttpHandler) GetHotel(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid hotel Id format"))
		return
	}

	hotel, err := hh.service.GetHotel(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(hotel, "Hotel retrieved successfully"))

}

func (hh *HotelHttpHandler) GetAll(ctx *gin.Context) {
	hotels, err := hh.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get all Hotels"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(hotels, "Hotels retrieved successfully"))
}

type CreateHotelRequest struct {
	Name string `json:"name"`
}

func (hh *HotelHttpHandler) CreateHotel(ctx *gin.Context) {
	hotelReq := &CreateHotelRequest{}

	err := ctx.ShouldBindJSON(hotelReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
		return
	}

	err = hh.service.CreateHotel(ctx, hotelReq.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Hotel created successfully"))
}
