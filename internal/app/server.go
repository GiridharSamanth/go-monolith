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

	"go-monolith/internal/app/container"
	"go-monolith/internal/app/middleware"
	routes "go-monolith/internal/bff/route"
	"go-monolith/pkg/auth"
	appctx "go-monolith/pkg/context"
	"go-monolith/pkg/metrics"
)

type Server struct {
	router    *gin.Engine
	container *container.Container
	server    *http.Server
}

func NewServer() *Server {
	// Initialize container first to get logger
	container := container.NewContainer()

	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New() // Use New() instead of Default() as we'll add our own middleware

	// Add middlewares
	router.Use(middleware.Recovery(container.Logger)) // Recovery should be first to catch all panics
	router.Use(appctx.Middleware())                   // Context middleware early in chain
	if container.Config.Server.EnableHTTPLogs {
		router.Use(middleware.RequestLogger(container.Logger))
	}
	router.Use(metrics.HTTPMetricsMiddleware(container.Metrics)) // Add metrics middleware
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	server := &http.Server{
		Addr:    container.Config.Server.Port,
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
	permissionVerifier := auth.NewMockPermissionVerifier()

	// Setup public routes (no authentication required)
	routes.SetupPublicRoutes(s.router)

	// Setup protected routes (authentication required)
	protectedRouter := s.router.Group("")
	protectedRouter.Use(auth.AuthMiddleware())
	routes.SetupProtectedRoutes(protectedRouter, s.container.Handlers, permissionVerifier)
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
