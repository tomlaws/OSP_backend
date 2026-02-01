package handlers

import (
	"fmt"
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	insight, err := h.insightService.CreateInsight(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("Error creating insight:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create insight",
		})
		return
	}

	c.JSON(http.StatusCreated, insight)
}

func (h *InsightHandler) GetInsights(c *gin.Context) {
	// Get query parameters for pagination
	var req models.GetInsightsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	insights, err := h.insightService.GetInsights(c.Request.Context(), req.Offset, req.Limit, req.SurveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve insights",
		})
		return
	}
	c.JSON(http.StatusOK, insights)
}
