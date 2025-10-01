package user_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	user_domain "github.com/098765432m/grpc-kafka/user/internal/domain"
	user_repo_mapping "github.com/098765432m/grpc-kafka/user/internal/infrastructure/repository/sqlc"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/infrastructure/repository/sqlc/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *user_repo.Queries
}

func NewUserService(repo *user_repo.Queries) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetUsers(ctx context.Context) ([]user_domain.User, error) {

	users, err := us.repo.GetUsers(ctx)
	if err != nil {
		zap.S().Errorln("Failed to get Users: ", err)
		return nil, err
	}

	result := user_repo_mapping.FromUsersRepoToUsersDomain(users)

	return result, nil
}

func (us *UserService) GetUserById(ctx context.Context, id pgtype.UUID) (*user_domain.User, error) {
	user, err := us.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Infoln("User not found")
			return nil, common_error.ErrNoRows
		}

		zap.S().Error("Failed to get User by id: ", err)
		return nil, err
	}

	result := user_repo_mapping.FromUserRepoToUserDomain(user)

	return &result, nil
}

func (us *UserService) GetUsersByIds(ctx context.Context, ids []pgtype.UUID) ([]user_domain.User, error) {

	users, err := us.repo.GetUsersByIds(ctx, ids)
	if err != nil {
		zap.S().Errorln("Failed to Get Users By Ids: ", err)
		return nil, err
	}

	return user_repo_mapping.FromUsersRepoToUsersDomain(users), nil
}

func (us *UserService) CreateUser(ctx context.Context, newUser *user_repo.CreateUserParams) error {

	// Hashed new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Errorln("Failed to hash password: ", err)
		return err
	}

	// Create User
	err = us.repo.CreateUser(ctx, user_repo.CreateUserParams{
		Username:    newUser.Username,
		Password:    string(hashedPassword),
		Address:     newUser.Address,
		Email:       newUser.Email,
		FullName:    newUser.FullName,
		PhoneNumber: newUser.PhoneNumber,
		Role:        newUser.Role,
		HotelID:     newUser.HotelID,
	})

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				zap.S().Info("Duplicated User: ", err)
				return common_error.ErrDuplicateRecord
			}
		}

		zap.S().Error("Failed to create User :", err)
		return err
	}

	return nil
}

func (us *UserService) UpdateUserById(ctx context.Context, params *user_repo.UpdateUserByIdParams) error {
	isUserExisted, err := us.repo.CheckUserExistsById(ctx, params.ID)
	if err != nil {
		zap.S().Error("Cannot check user exist by id")
		return err
	}

	if !isUserExisted {
		zap.S().Errorln("User is not existed to update")
		return common_error.ErrNoRows
	}
	err = us.repo.UpdateUserById(ctx, *params)
	if err != nil {

		zap.S().Errorln("Failed to update User by id: ", err)
		return err
	}

	return nil
}

func (us *UserService) DeleteUserById(ctx context.Context, id pgtype.UUID) error {

	const DEFAULT_ERR_MSG = "loi khong the tao tai khoan"
	isUserExisted, err := us.repo.CheckUserExistsById(ctx, id)
	if err != nil {
		zap.S().Errorln("Failed to check User exists by Id: ", err)
		return fmt.Errorf("%s", DEFAULT_ERR_MSG)
	}

	if !isUserExisted {
		return common_error.ErrNoRows
	}

	if err := us.repo.DeleteUserById(ctx, id); err != nil {
		zap.S().Errorln("Failed to delete User by id: ", err)
		return fmt.Errorf("%s", DEFAULT_ERR_MSG)
	}

	return nil
}

type SignInResult struct {
	Jwt      string
	Id       string
	Username string
	Email    string
	Role     string
}

func (us *UserService) SignIn(ctx context.Context, username string, password string) (*SignInResult, error) {

	// Check User by username
	checkUser, err := us.repo.CheckUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Info("user with this username is not existed: ", err)
			return nil, common_error.ErrNoRows
		}
		return nil, err
	}

	// Check Password and Hashed password
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	if err != nil {
		zap.S().Errorln("Wrong password!: ", err)
		return nil, common_error.ErrNoRows
	}

	// Get secret key from env
	secretKey := viper.GetString("JWT_SECRET_KEY")
	if secretKey == "" {
		zap.S().Errorln("Failed to get JWT SECRET KEY from env: ", err)
		return nil, err
	}

	exp := time.Now().Add(time.Hour * 24).Unix()

	// create jwt
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "http://localhost:3102",
		"sub":  checkUser.ID.String(),
		"role": checkUser.Role,
		"exp":  exp,
	})

	signedStr, err := t.SignedString([]byte(secretKey))
	if err != nil {
		zap.S().Errorln("Failed to return JWT: ", err)
		return nil, err
	}

	return &SignInResult{
		Jwt:      signedStr,
		Id:       checkUser.ID.String(),
		Username: checkUser.Username,
		Email:    checkUser.Email,
		Role:     string(checkUser.Role),
	}, nil
}
