package routes

import (
	"osp/internal/config"
	"osp/internal/handlers"
	"osp/internal/middleware"
	"osp/internal/services"

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
	// Initialize MongoDB collection
	db := client.Database("osp")
	surveysCollection := db.Collection("surveys")
	submissionsCollection := db.Collection("submissions")

	// Initialize services and handlers
	surveyService := services.NewSurveyService(surveysCollection)
	surveyHandler := handlers.NewSurveyHandler(surveyService)

	// Surveys routes
	surveys := api.Group("/surveys")
	{
		surveys.POST("", surveyHandler.CreateSurvey)
		surveys.GET("/:token", surveyHandler.GetSurvey)
	}
	// Submissions routes
	submissionService := services.NewSubmissionService(submissionsCollection)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)
	submissions := api.Group("/submissions")
	{
		submissions.POST("", submissionHandler.CreateSubmission)
	}
}
