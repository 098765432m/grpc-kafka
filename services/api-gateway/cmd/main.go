package main

import (
	"fmt"

	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	utils.Init()
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	zap.S().Infoln("Running api-gateway on port 5025")

	if err := router.Run(fmt.Sprintf(":%d", 5025)); err != nil {
		zap.S().Fatalf("Failed to start HTTP server: %v\n", err)

	}
}
