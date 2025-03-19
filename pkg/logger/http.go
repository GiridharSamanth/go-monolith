package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// HTTPRequest represents a generic HTTP request that can be implemented by different web frameworks
type HTTPRequest interface {
	Method() string
	Path() string
	Query() string
	Status() int
	ClientIP() string
	UserAgent() string
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	GetContext() context.Context
}

// HTTPLogger provides HTTP request logging functionality
type HTTPLogger struct {
	logger Logger
}

// NewHTTPLogger creates a new HTTP logger instance
func NewHTTPLogger(logger Logger) *HTTPLogger {
	return &HTTPLogger{
		logger: logger,
	}
}

// LogRequest logs an HTTP request with standard fields
func (h *HTTPLogger) LogRequest(req HTTPRequest, start time.Time) {
	// Get request ID from context or generate new one
	requestID, exists := req.Get("request_id")
	if !exists {
		requestID = uuid.New().String()
		req.Set("request_id", requestID)
	}

	// Calculate latency
	latency := time.Since(start)

	// Build path with query if present
	path := req.Path()
	if query := req.Query(); query != "" {
		path = path + "?" + query
	}

	// Log request details
	h.logger.Info(req.GetContext(), "request completed",
		String("request_id", requestID.(string)),
		String("method", req.Method()),
		String("path", path),
		Int("status", req.Status()),
		String("latency", latency.String()),
		String("client_ip", req.ClientIP()),
		String("user_agent", req.UserAgent()),
	)
}
