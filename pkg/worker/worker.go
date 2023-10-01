package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/downloader"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/parser"
	"github.com/avyukth/search-app/pkg/queue"
)

type Worker interface {
	Process(ctx context.Context, task queue.Task) error
}

type taskWorker struct {
	downloader downloader.Downloader
	parser     *parser.Parser
	dbClient   *mongo.Database
	indexer    *indexer.SearchEngine
}

func NewWorker(d downloader.Downloader, p *parser.Parser, db *mongo.Database, i *indexer.SearchEngine) Worker {
	return &taskWorker{
		downloader: d,
		parser:     p,
		dbClient:   db,
		indexer:    i,
	}
}

func (w *taskWorker) Process(ctx context.Context, task queue.Task) error {
	switch task.Type {
	case queue.DownloadAndProcess:
		return w.DownExtractAndProcess(ctx, task)
	case queue.WalkAndProcess:
		return w.walkDir(task.FilePath)
	default:
		return fmt.Errorf("unsupported task type: %v", task.Type)
	}
}

func (w *taskWorker) DownExtractAndProcess(ctx context.Context, task queue.Task) error {
	log.Printf("Starting processing for task: %+v", task)
	filePath, err := w.downloader.Download(ctx, task.FilePath)
	if err != nil {
		return err
	}
	log.Printf("Successfully downloaded file to: %s", filePath)

	extractedPath, err := w.downloader.ExtractTarGz(filePath)
	if err != nil {
		return err
	}
	return w.walkDir(extractedPath)
}

func (w *taskWorker) walkDir(dirPath string) error {
	var wg sync.WaitGroup

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}

		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			if err := w.processFile(filePath); err != nil {
				log.Printf("Error processing file at %s: %v", filePath, err)
			}
		}(path)

		return nil
	})

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (w *taskWorker) processFile(filePath string) error {
	parsedData, err := w.parser.Parse(filePath)
	if err != nil {
		return err
	}

	xmlID, err := w.dbClient.StoreXML(parsedData)
	if err != nil {
		return err
	}

	patent, err := w.parser.ParseToStruct(filePath, xmlID)
	if err != nil {
		return err
	}
	_, err = w.dbClient.StorePatent(patent)
	if err != nil {
		return err
	}

	if err := w.indexer.IndexPatent(patent); err != nil {
		return err
	}
	log.Printf("Successfully indexed the patent")
	return nil
}
