package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPMetricsMiddleware creates a middleware that tracks HTTP metrics
func HTTPMetricsMiddleware(metrics *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Prepare tags
		tags := []string{
			fmt.Sprintf("endpoint:%s", c.Request.URL.Path),
			fmt.Sprintf("method:%s", c.Request.Method),
			fmt.Sprintf("status:%d", c.Writer.Status()),
		}

		// Record metrics
		metrics.RecordHTTPRequestDuration(duration, tags)
		metrics.IncrementAPICalls(c.Request.URL.Path, c.Request.Method, nil)

		// Track failed requests (4xx and 5xx)
		if c.Writer.Status() >= 400 {
			metrics.IncrementFailedRequests(tags)
		}

		// Track errors (5xx)
		if c.Writer.Status() >= 500 {
			metrics.IncrementErrors(tags)
		}
	}
}
