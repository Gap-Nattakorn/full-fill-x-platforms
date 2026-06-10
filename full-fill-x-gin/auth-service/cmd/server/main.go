package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/config"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/handler"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/infra/postgres"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/middleware"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/repository"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/routes"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/security"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database, err := postgres.Open(postgres.Config{
		DSN:             cfg.DatabaseURL,
		MaxOpenConns:    cfg.DatabaseMaxOpenConns,
		MaxIdleConns:    cfg.DatabaseMaxIdleConns,
		ConnMaxLifetime: cfg.DatabaseConnMaxLife,
		SlowThreshold:   cfg.DatabaseSlowQuery,
		LogLevel:        cfg.DatabaseLogLevel,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewGormUserRepository(database.Gorm())
	tokenService := service.NewTokenService(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	passwordHasher := security.NewBcryptPasswordHasher(cfg.PasswordHashCost)
	authService := service.NewAuthService(userRepository, tokenService, passwordHasher)
	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler(database)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CorsMiddleware())

	routes.Register(router, authHandler, healthHandler)

	log.Println("auth-service listening on " + cfg.Port)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

	if err := database.Close(); err != nil {
		log.Fatal(err)
	}
}
