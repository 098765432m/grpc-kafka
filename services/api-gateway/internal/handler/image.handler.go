package api_handler

import (
	"net/http"

	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageClient image_pb.ImageServiceClient
}

func NewImageHandler(imageClient image_pb.ImageServiceClient) *ImageHandler {
	return &ImageHandler{
		imageClient: imageClient,
	}
}

func (ih *ImageHandler) RegisterRoutes(router *gin.RouterGroup) {
	imageHandler := router.Group("/images")

	imageHandler.DELETE("/:id", ih.DeleteImageById)
}

func (ih *ImageHandler) DeleteImageById(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := ih.imageClient.DeleteImageById(ctx, &image_pb.DeleteImageByIdRequest{
		Id: id,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Xoa tai khoan thanh cong"))
}

type DeleteImagesByIdsParams struct {
	Ids []string
}

func (ih *ImageHandler) DeleteImagesByIds(ctx *gin.Context) {
	var req DeleteImagesByIdsParams
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	_, err := ih.imageClient.DeleteImagesByIds(ctx, &image_pb.DeleteImagesByIdsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorApiResponse("Loi he thong"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessApiResponse(nil, "Xoa thanh cong"))
}
