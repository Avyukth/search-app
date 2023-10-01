package queue

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a unit of work to be processed.
type TaskType int

const (
	DownloadAndProcess TaskType = iota
	WalkAndProcess
)

type Task struct {
	FilePath string
	Type     TaskType
}

// TaskProcessor is an interface that represents the ability to process tasks.
type TaskProcessor interface {
	Process(ctx context.Context, task Task) error
}

// TaskQueue manages a queue of tasks and processes them using the provided TaskProcessor.
type TaskQueue struct {
	tasks     chan Task
	wg        *sync.WaitGroup
	processor TaskProcessor
	resume    chan struct{}
}

// NewTaskQueue creates a new TaskQueue with the given TaskProcessor and size.
func NewTaskQueue(size int, processor TaskProcessor) *TaskQueue {
	log.Println("Initializing TaskQueue with size:", size)
	q := &TaskQueue{
		tasks:     make(chan Task, size),
		wg:        &sync.WaitGroup{},
		processor: processor,
		resume:    make(chan struct{}),
	}
	log.Println("TaskQueue Initialized.")
	return q
}

// Enqueue adds a new task to the queue and resumes a worker.
func (q *TaskQueue) Enqueue(task Task) {
	log.Printf("Enqueueing task: %+v\n", task)
	q.tasks <- task
	log.Printf("Task: %+v enqueued.\n", task)
	log.Println("Sending signal to resume a worker.")
	q.resume <- struct{}{}
	log.Println("Signal sent to resume a worker.")
}

// Start initializes workers to process tasks.
func (q *TaskQueue) Start(ctx context.Context) {
	log.Println("Starting workers")
	for i := 0; i < cap(q.tasks); i++ {
		q.wg.Add(1)
		go q.worker(ctx)
		log.Printf("Worker %d started.\n", i)
	}
}

// Stop waits for all workers to finish processing and closes the tasks channel.
func (q *TaskQueue) Stop() {
	log.Println("Stopping TaskQueue, closing tasks channel.")
	close(q.tasks)
	q.wg.Wait()
	log.Println("All workers have finished processing, TaskQueue stopped.")
}

// worker is a goroutine that processes tasks from the queue.
func (q *TaskQueue) worker(ctx context.Context) {
	log.Println("Worker goroutine is running.")
	defer q.wg.Done()
	for {
		select {
		case task, ok := <-q.tasks:
			if !ok {
				log.Println("Tasks channel closed, exiting worker.")
				return
			}
			log.Printf("Processing task: %+v\n", task)
			taskCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
			if err := q.processor.Process(taskCtx, task); err != nil {
				log.Printf("Error processing task %+v: %v\n", task, err)
			} else {
				log.Printf("Task %+v processed successfully.\n", task)
			}
			cancel()
		case <-q.resume:
			log.Println("Worker resumed.")
			continue
		case <-ctx.Done():
			log.Println("Context done, exiting worker.")
			return
		}
	}
}
