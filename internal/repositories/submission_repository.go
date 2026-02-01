package repositories

import (
	"context"
	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	GetBySurveyID(ctx context.Context, surveyID interface{}) ([]models.Submission, error)
}

type MongoSubmissionRepository struct {
	collection *mongo.Collection
}

func NewMongoSubmissionRepository(collection *mongo.Collection) *MongoSubmissionRepository {
	return &MongoSubmissionRepository{
		collection: collection,
	}
}

func (r *MongoSubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	_, err := r.collection.InsertOne(ctx, submission)
	return err
}

func (r *MongoSubmissionRepository) GetBySurveyID(ctx context.Context, surveyID interface{}) ([]models.Submission, error) {
	cursor, err := r.collection.Find(ctx, map[string]interface{}{"survey_id": surveyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []models.Submission
	if err := cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}
	return submissions, nil
}
