package routes

import (
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/config"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/handler"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/middleware"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, cfg *config.Config) {
	registerHealthRoutes(router)
	registerAuthRoutes(router, cfg)
	registerProductRoutes(router, cfg)
	registerCartRoutes(router, cfg)
	registerOrderRoutes(router, cfg)
	registerAdminRoutes(router, cfg)
}

func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)

}

func registerAuthRoutes(router *gin.Engine, cfg *config.Config) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", proxy.Forward(cfg.AuthServiceURL))
		auth.POST("/login", proxy.Forward(cfg.AuthServiceURL))
		auth.POST("/refresh", proxy.Forward(cfg.AuthServiceURL))
		auth.POST("/logout", middleware.AuthRequiredMiddleware(), proxy.Forward(cfg.AuthServiceURL))
	}
}

func registerProductRoutes(router *gin.Engine, cfg *config.Config) {
	products := router.Group("/products")
	{
		products.GET("", proxy.Forward(cfg.CatalogServiceURL))
		products.GET("/:id", proxy.Forward(cfg.CatalogServiceURL))

		adminProducts := products.Group("")
		adminProducts.Use(middleware.AuthRequiredMiddleware(), middleware.AdminRequiredMiddleware())
		{
			adminProducts.POST("", proxy.Forward(cfg.CatalogServiceURL))
			adminProducts.PATCH("/:id", proxy.Forward(cfg.CatalogServiceURL))
			adminProducts.DELETE("/:id", proxy.Forward(cfg.CatalogServiceURL))
		}
	}
}

func registerCartRoutes(router *gin.Engine, cfg *config.Config) {
	cart := router.Group("/cart", middleware.AuthRequiredMiddleware())
	{
		cart.GET("", proxy.Forward(cfg.OrderServiceURL))
		cart.POST("/items", proxy.Forward(cfg.OrderServiceURL))
		cart.PATCH("/items/:id", proxy.Forward(cfg.OrderServiceURL))
		cart.DELETE("/items/:id", proxy.Forward(cfg.OrderServiceURL))
	}
}

func registerOrderRoutes(router *gin.Engine, cfg *config.Config) {
	orders := router.Group("/orders", middleware.AuthRequiredMiddleware())
	{
		orders.POST("/checkout", proxy.Forward(cfg.OrderServiceURL))
		orders.GET("", proxy.Forward(cfg.OrderServiceURL))
		orders.GET("/:id", proxy.Forward(cfg.OrderServiceURL))
	}

}

func registerAdminRoutes(router *gin.Engine, cfg *config.Config) {
	admin := router.Group("/admin", middleware.AuthRequiredMiddleware(), middleware.AdminRequiredMiddleware())
	{
		admin.GET("/metrics", proxy.Forward(cfg.AnalyticsServiceURL))
		admin.GET("/orders", proxy.Forward(cfg.OrderServiceURL))
		admin.GET("/inventory", proxy.Forward(cfg.InventoryServiceURL))
		admin.GET("/system-health", handler.Ready)
	}
}
