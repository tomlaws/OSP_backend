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

func (m *MockSurveyService) ListSurveys(ctx context.Context, offset, limit int64) ([]*models.Survey, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Survey), args.Error(1)
}

func (m *MockSurveyService) GetSurveyByToken(ctx context.Context, token string) (*models.Survey, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyService) GetSurveyByID(ctx context.Context, id bson.ObjectID) (*models.Survey, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyService) DeleteSurvey(ctx context.Context, id bson.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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
		router.GET("/surveys/:token", handler.GetSurveyByToken)

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
		router.GET("/surveys/:token", handler.GetSurveyByToken)
		mockService.On("GetSurveyByToken", mock.Anything, "invalid").Return(nil, errors.New("not found"))

		req, _ := http.NewRequest("GET", "/surveys/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListSurveys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.GET("/surveys", handler.ListSurveys)
		expectedSurveys := []*models.Survey{
			{Name: "Survey1"},
			{Name: "Survey2"},
		}
		mockService.On("ListSurveys", mock.Anything, int64(0), int64(10)).Return(expectedSurveys, nil)
		req, _ := http.NewRequest("GET", "/surveys?offset=0&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.GET("/surveys", handler.ListSurveys)
		mockService.On("ListSurveys", mock.Anything, int64(0), int64(10)).Return(nil, errors.New("db error"))
		req, _ := http.NewRequest("GET", "/surveys?offset=0&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteSurvey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.DELETE("/surveys/:id", handler.DeleteSurvey)
		surveyID := bson.NewObjectID()
		mockService.On("DeleteSurvey", mock.Anything, surveyID).Return(nil)
		req, _ := http.NewRequest("DELETE", "/surveys/"+surveyID.Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockSurveyService)
		handler := NewSurveyHandler(mockService)
		router := gin.Default()
		router.DELETE("/surveys/:id", handler.DeleteSurvey)
		surveyID := bson.NewObjectID()
		mockService.On("DeleteSurvey", mock.Anything, surveyID).Return(errors.New("db error"))
		req, _ := http.NewRequest("DELETE", "/surveys/"+surveyID.Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
