package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/avyukth/search-app/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf(cfg.MongoDBConfig.URI, cfg.MongoDBConfig.Username, cfg.MongoDBConfig.Password)
	log.Println("Connecting to MongoDB at: ", uri)
	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(cfg.MongoDBConfig.MaxPoolSize)
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
