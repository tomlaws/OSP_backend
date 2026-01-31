package services

import (
	"context"
	"fmt"
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
	return &InsightService{
		collection: collection,
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
	_, err := s.collection.InsertOne(ctx, insight)
	if err != nil {
		return nil, err
	}
	err = s.initializeInsightBatches(insight.ID)
	if err != nil {
		return nil, err
	}
	err = s.processInsightBatches(insight.ID)
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

func (s *InsightService) initializeInsightBatches(insightID bson.ObjectID) error {
	// Get insight
	var insight models.Insight
	err := s.collection.FindOne(context.TODO(), bson.M{"_id": insightID}).Decode(&insight)
	if err != nil {
		return err
	}
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
	answerBatches := []*models.InsightBatch{}
	for _, question := range survey.Questions {
		currentBatch := &models.InsightBatch{
			BatchNumber:      len(answerBatches) + 1,
			Question:         question,
			AggregatedAnswer: make(map[string]int),
			TextualAnswers:   []string{},
		}
		answerBatches = append(answerBatches, currentBatch)
		answers := responseMap[question.ID]
		textLength := 0
		for _, answer := range answers {
			switch question.Type {
			case models.QuestionTypeMultipleChoice:
				currentBatch.AggregatedAnswer[answer]++
			case models.QuestionTypeLikert:
				currentBatch.AggregatedAnswer[answer]++
			default:
				textLength += len(answer)
				if textLength > 4000 {
					// Start a new batch for this question
					currentBatch = &models.InsightBatch{
						BatchNumber:      len(answerBatches) + 1,
						Question:         question,
						AggregatedAnswer: make(map[string]int),
						TextualAnswers:   []string{},
					}
					answerBatches = append(answerBatches, currentBatch)
					textLength = 0
				}
				currentBatch.TextualAnswers = append(currentBatch.TextualAnswers, answer)
			}
		}
	}
	// print number of batches
	update := bson.M{
		"$set": bson.M{
			"batches": answerBatches,
			"status":  models.InsightProcessing,
		},
	}
	_, err = s.collection.UpdateByID(context.TODO(), insight.ID, update)
	if err != nil {
		fmt.Println("Failed to update insight with batches:", err)
		return err
	}
	updatedInsight := models.Insight{}
	err = s.collection.FindOne(context.TODO(), bson.M{"_id": insight.ID}).Decode(&updatedInsight)
	return err
}

func (s *InsightService) processInsightBatches(insightID bson.ObjectID) error {
	// Update insight status
	insight := &models.Insight{}
	err := s.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": insightID}, bson.M{
		"$set": bson.M{
			"status":     models.InsightProcessing,
			"updated_at": time.Now(),
		},
	}).Decode(insight)
	if err != nil {
		return err
	}
	// Process each batch using LLM
	cursor, err := s.collection.Find(context.TODO(), bson.M{"status": models.InsightProcessing})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var insight models.Insight
		if err := cursor.Decode(&insight); err != nil {
			fmt.Println("Error decoding insight:", err)
			continue
		}
		// Process each batch
		for i, batch := range insight.Batches {
			if batch.Summary != nil {
				// Already processed
				continue
			}
			summary, err := s.processInsightBatch(batch)
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

func (s *InsightService) processInsightBatch(batch models.InsightBatch) (string, error) {
	// Build prompt
	prompt := "Summarize the following answers to the question: " + batch.Question.Text + "\n\n"
	for _, answer := range batch.TextualAnswers {
		prompt += "- " + answer + "\n"
	}
	// LLM processing
	summary := "This is a simulated summary for the question: " + batch.Question.Text
	return summary, nil
}
