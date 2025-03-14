package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-monolith/internal/bff/handler"
)

// SetupPublicRoutes configures all the public routes for the BFF service
func SetupPublicRoutes(router gin.IRouter, handlers *handler.Handlers) {
	// v1.2 routes
	v1_2Group := router.Group("/v1.2")
	{
		// Story routes
		storyRoutes := v1_2Group.Group("/stories")
		{
			storyRoutes.GET("/:id", handlers.V1_2StoryHandler.GetStory)
		}
	}

	// v2.0 routes
	v2_0Group := router.Group("/v2.0")
	{
		// Story routes with enhanced functionality
		storyRoutes := v2_0Group.Group("/stories")
		{
			storyRoutes.GET("/:id", handlers.V2_0StoryHandler.GetStory)
		}
	}
}

// SetupProtectedRoutes configures all the protected routes for the BFF service
func SetupProtectedRoutes(router gin.IRouter, handlers *handler.Handlers) {
	// v1.2 protected routes
	v1_2Group := router.Group("/v1.2")
	{
		// Author routes
		authorRoutes := v1_2Group.Group("/authors")
		{
			// Placeholder endpoint until actual implementation
			authorRoutes.POST("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
			})
		}

		// Protected story routes
		storyRoutes := v1_2Group.Group("/stories")
		{
			// Placeholder endpoint until actual implementation
			storyRoutes.POST("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
			})
		}
	}

	// v2.0 protected routes
	v2_0Group := router.Group("/v2.0")
	{
		// Author routes
		authorRoutes := v2_0Group.Group("/authors")
		{
			// Placeholder endpoint until actual implementation
			authorRoutes.POST("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
			})
		}

		// Protected story routes
		storyRoutes := v2_0Group.Group("/stories")
		{
			// Placeholder endpoint until actual implementation
			storyRoutes.POST("", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
			})
		}
	}
}
