package main

import (
	"log"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/config"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/middleware"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	router := gin.New() // Create a new Gin router

	router.Use(gin.Recovery()) // Add recovery middleware to handle panics gracefully
	router.Use(middleware.RequestIDMiddleware()) // Add custom middleware to generate and attach a request ID
	router.Use(middleware.LoggerMiddleware()) // Add custom middleware to log request details
	router.Use(middleware.CorsMiddleware()) // Add custom middleware to handle CORS

	routes.Register(router)

	
	log.Println("api-gateway listening on " + cfg.Port)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}

}


