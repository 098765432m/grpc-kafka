package hotel_handler

import (
	"fmt"
	"net/http"

	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type HotelHttpHandler struct {
	service         *hotel_service.HotelService
	roomTypeService *room_type_service.RoomTypeService
	roomService     *room_service.RoomService
}

func NewHotelHttpHandler(service *hotel_service.HotelService, roomTypeService *room_type_service.RoomTypeService, roomService *room_service.RoomService) *HotelHttpHandler {
	return &HotelHttpHandler{
		service:         service,
		roomTypeService: roomTypeService,
		roomService:     roomService,
	}
}

func (hh *HotelHttpHandler) RegisterRoutes(handler *gin.RouterGroup) {
	hotels := handler.Group("/hotels")

	hotels.GET("/", hh.GetAll)
	hotels.POST("/", hh.CreateHotel)

	hotels.GET("/:id", hh.GetHotel)
	hotels.PUT("/:id", hh.UpdateHotel)
	hotels.DELETE("/:id", hh.DeleteHotel)
	hotels.GET("/:id/room-types", hh.GetRoomTypesByHotelId)
	hotels.GET("/:id/rooms", hh.GetRoomsByHotelId)
}

func (hh *HotelHttpHandler) GetHotel(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorApiResponse("Invalid hotel Id format"))
		return
	}

	hotel, err := hh.service.GetHotelById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to get Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(hotel, "Hotel retrieved successfully"))

}

func (hh *HotelHttpHandler) GetAll(ctx *gin.Context) {
	hotels, err := hh.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to get all Hotels"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(hotels, "Hotels retrieved successfully"))
}

type CreateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (hh *HotelHttpHandler) CreateHotel(ctx *gin.Context) {
	hotelReq := &CreateHotelRequest{}

	err := ctx.ShouldBindJSON(hotelReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Invalid request body"))
		return
	}

	err = hh.service.CreateHotel(ctx, &hotel_repo.CreateHotelParams{
		Name:    hotelReq.Name,
		Address: hotelReq.Address,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to create Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Hotel created successfully"))
}

type UpdateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (hh *HotelHttpHandler) UpdateHotel(ctx *gin.Context) {
	hotelReq := &UpdateHotelRequest{}
	if err := ctx.ShouldBindJSON(hotelReq); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorApiResponse("Failed to Bind Request"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorApiResponse("Failed to convert UUID"))
		return
	}

	if err := hh.service.UpdateHotelById(ctx, &hotel_repo.UpdateHotelByIdParams{
		ID:      id,
		Name:    hotelReq.Name,
		Address: hotelReq.Address,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(fmt.Sprintf("Failed to update Hotel: %v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Updated Hotel successfully"))
}

func (hh *HotelHttpHandler) DeleteHotel(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorApiResponse("Failed to convert UUID"))
		return
	}

	if err := hh.service.DeleteHotelById(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Failed to delete Hotel"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Deleted Hotel successfully"))
}

func (hh *HotelHttpHandler) GetRoomTypesByHotelId(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Failed to convert Hotel UUID format"))
		return
	}

	roomTypes, err := hh.roomTypeService.GetRoomTypesByHotelId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to Room Types by Hotel Id")
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomTypes, "Retrieved Room Types successfully"))
}

func (hh *HotelHttpHandler) GetRoomsByHotelId(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Failed to convert Hotel UUID format"))
		return
	}

	rooms, err := hh.roomService.GetRoomsByHotelId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to Rooms by Hotel Id")
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(rooms, "Retrieved Rooms successfully"))
}
