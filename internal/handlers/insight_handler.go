package handlers

import (
	"fmt"
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InsightHandler struct {
	insightService services.IInsightService
}

func NewInsightHandler(insightService services.IInsightService) *InsightHandler {
	return &InsightHandler{
		insightService: insightService,
	}
}

func (h *InsightHandler) CreateInsight(c *gin.Context) {
	var req models.CreateInsightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.CreateInsightResponse{
			Error: err.Error(),
		})
		return
	}

	insight, err := h.insightService.CreateInsight(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("Error creating insight:", err)
		c.JSON(http.StatusInternalServerError, &models.CreateInsightResponse{
			Error: "Failed to create insight",
		})
		return
	}
	c.JSON(http.StatusCreated, &models.CreateInsightResponse{
		Data: insight,
	})
}

func (h *InsightHandler) GetInsights(c *gin.Context) {
	// Get query parameters for pagination
	var req models.GetInsightsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.GetInsightsResponse{
			Error: "Invalid query parameters",
		})
		return
	}
	// Convert surveyID string to bson.ObjectID pointer if provided
	var surveyID *bson.ObjectID
	if req.SurveyID != nil {
		id, err := bson.ObjectIDFromHex(*req.SurveyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, &models.GetInsightsResponse{
				Error: "Invalid survey ID",
			})
			return
		}
		surveyID = &id
	}
	insights, err := h.insightService.GetInsights(c.Request.Context(), req.Offset, req.Limit, surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.GetInsightsResponse{
			Error: "Failed to retrieve insights",
		})
		return
	}
	c.JSON(http.StatusOK, &models.GetInsightsResponse{
		Data: insights,
	})
}

func (h *InsightHandler) GetInsight(c *gin.Context) {
	var req models.GetInsightRequest
	if err := c.BindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.GetInsightResponse{
			Error: "Invalid query parameters",
		})
		return
	}
	var err error
	var insightID bson.ObjectID
	insightID, err = bson.ObjectIDFromHex(req.ID)
	insight, err := h.insightService.GetInsight(c.Request.Context(), insightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.GetInsightResponse{
			Error: "Failed to retrieve insight",
		})
		return
	}
	c.JSON(http.StatusOK, &models.GetInsightResponse{
		Data: insight,
	})
}
