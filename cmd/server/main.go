package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-monolith/internal/app"
	"go-monolith/internal/app/config"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize GORM
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create and setup server
	server := app.NewServer(db)
	server.SetupRoutes()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
