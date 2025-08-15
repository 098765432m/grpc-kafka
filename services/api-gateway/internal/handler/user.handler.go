package api_handler

import (
	"errors"
	"net/http"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (uh *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	userHandler := router.Group("/users")

	userHandler.POST("/sign-in", uh.SignIn)
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uh *UserHandler) SignIn(ctx *gin.Context) {
	// Get request body
	signInReq := &SignInRequest{}
	err := ctx.ShouldBindJSON(signInReq)
	if err != nil {
		zap.S().Errorln("Failed to get sign in request: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	signInResult, err := uh.userClient.SignIn(ctx, &user_pb.SignInRequest{
		Username: signInReq.Username,
		Password: signInReq.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan hoac mat khau khong dung"))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong the dang nhap"))
			return
		}
	}

	ctx.SetSameSite(http.SameSiteLaxMode)

	cookie, err := ctx.Cookie("sign_in")
	if err != nil {
		zap.S().Errorln("Cookie err: ", err)
		cookie = "NotSet"

		ctx.SetCookie("sign_in", signInResult.Jwt, 3600, "/", "localhost", false, true)
	}

	zap.S().Infoln("Cookie setting")

	zap.S().Infoln(cookie)

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}{
		signInResult.Id,
		signInResult.Username,
		signInResult.Email,
		signInResult.Role,
	}, "Dang nhap thanh cong"))
}
