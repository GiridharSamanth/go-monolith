package auth

import (
	"net/http"

	appctx "go-monolith/pkg/context"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a new authentication middleware
func AuthMiddleware(tokenExtractor TokenExtractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		// Parse token and type from header
		token, tokenType := ParseAuthHeader(authHeader)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		// Extract and validate token
		claims, err := tokenExtractor.ExtractAndValidate(token, tokenType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Update context with user ID
		ctx := appctx.FromContext(c.Request.Context()).WithUserID(claims.UserID)
		c.Request = c.Request.WithContext(ctx.ToContext(c.Request.Context()))

		c.Next()
	}
}

// RequirePermission creates a new authorization middleware for a specific action and resource
func RequirePermission(verifier PermissionVerifier, action string, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get context
		ctx := appctx.FromContext(c.Request.Context())

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
