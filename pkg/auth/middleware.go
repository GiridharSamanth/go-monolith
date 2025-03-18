package auth

import (
	"fmt"
	"net/http"

	appctx "go-monolith/pkg/context"

	"github.com/gin-gonic/gin"
)

// mockTokenCache represents a mock cache of access tokens to user IDs
var mockTokenCache = map[string]string{
	"550e8400-e29b-41d4-a716-446655440000": "user_1",
	"6ba7b810-9dad-11d1-80b4-00c04fd430c8": "user_2",
	"6ba7b811-9dad-11d1-80b4-00c04fd430c8": "user_3",
	"6ba7b812-9dad-11d1-80b4-00c04fd430c8": "user_4",
	"6ba7b813-9dad-11d1-80b4-00c04fd430c8": "user_5",
}

// AuthMiddleware creates a new authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access-token header
		accessToken := c.GetHeader("access-token")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "access-token header is required",
			})
			return
		}

		// Look up user ID in mock cache
		userID, exists := mockTokenCache[accessToken]
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid access token",
			})
			return
		}

		// Update context with user ID
		ctx := appctx.FromContext(c.Request.Context()).WithUserID(userID)
		c.Request = c.Request.WithContext(ctx.ToContext(c.Request.Context()))

		c.Next()
	}
}

// RequirePermission creates a new authorization middleware for a specific action and resource
func RequirePermission(verifier PermissionVerifier, action string, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get context
		ctx := appctx.FromContext(c.Request.Context())

		fmt.Println("RequirePermission", "action", action, "resource", resource, "userID", ctx.UserID())

		// Check user ID
		userID := ctx.UserID()
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not authenticated",
			})
			return
		}

		// Get resource ID from URL parameter if it exists
		resourceID := c.Param("id")

		// Verify permission
		if !verifier.Verify(c.Request.Context(), action, resource, resourceID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
			})
			return
		}

		c.Next()
	}
}
