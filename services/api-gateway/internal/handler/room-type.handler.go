package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomTypeHandler struct {
	roomTypeClient room_type_pb.RoomTypeServiceClient
	roomClient     room_pb.RoomServiceClient
}

func NewRoomTypeHandler(roomTypeClient room_type_pb.RoomTypeServiceClient,
	roomClient room_pb.RoomServiceClient) *RoomTypeHandler {
	return &RoomTypeHandler{
		roomTypeClient: roomTypeClient,
		roomClient:     roomClient,
	}
}

func (rth *RoomTypeHandler) RegisterRoutes(router *gin.RouterGroup) {
	roomTypeHandler := router.Group("/room-types")

	roomTypeHandler.POST("/", rth.CreateRoomType)

	roomTypeHandler.GET("/:id", rth.GetRoomTypeById)
	roomTypeHandler.GET("/:id/rooms", rth.GetRoomsByRoomTypeId)
}

func (rth *RoomTypeHandler) GetRoomTypeById(ctx *gin.Context) {
	id := ctx.Param("id")

	roomTypeGrpcResult, err := rth.roomTypeClient.GetRoomTypeById(ctx, &room_type_pb.GetRoomTypeByIdRequest{
		Id: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, utils.ErrorApiResponse("Loai phong khong ton tai"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc loai phong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomTypeGrpcResult.RoomType, "Lay loai phong thanh cong"))
}

func (rth *RoomTypeHandler) GetRoomsByRoomTypeId(ctx *gin.Context) {
	roomTypeId := ctx.Param("id")

	roomsResult, err := rth.roomClient.GetRoomsByRoomTypeId(ctx, &room_pb.GetRoomsByRoomTypeIdRequest{
		RoomTypeId: roomTypeId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach phong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomsResult.Rooms, "Lay danh sach phong thanh cong"))
}

// Create Room Type
type CreateRoomTypeBody struct {
	RoomTypeName string `json:"room_type_name"`
	Price        int    `json:"price"`
	HotelId      string `json:"hotel_id"`
}

func (rth *RoomTypeHandler) CreateRoomType(ctx *gin.Context) {

	var reqBody CreateRoomTypeBody
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		zap.S().Infoln("Failed to get request body: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi khong tao duoc loai phong"))
		return
	}

	_, err := rth.roomTypeClient.CreateRoomType(ctx, &room_type_pb.CreateRoomTypeRequest{
		Name:    reqBody.RoomTypeName,
		Price:   int32(reqBody.Price),
		HotelId: reqBody.HotelId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong tao duoc danh sach phong"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessApiResponse(nil, "Tao thanh cong"))
}
