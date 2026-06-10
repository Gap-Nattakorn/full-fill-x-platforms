package response

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

func Error(c *gin.Context, status int, code string, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
