package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-monolith/internal/bff/handler"
	"go-monolith/pkg/auth"
)

// SetupPublicRoutes configures all public routes that don't require authentication
func SetupPublicRoutes(router gin.IRouter) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Access token endpoint
	router.GET("/get-accesstoken", func(c *gin.Context) {
		// TODO: Implement token generation logic
		c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
	})
}

// SetupProtectedRoutes configures all routes that require authentication
func SetupProtectedRoutes(router gin.IRouter, handlers *handler.Handlers, permissionVerifier auth.PermissionVerifier) {
	// v1.2 routes
	router.POST("/v1.2/authors",
		auth.RequirePermission(permissionVerifier, "create", "author"),
		func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
		},
	)

	router.GET("/v1.2/stories",
		auth.RequirePermission(permissionVerifier, "get", "story"),
		handlers.V1_2StoryHandler.GetStory,
	)

	// v2.0 routes
	router.GET("/v2.0/stories/:id",
		auth.RequirePermission(permissionVerifier, "get", "story"),
		handlers.V2_0StoryHandler.GetStory,
	)

	router.DELETE("/v2.0/stories/:id",
		auth.RequirePermission(permissionVerifier, "delete", "story"),
		func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
		},
	)
}
