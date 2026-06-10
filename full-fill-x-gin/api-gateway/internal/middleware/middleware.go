package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// requestIDMiddleware generates a unique request ID for each incoming request and attaches it to the context and response headers.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = strconv.FormatInt(time.Now().UnixNano(), 36)
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// loggerMiddleware logs the HTTP method, path, status code, latency, and request ID for each request.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf(
			"method=%s path=%s status=%d latency=%s request_id=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
			c.GetString("request_id"),
		)
	}
}

// corsMiddleware sets the necessary CORS headers to allow cross-origin requests and handles preflight OPTIONS requests.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// authRequiredMiddleware checks if the request includes a valid authorization header.
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			abortWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
			return
		}

		c.Set("user_id", "demo-user")
		c.Next()
	}
}

// adminRequiredMiddleware checks if the user has an admin role based on the "X-User-Role" header.
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")

		if role != "admin" {
			abortWithError(c, http.StatusForbidden, "FORBIDDEN", "Admin access required")
			return
		}

		c.Next()
	}
}

// mockRoute returns a handler function that simulates a response from a backend service, including the service name, action, and request ID.
func mockRoute(service string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":    service,
			"action":     action,
			"request_id": c.GetString("request_id"),
		})
	}
}

// abortWithError is a helper function that aborts the request and returns a JSON error response with the specified status code, error code, message, and request ID.
func abortWithError(c *gin.Context, status int, code string, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": gin.H{
			"code":       code,
			"message":    message,
			"request_id": c.GetString("request_id"),
		},
	})
}
