package handler

import (
	"net/http"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/response"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/service"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req validation.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, tokenPair, err := h.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "REGISTER_FAILED", "Unable to register user")
		return
	}

	response.OK(c, http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
		"tokens": tokenPair,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req validation.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, tokenPair, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	response.OK(c, http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
		"tokens": tokenPair,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{"message": "refresh token placeholder"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{"message": "logout placeholder"})
}
