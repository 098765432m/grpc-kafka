package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/room_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomGrpcHandler struct {
	roomClient room_pb.RoomServiceClient
}

func NewRoomGrpchandler(roomClient room_pb.RoomServiceClient) *RoomGrpcHandler {
	return &RoomGrpcHandler{
		roomClient: roomClient,
	}
}

func (rh *RoomGrpcHandler) RegisterRoutes(router *gin.RouterGroup) {
	roomHandler := router.Group("/rooms")

	roomHandler.GET("/:id", rh.GetRoomById)
}

func (rh *RoomGrpcHandler) GetRoomById(ctx *gin.Context) {
	id := ctx.Param("id")

	roomGrpcResult, err := rh.roomClient.GetRoomById(ctx, &room_pb.GetRoomByIdRequest{
		Id: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, utils.ErrorApiResponse("Phong khong ton tai"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc phong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(roomGrpcResult.Room, "Lay loai phong thanh cong"))
}
