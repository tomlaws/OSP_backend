package routes

import (
	"osp/internal/config"
	"osp/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRouter(cfg *config.Config, client *mongo.Client) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	setupAPIRoutes(api, cfg, client)

	return router
}

func setupAPIRoutes(api *gin.RouterGroup, cfg *config.Config, client *mongo.Client) {
}
