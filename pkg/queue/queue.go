package queue

import (
	"context"
	"log"
	"sync"
)

// Task represents a unit of work to be processed.
type Task struct {
	FilePath string
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
	resume    chan struct{} // Channel to resume workers
}

// NewTaskQueue creates a new TaskQueue with the given TaskProcessor and size.
func NewTaskQueue(size int, processor TaskProcessor) *TaskQueue {
	return &TaskQueue{
		tasks:     make(chan Task, size),
		wg:        &sync.WaitGroup{},
		processor: processor,
		resume:    make(chan struct{}), // Initialize the resume channel
	}
}

// Enqueue adds a new task to the queue and resumes a worker.
func (q *TaskQueue) Enqueue(task Task) {
	q.tasks <- task
	q.resume <- struct{}{} // Send signal to resume a worker
}

// Start initializes workers to process tasks.
func (q *TaskQueue) Start(ctx context.Context) {
	for i := 0; i < cap(q.tasks); i++ {
		q.wg.Add(1)
		go q.worker(ctx)
	}
}

// Stop waits for all workers to finish processing and closes the tasks channel.
func (q *TaskQueue) Stop() {
	close(q.tasks)
	q.wg.Wait()
}

// worker is a goroutine that processes tasks from the queue.
func (q *TaskQueue) worker(ctx context.Context) {
	defer q.wg.Done()
	for {
		select {
		case task, ok := <-q.tasks:
			if !ok {
				return // exit if the tasks channel is closed
			}
			if err := q.processor.Process(ctx, task); err != nil {
				log.Printf("Error processing task %v: %v", task, err)
			}
		case <-q.resume: // resume when a signal is received
			continue
		case <-ctx.Done():
			return // exit if context is done
		}
	}
}
