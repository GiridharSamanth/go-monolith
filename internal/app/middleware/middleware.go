package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// ContextKeyRequestID is the key for request ID in context
	ContextKeyRequestID = "request_id"
	// ContextKeyUserID is the key for user ID in context
	ContextKeyUserID = "user_id"
)

// RequestLogger returns a gin middleware for logging requests using zap
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Get request ID from header or generate new one
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(ContextKeyRequestID, requestID)

		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		status := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// CORS returns a gin middleware for handling CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Authorization returns a gin middleware for handling JWT authorization
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := parts[1]
		// TODO: Implement actual JWT validation
		if !isValidToken(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Set user ID in context after successful validation
		// This is a placeholder - actual implementation would extract user ID from JWT claims
		c.Set(ContextKeyUserID, "user_id_from_token")

		c.Next()
	}
}

// RequestContext adds request-scoped values to the context
func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not already set
		if requestID, exists := c.Get(ContextKeyRequestID); !exists {
			requestID = uuid.New().String()
			c.Set(ContextKeyRequestID, requestID)
		}

		// Set any other context values needed throughout the request lifecycle

		c.Next()
	}
}

// SecurityHeaders adds security-related headers to the response
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strict-Transport-Security
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")
		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")
		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// Content-Security-Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Referrer-Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Permissions-Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// Recovery returns a gin middleware for recovering from panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", c.GetString(ContextKeyRequestID)),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

// isValidToken is a placeholder for actual JWT validation
func isValidToken(token string) bool {
	// TODO: Implement actual JWT validation
	return token != ""
}
