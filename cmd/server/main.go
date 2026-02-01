package main

import (
	"context"
	"log"
	"osp/internal/config"
	"osp/internal/database"
	"osp/internal/models"
	"osp/internal/routes"

	"github.com/hibiken/asynq"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	dbClient, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer dbClient.Disconnect(context.Background())

	// Create asynq client
	redisOpt, err := asynq.ParseRedisURI(cfg.RedisUri)
	if err != nil {
		log.Fatalf("Failed to parse Redis URI: %v", err)
	}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// Start asynq worker
	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"insights": 1,
		},
	})

	mux := asynq.NewServeMux()
	jobSystem := &models.JobSystem{
		Client: asynqClient,
		Server: asynqServer,
		Mux:    mux,
	}

	// Create Gin router
	router := routes.SetupRouter(cfg, dbClient, jobSystem)

	go func() {
		if err := asynqServer.Run(mux); err != nil {
			log.Printf("asynq server stopped: %v", err)
		}
	}()

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
