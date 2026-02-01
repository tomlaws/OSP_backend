package services

import (
	"context"
	"errors"
	"fmt"
	"osp/internal/models"
	"osp/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ISubmissionService interface {
	CreateSubmission(ctx context.Context, req *models.CreateSubmissionRequest) (*models.Submission, error)
	GetSubmissions(ctx context.Context, offset int64, limit int64, surveyID *bson.ObjectID) ([]*models.Submission, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

type SubmissionService struct {
	submissionRepo repositories.SubmissionRepository
	surveyRepo     repositories.SurveyRepository
}

func NewSubmissionService(submissionRepo repositories.SubmissionRepository, surveyRepo repositories.SurveyRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: submissionRepo,
		surveyRepo:     surveyRepo,
	}
}

func (s *SubmissionService) CreateSubmission(ctx context.Context, req *models.CreateSubmissionRequest) (*models.Submission, error) {
	survey, err := s.surveyRepo.GetByToken(ctx, req.SurveyToken)
	if err != nil {
		return nil, errors.New("Survey not found")
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
	err = s.submissionRepo.Create(ctx, submission)
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (s *SubmissionService) GetSubmissions(ctx context.Context, offset int64, limit int64, surveyID *bson.ObjectID) ([]*models.Submission, error) {
	return s.submissionRepo.GetSubmissions(ctx, offset, limit, surveyID)
}

func (s *SubmissionService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.submissionRepo.Delete(ctx, id)
}
