package user_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/repository/user"
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

func (us *UserService) GetUsers(ctx context.Context) ([]user_repo.User, error) {

	users, err := us.repo.GetUsers(ctx)
	if err != nil {
		zap.S().Errorln("Failed to get Users: ", err)
		return nil, err
	}

	return users, nil
}

func (us *UserService) GetUserById(ctx context.Context, id pgtype.UUID) (*user_repo.User, error) {
	user, err := us.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			zap.S().Infoln("User not found")
			return nil, common_error.ErrNoRows
		}

		zap.S().Error("Failed to get Hotel by id: ", err)
		return nil, err
	}

	return &user, nil
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
				zap.S().Errorln("Duplicated User", err)
				return common_error.ErrDuplicateRecord
			}
		}

		zap.S().Error("Failed to create User :", err)
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

func (us *UserService) SignIn(ctx context.Context, username string, password string) (string, error) {

	// Check User by username
	checkUser, err := us.repo.CheckUserByUsername(ctx, username)
	if err != nil {
		zap.S().Errorln("Failed to get User by username: ", err)
		return "", common_error.ErrNoRows
	}

	// Check Password and Hashed password
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	if err != nil {
		zap.S().Errorln("Wrong password!: ", err)
		return "", common_error.ErrNoRows
	}

	// Get secret key from env
	secretKey := viper.GetString("JWT_SECRET_KEY")
	if secretKey == "" {
		zap.S().Errorln("Failed to get JWT SECRET KEY from env: ", err)
		return "", err
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
		return "", err
	}

	return signedStr, nil
}
