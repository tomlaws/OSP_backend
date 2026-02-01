package routes

import (
	"osp/internal/config"
	"osp/internal/handlers"
	"osp/internal/middleware"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRouter(cfg *config.Config, client *mongo.Client, jobSystem *models.JobSystem) *gin.Engine {
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
	setupAPIRoutes(api, cfg, client, jobSystem)

	return router
}

func setupAPIRoutes(api *gin.RouterGroup, cfg *config.Config, client *mongo.Client, jobSystem *models.JobSystem) {
	// Initialize MongoDB collection
	db := client.Database("osp")
	surveysCollection := db.Collection("surveys")
	insightsCollection := db.Collection("insights")
	submissionsCollection := db.Collection("submissions")

	// Initialize services and handlers
	surveyService := services.NewSurveyService(surveysCollection)
	surveyHandler := handlers.NewSurveyHandler(surveyService)

	chatCompletionService := services.NewChatCompletionService(db.Collection("chat_completion_logs"))
	insightService := services.NewInsightService(insightsCollection, chatCompletionService, jobSystem)
	insightHandler := handlers.NewInsightHandler(insightService)
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
	// Admin routes (secured)
	admin := api.Group("/admin")
	admin.Use(middleware.AdminBearerAuth(cfg.RootToken))
	{
		insights := admin.Group("/insights")
		{
			insights.POST("", insightHandler.CreateInsight)
			insights.GET("", insightHandler.GetInsights)
		}
	}
}
