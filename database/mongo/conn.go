package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/avyukth/search-app/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
)

// GetMongoClient returns a MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		cfg := config.LoadConfig()

		// Set client options
		clientOptions := options.Client().ApplyURI(cfg.MongoURI)
		clientOptions.SetMaxPoolSize(cfg.MongoMaxPoolSize)

		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceErr = err
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceErr = err
		}

		if clientInstanceErr == nil {
			clientInstance = client
			log.Println("Connected to MongoDB!")
		}
	})

	return clientInstance, clientInstanceErr
}

// GetMongoDB returns a MongoDB database instance
func GetMongoDB() *mongo.Database {
	cfg := config.LoadConfig()
	client, err := GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(cfg.MongoDatabase)
}

// CloseMongoClient closes the existing MongoDB client
func CloseMongoClient() {
	client, err := GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}
