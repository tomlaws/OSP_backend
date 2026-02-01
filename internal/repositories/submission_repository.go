package repositories

import (
	"context"
	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	GetAllSubmissions(ctx context.Context, surveyID bson.ObjectID) ([]*models.Submission, error)
	GetSubmissions(ctx context.Context, offset int64, limit int64, surveyID *bson.ObjectID) ([]*models.Submission, error)
	Delete(ctx context.Context, id bson.ObjectID) error
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

func (r *MongoSubmissionRepository) GetAllSubmissions(ctx context.Context, surveyID bson.ObjectID) ([]*models.Submission, error) {
	filter := bson.M{"survey_id": surveyID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []*models.Submission
	for cursor.Next(ctx) {
		var submission models.Submission
		if err := cursor.Decode(&submission); err != nil {
			return nil, err
		}
		submissions = append(submissions, &submission)
	}
	return submissions, nil
}

func (r *MongoSubmissionRepository) GetSubmissions(ctx context.Context, offset int64, limit int64, surveyID *bson.ObjectID) ([]*models.Submission, error) {
	filter := bson.M{}
	if surveyID != nil {
		filter["survey_id"] = *surveyID
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "updated_at", Value: -1},
			{Key: "created_at", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var submissions []*models.Submission
	for cursor.Next(ctx) {
		var submission models.Submission
		if err := cursor.Decode(&submission); err != nil {
			return nil, err
		}
		submissions = append(submissions, &submission)
	}
	return submissions, nil

}

func (r *MongoSubmissionRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
