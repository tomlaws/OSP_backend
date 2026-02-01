package handlers

import (
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
)

type SurveyHandler struct {
	surveyService services.ISurveyService
}

func NewSurveyHandler(surveyService services.ISurveyService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
	}
}

func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	var req models.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Specfication validation
	for _, question := range req.Questions {
		if question.Type == models.QuestionTypeMultipleChoice && question.Specification.MultipleChoiceSpecification == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "MultipleChoiceSpecification is required for MULTIPLE_CHOICE question type",
			})
			return
		}
		if question.Type == models.QuestionTypeLikert && question.Specification.LikertSpecification == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "LikertSpecification is required for LIKERT question type",
			})
			return
		}
		if question.Type == models.QuestionTypeTextbox && question.Specification.TextboxSpecification == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "TextboxSpecification is required for TEXTBOX question type",
			})
			return
		}
	}

	item, err := h.surveyService.CreateSurvey(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create item",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	token := c.Param("token")
	survey, err := h.surveyService.GetSurveyByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID",
		})
		return
	}

	c.JSON(http.StatusOK, survey)
}
