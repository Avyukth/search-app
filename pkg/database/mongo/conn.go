package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/avyukth/search-app/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBConfig.URI).SetMaxPoolSize(cfg.MongoDBConfig.MaxPoolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := &Database{
		Client:     client,
		Config:     cfg,
		Collection: client.Database(cfg.MongoDBConfig.Database).Collection(cfg.MongoDBConfig.LinkCollectionName),
	}

	return db, nil
}
