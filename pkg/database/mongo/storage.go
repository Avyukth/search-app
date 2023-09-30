package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *Database) StoreXML(data map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Client.Database(db.Config.MongoDBConfig.Database).Collection(db.Config.MongoDBConfig.StorageCollectionName)
	// Insert the parsed XML data into the MongoDB collection
	result, err := collection.InsertOne(ctx, bson.M(data))
	if err != nil {
		return "", fmt.Errorf("error storing XML data to MongoDB: %v", err)
	}

	// Return the inserted ID
	return fmt.Sprintf("%v", result.InsertedID), nil
}

func (db *Database) StorePatent(patent *Patent) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.Client.Database(db.Config.MongoDBConfig.Database).Collection(db.Config.MongoDBConfig.IndexCollectionName)
	// Insert the parsed Patent data into the MongoDB collection
	result, err := collection.InsertOne(ctx, patent)
	if err != nil {
		return "", fmt.Errorf("error storing Patent data to MongoDB: %v", err)
	}

	// Return the inserted ID
	return fmt.Sprintf("%v", result.InsertedID), nil
}

func (db *Database) RetrievePatent(patentStorageID string) (*Patent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Client.Database(db.Config.MongoDBConfig.Database).Collection(db.Config.MongoDBConfig.IndexCollectionName)

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(patentStorageID)
	if err != nil {
		return nil, fmt.Errorf("error converting string ID to ObjectID: %v", err)
	}

	// Create a filter to match the document
	filter := bson.M{"_id": objID}

	var patent Patent
	err = collection.FindOne(ctx, filter).Decode(&patent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no patent found with ID: %s", patentStorageID)
		}
		return nil, fmt.Errorf("error retrieving patent from MongoDB: %v", err)
	}

	return &patent, nil
}

func (db *Database) RetrieveXML(xmlStorageID string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Client.Database(db.Config.MongoDBConfig.Database).Collection(db.Config.MongoDBConfig.StorageCollectionName)

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(xmlStorageID)
	if err != nil {
		return nil, fmt.Errorf("error converting string ID to ObjectID: %v", err)
	}

	// Create a filter to match the document
	filter := bson.M{"_id": objID}

	var xmlData map[string]interface{}
	err = collection.FindOne(ctx, filter).Decode(&xmlData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no XML data found with ID: %s", xmlStorageID)
		}
		return nil, fmt.Errorf("error retrieving XML data from MongoDB: %v", err)
	}

	return xmlData, nil
}
