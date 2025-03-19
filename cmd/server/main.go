package main

import (
	"log"
	"os"

	"go-monolith/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env.local file only in development/local environment
	env := os.Getenv("APP_ENV")
	if env == "" || env == "development" || env == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Printf("Warning: Error loading .env.local file: %v", err)
		}
	}

	// Create and setup server
	server := app.NewServer()
	server.SetupRoutes()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
