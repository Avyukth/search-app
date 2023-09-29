package mongo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckAndSetLinkStatus checks the link status and sets it to processed if not already processed or completed
func (db *Database) CheckAndSetLinkStatus(link string) (bool, error) {
	collection := db.Client.Database(db.Config.MongoDatabase).Collection(db.Config.MongoDBLinkCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate hash of the link
	hash := sha256.Sum256([]byte(link))
	linkHash := hex.EncodeToString(hash[:])

	// Check if link is already processed or completed
	var linkStatus LinkStatus
	err := collection.FindOne(ctx, map[string]interface{}{"linkHash": linkHash, "status": map[string]interface{}{"$in": []string{"processed", "completed"}}}).Decode(&linkStatus)
	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}
	if err == nil {
		return false, nil // Link is already processed or completed
	}

	// Set link to processed state
	_, err = collection.InsertOne(ctx, LinkStatus{
		LinkHash:  linkHash,
		Status:    "processed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
