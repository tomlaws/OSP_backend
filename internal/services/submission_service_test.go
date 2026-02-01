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

// MockSubmissionRepository is a mock implementation of SubmissionRepository
type MockSubmissionRepository struct {
	mock.Mock
}

func (m *MockSubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

func (m *MockSubmissionRepository) GetAllSubmissions(ctx context.Context, surveyID bson.ObjectID) ([]*models.Submission, error) {
	args := m.Called(ctx, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetSubmissions(ctx context.Context, offset int64, limit int64, surveyID *bson.ObjectID) ([]*models.Submission, error) {
	args := m.Called(ctx, offset, limit, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetBySurveyID(ctx context.Context, surveyID bson.ObjectID) ([]models.Submission, error) {
	args := m.Called(ctx, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateSubmission(t *testing.T) {
	t.Run("SurveyNotFound", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		mockSurveyRepo.On("GetByToken", mock.Anything, "invalid").Return(nil, errors.New("not found"))

		req := &models.CreateSubmissionRequest{SurveyToken: "invalid"}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "Survey not found", err.Error())
	})

	t.Run("InvalidQuestionID", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		surveyID := bson.NewObjectID()
		survey := &models.Survey{ID: surveyID, Questions: []models.Question{}}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: bson.NewObjectID(), Answer: "ans"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid question ID")
	})

	t.Run("Validation_Textbox_Success", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{
						TextboxSpecification: &models.TextboxSpecification{MaxLength: 10},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)
		mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "short"},
			},
		}
		submission, err := service.CreateSubmission(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, submission)
	})

	t.Run("Validation_Textbox_Fail", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{
						TextboxSpecification: &models.TextboxSpecification{MaxLength: 3},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "long answer"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid answer")
	})

	t.Run("Validation_MultipleChoice_Success", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeMultipleChoice,
					Specification: models.QuestionSpecification{
						MultipleChoiceSpecification: &models.MultipleChoiceSpecification{
							Options: []string{"A", "B"},
						},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)
		mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "A"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("Validation_MultipleChoice_Fail", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeMultipleChoice,
					Specification: models.QuestionSpecification{
						MultipleChoiceSpecification: &models.MultipleChoiceSpecification{
							Options: []string{"A", "B"},
						},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "C"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
	})

	t.Run("Validation_Likert_Success", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeLikert,
					Specification: models.QuestionSpecification{
						LikertSpecification: &models.LikertSpecification{
							Min: 1, Max: 5,
						},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)
		mockSubmissionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "3"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("Validation_Likert_Fail", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:   qID,
					Type: models.QuestionTypeLikert,
					Specification: models.QuestionSpecification{
						LikertSpecification: &models.LikertSpecification{
							Min: 1, Max: 5,
						},
					},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)

		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID, Answer: "6"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
	})

	t.Run("MissingResponse", func(t *testing.T) {
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)

		qID1 := bson.NewObjectID()
		qID2 := bson.NewObjectID()
		survey := &models.Survey{
			ID: bson.NewObjectID(),
			Questions: []models.Question{
				{
					ID:            qID1,
					Type:          models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{TextboxSpecification: &models.TextboxSpecification{MaxLength: 100}},
				},
				{
					ID:            qID2,
					Type:          models.QuestionTypeTextbox,
					Specification: models.QuestionSpecification{TextboxSpecification: &models.TextboxSpecification{MaxLength: 100}},
				},
			},
		}
		mockSurveyRepo.On("GetByToken", mock.Anything, "token").Return(survey, nil)

		// Only answering qID1
		req := &models.CreateSubmissionRequest{
			SurveyToken: "token",
			Responses: []models.SubmissionResponse{
				{QuestionID: qID1, Answer: "A"},
			},
		}
		_, err := service.CreateSubmission(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing answer")
	})
}

func TestService_GetSubmissionsBySurveyID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSubmissionRepo := new(MockSubmissionRepository)
		mockSurveyRepo := new(MockSurveyRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)
		surveyID := bson.NewObjectID()
		expectedSubmissions := []*models.Submission{
			{ID: bson.NewObjectID()},
			{ID: bson.NewObjectID()},
		}
		mockSubmissionRepo.On("GetSubmissions", mock.Anything, int64(0), int64(10), &surveyID).Return(expectedSubmissions, nil)
		submissions, err := service.GetSubmissions(context.Background(), 0, 10, &surveyID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSubmissions, submissions)
	})
}

func TestService_DeleteSubmission(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSubmissionRepo := new(MockSubmissionRepository)
		mockSurveyRepo := new(MockSurveyRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)
		submissionID := bson.NewObjectID()
		mockSubmissionRepo.On("Delete", mock.Anything, submissionID).Return(nil)
		err := service.Delete(context.Background(), submissionID)
		assert.NoError(t, err)
	})
	t.Run("RepoError", func(t *testing.T) {
		mockSubmissionRepo := new(MockSubmissionRepository)
		mockSurveyRepo := new(MockSurveyRepository)
		service := NewSubmissionService(mockSubmissionRepo, mockSurveyRepo)
		submissionID := bson.NewObjectID()
		mockSubmissionRepo.On("Delete", mock.Anything, submissionID).Return(errors.New("db error"))
		err := service.Delete(context.Background(), submissionID)
		assert.Error(t, err)
	})
}
