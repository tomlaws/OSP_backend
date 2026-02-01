package services

import (
	"context"
	"testing"

	"osp/internal/models"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockInsightRepository struct {
	mock.Mock
}

func (m *MockInsightRepository) Create(ctx context.Context, insight *models.Insight) error {
	args := m.Called(ctx, insight)
	return args.Error(0)
}
func (m *MockInsightRepository) GetByID(ctx context.Context, id interface{}) (*models.Insight, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Insight), args.Error(1)
}
func (m *MockInsightRepository) GetInsights(ctx context.Context, offset, limit int64, surveyID *string) ([]*models.Insight, error) {
	args := m.Called(ctx, offset, limit, surveyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Insight), args.Error(1)
}
func (m *MockInsightRepository) Update(ctx context.Context, id interface{}, update interface{}) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

type MockChatCompletionService struct {
	mock.Mock
}

func (m *MockChatCompletionService) NewRequest(reqBody models.ChatCompletionRequest, reference *string) (*string, error) {
	args := m.Called(reqBody, reference)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

type MockJobEnqueuer struct {
	mock.Mock
}

func (m *MockJobEnqueuer) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	// For variadic functions, we can just pass them as is to Called
	// But since opts is a slice here, we pass it as a single argument
	args := m.Called(task, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*asynq.TaskInfo), args.Error(1)
}

func TestService_CreateInsight(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockInsightRepo := new(MockInsightRepository)
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		mockChat := new(MockChatCompletionService)
		mockEnqueuer := new(MockJobEnqueuer)

		service := NewInsightService(mockInsightRepo, mockSurveyRepo, mockSubmissionRepo, mockChat, mockEnqueuer)

		surveyID := bson.NewObjectID()
		survey := &models.Survey{
			ID:        surveyID,
			Questions: []models.Question{},
		}

		submissions := []models.Submission{}

		mockSurveyRepo.On("GetByID", mock.Anything, surveyID).Return(survey, nil)
		mockSubmissionRepo.On("GetBySurveyID", mock.Anything, surveyID).Return(submissions, nil)

		mockInsightRepo.On("Create", mock.Anything, mock.MatchedBy(func(i *models.Insight) bool {
			return i.SurveyID == surveyID
		})).Return(nil)

		// Expect Enqueue with any task and any opts
		mockEnqueuer.On("Enqueue", mock.Anything, mock.Anything).Return(&asynq.TaskInfo{}, nil)

		// CreateInsight calls GetByID at the end
		mockInsightRepo.On("GetByID", mock.Anything, mock.Anything).Return(&models.Insight{ID: bson.NewObjectID(), SurveyID: surveyID}, nil)

		req := &models.CreateInsightRequest{SurveyID: surveyID}
		insight, err := service.CreateInsight(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, insight)
		mockInsightRepo.AssertExpectations(t)
		mockEnqueuer.AssertExpectations(t)
	})

	t.Run("ProcessInsight_Success", func(t *testing.T) {
		mockInsightRepo := new(MockInsightRepository)
		mockSurveyRepo := new(MockSurveyRepository)
		mockSubmissionRepo := new(MockSubmissionRepository)
		mockChat := new(MockChatCompletionService)
		mockEnqueuer := new(MockJobEnqueuer)

		service := NewInsightService(mockInsightRepo, mockSurveyRepo, mockSubmissionRepo, mockChat, mockEnqueuer)

		insightID := bson.NewObjectID()
		insight := &models.Insight{
			ID:          insightID,
			SurveyID:    bson.NewObjectID(),
			ContextType: models.ContextType("COURSE_FEEDBACK"),
			Batches: []models.InsightBatch{
				{
					BatchNumber:    1,
					Question:       models.Question{Type: models.QuestionTypeTextbox, Text: "Q1"},
					TextualAnswers: &[]string{"A1", "A2"},
				},
			},
		}

		mockInsightRepo.On("GetByID", mock.Anything, insightID).Return(insight, nil)

		// 1. Update status to Processing
		mockInsightRepo.On("Update", mock.Anything, insightID, mock.MatchedBy(func(u interface{}) bool {
			m, ok := u.(bson.M)
			if !ok {
				return false
			}
			set, ok := m["$set"].(bson.M)
			if !ok {
				return false
			}
			return set["status"] == models.InsightProcessing
		})).Return(nil)

		// 2. Chat completion for batch
		summary := "Summary 1"
		mockChat.On("NewRequest", mock.Anything, mock.Anything).Return(&summary, nil).Once()

		// 3. Update batch with summary
		mockInsightRepo.On("Update", mock.Anything, insightID, mock.MatchedBy(func(u interface{}) bool {
			m, ok := u.(bson.M)
			if !ok {
				return false
			}
			set, ok := m["$set"].(bson.M)
			if !ok {
				return false
			}
			_, hasBatches := set["batches"]
			return hasBatches
		})).Return(nil)

		// 4. Chat completion for Meta Summary
		metaAnalysis := "Meta Analysis"
		mockChat.On("NewRequest", mock.Anything, mock.Anything).Return(&metaAnalysis, nil).Once()

		// 5. Final update
		mockInsightRepo.On("Update", mock.Anything, insightID, mock.MatchedBy(func(u interface{}) bool {
			m, ok := u.(bson.M)
			if !ok {
				return false
			}
			set, ok := m["$set"].(bson.M)
			if !ok {
				return false
			}
			return set["status"] == models.InsightCompleted && set["analysis"] == "Meta Analysis"
		})).Return(nil)

		err := service.ProcessInsight(insightID)
		assert.NoError(t, err)
		mockInsightRepo.AssertExpectations(t)
		mockChat.AssertExpectations(t)
	})
}
