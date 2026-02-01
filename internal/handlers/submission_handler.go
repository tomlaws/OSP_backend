package handlers

import (
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SubmissionHandler struct {
	submissionService services.ISubmissionService
}

func NewSubmissionHandler(submissionService services.ISubmissionService) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
	}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	var req models.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.CreateSubmissionResponse{
			Error: err.Error(),
		})
		return
	}
	submission, err := h.submissionService.CreateSubmission(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.CreateSubmissionResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.CreateSubmissionResponse{Data: submission})
}

func (h *SubmissionHandler) GetSubmissions(c *gin.Context) {
	// Get query parameters for pagination
	var req models.GetSubmissionsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &models.GetSubmissionsResponse{
			Error: "Invalid query parameters",
		})
		return
	}
	// Convert surveyID string to bson.ObjectID pointer if provided
	var surveyID *bson.ObjectID
	if req.SurveyID != nil {
		id, err := bson.ObjectIDFromHex(*req.SurveyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, &models.GetSubmissionsResponse{
				Error: "Invalid survey ID",
			})
			return
		}
		surveyID = &id
	}
	submissions, err := h.submissionService.GetSubmissions(c.Request.Context(), req.Offset, req.Limit, surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.GetSubmissionsResponse{
			Error: "Failed to retrieve submissions",
		})
		return
	}
	c.JSON(http.StatusOK, &models.GetSubmissionsResponse{
		Data: submissions,
	})
}

func (h *SubmissionHandler) DeleteSubmission(c *gin.Context) {
	var uriReq models.DeleteSubmissionRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, &models.DeleteSubmissionResponse{
			Error: err.Error(),
		})
		return
	}
	id, err := bson.ObjectIDFromHex(uriReq.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.DeleteSubmissionResponse{
			Error: "Invalid submission ID",
		})
		return
	}
	err = h.submissionService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.DeleteSubmissionResponse{
			Error: "Failed to delete submission",
		})
		return
	}
	c.JSON(http.StatusOK, &models.DeleteSubmissionResponse{
		Error: "",
	})
}
