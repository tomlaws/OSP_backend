package services

import (
	"context"
	"errors"
	"testing"

	"osp/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSurveyRepository is a mock implementation of SurveyRepository
type MockSurveyRepository struct {
	mock.Mock
}

func (m *MockSurveyRepository) Create(ctx context.Context, survey *models.Survey) error {
	args := m.Called(ctx, survey)
	return args.Error(0)
}

func (m *MockSurveyRepository) GetByToken(ctx context.Context, token string) (*models.Survey, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyRepository) GetByID(ctx context.Context, id interface{}) (*models.Survey, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func TestService_CreateSurvey(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)

		req := &models.CreateSurveyRequest{
			Name: "Test Survey",
			Questions: []models.QuestionInput{
				{
					Text: "Q1",
					Type: models.QuestionTypeTextbox,
				},
			},
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *models.Survey) bool {
			return s.Name == "Test Survey" && len(s.Questions) == 1 && s.Token != ""
		})).Return(nil)

		survey, err := service.CreateSurvey(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, survey)
		assert.Equal(t, "Test Survey", survey.Name)
		assert.NotEmpty(t, survey.Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)

		req := &models.CreateSurveyRequest{Name: "Test"}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		survey, err := service.CreateSurvey(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, survey)
	})
}

func TestService_GetSurveyByToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)

		expectedSurvey := &models.Survey{Token: "abc"}
		mockRepo.On("GetByToken", mock.Anything, "abc").Return(expectedSurvey, nil)

		survey, err := service.GetSurveyByToken(context.Background(), "abc")

		assert.NoError(t, err)
		assert.Equal(t, expectedSurvey, survey)
	})
}
