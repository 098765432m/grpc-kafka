package hotel_handler

import (
	"fmt"
	"net/http"

	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
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
	hotels.PUT("/:id", hh.UpdateHotel)
	hotels.DELETE("/:id", hh.DeleteHotel)
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
	Name    string `json:"name"`
	Address string `json:"address"`
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

type UpdateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (hh *HotelHttpHandler) UpdateHotel(ctx *gin.Context) {
	hotelReq := &UpdateHotelRequest{}
	if err := ctx.ShouldBindJSON(hotelReq); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to Bind Request"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to convert UUID"))
		return
	}

	if err := hh.service.UpdateHotel(ctx, &hotel_repo.UpdateHotelByIdParams{
		ID:      id,
		Name:    hotelReq.Name,
		Address: hotelReq.Address,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Sprintf("Failed to update Hotel: %v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Updated Hotel successfully"))
}

func (hh *HotelHttpHandler) DeleteHotel(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to convert UUID"))
		return
	}

	if err := hh.service.DeleteHotel(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Deleted Hotel successfully"))
}
