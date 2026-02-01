package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"osp/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InsightService struct {
	collection            *mongo.Collection
	chatCompletionService *ChatCompletionService
}

func NewInsightService(collection *mongo.Collection, chatCompletionService *ChatCompletionService) *InsightService {
	insightService := &InsightService{
		collection:            collection,
		chatCompletionService: chatCompletionService,
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
		summary, err := s.processInsightBatch(context.TODO(), insight.ID, insight.ContextType, batch)
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
			_, _ = s.collection.UpdateByID(context.TODO(), insight.ID, update)
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
		_, err = s.collection.UpdateByID(context.TODO(), insight.ID, finalUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *InsightService) processInsightBatch(ctx context.Context, insightID bson.ObjectID, contextType models.ContextType, batch models.InsightBatch) (*string, error) {
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

	ref := fmt.Sprintf("insight:%s batch:%d", insightID.Hex(), batch.BatchNumber)
	resp, err := s.chatCompletionService.NewRequest(reqBody, &ref)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *InsightService) generateMetaSummary(insight models.Insight) (*string, error) {
	meta := "Create an overall analysis across all questions.\n"
	meta += "Provide:\n"
	meta += "- key themes\n- strongest signals\n- notable outliers\n- actionable recommendations\n"
	meta += "\nBatches:\n"

	for _, batch := range insight.Batches {
		meta += fmt.Sprintf("\nBatch %d\nQuestion: %s\nType: %s\n", batch.BatchNumber, batch.Question.Text, batch.Question.Type)
		if batch.Summary != nil {
			meta += fmt.Sprintf("Summary: %s\n", *batch.Summary)
		}
		if batch.AggregatedAnswer != nil && len(*batch.AggregatedAnswer) > 0 {
			payloadBytes, _ := json.Marshal(batch.AggregatedAnswer)
			meta += fmt.Sprintf("Aggregated: %s\n", string(payloadBytes))
		}
		if batch.TextualAnswers != nil && len(*batch.TextualAnswers) > 0 {
			meta += fmt.Sprintf("Text answers count: %d\n", len(*batch.TextualAnswers))
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
		return nil, err
	}
	return resp, nil
}
