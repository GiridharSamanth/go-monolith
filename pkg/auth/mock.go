package auth

import (
	"context"
	"strings"
)

// MockTokenExtractor is a mock implementation of TokenExtractor
type MockTokenExtractor struct{}

// ExtractAndValidate implements TokenExtractor interface
func (m *MockTokenExtractor) ExtractAndValidate(token string, tokenType TokenType) (*UserClaims, error) {
	// TODO: Replace with actual implementation
	// This is just a mock that assumes the token is the user ID
	return &UserClaims{
		UserID: token,
	}, nil
}

// MockPermissionVerifier is a mock implementation of PermissionVerifier
type MockPermissionVerifier struct{}

// Verify implements PermissionVerifier interface
func (m *MockPermissionVerifier) Verify(ctx context.Context, action string, resource string, resourceID string) bool {
	// TODO: Replace with actual implementation
	// This mock implementation always returns true
	return true
}

// NewMockTokenExtractor creates a new mock token extractor
func NewMockTokenExtractor() TokenExtractor {
	return &MockTokenExtractor{}
}

// NewMockPermissionVerifier creates a new mock permission verifier
func NewMockPermissionVerifier() PermissionVerifier {
	return &MockPermissionVerifier{}
}

// ParseAuthHeader parses the Authorization header and returns token and type
func ParseAuthHeader(header string) (string, TokenType) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", ""
	}

	switch strings.ToLower(parts[0]) {
	case "bearer":
		// Check if it looks like a JWT (contains two dots)
		if strings.Count(parts[1], ".") == 2 {
			return parts[1], TokenTypeJWT
		}
		return parts[1], TokenTypeAccessToken
	default:
		return "", ""
	}
}
