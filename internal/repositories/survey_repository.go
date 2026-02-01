package repositories

import (
	"context"
	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SurveyRepository interface {
	Create(ctx context.Context, survey *models.Survey) error
	List(ctx context.Context, offset, limit int64) ([]*models.Survey, error)
	GetByToken(ctx context.Context, token string) (*models.Survey, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.Survey, error)
	Delete(ctx context.Context, id bson.ObjectID) error
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

func (r *MongoSurveyRepository) List(ctx context.Context, offset, limit int64) ([]*models.Survey, error) {
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var surveys []*models.Survey
	if err := cursor.All(ctx, &surveys); err != nil {
		return nil, err
	}
	return surveys, nil
}

func (r *MongoSurveyRepository) GetByToken(ctx context.Context, token string) (*models.Survey, error) {
	var survey models.Survey
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *MongoSurveyRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.Survey, error) {
	var survey models.Survey
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&survey)
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *MongoSurveyRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
