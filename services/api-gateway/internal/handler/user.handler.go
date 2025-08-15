package api_handler

import (
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
)

type UserHandler struct {
	userClient  user_pb.HotelServiceClient
	imageClient image_pb.ImageServiceClient
}

func NewUserHandler(
	userClient user_pb.HotelServiceClient,
	imageClient image_pb.ImageServiceClient,
) *UserHandler {
	return &UserHandler{
		userClient:  userClient,
		imageClient: imageClient,
	}
}

func RegisterRoutes() {

}
