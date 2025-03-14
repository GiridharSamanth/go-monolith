package auth

import "context"

const (
	// ContextKeyUserID is the key for user ID in context
	ContextKeyUserID = "user_id"
	// ContextKeyPermissions is the key for user permissions in context
	ContextKeyPermissions = "permissions"
)

// TokenType represents the type of token
type TokenType string

const (
	// TokenTypeJWT represents a JWT token
	TokenTypeJWT TokenType = "jwt"
	// TokenTypeAccessToken represents an access token
	TokenTypeAccessToken TokenType = "access_token"
)

// UserClaims represents the claims we extract from the token
type UserClaims struct {
	UserID string
	// Add other claims as needed
}

// PermissionVerifier defines the interface for permission checking
type PermissionVerifier interface {
	// Verify checks if the user has permission to perform the action on the resource
	Verify(ctx context.Context, action string, resource string, resourceID string) bool
}

// TokenExtractor defines the interface for token extraction and validation
type TokenExtractor interface {
	// ExtractAndValidate extracts and validates the token, returning the user claims
	ExtractAndValidate(token string, tokenType TokenType) (*UserClaims, error)
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// NewContextWithUserID creates a new context with user ID
func NewContextWithUserID(ctx context.Context, userID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ContextKeyUserID, userID)
}
