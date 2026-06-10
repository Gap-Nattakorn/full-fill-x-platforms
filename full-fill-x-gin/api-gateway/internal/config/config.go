package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)


type Config struct {
	Port string
	AuthServiceURL    string
	CatalogServiceURL  string
	OrderServiceURL    string
	InventoryServiceURL string
	AnalyticsServiceURL string
	JWTSecret          string
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")
	return &Config{
		Port:               getEnv("PORT", ":5001"),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:5002"),
		CatalogServiceURL:  getEnv("CATALOG_SERVICE_URL", "http://localhost:5003"),
		OrderServiceURL:    getEnv("ORDER_SERVICE_URL", "http://localhost:5004"),
		InventoryServiceURL: getEnv("INVENTORY_SERVICE_URL", "http://localhost:5005"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:5006"),
		JWTSecret:          getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

// func getDurationEnv(key string, fallback time.Duration) time.Duration {
// 	value := strings.TrimSpace(os.Getenv(key))
// 	if value == "" {
// 		return fallback
// 	}

// 	duration, err := time.ParseDuration(value)
// 	if err != nil {
// 		return fallback
// 	}
// 	return duration
// }

// func getIntEnv(key string, fallback int) int {
// 	value := strings.TrimSpace(os.Getenv(key))
// 	if value == "" {
// 		return fallback
// 	}

// 	parsed, err := strconv.Atoi(value)
// 	if err != nil || parsed < 1 {
// 		return fallback
// 	}
// 	return parsed
// }
