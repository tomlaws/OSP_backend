package main

import (
	"context"
	"log"
	"osp/internal/config"
	"osp/internal/database"
	"osp/internal/routes"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	client, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Create Gin router
	router := routes.SetupRouter(cfg, client)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
