package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"osp/internal/models"
	"osp/internal/repositories"
	"time"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IInsightService interface {
	CreateInsight(ctx context.Context, req *models.CreateInsightRequest) (*models.Insight, error)
	GetInsights(ctx context.Context, offset, limit int64, surveyID *bson.ObjectID) ([]*models.Insight, error)
	GetInsight(ctx context.Context, id bson.ObjectID) (*models.Insight, error)
	ProcessInsight(insightID bson.ObjectID) error
	RegisterHandlers(mux *asynq.ServeMux)
}

type JobEnqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type InsightService struct {
	insightRepo           repositories.InsightRepository
	surveyRepo            repositories.SurveyRepository
	submissionRepo        repositories.SubmissionRepository
	chatCompletionService IChatCompletionService
	jobEnqueuer           JobEnqueuer
}

func NewInsightService(
	insightRepo repositories.InsightRepository,
	surveyRepo repositories.SurveyRepository,
	submissionRepo repositories.SubmissionRepository,
	chatCompletionService IChatCompletionService,
	jobEnqueuer JobEnqueuer,
) *InsightService {
	return &InsightService{
		insightRepo:           insightRepo,
		surveyRepo:            surveyRepo,
		submissionRepo:        submissionRepo,
		chatCompletionService: chatCompletionService,
		jobEnqueuer:           jobEnqueuer,
	}
}

func (s *InsightService) RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeProcessInsight, func(ctx context.Context, task *asynq.Task) error {
		insightID, err := parseProcessInsightPayload(task)
		if err != nil {
			return err
		}
		log.Printf("asynq: processing insight %s", insightID.Hex())
		return s.ProcessInsight(insightID)
	})
}

