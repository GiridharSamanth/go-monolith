package logger

import (
	"context"

	"github.com/gin-gonic/gin"
)

// GinRequest implements HTTPRequest interface for Gin framework
type GinRequest struct {
	c *gin.Context
}

// NewGinRequest creates a new GinRequest adapter
func NewGinRequest(c *gin.Context) *GinRequest {
	return &GinRequest{c: c}
}

// Method returns the HTTP method
func (g *GinRequest) Method() string {
	return g.c.Request.Method
}

// Path returns the request path
func (g *GinRequest) Path() string {
	return g.c.Request.URL.Path
}

// Query returns the raw query string
func (g *GinRequest) Query() string {
	return g.c.Request.URL.RawQuery
}

// Status returns the response status code
func (g *GinRequest) Status() int {
	return g.c.Writer.Status()
}

// ClientIP returns the client IP address
func (g *GinRequest) ClientIP() string {
	return g.c.ClientIP()
}

// UserAgent returns the user agent string
func (g *GinRequest) UserAgent() string {
	return g.c.Request.UserAgent()
}

// Set stores a value in the context
func (g *GinRequest) Set(key string, value interface{}) {
	g.c.Set(key, value)
}

// Get retrieves a value from the context
func (g *GinRequest) Get(key string) (interface{}, bool) {
	return g.c.Get(key)
}

// GetContext returns the request context
func (g *GinRequest) GetContext() context.Context {
	return g.c.Request.Context()
}
