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
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockInsightService is a mock implementation of IInsightService
type MockInsightService struct {
	mock.Mock
}

func (m *MockInsightService) CreateInsight(ctx context.Context, req *models.CreateInsightRequest) (*models.Insight, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Insight), args.Error(1)
}

func (m *MockInsightService) GetInsights(ctx context.Context, offset, limit int64, surveyID *string) ([]*models.Insight, error) {
	args := m.Called(ctx, offset, limit, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Insight), args.Error(1)
}

func (m *MockInsightService) ProcessInsight(insightID bson.ObjectID) error {
	args := m.Called(insightID)
	return args.Error(0)
}

func (m *MockInsightService) RegisterHandlers(mux *asynq.ServeMux) {
	m.Called(mux)
}

func TestCreateInsight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockInsightService)
		handler := NewInsightHandler(mockService)
		router := gin.Default()
		router.POST("/insights", handler.CreateInsight)

		surveyID := bson.NewObjectID()
		reqBody := models.CreateInsightRequest{
			SurveyID:    surveyID,
			ContextType: models.ContextType("COURSE_FEEDBACK"),
		}

		expectedInsight := &models.Insight{
			ID:       bson.NewObjectID(),
			SurveyID: surveyID,
			Status:   models.InsightPending,
		}

		mockService.On("CreateInsight", mock.Anything, mock.MatchedBy(func(req *models.CreateInsightRequest) bool {
			return req.SurveyID == surveyID
		})).Return(expectedInsight, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/insights", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockInsightService)
		handler := NewInsightHandler(mockService)
		router := gin.Default()
		router.POST("/insights", handler.CreateInsight)

		surveyID := bson.NewObjectID()
		reqBody := models.CreateInsightRequest{
			SurveyID:    surveyID,
			ContextType: models.ContextType("COURSE_FEEDBACK"),
		}

		mockService.On("CreateInsight", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/insights", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetInsights(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockInsightService)
		handler := NewInsightHandler(mockService)
		router := gin.Default()
		router.GET("/insights", handler.GetInsights)

		expectedInsights := []*models.Insight{
			{ID: bson.NewObjectID()},
		}

		mockService.On("GetInsights", mock.Anything, int64(0), int64(10), (*string)(nil)).Return(expectedInsights, nil)

		req, _ := http.NewRequest("GET", "/insights?offset=0&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
