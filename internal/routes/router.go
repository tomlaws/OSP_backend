package routes

import (
	"osp/internal/config"
	"osp/internal/handlers"
	"osp/internal/middleware"
	"osp/internal/models"
	"osp/internal/repositories"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRouter(cfg *config.Config, client *mongo.Client, jobSystem *models.JobSystem) *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())

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
	surveyRepo := repositories.NewMongoSurveyRepository(surveysCollection)
	submissionRepo := repositories.NewMongoSubmissionRepository(submissionsCollection)
	insightRepo := repositories.NewMongoInsightRepository(insightsCollection)

	surveyService := services.NewSurveyService(surveyRepo)
	surveyHandler := handlers.NewSurveyHandler(surveyService)

	chatCompletionService := services.NewChatCompletionService(db.Collection("chat_completion_logs"))
	insightService := services.NewInsightService(insightRepo, surveyRepo, submissionRepo, chatCompletionService, jobSystem.Client)
	insightService.RegisterHandlers(jobSystem.Mux)
	insightHandler := handlers.NewInsightHandler(insightService)

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Surveys routes
	surveys := api.Group("/surveys")
	{
		surveys.GET("/:token", surveyHandler.GetSurveyByToken)
	}
	// Submissions routes
	submissionService := services.NewSubmissionService(submissionRepo, surveyRepo)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)
	submissions := api.Group("/submissions")
	{
		submissions.POST("", submissionHandler.CreateSubmission)
	}
	// Admin routes (secured)
	admin := api.Group("/admin")
	admin.Use(middleware.AdminBearerAuth(cfg.RootToken))
	{
		surveys := admin.Group("/surveys")
		{
			surveys.POST("", surveyHandler.CreateSurvey)
			surveys.GET("", surveyHandler.ListSurveys)
			surveys.GET("/:id", surveyHandler.GetSurvey)
			surveys.DELETE("/:id", surveyHandler.DeleteSurvey)
		}
		submissions := admin.Group("/submissions")
		{
			submissions.GET("/", submissionHandler.GetSubmissions)
			submissions.DELETE("/:id", submissionHandler.DeleteSubmission)
		}
		insights := admin.Group("/insights")
		{
			insights.POST("", insightHandler.CreateInsight)
			insights.GET("", insightHandler.GetInsights)
			insights.GET("/:id", insightHandler.GetInsight)
		}
	}
}
