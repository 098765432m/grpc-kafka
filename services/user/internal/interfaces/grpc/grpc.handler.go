package user_handler

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	user_service "github.com/098765432m/grpc-kafka/user/internal/application"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/infrastructure/repository/sqlc/user"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcHandler struct {
	user_pb.UnimplementedUserServiceServer
	service *user_service.UserService
}

func NewUserGrpcHandler(service *user_service.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		service: service,
	}
}

func (ug *UserGrpcHandler) GetUserById(ctx context.Context, req *user_pb.GetUserByIdRequest) (*user_pb.GetUserByIdResponse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Infoln("Invalid user UUID: ", err)
		return nil, status.Errorf(codes.InvalidArgument, "UUID khong hop le")
	}

	user, err := ug.service.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, common_error.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "Tai khoan khong ton tai")
		}

		return nil, status.Errorf(codes.Internal, "Loi he thong")
	}

	return &user_pb.GetUserByIdResponse{
		User: &user_pb.User{
			Id:          user.Id,
			Username:    user.Username,
			Password:    user.Password,
			Address:     user.Address,
			Email:       user.Email,
			FullName:    user.Address,
			PhoneNumber: user.PhoneNumber,
			Role:        string(user.Role),
			HotelId:     user.HotelId,
		},
	}, nil
}

func (ug *UserGrpcHandler) GetUsersByIds(ctx context.Context, req *user_pb.GetUsersByIdsRequest) (*user_pb.GetUsersByIdsResponse, error) {
	ids := make([]pgtype.UUID, 0, len(req.GetIds()))
	for _, userIdReq := range req.Ids {
		var id pgtype.UUID
		// id := &pgtype.UUID{}
		if err := id.Scan(userIdReq); err != nil {
			zap.S().Info("Invalid User UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi user UUID")
		}

		ids = append(ids, id)
	}

	users, err := ug.service.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	usersGrpcResult := make([]*user_pb.User, 0, len(users))
	for _, user := range users {
		userGrpcResult := &user_pb.User{
			Id:          user.Id,
			Username:    user.Username,
			Password:    user.Password,
			Address:     user.Address,
			Email:       user.Username,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
			HotelId:     user.HotelId,
			Role:        string(user.Role),
		}
		usersGrpcResult = append(usersGrpcResult, userGrpcResult)
	}

	return &user_pb.GetUsersByIdsResponse{
		Users: usersGrpcResult,
	}, nil
}

func (ug *UserGrpcHandler) CreateUser(ctx context.Context, req *user_pb.CreateUserRequest) (*user_pb.CreateUserResponse, error) {

	var hotelId pgtype.UUID
	if req.HotelId != "" {
		if err := hotelId.Scan(req.HotelId); err != nil {
			zap.S().Infoln("Invalid User UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi he thong")
		}
	}

	err := ug.service.CreateUser(ctx, &user_repo.CreateUserParams{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
		Role:        user_repo.RoleEnum(req.Role),
		HotelID:     hotelId,
	})
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrDuplicateRecord):
			return nil, status.Error(codes.AlreadyExists, "Tai khoan da ton tai")
		}

		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	return &user_pb.CreateUserResponse{}, nil
}

func (ug *UserGrpcHandler) UpdateUserById(ctx context.Context, req *user_pb.UpdateUserByIdRequest) (*user_pb.UpdateUserByIdResponse, error) {

	var id pgtype.UUID
	if req.GetUser().Id == "" {
		if err := id.Scan(req.GetUser().Id); err != nil {
			zap.S().Infoln("Invalid User UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi he thong")
		}
	}

	var hotelId pgtype.UUID
	if req.GetUser().GetHotelId() == "" {
		if err := hotelId.Scan(req.GetUser().GetHotelId()); err != nil {
			zap.S().Infoln("Invalid Hotel UUID: ", err)
			return nil, status.Error(codes.InvalidArgument, "Loi he thong")
		}
	}

	// Default Role if empty
	var newRole user_repo.RoleEnum
	if req.User.GetRole() == "" {
		newRole = user_repo.RoleEnumGUEST
	}

	err := ug.service.UpdateUserById(ctx, &user_repo.UpdateUserByIdParams{
		ID:          id,
		Username:    req.User.Username,
		Password:    req.User.Password,
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
		FullName:    req.User.FullName,
		Role:        newRole,
		HotelID:     hotelId,
	})
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):

			return nil, status.Error(codes.NotFound, "Tai khoan khong ton tai de cap nhat")
		}

		zap.S().Errorln(err)
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	return &user_pb.UpdateUserByIdResponse{}, nil
}

func (ug *UserGrpcHandler) DeleteUserById(ctx context.Context, req *user_pb.DeleteUserByIdRequest) (*user_pb.DeleteUserByIdResponse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Infoln("Invalid User UUID: ", err)
		return nil, status.Error(codes.InvalidArgument, "Loi he thong")
	}
	err := ug.service.DeleteUserById(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong the xoa tai khoan")
	}

	return &user_pb.DeleteUserByIdResponse{}, nil
}

func (ug *UserGrpcHandler) SignIn(ctx context.Context, req *user_pb.SignInRequest) (*user_pb.SignInResponse, error) {
	signInResult, err := ug.service.SignIn(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, common_error.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "%v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%v", err.Error())
	}

	return &user_pb.SignInResponse{
		Jwt:      signInResult.Jwt,
		Id:       signInResult.Id,
		Username: signInResult.Username,
		Email:    signInResult.Email,
		Role:     signInResult.Role,
	}, nil
}
