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
)

// MockSurveyService is a mock implementation of ISurveyService
type MockSurveyService struct {
	mock.Mock
}

func (m *MockSurveyService) CreateSurvey(ctx context.Context, req *models.CreateSurveyRequest) (*models.Survey, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyService) GetSurveyByToken(ctx context.Context, token string) (*models.Survey, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func TestCreateSurvey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.POST("/surveys", handler.CreateSurvey)

		reqBody := models.CreateSurveyRequest{
			Name: "Test Survey",
			Questions: []models.QuestionInput{
				{
					Text: "Q1",
					Type: models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{
						TextboxSpecification: &models.TextboxSpecification{
							MaxLength: 100,
						},
					},
				},
			},
		}

		expectedSurvey := &models.Survey{
			Name: "Test Survey",
		}

		mockService.On("CreateSurvey", mock.Anything, mock.MatchedBy(func(req *models.CreateSurveyRequest) bool {
			return req.Name == "Test Survey"
		})).Return(expectedSurvey, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/surveys", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.POST("/surveys", handler.CreateSurvey)

		reqBody := models.CreateSurveyRequest{
			// Name missing
			Questions: []models.QuestionInput{
				{
					Text: "Q1",
					Type: models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{
						TextboxSpecification: &models.TextboxSpecification{
							MaxLength: 100,
						},
					},
				},
			},
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/surveys", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateSurvey")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.POST("/surveys", handler.CreateSurvey)

		reqBody := models.CreateSurveyRequest{
			Name: "Test Survey",
			Questions: []models.QuestionInput{
				{
					Text: "Q1",
					Type: models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{
						TextboxSpecification: &models.TextboxSpecification{
							MaxLength: 100,
						},
					},
				},
			},
		}

		mockService.On("CreateSurvey", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/surveys", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetSurvey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.GET("/surveys/:token", handler.GetSurvey)

		expectedSurvey := &models.Survey{
			Name:  "Test Survey",
			Token: "abcde",
		}

		mockService.On("GetSurveyByToken", mock.Anything, "abcde").Return(expectedSurvey, nil)

		req, _ := http.NewRequest("GET", "/surveys/abcde", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.GET("/surveys/:token", handler.GetSurvey)

		mockService.On("GetSurveyByToken", mock.Anything, "invalid").Return(nil, errors.New("not found"))

		req, _ := http.NewRequest("GET", "/surveys/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
