package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-monolith/internal/app/container"
	"go-monolith/internal/app/middleware"
	routes "go-monolith/internal/bff/route"
	"go-monolith/pkg/auth"
	appctx "go-monolith/pkg/context"
	"go-monolith/pkg/logger"
)

type Server struct {
	router    *gin.Engine
	container *container.Container
	server    *http.Server
}

func NewServer(db *gorm.DB) *Server {
	// Initialize container first to get logger
	container := container.NewContainer(db)

	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New() // Use New() instead of Default() as we'll add our own middleware

	// Add middlewares
	router.Use(gin.Recovery())      // Use gin's built-in recovery middleware
	router.Use(appctx.Middleware()) // Add context middleware first
	router.Use(logger.HTTPLogMiddleware(container.Config.Logger, container.Config.Server.EnableHTTPLogs))
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &Server{
		router:    router,
		container: container,
		server:    server,
	}
}

func (s *Server) SetupRoutes() {
	// Initialize auth components (using mocks for now)
	tokenExtractor := auth.NewMockTokenExtractor()
	permissionVerifier := auth.NewMockPermissionVerifier()

	// Public routes (no authentication required)
	public := s.router.Group("")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Setup public routes with handlers from container
		routes.SetupPublicRoutes(public, s.container.Handlers)
	}

	// Protected routes (authentication required)
	protected := s.router.Group("")
	protected.Use(auth.AuthMiddleware(tokenExtractor))
	{
		// Story routes
		storyRoutes := protected.Group("/stories")
		{
			// Example of using permission middleware
			storyRoutes.POST("", auth.RequirePermission(permissionVerifier, "create", "story"))
			storyRoutes.PUT("/:id", auth.RequirePermission(permissionVerifier, "update", "story"))
			storyRoutes.DELETE("/:id", auth.RequirePermission(permissionVerifier, "delete", "story"))
		}

		// Setup other protected routes with handlers from container
		routes.SetupProtectedRoutes(protected, s.container.Handlers)
	}
}

func (s *Server) Start() error {
	// Start the server
	go func() {
		log.Printf("Server is starting on %s\n", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// We'll catch SIGINT (Ctrl+C) and SIGTERM (container termination)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
