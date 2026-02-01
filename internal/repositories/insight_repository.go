package repositories

import (
	"context"
	"osp/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InsightRepository interface {
	Create(ctx context.Context, insight *models.Insight) error
	GetByID(ctx context.Context, id bson.ObjectID) (*models.Insight, error)
	GetInsights(ctx context.Context, offset, limit int64, surveyID *bson.ObjectID) ([]*models.Insight, error)
	Update(ctx context.Context, id bson.ObjectID, update interface{}) error
}

type MongoInsightRepository struct {
	collection *mongo.Collection
}

func NewMongoInsightRepository(collection *mongo.Collection) *MongoInsightRepository {
	return &MongoInsightRepository{
		collection: collection,
	}
}

func (r *MongoInsightRepository) Create(ctx context.Context, insight *models.Insight) error {
	_, err := r.collection.InsertOne(ctx, insight)
	return err
}

func (r *MongoInsightRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.Insight, error) {
	var insight models.Insight
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&insight)
	if err != nil {
		return nil, err
	}
	return &insight, nil
}

func (r *MongoInsightRepository) GetInsights(ctx context.Context, offset, limit int64, surveyID *bson.ObjectID) ([]*models.Insight, error) {
	filter := bson.M{}
	if surveyID != nil {
		filter["survey_id"] = *surveyID
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "completed_at", Value: -1},
			{Key: "updated_at", Value: -1},
			{Key: "created_at", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
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
	return insights, nil
}

func (r *MongoInsightRepository) Update(ctx context.Context, id bson.ObjectID, update interface{}) error {
	_, err := r.collection.UpdateByID(ctx, id, update)
	return err
}
