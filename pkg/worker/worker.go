// worker.go
package worker

import (
	"context"
	"log"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/downloader"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/parser"
	"github.com/avyukth/search-app/pkg/queue"
)

// Worker is an interface that represents the ability to process tasks.
type Worker interface {
	Process(ctx context.Context, task queue.Task) error
}

// taskWorker is a concrete implementation of Worker that can process tasks.
type taskWorker struct {
	downloader downloader.Downloader
	parser     *parser.Parser
	dbClient   *mongo.Database
	indexer    *indexer.SearchEngine
}

// NewWorker acts as a constructor and returns an instance of Worker.
func NewWorker(d downloader.Downloader, p *parser.Parser, db *mongo.Database, i *indexer.SearchEngine) Worker {
	return &taskWorker{
		downloader: d,
		parser:     p,
		dbClient:   db,
		indexer:    i,
	}
}

// Process processes the given task.
func (w *taskWorker) Process(ctx context.Context, task queue.Task) error {
	// Download the file
	filePath, err := w.downloader.Download(ctx, task.FilePath)
	if err != nil {
		return err
	}

	// Parse the downloaded file
	parsedData, err := w.parser.Parse(filePath)
	if err != nil {
		log.Printf("Error parsing XML: %v", err)
		return err
	}

	// Store the parsed data to MongoDB
	xmlID, err := w.dbClient.StoreXML(parsedData)
	if err != nil {
		log.Printf("Error storing XML to MongoDB: %v", err)
		return err
	}

	// Build and store the patent
	patent, err := w.parser.BuildPatent(parsedData, xmlID)
	if err != nil {
		log.Printf("Error building patent: %v", err)
		return err
	}

	_, err = w.dbClient.StorePatent(patent)
	if err != nil {
		log.Printf("Error storing patent to MongoDB: %v", err)
		return err
	}

	// Index the patent
	err = w.indexer.IndexPatent(patent)
	if err != nil {
		log.Printf("Error indexing patent: %v", err)
		return err
	}

	return nil
}