package utils

import "github.com/gin-gonic/gin"

func ErrorResponse(msg string) gin.H {
	return gin.H{
		"success": false,
		"error":   msg,
	}
}

func SuccessResponse(result any, msg string) gin.H {
	return gin.H{
		"success": true,
		"result":  result,
		"message": msg,
	}
}
