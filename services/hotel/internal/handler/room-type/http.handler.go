package room_type_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/utils"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomTypeHttpHandler struct {
	service *room_type_service.RoomTypeService
}

func NewHotelHttpHandler(service *room_type_service.RoomTypeService) *RoomTypeHttpHandler {
	return &RoomTypeHttpHandler{
		service: service,
	}
}

func (rth *RoomTypeHttpHandler) RegisterRoutes(handler *gin.RouterGroup) {
	roomTypes := handler.Group("/room-types")

	roomTypes.POST("/", rth.CreateRoomType)
	roomTypes.GET("/:id", rth.GetRoomType)
	roomTypes.DELETE("/:id", rth.DeleteRoomType)
}

func (rth *RoomTypeHttpHandler) GetRoomType(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid Room Type Id format"))
		return
	}

	roomType, err := rth.service.GetRoomTypeById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get Room Type"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(roomType, "Room Type retrieved successfully"))

}

type CreateRoomTypeRequest struct {
	Name    string `json:"name"`
	Price   int    `json:"price"`
	HotelId string `json:"hotel_id,omitempty"`
}

func (rth *RoomTypeHttpHandler) CreateRoomType(ctx *gin.Context) {
	roomTypeReq := &CreateRoomTypeRequest{}

	var hotelId pgtype.UUID
	if err := hotelId.Scan(hotelId); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid hotel UUID format"))
	}

	err := ctx.ShouldBindJSON(roomTypeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
		return
	}

	err = rth.service.CreateRoomType(ctx, &room_type_repo.CreateRoomTypeParams{
		Name:    roomTypeReq.Name,
		Price:   int32(roomTypeReq.Price),
		HotelID: hotelId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create Room Type"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Room Type created successfully"))
}

func (rth *RoomTypeHttpHandler) DeleteRoomType(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Failed to convert UUID"))
		return
	}

	if err := rth.service.DeleteRoomTypeById(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete Room Type"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Deleted Room Type successfully"))
}
