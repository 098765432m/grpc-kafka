package rating_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/utils"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RatingHttpHandler struct {
	service *rating_service.RatingService
}

func NewHotelHttpHandler(service *rating_service.RatingService) *RatingHttpHandler {
	return &RatingHttpHandler{
		service: service,
	}
}

func (rh *RatingHttpHandler) RegisterRoutes(router *gin.RouterGroup) {
	ratings := router.Group("/ratings")

	ratings.POST("/", rh.CreateRating)
}

type CreateRatingParams struct {
	Rating  int    `json:"rating"`
	HotelId string `json:"hotel_id"`
	UserId  string `json:"user_id"`
	Comment string `json:"comment,omitempty"`
}

func (rh *RatingHttpHandler) CreateRating(ctx *gin.Context) {

	createRatingReq := &CreateRatingParams{}
	if err := ctx.ShouldBindJSON(createRatingReq); err != nil {
		zap.S().Errorln("Failed to get request body: ", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Failed to get request body"))
	}

	var hotelId pgtype.UUID
	if err := hotelId.Scan(createRatingReq.HotelId); err != nil {
		errMsg := "Failed to convert hotel UUID"
		zap.S().Errorln(errMsg, err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse(errMsg))
	}

	var userId pgtype.UUID
	if err := userId.Scan(createRatingReq.UserId); err != nil {
		errMsg := "Failed to convert user UUID"
		zap.S().Errorln(errMsg, err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse(errMsg))
	}

	err := rh.service.CreateRating(ctx, &rating_repo.CreateRatingParams{
		Rating:  int32(createRatingReq.Rating),
		HotelID: hotelId.String(),
		UserID:  userId.String(),
		Comment: createRatingReq.Comment,
	})

	if err != nil {
		errMsg := "Failed to create Rating"
		zap.S().Errorln(errMsg, err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse(errMsg))
	}

	ctx.JSON(http.StatusCreated, utils.SuccessApiResponse(nil, "Created rating successfully"))
}
