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

// SetupDatabase sets up the database connection using the provided configuration
func SetupDatabase(cfg *config.Config) (*mongo.Database, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf(cfg.MongoURI, cfg.MongoUsername, cfg.MongoPassword)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB Ping Error: %v", err)
		return nil, nil, err
	}

	db := client.Database(cfg.MongoDatabase)
	return db, client, nil
}
