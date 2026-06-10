package routes

import (
	"net/http"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	registerHealthRoutes(router)
	registerAuthRoutes(router)
	registerProductRoutes(router)
	registerCartRoutes(router)
	registerOrderRoutes(router)
	registerAdminRoutes(router)
}

func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})
}

func registerAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", mockRoute("auth-service", "register"))
		auth.POST("/login", mockRoute("auth-service", "login"))
		auth.POST("/refresh", mockRoute("auth-service", "refresh token"))

		auth.POST("/logout", middleware.AuthRequiredMiddleware(), mockRoute("auth-service", "logout"))
	}
}

func registerProductRoutes(router *gin.Engine) {
	products := router.Group("/products")
	{
		products.GET("", mockRoute("catalog-service", "list products"))
		products.GET("/:id", mockRoute("catalog-service", "get product"))

		adminProducts := products.Group("")
		adminProducts.Use(middleware.AuthRequiredMiddleware(), middleware.AdminRequiredMiddleware())
		{
			adminProducts.POST("", mockRoute("catalog-service", "create product"))
			adminProducts.PATCH("/:id", mockRoute("catalog-service", "update product"))
			adminProducts.DELETE("/:id", mockRoute("catalog-service", "delete product"))
		}
	}
}

func registerCartRoutes(router *gin.Engine) {
	cart := router.Group("/cart")
	cart.Use(middleware.AuthRequiredMiddleware())
	{
		cart.GET("", mockRoute("order-service", "get cart"))
		cart.POST("/items", mockRoute("order-service", "add cart item"))
		cart.PATCH("/items/:id", mockRoute("order-service", "update cart item"))
		cart.DELETE("/items/:id", mockRoute("order-service", "remove cart item"))
	}
}

func registerOrderRoutes(router *gin.Engine) {
	orders := router.Group("/orders")
	orders.Use(middleware.AuthRequiredMiddleware())
	{
		orders.POST("/checkout", mockRoute("order-service", "checkout"))
		orders.GET("", mockRoute("order-service", "list orders"))
		orders.GET("/:id", mockRoute("order-service", "get order"))
	}
}

func registerAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthRequiredMiddleware(), middleware.AdminRequiredMiddleware())
	{
		admin.GET("/metrics", mockRoute("analytics-service", "admin metrics"))
		admin.GET("/orders", mockRoute("order-service", "admin orders"))
		admin.GET("/inventory", mockRoute("inventory-service", "admin inventory"))
		admin.GET("/system-health", mockRoute("api-gateway", "system health"))
	}
}

func mockRoute(service string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":    service,
			"action":     action,
			"request_id": c.GetString("request_id"),
		})
	}
}