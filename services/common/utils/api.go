package utils

import "github.com/gin-gonic/gin"

func ErrorApiResponse(msg string) gin.H {
	return gin.H{
		"success": false,
		"error":   msg,
	}
}

func SuccessApiResponse(result any, msg string) gin.H {
	return gin.H{
		"success": true,
		"result":  result,
		"message": msg,
	}
}
