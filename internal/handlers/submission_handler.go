package handlers

import (
	"net/http"
	"osp/internal/models"
	"osp/internal/services"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	submission, err := h.submissionService.CreateSubmission(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.CreateSubmissionResponse{Data: submission})
}
