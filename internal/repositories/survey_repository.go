package repositories

import (
	"context"
	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveyRepository interface {
	Create(ctx context.Context, survey *models.Survey) error
	GetByToken(ctx context.Context, token string) (*models.Survey, error)
	GetByID(ctx context.Context, id interface{}) (*models.Survey, error)
}

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository(collection *mongo.Collection) *MongoSurveyRepository {
	return &MongoSurveyRepository{
		collection: collection,
	}
}

func (r *MongoSurveyRepository) Create(ctx context.Context, survey *models.Survey) error {
	_, err := r.collection.InsertOne(ctx, survey)
	return err
}

func (r *MongoSurveyRepository) GetByToken(ctx context.Context, token string) (*models.Survey, error) {
	var survey models.Survey
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *MongoSurveyRepository) GetByID(ctx context.Context, id interface{}) (*models.Survey, error) {
	var survey models.Survey
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	return &survey, nil
}
