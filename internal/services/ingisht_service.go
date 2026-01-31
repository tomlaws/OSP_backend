package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"osp/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InsightService struct {
	collection *mongo.Collection
}

func NewInsightService(collection *mongo.Collection) *InsightService {
	insightService := &InsightService{
		collection: collection,
	}
	go insightService.startInsightProcessingWorker()
	return insightService
}

func (s *InsightService) startInsightProcessingWorker() {
	// Mark all PENDING insights as FAILED from previous runs
	_, err := s.collection.UpdateMany(
		context.TODO(),
		bson.M{"status": models.InsightPending},
		bson.M{"$set": bson.M{"status": models.InsightFailed}},
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			{
				Key: "$match",
				Value: bson.D{
					{Key: "operationType", Value: "insert"},
					{Key: "fullDocument.status", Value: "PENDING"},
				},
			},
		}}
	stream, err := s.collection.Watch(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close(ctx)
	fmt.Println("Started insight processing watcher...")
	for stream.Next(ctx) {
		var event bson.M
		if err := stream.Decode(&event); err != nil {
			panic(err)
		}
		output, err := json.MarshalIndent(event["fullDocument"], "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
		var insight models.Insight
		bsonBytes, _ := bson.Marshal(event["fullDocument"])
		bson.Unmarshal(bsonBytes, &insight)
		fmt.Printf("Processing insight ID: %s\n", insight.ID.Hex())
		err = s.ProcessInsight(insight)
	}
	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}
}

func (s *InsightService) CreateInsight(ctx context.Context, req *models.CreateInsightRequest) (*models.Insight, error) {
	insight := &models.Insight{
		ID:          bson.NewObjectID(),
		SurveyID:    req.SurveyID,
		ContextType: req.ContextType,
		Status:      models.InsightPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.preprocessInsight(insight)
	_, err := s.collection.InsertOne(ctx, insight)
	if err != nil {
		return nil, err
	}
	err = s.collection.FindOne(ctx, bson.M{"_id": insight.ID}).Decode(&insight)
	if err != nil {
		return nil, err
	}
	return insight, nil
}

// Retrieve insights by using offset and limit for pagination
func (s *InsightService) GetInsights(ctx context.Context, offset, limit int64, surveyID *string) ([]*models.Insight, error) {
	filter := bson.M{}
	if surveyID != nil {
		bsonSurveyID, err := bson.ObjectIDFromHex(*surveyID)
		if err != nil {
			return nil, err
		}
		filter["survey_id"] = bsonSurveyID
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []*models.Insight
	for cursor.Next(ctx) {
		var insight models.Insight
		if err := cursor.Decode(&insight); err != nil {
			return nil, err
		}
		insights = append(insights, &insight)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return insights, nil
}

func (s *InsightService) preprocessInsight(insight *models.Insight) error {
	// Get all submissions for the survey
	var submissions []models.Submission
	cursor, err := s.collection.Database().Collection("submissions").Find(context.TODO(), bson.M{"survey_id": insight.SurveyID})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var submission models.Submission
		if err := cursor.Decode(&submission); err != nil {
			insight.Status = models.InsightFailed
			s.collection.UpdateByID(context.TODO(), insight.ID, insight)
			return err
		}
		submissions = append(submissions, submission)
	}
	if err := cursor.Err(); err != nil {
		insight.Status = models.InsightFailed
		s.collection.UpdateByID(context.TODO(), insight.ID, insight)
		return err
	}
	// Get the survey
	var survey models.Survey
	err = s.collection.Database().Collection("surveys").FindOne(context.TODO(), bson.M{"_id": insight.SurveyID}).Decode(&survey)
	if err != nil {
		insight.Status = models.InsightFailed
		s.collection.UpdateByID(context.TODO(), insight.ID, insight)
		return err
	}
	// Build map of question ID to question
	questionMap := make(map[bson.ObjectID]models.Question)
	for _, question := range survey.Questions {
		questionMap[question.ID] = question
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
		currentBatch := models.InsightBatch{
			BatchNumber:      len(insightBatches) + 1,
			Question:         question,
			AggregatedAnswer: make(map[string]int),
			TextualAnswers:   &[]string{},
		}
		insightBatches = append(insightBatches, currentBatch)
		answers := responseMap[question.ID]
		textLength := 0
		for _, answer := range answers {
			switch question.Type {
			case models.QuestionTypeMultipleChoice:
				currentBatch.AggregatedAnswer[answer]++
			case models.QuestionTypeLikert:
				currentBatch.AggregatedAnswer[answer]++
			default:
				*currentBatch.TextualAnswers = append(*currentBatch.TextualAnswers, answer)
				textLength += len(answer)
				if textLength > 4000 {
					// Start a new batch for this question
					currentBatch = models.InsightBatch{
						BatchNumber:      len(insightBatches) + 1,
						Question:         question,
						AggregatedAnswer: make(map[string]int),
						TextualAnswers:   &[]string{},
					}
					insightBatches = append(insightBatches, currentBatch)
					textLength = 0
				}
			}
		}
	}
	insight.Batches = insightBatches
	return nil
}

func (s *InsightService) ProcessInsight(insight models.Insight) error {
	// Update insight status
	update := bson.M{
		"$set": bson.M{
			"status":     models.InsightProcessing,
			"updated_at": time.Now(),
		},
	}
	_, err := s.collection.UpdateByID(context.TODO(), insight.ID, update)
	if err != nil {
		return err
	}
	for i, batch := range insight.Batches {
		if batch.Summary != nil {
			// Already processed
			continue
		}
		summary, err := s.processInsightBatch(insight.ContextType, batch)
		insight.Batches[i].Summary = &summary
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
		_, err = s.collection.UpdateByID(context.TODO(), insight.ID, update)
		if err != nil {
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
		finalUpdate := bson.M{
			"$set": bson.M{
				"status":       models.InsightCompleted,
				"completed_at": time.Now(),
				"updated_at":   time.Now(),
			},
		}
		_, err = s.collection.UpdateByID(context.TODO(), insight.ID, finalUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *InsightService) processInsightBatch(contextType models.ContextType, batch models.InsightBatch) (string, error) {
	// Build prompt
	prompt := "Summarize the following answers to the question: " + batch.Question.Text + "\n\n"
	for _, answer := range *batch.TextualAnswers {
		prompt += "- " + answer + "\n"
	}
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
		payload = fmt.Sprintf("Likert scale distribution: %s", string(payloadBytes))
	default:
		payload = "No valid answers provided."
	}

	reqBody := models.ChatCompletionRequest{
		Messages: []models.ChatCompletionMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a helpful assistant. Analyze survey responses in the context of %s.", contextType),
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Batch %d:\nQuestion: %s\n%s\n\nPlease generate a concise summary.", batch.BatchNumber, batch.Question.Text, payload),
			},
		},
		Temperature: 0.7,
		TopP:        1.0,
		MaxTokens:   500,
		Model:       "openai/gpt-4o-mini",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://models.github.ai/inference/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var chatCompletionResponse models.ChatCompletionResponse
	if err := json.Unmarshal(body, &chatCompletionResponse); err != nil {
		return "", err
	}

	if len(chatCompletionResponse.Choices) > 0 {
		return chatCompletionResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no summary returned")
}
