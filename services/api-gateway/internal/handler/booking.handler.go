package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingClient booking_pb.BookingServiceClient
}

func NewBookingHandler(bookingClient booking_pb.BookingServiceClient) *BookingHandler {
	return &BookingHandler{
		bookingClient: bookingClient,
	}
}

func (bh *BookingHandler) RegisterRoutes(router *gin.RouterGroup) {
	bookingHandler := router.Group("/bookings")

	bookingHandler.POST("/", bh.BookingRooms)
}

type BookingRoomsRequest struct {
	RoomTypeBooked []string `json:"room_type_booked"`
	CheckInDate    string   `json:"check_in_date"`
	CheckOutDate   string   `json:"check_out_date"`
	// NumOfGuests int32 `json:"num_of_guests"`
	UserId string `json:"user_id"`
}

func (bh *BookingHandler) BookingRooms(ctx *gin.Context) {
	var bookingReq *BookingRoomsRequest
	if err := ctx.ShouldBindJSON(&bookingReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
		return
	}

	// bookingRes, err := bh.bookingClient.
}
