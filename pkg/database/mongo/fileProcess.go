package mongo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/avyukth/search-app/pkg/downloader"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var _ downloader.LinkProcessor = &Database{}

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

func (db *Database) IsLinkProcessed(ctx context.Context, id string) (bool, error) {
	var result LinkStatus
	collection := db.Client.Database(db.Config.MongoDatabase).Collection(db.Config.MongoDBLinkCollectionName)

	err := collection.FindOne(ctx, bson.M{"_id": id, "status": "processing"}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *Database) MarkLinkAsCompleted(ctx context.Context, id string) error {
	collection := db.Client.Database(db.Config.MongoDatabase).Collection(db.Config.MongoDBLinkCollectionName)

	_, err := collection.InsertOne(ctx, bson.M{"_id": id, "status": "completed", "updatedAt": time.Now()})
	return err
}

func (db *Database) MarkLinkAsProcessing(ctx context.Context, id string) error {
	collection := db.Client.Database(db.Config.MongoDatabase).Collection(db.Config.MongoDBLinkCollectionName)

	_, err := collection.InsertOne(ctx, bson.M{"_id": id, "status": "processing", "updatedAt": time.Now()})
	return err
}

func (d *Database) Download(ctx context.Context, url string, destPath string) error {
	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server return non-200 status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
