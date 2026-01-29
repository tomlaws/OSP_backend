package database

import (
	"context"
	"fmt"
	"osp/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(cfg *config.Config) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(cfg.DBUri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB successfully")
	return client, nil
}
