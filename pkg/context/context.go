package context

import (
	"context"
	"time"
)

type contextKey string

const (
	contextKeyApp = contextKey("app_context")
)

// Context holds request-scoped data
type Context struct {
	traceID    string
	userID     string
	timestamp  time.Time
	clientIP   string
	userAgent  string
	apiVersion string
	locale     string
}

// New creates a new Context with initial values
func New() *Context {
	return &Context{
		timestamp: time.Now(),
	}
}

// WithTraceID sets the trace ID
func (c *Context) WithTraceID(traceID string) *Context {
	c.traceID = traceID
	return c
}

// WithUserID sets the user ID
func (c *Context) WithUserID(userID string) *Context {
	c.userID = userID
	return c
}

// WithClientInfo sets client-related information
func (c *Context) WithClientInfo(clientIP, userAgent string) *Context {
	c.clientIP = clientIP
	c.userAgent = userAgent
	return c
}

// WithAPIVersion sets the API version
func (c *Context) WithAPIVersion(version string) *Context {
	c.apiVersion = version
	return c
}

// WithLocale sets the locale
func (c *Context) WithLocale(locale string) *Context {
	c.locale = locale
	return c
}

// TraceID returns the trace ID
func (c *Context) TraceID() string {
	return c.traceID
}

// UserID returns the user ID
func (c *Context) UserID() string {
	return c.userID
}

// Timestamp returns the context creation timestamp
func (c *Context) Timestamp() time.Time {
	return c.timestamp
}

// ClientIP returns the client IP address
func (c *Context) ClientIP() string {
	return c.clientIP
}

// UserAgent returns the user agent string
func (c *Context) UserAgent() string {
	return c.userAgent
}

// APIVersion returns the API version
func (c *Context) APIVersion() string {
	return c.apiVersion
}

// Locale returns the locale
func (c *Context) Locale() string {
	return c.locale
}

// ToContext adds the Context to a context.Context
func (c *Context) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyApp, c)
}

// FromContext retrieves the Context from a context.Context
func FromContext(ctx context.Context) *Context {
	if ctx == nil {
		return New()
	}
	if c, ok := ctx.Value(contextKeyApp).(*Context); ok {
		return c
	}
	return New()
}

// FromRequest creates a new Context from a request context
func FromRequest(ctx context.Context) *Context {
	return FromContext(ctx)
}
