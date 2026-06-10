package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	DatabaseMaxOpenConns int
	DatabaseMaxIdleConns int
	DatabaseConnMaxLife  time.Duration
	DatabaseSlowQuery    time.Duration
	DatabaseLogLevel     string
	RedisURL             string
	JWTAccessSecret      string
	JWTRefreshSecret     string
	PasswordHashCost     int
	ShutdownTimeout      time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	return &Config{
		Port:                 getEnv("AUTH_SERVICE_PORT", ":5002"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"),
		DatabaseMaxOpenConns: getIntEnv("DATABASE_MAX_OPEN_CONNS", 25),
		DatabaseMaxIdleConns: getIntEnv("DATABASE_MAX_IDLE_CONNS", 10),
		DatabaseConnMaxLife:  getDurationEnv("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute),
		DatabaseSlowQuery:    getDurationEnv("DATABASE_SLOW_QUERY_THRESHOLD", 200*time.Millisecond),
		DatabaseLogLevel:     getEnv("DATABASE_LOG_LEVEL", "warn"),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTAccessSecret:      getEnv("JWT_ACCESS_SECRET", "change-me-access"),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", "change-me-refresh"),
		PasswordHashCost:     getIntEnv("PASSWORD_HASH_COST", 12),
		ShutdownTimeout:      getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
