package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookingHandler struct {
	bookingClient booking_pb.BookingServiceClient
	roomClient    room_pb.RoomServiceClient
}

func NewBookingHandler(bookingClient booking_pb.BookingServiceClient,
	roomClient room_pb.RoomServiceClient) *BookingHandler {
	return &BookingHandler{
		bookingClient: bookingClient,
		roomClient:    roomClient,
	}
}

func (bh *BookingHandler) RegisterRoutes(router *gin.RouterGroup) {
	bookingHandler := router.Group("/bookings")

	bookingHandler.POST("/", bh.BookingRooms)
	bookingHandler.POST("/delete", bh.DeleteBookingsByIds)
}

type BookedRooms struct {
	RoomTypeBookedId string `json:"room_type_booked"`
	NumberOfRooms    int    `json:"number_of_rooms"`
}

type BookingRoomsRequest struct {
	BookedRooms  []BookedRooms `json:"booked_rooms"`
	CheckInDate  string        `json:"check_in_date"`
	CheckOutDate string        `json:"check_out_date"`
	Total        int           `json:"total"`
	// NumOfGuests int32 `json:"num_of_guests"`
	UserId string `json:"user_id"`
}

// TODO: Check this function can this use with waitGroup
func (bh *BookingHandler) BookingRooms(ctx *gin.Context) {
	var bookingReq *BookingRoomsRequest
	if err := ctx.ShouldBindJSON(&bookingReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
		return
	}

	type BookedRoomsReq struct {
		RoomTypeId string
		RoomIds    []string
	}

	// Save room Type id and room ids that assign to booking
	bookedRoomsReq := make([]BookedRoomsReq, 0, len(bookingReq.BookedRooms))

	// Check for available Rooms for Room Type
	for _, booking := range bookingReq.BookedRooms {
		listRoomsGrpcResult, err := bh.roomClient.GetListOfAvailableRoomsByRoomTypeId(ctx, &room_pb.GetListOfAvailableRoomsByRoomTypeIdRequest{
			RoomTypeId:    booking.RoomTypeBookedId,
			NumberOfRooms: int32(booking.NumberOfRooms),
		})
		// TODO Bat loi Khong co phong AVAILABLE

		if err != nil {
			zap.S().Info("Cannot get available rooms: ", err)
			ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong dat duoc phong"))
		}

		bookedRoomsReq = append(bookedRoomsReq, BookedRoomsReq{
			RoomTypeId: booking.RoomTypeBookedId,
			RoomIds:    listRoomsGrpcResult.GetRoomIds(),
		})
	}

	// Create booking param for each room
	var newBookings []*booking_pb.BookingParam
	for _, bookedRoom := range bookedRoomsReq {
		for _, roomId := range bookedRoom.RoomIds {
			newBooking := &booking_pb.BookingParam{
				CheckIn:    bookingReq.CheckInDate,
				CheckOut:   bookingReq.CheckOutDate,
				Total:      int32(bookingReq.Total),
				RoomTypeId: bookedRoom.RoomTypeId,
				RoomId:     roomId,
				UserId:     bookingReq.UserId,
			}
			newBookings = append(newBookings, newBooking)

		}

		// TODO: Set status rooms to BOOKED
		_, err := bh.roomClient.SetOccupiedStatusToRooms(ctx, &room_pb.SetOccupiedStatusToRoomsRequest{
			RoomIds: bookedRoom.RoomIds,
		})
		if err != nil {
			zap.S().Info("Cannot set status booked to Rooms: ", err)
			ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong dat duoc phong"))
		}
	}

	// Create Bookings with many rooms
	_, err := bh.bookingClient.CreateBookings(ctx, &booking_pb.NewBookingParam{
		NewBookings: newBookings,
	})
	if err != nil {
		zap.S().Info("Cannot booking rooms: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong dat duoc phong"))
	}

	ctx.JSON(http.StatusCreated, utils.SuccessApiResponse(nil, "Dat phong thanh cong"))
}

type DeleteBookingsRequest struct {
	BookingIds []string `json:"booking_ids"`
}

func (bh *BookingHandler) DeleteBookingsByIds(ctx *gin.Context) {
	var req *DeleteBookingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Request khong hop le"))
		return
	}

	_, err := bh.bookingClient.DeleteBookingsByIds(ctx, &booking_pb.DeleteBookingRequest{
		BookingIds: req.BookingIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong xoa duoc booking"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Xoa bookings thanh cong"))
}
