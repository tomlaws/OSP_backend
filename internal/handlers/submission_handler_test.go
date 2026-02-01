package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"osp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockSubmissionService is a mock implementation of ISubmissionService
type MockSubmissionService struct {
	mock.Mock
}

func (m *MockSubmissionService) CreateSubmission(ctx context.Context, req *models.CreateSubmissionRequest) (*models.Submission, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func TestCreateSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)
		router := gin.Default()
		router.POST("/submissions", handler.CreateSubmission)

		qID := bson.NewObjectID()
		reqBody := models.CreateSubmissionRequest{
			SurveyToken: "abcde",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "Answer"},
			},
		}

		expectedSubmission := &models.Submission{
			SurveyID: bson.NewObjectID(),
		}

		mockService.On("CreateSubmission", mock.Anything, mock.MatchedBy(func(req *models.CreateSubmissionRequest) bool {
			return req.SurveyToken == "abcde"
		})).Return(expectedSubmission, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)
		router := gin.Default()
		router.POST("/submissions", handler.CreateSubmission)

		// Invalid JSON
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBufferString("{invalid"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)
		router := gin.Default()
		router.POST("/submissions", handler.CreateSubmission)

		qID := bson.NewObjectID()
		reqBody := models.CreateSubmissionRequest{
			SurveyToken: "abcde",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "Answer"},
			},
		}

		mockService.On("CreateSubmission", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/submissions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
