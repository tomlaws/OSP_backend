package handlers

import (
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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
		c.JSON(http.StatusBadRequest, &models.CreateSurveyResponse{
			Error: err.Error(),
		})
		return
	}
	// Specfication validation
	for _, question := range req.Questions {
		if question.Type == models.QuestionTypeMultipleChoice && question.Specification.MultipleChoiceSpecification == nil {
			c.JSON(http.StatusBadRequest, &models.CreateSurveyResponse{
				Error: "MultipleChoiceSpecification is required for MULTIPLE_CHOICE question type",
			})
			return
		}
		if question.Type == models.QuestionTypeLikert && question.Specification.LikertSpecification == nil {
			c.JSON(http.StatusBadRequest, &models.CreateSurveyResponse{
				Error: "LikertSpecification is required for LIKERT question type",
			})
			return
		}
		if question.Type == models.QuestionTypeTextbox && question.Specification.TextboxSpecification == nil {
			c.JSON(http.StatusBadRequest, &models.CreateSurveyResponse{
				Error: "TextboxSpecification is required for TEXTBOX question type",
			})
			return
		}
	}

	survey, err := h.surveyService.CreateSurvey(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.CreateSurveyResponse{
			Error: "Failed to create survey",
		})
		return
	}

	c.JSON(http.StatusCreated, &models.CreateSurveyResponse{
		Data: survey,
	})
}

func (h *SurveyHandler) ListSurveys(c *gin.Context) {
	var req models.ListSurveysRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.ListSurveysResponse{
			Error: err.Error(),
		})
		return
	}
	surveys, err := h.surveyService.ListSurveys(c.Request.Context(), req.Offset, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.ListSurveysResponse{
			Error: "Failed to retrieve surveys",
		})
		return
	}
	c.JSON(http.StatusOK, models.ListSurveysResponse{Data: surveys})
}

func (h *SurveyHandler) GetSurveyByToken(c *gin.Context) {
	var uriReq models.GetSurveyByTokenRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, &models.GetSurveyByTokenResponse{
			Error: err.Error(),
		})
		return
	}
	survey, err := h.surveyService.GetSurveyByToken(c.Request.Context(), uriReq.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.GetSurveyByTokenResponse{
			Error: "Invalid survey token",
		})
		return
	}

	c.JSON(http.StatusOK, &models.GetSurveyByTokenResponse{
		Data: survey,
	})
}

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	var uriReq models.GetSurveyRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, &models.GetSurveyResponse{
			Error: err.Error(),
		})
		return
	}
	var surveyID bson.ObjectID
	surveyID, err := bson.ObjectIDFromHex(uriReq.ID)
	survey, err := h.surveyService.GetSurveyByID(c.Request.Context(), surveyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.GetSurveyResponse{
			Error: "Invalid survey ID",
		})
		return
	}
	c.JSON(http.StatusOK, survey)
}

func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
	var uriReq models.DeleteSurveyRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, &models.DeleteSurveyResponse{
			Error: err.Error(),
		})
		return
	}
	surveyID, err := bson.ObjectIDFromHex(uriReq.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.DeleteSurveyResponse{
			Error: "Invalid survey ID",
		})
		return
	}
	err = h.surveyService.DeleteSurvey(c.Request.Context(), surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.DeleteSurveyResponse{
			Error: "Failed to delete survey",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
