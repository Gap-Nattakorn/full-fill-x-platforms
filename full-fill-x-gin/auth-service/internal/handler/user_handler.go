package handler

import (
	"net/http"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/response"
	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{
		"id":    c.GetString("user_id"),
		"email": c.GetString("email"),
		"role":  c.GetString("role"),
	})
}