func (s *InsightService) CreateInsight(ctx context.Context, req *models.CreateInsightRequest) (*models.Insight, error) {
	// Check survey existence
	_, err := s.surveyRepo.GetByID(ctx, req.SurveyID)
	if err != nil {
		return nil, fmt.Errorf("survey not found")
	}

	insight := &models.Insight{
		ID:          bson.NewObjectID(),
		SurveyID:    req.SurveyID,
		ContextType: req.ContextType,
		Status:      models.InsightPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Preprocess (load data)
	if err := s.preprocessInsight(ctx, insight); err != nil {
		return nil, err
	}

	err = s.insightRepo.Create(ctx, insight)
	if err != nil {
		return nil, err
	}

	// Enqueue background processing.
	if s.jobEnqueuer != nil {
		task, err := newProcessInsightTask(insight.ID)
		if err != nil {
			return nil, err
		}
		if _, err := s.jobEnqueuer.Enqueue(task, asynq.Queue("insights"), asynq.MaxRetry(10)); err != nil {
			return nil, err
		}
	}

	return s.insightRepo.GetByID(ctx, insight.ID)
}

// Retrieve insights by using offset and limit for pagination
func (s *InsightService) GetInsights(ctx context.Context, offset, limit int64, surveyID *bson.ObjectID) ([]*models.Insight, error) {
	return s.insightRepo.GetInsights(ctx, offset, limit, surveyID)
}

func (s *InsightService) GetInsight(ctx context.Context, id bson.ObjectID) (*models.Insight, error) {
	return s.insightRepo.GetByID(ctx, id)
}

func (s *InsightService) preprocessInsight(ctx context.Context, insight *models.Insight) error {
	// Get all submissions for the survey
	submissions, err := s.submissionRepo.GetAllSubmissions(ctx, insight.SurveyID)
	if err != nil {
		return err
	}

	// Get the survey
	survey, err := s.surveyRepo.GetByID(ctx, insight.SurveyID)
	if err != nil {
		return err
	}

	// Build map of question ID to responses
	responseMap := make(map[bson.ObjectID][]string)
	for _, submission := range submissions {
		for _, response := range submission.Responses {
			responseMap[response.QuestionID] = append(responseMap[response.QuestionID], response.Answer)
		}
	}

	// Build answer batches
	insightBatches := []models.InsightBatch{}
	for _, question := range survey.Questions {
		insightBatch := &models.InsightBatch{
			BatchNumber: len(insightBatches) + 1,
			Question:    question,
		}
		if question.Type == models.QuestionTypeMultipleChoice || question.Type == models.QuestionTypeLikert {
			aggMap := make(map[string]int)
			insightBatch.AggregatedAnswer = &aggMap
		} else {
			insightBatch.TextualAnswers = &[]string{}
		}
		insightBatches = append(insightBatches, *insightBatch)
		currentBatch := &insightBatches[len(insightBatches)-1]

		answers := responseMap[question.ID]
		textLength := 0
		for _, answer := range answers {
			switch question.Type {
			case models.QuestionTypeMultipleChoice:
				(*currentBatch.AggregatedAnswer)[answer]++
			case models.QuestionTypeLikert:
				(*currentBatch.AggregatedAnswer)[answer]++
			default:
				// If this answer would exceed the limit, start a new batch BEFORE appending.
				if textLength > 0 && textLength+len(answer) > 4000 {
					newAggMap := make(map[string]int)
					insightBatches = append(insightBatches, models.InsightBatch{
						BatchNumber:      len(insightBatches) + 1,
						Question:         question,
						AggregatedAnswer: &newAggMap,
						TextualAnswers:   &[]string{},
					})
					currentBatch = &insightBatches[len(insightBatches)-1]
					textLength = 0
				}

				*currentBatch.TextualAnswers = append(*currentBatch.TextualAnswers, answer)
				textLength += len(answer)
			}
		}
	}
	insight.Batches = insightBatches
	return nil
}

func (s *InsightService) ProcessInsight(insightID bson.ObjectID) error {
	ctx := context.TODO() // Background context for async task

	// Retrieve the insight
	insight, err := s.insightRepo.GetByID(ctx, insightID)
	if err != nil {
		return err
	}

	// Update insight status
	update := bson.M{
		"$set": bson.M{
			"status":     models.InsightProcessing,
			"updated_at": time.Now(),
		},
	}
	if err := s.insightRepo.Update(ctx, insight.ID, update); err != nil {
		return err
	}

	for i, batch := range insight.Batches {
		if batch.Summary != nil {
			// Already processed
			continue
		}
		summary, err := s.processInsightBatch(insight.ID, insight.ContextType, batch)
		insight.Batches[i].Summary = summary
		if err != nil {
			errMsg := err.Error()
			insight.Batches[i].ErrorLog = &errMsg
		}
		// Update the insight in the database
		update := bson.M{
			"$set": bson.M{
				"batches":    insight.Batches,
				"updated_at": time.Now(),
			},
		}
		if err := s.insightRepo.Update(ctx, insight.ID, update); err != nil {
			fmt.Println("Failed to update insight batch:", err)
			continue
		}
	}
	// Check if all batches are processed and update insight status
	allProcessed := true
	for _, batch := range insight.Batches {
		if batch.Summary == nil && batch.ErrorLog == nil {
			allProcessed = false
			break
		}
	}
	if allProcessed {
		// Meta-summary (overall analysis) after all batches are processed.
		analysis, analysisErr := s.generateMetaSummary(insight)
		if analysisErr != nil {
			errMsg := analysisErr.Error()
			update := bson.M{
				"$set": bson.M{
					"analysis":     "",
					"status":       models.InsightFailed,
					"updated_at":   time.Now(),
					"completed_at": time.Now(),
					"error_log":    errMsg,
				},
			}
			_ = s.insightRepo.Update(ctx, insight.ID, update)
			return analysisErr
		}

		finalUpdate := bson.M{
			"$set": bson.M{
				"analysis":     analysis,
				"status":       models.InsightCompleted,
				"completed_at": time.Now(),
				"updated_at":   time.Now(),
			},
		}
		if err := s.insightRepo.Update(ctx, insight.ID, finalUpdate); err != nil {
			return err
		}
	}
	return nil
}

func (s *InsightService) processInsightBatch(insightID bson.ObjectID, contextType models.ContextType, batch models.InsightBatch) (*string, error) {
	// LLM processing
	var payload string

	switch batch.Question.Type {
	case "TEXTBOX":
		payload = fmt.Sprintf("Answers: %v", batch.TextualAnswers)
	case "MULTIPLE_CHOICE":
		payloadBytes, _ := json.Marshal(batch.AggregatedAnswer)
		payload = fmt.Sprintf("Aggregated answers: %s", string(payloadBytes))
	case "LIKERT":
		payloadBytes, _ := json.Marshal(batch.AggregatedAnswer)
		payload = fmt.Sprintf("Aggregated answers: %s", string(payloadBytes))
	}

	reqBody := models.ChatCompletionRequest{
		Messages: []models.ChatCompletionMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a helpful assistant. Summarize the following survey responses in the context of %s.", contextType),
			},
			{
				Role:    "user",
				Content: payload,
			},
		},
		Temperature: 0.5,
		TopP:        1.0,
		MaxTokens:   800,
		Model:       "openai/gpt-4o-mini",
	}

	ref := fmt.Sprintf("insight:%s batch:%d", insightID.Hex(), batch.BatchNumber)
	resp, err := s.chatCompletionService.NewRequest(reqBody, &ref)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *InsightService) generateMetaSummary(insight *models.Insight) (string, error) {
	meta := "Here are the summaries of different batches of answers:\n"
	for _, batch := range insight.Batches {
		if batch.Summary != nil {
			meta += fmt.Sprintf("Batch %d (Question: %s): %s\n", batch.BatchNumber, batch.Question.Text, *batch.Summary)
		}
		if batch.ErrorLog != nil {
			meta += fmt.Sprintf("Error: %s\n", *batch.ErrorLog)
		}
	}

	reqBody := models.ChatCompletionRequest{
		Messages: []models.ChatCompletionMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a helpful assistant. Analyze survey responses in the context of %s.", insight.ContextType),
			},
			{
				Role:    "user",
				Content: meta,
			},
		},
		Temperature: 0.5,
		TopP:        1.0,
		MaxTokens:   800,
		Model:       "openai/gpt-4o-mini",
	}

	ref := fmt.Sprintf("insight:%s meta", insight.ID.Hex())
	resp, err := s.chatCompletionService.NewRequest(reqBody, &ref)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("empty response")
	}
	return *resp, nil
}

// Asynq task definitions
const TypeProcessInsight = "insight:process"

type ProcessInsightPayload struct {
	InsightID string `json:"insight_id"`
}

func newProcessInsightTask(insightID bson.ObjectID) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessInsightPayload{InsightID: insightID.Hex()})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProcessInsight, payload), nil
}

func parseProcessInsightPayload(task *asynq.Task) (bson.ObjectID, error) {
	var payload ProcessInsightPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return bson.ObjectID{}, err
	}
	if payload.InsightID == "" {
		return bson.ObjectID{}, fmt.Errorf("missing insight_id")
	}
	return bson.ObjectIDFromHex(payload.InsightID)
}
