package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Middleware creates a new Context for each request
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new context with trace ID
		ctx := New().
			WithTraceID(uuid.New().String()).
			WithClientInfo(c.ClientIP(), c.Request.UserAgent())

		// Add API version if present in headers
		if version := c.GetHeader("X-API-Version"); version != "" {
			ctx = ctx.WithAPIVersion(version)
		}

		// Add locale if present in headers
		if locale := c.GetHeader("Accept-Language"); locale != "" {
			ctx = ctx.WithLocale(locale)
		}

		// Store context in request
		c.Request = c.Request.WithContext(ctx.ToContext(c.Request.Context()))

		c.Next()
	}
}
