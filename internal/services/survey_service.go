package services

import (
	"context"
	"math/rand"
	"osp/internal/models"
	"osp/internal/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ISurveyService defines the business logic for survey operations
type ISurveyService interface {
	CreateSurvey(ctx context.Context, req *models.CreateSurveyRequest) (*models.Survey, error)
	GetSurveyByToken(ctx context.Context, token string) (*models.Survey, error)
}

type SurveyService struct {
	repo repositories.SurveyRepository
}

func NewSurveyService(repo repositories.SurveyRepository) *SurveyService {
	return &SurveyService{
		repo: repo,
	}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, req *models.CreateSurveyRequest) (*models.Survey, error) {
	survey := &models.Survey{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Token:     generateRandomToken(5),
		Questions: make([]models.Question, len(req.Questions)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for i, qInput := range req.Questions {
		survey.Questions[i] = models.Question{
			ID:            bson.NewObjectID(),
			Text:          qInput.Text,
			Type:          qInput.Type,
			Specification: qInput.Specification,
		}
	}

	err := s.repo.Create(ctx, survey)
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func generateRandomToken(i int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for j := 0; j < i; j++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (s *SurveyService) GetSurveyByToken(ctx context.Context, token string) (*models.Survey, error) {
	return s.repo.GetByToken(ctx, token)
}
