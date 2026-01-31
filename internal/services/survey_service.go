package services

import (
	"context"
	"osp/internal/models"
	"time"

	"math/rand"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveyService struct {
	collection *mongo.Collection
}

func NewSurveyService(collection *mongo.Collection) *SurveyService {
	return &SurveyService{
		collection: collection,
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

	_, err := s.collection.InsertOne(ctx, survey)
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

func (s *SurveyService) GetSurveyByToken(context context.Context, token string) (any, any) {
	var survey models.Survey
	err := s.collection.FindOne(context, bson.M{"token": token}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	return &survey, nil
}
