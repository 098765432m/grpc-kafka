package api_handler

import (
	"errors"
	"net/http"

	api_dto "github.com/098765432m/grpc-kafka/api-gateway/internal/dto"
	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (uh *UserHandler) GetUserById(ctx *gin.Context) {

	id := ctx.Param("id")

	userGrpc, err := uh.userClient.GetUserById(ctx, &user_pb.GetUserByIdRequest{
		Id: id,
	})
	user := userGrpc.User

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(int(st.Code()), utils.ErrorApiResponse(st.Message()))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
		return
	}

	imageGrpc, err := uh.imageClient.GetImagesByUserId(ctx, &image_pb.GetImagesByUserIdRequest{
		UserId: user.GetId(),
	})
	image := imageGrpc.Image
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:

			default:
				ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
			return

		}

	}

	resp := api_dto.UserResponse{
		Id:          user.GetId(),
		Username:    user.GetUsername(),
		Email:       user.GetEmail(),
		PhoneNumber: user.GetPhoneNumber(),
		FullName:    user.GetFullName(),
		Role:        user.GetRole(),
		HotelId:     user.GetHotelId(),
		Image: api_dto.UserImage{
			Id:       image.GetId(),
			PublicId: image.GetPublicId(),
			Format:   image.GetFormat(),
		},
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(resp, "Thanh cong"))
}

type UpdateUserParams struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Role        string `json:"role,omitempty"`
	HotelId     string `json:"hotel_id,omitempty"`
}

func (uh *UserHandler) UpdateUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	updateReq := &UpdateUserParams{}
	if err := ctx.ShouldBindJSON(updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Khong cap nhat duoc tai khoan"))
		return
	}

	_, err := uh.userClient.UpdateUserById(ctx, &user_pb.UpdateUserByIdRequest{
		User: &user_pb.User{
			Id:          id,
			Username:    updateReq.Username,
			Password:    updateReq.Password,
			Address:     updateReq.Address,
			Email:       updateReq.Email,
			PhoneNumber: updateReq.PhoneNumber,
			FullName:    updateReq.FullName,
			Role:        updateReq.Role,
			HotelId:     updateReq.HotelId,
		},
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}
}

func (uh *UserHandler) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uh.userClient.DeleteUserById(ctx, &user_pb.DeleteUserByIdRequest{
		Id: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}
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
