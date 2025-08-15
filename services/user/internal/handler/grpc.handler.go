package user_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	user_service "github.com/098765432m/grpc-kafka/user/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserGrpcHandler struct {
	user_pb.UnimplementedHotelServiceServer
	service *user_service.UserService
}

func NewUserGrpcHandler(service *user_service.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		service: service,
	}
}

func (ug *UserGrpcHandler) GetUserById(ctx context.Context, req *user_pb.GetUserByIdRequest) (*user_pb.GetUserByIdResponse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.GetId()); err != nil {
		return nil, err
	}

	return &user_pb.GetUserByIdResponse{}, nil
}
