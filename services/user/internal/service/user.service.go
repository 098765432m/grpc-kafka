package user_service

import (
	"context"

	user_repo "github.com/098765432m/grpc-kafka/user/internal/repository/user"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UserService struct {
	repo *user_repo.Queries
}

func NewUserService(repo *user_repo.Queries) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetUser(ctx context.Context, id pgtype.UUID) (*user_repo.User, error) {
	user, err := us.repo.GetUserById(ctx, id)
	if err != nil {
		zap.S().Error("Failed to get Hotel by id")
		return nil, err
	}

	return &user, nil
}

func (us *UserService) CreateUser(ctx context.Context, newUser *user_repo.CreateUserParams) error {

	err := us.repo.CreateUser(ctx, *newUser)
	if err != nil {
		zap.S().Error("Failed to create User :", err)
		return err
	}

	return nil
}
