package user_handler

import (
	"errors"
	"net/http"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/utils"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/repository/user"
	user_service "github.com/098765432m/grpc-kafka/user/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UserHttpHandler struct {
	service *user_service.UserService
}

func NewUserHttpHandler(service *user_service.UserService) *UserHttpHandler {
	return &UserHttpHandler{
		service: service,
	}
}

func (uh *UserHttpHandler) RegisterRoutes(handler *gin.RouterGroup) {
	users := handler.Group("/users")

	users.GET("/", uh.GetUsers)
	users.POST("/", uh.CreateUser)

	users.GET("/:id", uh.GetUserById)
	users.DELETE("/:id", uh.DeleteUserById)

	users.POST("/sign-in", uh.SignIn)
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DTO Response
type UserDtoResponse struct {
	Id          string
	Username    string
	Email       string
	PhoneNumber string
	FullName    string
	Role        string
	HotelId     string
}

func (uh *UserHttpHandler) GetUsers(ctx *gin.Context) {

	users, err := uh.service.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi khong lay duoc danh sach tai khoan"))
		return
	}

	var usersDtoRes []*UserDtoResponse
	for _, user := range users {
		userDtoRes := &UserDtoResponse{
			Id:          user.ID.String(),
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
			Role:        string(user.Role),
			HotelId:     user.HotelID.String(),
		}

		usersDtoRes = append(usersDtoRes, userDtoRes)
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(usersDtoRes, "Thanh cong"))
}

func (uh *UserHttpHandler) GetUserById(ctx *gin.Context) {

	// Get Id Param
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		zap.S().Errorln("Invalid UUID format: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Loi he thong"))
		return
	}

	user, err := uh.service.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, common_error.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tai khoan khong ton tai"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi khong tim duoc tai khoan"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(UserDtoResponse{
		Id:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		Role:        string(user.Role),
		HotelId:     user.HotelID.String(),
	}, "Thanh cong"))
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role,omitempty"`
	HotelId     string `json:"hotel_id,omitempty"`
}

func (uh *UserHttpHandler) CreateUser(ctx *gin.Context) {

	// Get request body
	createUserReq := &CreateUserRequest{}
	if err := ctx.ShouldBindJSON(createUserReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Loi he thong"))
		return
	}

	// Check hotel UUID
	var hotelId pgtype.UUID
	if err := hotelId.Scan(createUserReq.HotelId); err != nil {
		if createUserReq.HotelId == "" {
			hotelId.Valid = false
		} else {
			zap.S().Errorln("Invalid UUID: ", err)
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi he thong"))
			return
		}
	}

	// Check is Role enum valid
	var userRole user_repo.RoleEnum
	if err := userRole.Scan(createUserReq.Role); err != nil {

		zap.S().Errorln("Invalid User role enum")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi he thong"))
		return

	}

	if userRole == "" {
		userRole = user_repo.RoleEnumGUEST
	}

	// Create User
	err := uh.service.CreateUser(ctx, &user_repo.CreateUserParams{
		Username:    createUserReq.Username,
		Password:    createUserReq.Password,
		Address:     createUserReq.Address,
		Email:       createUserReq.Email,
		FullName:    createUserReq.FullName,
		PhoneNumber: createUserReq.PhoneNumber,
		Role:        userRole,
		HotelID:     hotelId,
	})

	if err != nil {

		if errors.Is(err, common_error.ErrDuplicateRecord) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Thong tin bi trung"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Khong the tao tai khoan"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(nil, "Tao tai khoan thanh cong"))
}

func (uh *UserHttpHandler) DeleteUserById(ctx *gin.Context) {

	// Get Id Param
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		zap.S().Errorln("Invalid UUID: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Loi ID tai khoan khong hop le"))
		return
	}

	if err := uh.service.DeleteUserById(ctx, id); err != nil {
		switch {
		// Catch no rows found
		case errors.Is(err, common_error.ErrNoRows):
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tai khoan khong ton tai"))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi khong the xoa tai khoan"))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(nil, "Xoa tai khoan thanh cong"))
}

func (uh *UserHttpHandler) SignIn(ctx *gin.Context) {

	// Get request body
	signInReq := &SignInRequest{}
	err := ctx.ShouldBindJSON(signInReq)
	if err != nil {
		zap.S().Errorln("Failed to get sign in request: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Loi he thong"))
		return
	}

	jwt, err := uh.service.SignIn(ctx, signInReq.Username, signInReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Tai khoan hoac mat khau khong dung"))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Loi khong the dang nhap"))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(jwt, "Dang nhap thanh cong"))
}
