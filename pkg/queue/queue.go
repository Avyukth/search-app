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
	}
}

func (q *TaskQueue) Enqueue(task Task) {
	q.tasks <- task
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
	for task := range q.tasks {
		if err := q.processTask(ctx, task); err != nil {
			log.Printf("Error processing task %v: %v", task, err)
		}
	}
}

func (q *TaskQueue) processTask(ctx context.Context, task Task) error {
	// 1. Parse the XML file
	parsedData, err := q.parser.Parse(task.FilePath)
	if err != nil {
		return err
	}

	// 2. Store parsed XML to MongoDB
	xmlID, err := q.dbClient.StoreXML(parsedData)
	if err != nil {
		return err
	}

	// 3. Build Patent Object from parsedData and xmlID
	patent, err := q.parser.BuildPatent(parsedData, xmlID)
	if err != nil {
		return err
	}

	// 4. Store Patent Object to MongoDB
	_, err = q.dbClient.StorePatent(patent)
	if err != nil {
		return err
	}

	// 5. Build Search Index
	err = q.indexer.IndexPatent(patent)
	if err != nil {
		return err
	}

	return nil
}
