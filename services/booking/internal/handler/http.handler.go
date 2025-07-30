package booking_handler

import (
	"fmt"
	"net/http"

	booking_repo "github.com/098765432m/grpc-kafka/booking/internal/repository/booking"
	booking_service "github.com/098765432m/grpc-kafka/booking/internal/service"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookingHttpHandler struct {
	service *booking_service.BookingService
}

func NewBookingHttpHandler(service *booking_service.BookingService) *BookingHttpHandler {
	return &BookingHttpHandler{
		service: service,
	}
}

func (bh *BookingHttpHandler) RegisterRoutes(handler *gin.RouterGroup) {
	bookings := handler.Group("/bookings")

	bookings.GET("/:id", bh.GetBooking)
	bookings.POST("/", bh.BookingRoom)
	bookings.DELETE("/:id", bh.DeleteBooking)
}

type BookingRoomRequest struct {
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	Total    int    `json:"total"`
	Status   string `json:"status"`
}

func (bh *BookingHttpHandler) GetBooking(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(id); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Booking UUID"))
		return
	}

	booking, err := bh.service.GetBookingById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Sprintf("Cannot get Booking: %v", err)))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(booking, "Get Booking successfully"))
}

func (bh *BookingHttpHandler) BookingRoom(ctx *gin.Context) {

	bookingReq := &BookingRoomRequest{}
	if err := ctx.ShouldBindJSON(bookingReq); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid request body"))
		return
	}

	checkIn, err := utils.ParsePgDate(bookingReq.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid date format"))
		return
	}

	checkOut, err := utils.ParsePgDate(bookingReq.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid date format"))
		return
	}

	var bookingStatus booking_repo.BookingStatus
	if err := bookingStatus.Scan(bookingReq.Status); err != nil {
		ctx.JSON(http.StatusBadGateway, utils.ErrorResponse("Invalid booking status"))
		return
	}

	err = bh.service.CreateBooking(ctx, &booking_repo.CreateBookingParams{
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Total:    int32(bookingReq.Total),
		Status:   bookingStatus,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Cannot create Booking"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(nil, "Booking Room successfully"))
}

func (bh *BookingHttpHandler) DeleteBooking(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(id); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid Booking UUID"))
		return
	}

	err := bh.service.DeleteBooking(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Cannot delete Booking"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Deleted Booking successfully"))
}
