package handler

import (
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/response"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}

func Ready(c *gin.Context) {
	response.OK(c, gin.H{"status": "ready"})
}