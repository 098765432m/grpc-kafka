package api_handler

import (
	"fmt"
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
	bookingHandler.DELETE("/", bh.DeleteBookingsById)
	bookingHandler.POST("/delete", bh.DeleteBookingsByIds)
}

type BookedRooms struct {
	RoomTypeBookedId string `json:"room_type_booked"`
	NumberOfRooms    int    `json:"number_of_rooms"`
}

type BookingRoomsRequest struct {
	BookedRooms  []BookedRooms `json:"booked_rooms"`
	CheckInDate  string        `json:"check_in_date" binding:"required"`
	CheckOutDate string        `json:"check_out_date" binding:"required"`
	Total        int           `json:"total" binding:"required"`
	// NumOfGuests int32 `json:"num_of_guests"`
	UserId string `json:"user_id" binding:"required"`
}

// Param struct for Create Bookings
type BookedRoomsReq struct {
	RoomTypeId string
	RoomIds    []string
}

func (bh *BookingHandler) BookingRooms(ctx *gin.Context) {
	var bookingReq *BookingRoomsRequest
	if err := ctx.ShouldBindJSON(&bookingReq); err != nil {
		zap.S().Info("Invalid Booking request: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
		return
	}

	// Check booked room != 0
	if len(bookingReq.BookedRooms) == 0 {
		zap.S().Info("Booked rooms are not found")
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
		return
	}
	// Save room Type id and room ids that assign to booking
	bookedRoomsReq := make([]BookedRoomsReq, 0, len(bookingReq.BookedRooms))

	// Check for available Rooms for Room Type
	for _, bookedRoom := range bookingReq.BookedRooms {

		// Get already booked Rooms in range of time
		roomsGrpcResult, err := bh.bookingClient.GetUnavailableRoomsByRoomTypeId(ctx, &booking_pb.GetUnavailableRoomsByRoomTypeIdRequest{
			RoomTypeId: bookedRoom.RoomTypeBookedId,
			CheckIn:    bookingReq.CheckInDate,
			CheckOut:   bookingReq.CheckOutDate,
		})
		if err != nil {
			zap.S().Infoln("Failed to get UNAVAILABLE Rooms: ", err)
			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
			return
		}

		unavailableRooms := roomsGrpcResult.GetRoomIds()

		roomIds, err := bh.roomClient.GetListOfRemainRooms(ctx, &room_pb.GetListOfRemainRoomsRequest{
			RoomTypeId:    bookedRoom.RoomTypeBookedId,
			BookedRoomIds: unavailableRooms,
			NumberOfRooms: int32(bookedRoom.NumberOfRooms),
		})
		if err != nil {
			zap.S().Infoln("Failed to get AVAILABLE Rooms: ", err)
			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong con phong trong"))
			return
		}

		bookedRoomsReq = append(bookedRoomsReq, BookedRoomsReq{
			RoomTypeId: bookedRoom.RoomTypeBookedId,
			RoomIds:    roomIds.GetRoomIds(),
		})
	}

	// Create booking param for each room
	var newBookings []*booking_pb.BookingParam
	for _, bookedRoom := range bookedRoomsReq {
		for i, roomId := range bookedRoom.RoomIds {
			newBooking := &booking_pb.BookingParam{
				CheckIn:    bookingReq.CheckInDate,
				CheckOut:   bookingReq.CheckOutDate,
				Total:      int32(bookingReq.Total),
				RoomTypeId: bookedRoom.RoomTypeId,
				RoomId:     roomId,
				UserId:     bookingReq.UserId,
			}

			zap.L().Info("Check booked room", zap.Any(fmt.Sprintf("Booked room %d", i), newBooking))

			newBookings = append(newBookings, newBooking)

		}
	}

	// Create Bookings with many rooms
	_, err := bh.bookingClient.CreateBookings(ctx, &booking_pb.NewBookingParam{
		NewBookings: newBookings,
	})
	if err != nil {
		zap.S().Info("Cannot booking rooms: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong dat duoc phong"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessApiResponse(nil, "Dat phong thanh cong"))
}

// TODO: Worker Pool to utilize waitGroup for performance
// Worker Pool to process AVAILABLE Rooms
// func (bh *BookingHandler) processBookedRooms(ctx context.Context, bookingReq *BookingRoomsRequest) ([]BookedRoomsReq, error) {
// 	// Save room Type id and room ids that assign to booking
// 	bookedRoomsReq := make([]BookedRoomsReq, 0, len(bookingReq.BookedRooms))

// 	// Check for available Rooms for Room Type
// 	for _, bookedRoom := range bookingReq.BookedRooms {

// 		// Get already booked Rooms in range of time
// 		roomsGrpcResult, err := bh.bookingClient.GetUnavailableRoomsByRoomTypeId(ctx, &booking_pb.GetUnavailableRoomsByRoomTypeIdRequest{
// 			RoomTypeId: bookedRoom.RoomTypeBookedId,
// 			CheckIn:    bookingReq.CheckInDate,
// 			CheckOut:   bookingReq.CheckOutDate,
// 		})
// 		if err != nil {
// 			zap.S().Infoln("Failed to get UNAVAILABLE Rooms: ", err)
// 			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong dat duoc phong"))
// 			return
// 		}

// 		unavailableRooms := roomsGrpcResult.GetRoomIds()

// 		roomIds, err := bh.roomClient.GetListOfRemainRooms(ctx, &room_pb.GetListOfRemainRoomsRequest{
// 			RoomTypeId:    bookedRoom.RoomTypeBookedId,
// 			BookedRoomIds: unavailableRooms,
// 			NumberOfRooms: int32(bookedRoom.NumberOfRooms),
// 		})
// 		if err != nil {
// 			zap.S().Infoln("Failed to get AVAILABLE Rooms: ", err)
// 			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong con phong trong"))
// 			return
// 		}

// 		bookedRoomsReq = append(bookedRoomsReq, BookedRoomsReq{
// 			RoomTypeId: bookedRoom.RoomTypeBookedId,
// 			RoomIds:    roomIds.GetRoomIds(),
// 		})
// 	}
// }

func (bh *BookingHandler) DeleteBookingsById(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		zap.S().Infoln("Invalid UUID: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi he thong"))
	}

	_, err := bh.bookingClient.DeleteBookingsById(ctx, &booking_pb.DeleteBookingByIdRequest{
		BookingId: id.String(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong xoa duoc booking"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Xoa thanh cong"))
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

	_, err := bh.bookingClient.DeleteBookingsByIds(ctx, &booking_pb.DeleteBookingByIdsRequest{
		BookingIds: req.BookingIds,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Khong xoa duoc booking"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Xoa bookings thanh cong"))
}
