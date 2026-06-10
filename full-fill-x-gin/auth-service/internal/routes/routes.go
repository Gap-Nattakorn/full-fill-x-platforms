package routes

import (
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/handler"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler) {
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", middleware.AuthRequiredMiddleware(), authHandler.Logout)
	}

	users := router.Group("/users", middleware.AuthRequiredMiddleware())
	{
		users.GET("/me", handler.Me)
	}
}
