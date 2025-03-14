package logger

import (
	"time"

	appctx "go-monolith/pkg/context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTPLogMiddleware creates a middleware for logging HTTP requests
func HTTPLogMiddleware(logger *Logger, enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if logging is disabled
		if !enabled {
			c.Next()
			return
		}

		// Get context
		ctx := appctx.FromContext(c.Request.Context())

		// Time request
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)

		// Prepare fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
		}

		// Add query parameters if present
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// Add error if present
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Log based on status code
		msg := "HTTP Request"
		switch {
		case c.Writer.Status() >= 500:
			logger.Error(ctx.ToContext(c.Request.Context()), msg, fields...)
		case c.Writer.Status() >= 400:
			logger.Warn(ctx.ToContext(c.Request.Context()), msg, fields...)
		default:
			logger.Info(ctx.ToContext(c.Request.Context()), msg, fields...)
		}
	}
}
