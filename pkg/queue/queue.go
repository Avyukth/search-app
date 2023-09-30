package queue

import (
	"context"
	"log"
	"sync"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/parser"
)

type TaskQueue struct {
	tasks    chan Task
	wg       *sync.WaitGroup
	parser   *parser.Parser
	dbClient *mongo.Database
	indexer  *indexer.SearchEngine
	resume   chan struct{} // Channel to resume workers
}

type Task struct {
	FilePath string
}

func NewTaskQueue(size int, parser *parser.Parser, dbClient *mongo.Database, indexer *indexer.SearchEngine) *TaskQueue {
	return &TaskQueue{
		tasks:    make(chan Task, size),
		wg:       &sync.WaitGroup{},
		parser:   parser,
		dbClient: dbClient,
		indexer:  indexer,
		resume:   make(chan struct{}), // Initialize the resume channel
	}
}

func (q *TaskQueue) Enqueue(task Task) {
	q.tasks <- task
	q.resume <- struct{}{} // Send signal to resume a worker
}

func (q *TaskQueue) Start(ctx context.Context) {
	for i := 0; i < cap(q.tasks); i++ {
		q.wg.Add(1)
		go q.worker(ctx)
	}
}

func (q *TaskQueue) Stop() {
	close(q.tasks)
	q.wg.Wait()
}

func (q *TaskQueue) worker(ctx context.Context) {
	defer q.wg.Done()
	for {
		select {
		case task, ok := <-q.tasks:
			if !ok {
				return // exit if the tasks channel is closed
			}
			if err := q.processTask(ctx, task); err != nil {
				log.Printf("Error processing task %v: %v", task, err)
			}
		case <-q.resume: // resume when a signal is received
			continue
		case <-ctx.Done():
			return // exit if context is done
		}
	}
}

func (q *TaskQueue) processTask(ctx context.Context, task Task) error {
	parsedData, err := q.parser.Parse(task.FilePath)
	if err != nil {
		log.Printf("Error parsing XML: %v", err)
		return err
	}

	xmlID, err := q.dbClient.StoreXML(parsedData)
	if err != nil {
		log.Printf("Error storing XML to MongoDB: %v", err)
		return err
	}

	patent, err := q.parser.BuildPatent(parsedData, xmlID)
	if err != nil {
		log.Printf("Error building patent: %v", err)
		return err
	}

	_, err = q.dbClient.StorePatent(patent)
	if err != nil {
		log.Printf("Error storing patent to MongoDB: %v", err)
		return err
	}

	err = q.indexer.IndexPatent(patent)
	if err != nil {
		log.Printf("Error indexing patent: %v", err)
		return err
	}

	return nil
}
