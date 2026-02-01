package services

import (
	"context"
	"errors"
	"testing"

	"osp/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockSurveyRepository is a mock implementation of SurveyRepository
type MockSurveyRepository struct {
	mock.Mock
}

func (m *MockSurveyRepository) Create(ctx context.Context, survey *models.Survey) error {
	args := m.Called(ctx, survey)
	return args.Error(0)
}

func (m *MockSurveyRepository) List(ctx context.Context, offset, limit int64) ([]*models.Survey, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Survey), args.Error(1)
}

func (m *MockSurveyRepository) GetByToken(ctx context.Context, token string) (*models.Survey, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.Survey, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockSurveyRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

func TestService_GetSurveyByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)

		expectedSurvey := &models.Survey{ID: bson.NewObjectID()}
		mockRepo.On("GetByID", mock.Anything, expectedSurvey.ID).Return(expectedSurvey, nil)

		survey, err := service.GetSurveyByID(context.Background(), expectedSurvey.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedSurvey, survey)
	})
}

func TestService_ListSurveys(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)
		expectedSurveys := []*models.Survey{
			{ID: bson.NewObjectID()},
			{ID: bson.NewObjectID()},
		}
		mockRepo.On("List", mock.Anything, int64(0), int64(10)).Return(expectedSurveys, nil)

		surveys, err := service.ListSurveys(context.Background(), 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedSurveys, surveys)
	})
}

func TestService_DeleteSurvey(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)
		surveyID := bson.NewObjectID()

		mockRepo.On("Delete", mock.Anything, surveyID).Return(nil)

		err := service.DeleteSurvey(context.Background(), surveyID)
		assert.NoError(t, err)
	})
	t.Run("RepoError", func(t *testing.T) {
		mockRepo := new(MockSurveyRepository)
		service := NewSurveyService(mockRepo)
		surveyID := bson.NewObjectID()
		mockRepo.On("Delete", mock.Anything, surveyID).Return(errors.New("db error"))

		err := service.DeleteSurvey(context.Background(), surveyID)
		assert.Error(t, err)
	})
}
