package api_handler

import (
	"net/http"
	"strconv"

	api_dto "github.com/098765432m/grpc-kafka/api-gateway/internal/dto"
	"github.com/098765432m/grpc-kafka/common/gen-proto/booking_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/gen-proto/user_pb"
	"github.com/098765432m/grpc-kafka/common/model"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userClient    user_pb.UserServiceClient
	imageClient   image_pb.ImageServiceClient
	bookingClient booking_pb.BookingServiceClient
}

func NewUserHandler(
	userClient user_pb.UserServiceClient,
	imageClient image_pb.ImageServiceClient,
	bookingClient booking_pb.BookingServiceClient,
) *UserHandler {
	return &UserHandler{
		userClient:    userClient,
		imageClient:   imageClient,
		bookingClient: bookingClient,
	}
}

func (uh *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	userHandler := router.Group("/users")

	userHandler.POST("/")

	userHandler.GET("/:id", uh.GetUserById)
	userHandler.PUT("/:id", uh.UpdateUserById)

	userHandler.GET("/:id/bookings", uh.GetBookingsByUserId)

	userHandler.POST("/sign-in", uh.SignIn)
	userHandler.POST("/sign-up", uh.SignUp)
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
	user := userGrpc.GetUser()

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan khong ton tai"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	imageGrpc, err := uh.imageClient.GetImageByUserId(ctx, &image_pb.GetImageByUserIdRequest{
		UserId: user.GetId(),
	})
	image := imageGrpc.GetImage()
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

	zap.S().Infoln("IMage: ", image)

	resp := api_dto.UserResponse{
		Id:          user.GetId(),
		Username:    user.GetUsername(),
		Email:       user.GetEmail(),
		PhoneNumber: user.GetPhoneNumber(),
		FullName:    user.GetFullName(),
		Role:        user.GetRole(),
	}

	// if image exist add image
	if image != nil {
		resp.Image = &api_dto.UserImage{
			Id:       image.GetId(),
			PublicId: image.GetPublicId(),
			Format:   image.GetFormat(),
		}
	}

	// if hotel ID exists, add hotel Id
	if user.GetHotelId() != "" {
		resp.HotelId = user.HotelId
	}

	zap.S().Infoln("User: ", resp)

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
		zap.S().Errorln("Cannot bind JSON Request: ", err)
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
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan khong ton tai de cap nhat"))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Cap nhat tai khoan thanh cong"))
}

// Return Bookings By UserId
func (uh *UserHandler) GetBookingsByUserId(ctx *gin.Context) {
	// Lay Request Param
	userId := ctx.Param("id")
	checkDateStart := ctx.Query("check_date_start")
	checkDateEnd := ctx.Query("check_date_end")
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		zap.S().Infoln("Failed to convert int of Size")
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi xay ra"))
		return
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		zap.S().Infoln("Failed to convert int of Offset")
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi xay ra"))
		return
	}

	zap.L().Info("Check Request of get bookings by userId", zap.Any("userId", userId), zap.Any("check Date Start", checkDateStart), zap.Any("Check Date End", checkDateEnd), zap.Any("Size", size), zap.Any("Offset", offset))

	// Lay danh sach Bookings
	resultBookingsByUserId, err := uh.bookingClient.GetBookingsByUserId(ctx, &booking_pb.GetBookingsByUserIdRequest{
		UserId:         userId,
		CheckDateStart: checkDateStart,
		CheckDateEnd:   checkDateEnd,
		Size:           int32(size),
		Offset:         int32(offset),
	})
	if err != nil {
		zap.S().Infoln("Failed to get list of Bookings by User ID: ", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong lay duoc danh sach dat phong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(resultBookingsByUserId.Bookings, "Thanh cong"))
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
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan khong ton tai de xoa"))
				return
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(err.Error()))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}
}

type SignUpRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {
	req := &SignUpRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap.S().Info("Failed to parsed req body: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi request body"))
		return
	}

	_, err := uh.userClient.CreateUser(ctx, &user_pb.CreateUserRequest{
		Username:    req.Username,
		Password:    req.Password,
		Address:     req.Address,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Role:        model.GUEST_ROLE,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan da ton tai"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Dang ky that bai"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Dang ky thanh cong"))
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
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Tai khoan hoac mat khau khong dung"))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi khong the dang nhap"))
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)

	cookie, err := ctx.Cookie("user")
	if err != nil {
		zap.S().Errorln("Cookie err: ", err)
		cookie = "NotSet"

		ctx.SetCookie("user", signInResult.Jwt, 3600, "/", "localhost", false, true)
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
