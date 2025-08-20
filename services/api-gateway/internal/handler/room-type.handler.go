package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/room_type_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomTypeGrpcHandler struct {
	roomTypeClient room_type_pb.RoomTypeServiceClient
	roomClient     room_pb.RoomServiceClient
}

func NewRoomTypeGrpchandler(roomTypeClient room_type_pb.RoomTypeServiceClient,
	roomClient room_pb.RoomServiceClient) *RoomTypeGrpcHandler {
	return &RoomTypeGrpcHandler{
		roomTypeClient: roomTypeClient,
		roomClient:     roomClient,
	}
}

func (rth *RoomTypeGrpcHandler) RegisterRoutes(router *gin.RouterGroup) {
	roomTypeHandler := router.Group("/room-types")

	roomTypeHandler.GET("/:id", rth.GetRoomTypeById)
	roomTypeHandler.GET("/:id/rooms", rth.GetRoomsByRoomTypeId)
}

func (rth *RoomTypeGrpcHandler) GetRoomTypeById(ctx *gin.Context) {
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

func (rth *RoomTypeGrpcHandler) GetRoomsByRoomTypeId(ctx *gin.Context) {
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
