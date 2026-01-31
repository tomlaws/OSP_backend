package services

import (
	"context"
	"errors"
	"fmt"
	"osp/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubmissionService struct {
	collection *mongo.Collection
}

func NewSubmissionService(collection *mongo.Collection) *SubmissionService {
	return &SubmissionService{
		collection: collection,
	}
}

func (s *SubmissionService) CreateSubmission(ctx context.Context, req *models.CreateSubmissionRequest) (*models.Submission, error) {
	var survey models.Survey
	err := s.collection.Database().Collection("surveys").FindOne(ctx, bson.M{"token": req.SurveyToken}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	// Map question ID to answer
	questionMap := make(map[bson.ObjectID]string)
	// Validate questions and responses
	for _, resp := range req.Responses {
		var question *models.Question
		for _, q := range survey.Questions {
			if q.ID == resp.QuestionID {
				question = &q
				break
			}
		}
		if question == nil {
			err := errors.New("invalid question ID: " + resp.QuestionID.Hex())
			return nil, err
		}
		// Validate answer based on question type
		var validAnswer bool = false
		if question.Type == models.QuestionTypeMultipleChoice {
			for _, option := range question.Specification.Options {
				if option == resp.Answer {
					validAnswer = true
					break
				}
			}
		}
		if question.Type == models.QuestionTypeLikert {
			num := 0
			_, err := fmt.Sscanf(resp.Answer, "%d", &num)
			if err == nil && num >= question.Specification.Min && num <= question.Specification.Max {
				validAnswer = true
			}
		}
		if question.Type == models.QuestionTypeTextbox {
			if len(resp.Answer) <= question.Specification.MaxLength {
				validAnswer = true
			}
		}
		if !validAnswer {
			err := errors.New("invalid answer for question ID: " + resp.QuestionID.Hex())
			return nil, err
		}
		questionMap[resp.QuestionID] = resp.Answer
	}
	validatedResponses := make([]models.SubmissionResponse, 0, len(survey.Questions))
	for _, question := range survey.Questions {
		answer, ok := questionMap[question.ID]
		if !ok {
			err := errors.New("missing answer for question ID: " + question.ID.Hex())
			return nil, err
		}
		validatedResponses = append(validatedResponses, models.SubmissionResponse{
			QuestionID: question.ID,
			Answer:     answer,
		})
	}
	submission := &models.Submission{
		ID:        bson.NewObjectID(),
		SurveyID:  survey.ID,
		Responses: validatedResponses,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.collection.InsertOne(ctx, submission)
	if err != nil {
		return nil, err
	}
	return submission, nil
}
