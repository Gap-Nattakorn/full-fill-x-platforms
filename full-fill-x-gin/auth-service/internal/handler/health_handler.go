package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/response"
	"github.com/gin-gonic/gin"
)

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	readiness ReadinessChecker
}

func NewHealthHandler(readiness ReadinessChecker) *HealthHandler {
	return &HealthHandler{readiness: readiness}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.readiness.Ping(ctx); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "Database is not ready")
		return
	}

	response.OK(c, http.StatusOK, gin.H{"status": "ready"})
}
