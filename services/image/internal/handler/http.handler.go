package image_handler

import (
	"net/http"

	image_service "github.com/098765432m/grpc-kafka/image/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ImageHttpHandler struct {
	service *image_service.ImageService
}

func NewImageHttpHandler(service *image_service.ImageService) *ImageHttpHandler {
	return &ImageHttpHandler{
		service: service,
	}
}

func (ih *ImageHttpHandler) RegisterRoutes(router *gin.RouterGroup) {
	images := router.Group("/images")

	images.GET("/:id", ih.GetImage)
	images.DELETE("/:id", ih.DeleteImage)
}

func (ih *ImageHttpHandler) GetImage(ctx *gin.Context) {

	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid image ID",
		})

	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "GetImage endpoint not implemented yet",
	})
}

func (ih *ImageHttpHandler) DeleteImage(ctx *gin.Context) {

	var id pgtype.UUID
	if err := id.Scan((ctx.Param("id"))); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid image ID",
		})
		return
	}

	err := ih.service.DeleteImage(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete image",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted image successfully",
	})
}
