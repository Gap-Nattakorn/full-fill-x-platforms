package response

import (
	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data gin.H) {
	c.JSON(200, gin.H{
		"data":       data,
		"request_id": c.GetString("request_id"),
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": gin.H{
			"code":       code,
			"message":    message,
			"request_id": c.GetString("request_id"),
		},
	})
}