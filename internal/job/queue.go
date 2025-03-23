package job

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Status represents the status of a job
type Status string

const (
	// StatusPending indicates the job is waiting to be processed
	StatusPending Status = "pending"
	// StatusRunning indicates the job is currently running
	StatusRunning Status = "running"
	// StatusCompleted indicates the job has completed successfully
	StatusCompleted Status = "completed"
	// StatusFailed indicates the job has failed
	StatusFailed Status = "failed"
)

// Job represents an asynchronous job
type Job struct {
	ID          string
	Type        string
	Payload     map[string]interface{}
	Status      Status
	Result      interface{}
	Error       error
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

// Handler is a function that processes a job
type Handler func(ctx context.Context, payload map[string]interface{}) (interface{}, error)

// Queue manages asynchronous jobs
type Queue struct {
	jobs     map[string]*Job
	handlers map[string]Handler
	mutex    sync.RWMutex
}

// NewQueue creates a new job queue
func NewQueue() *Queue {
	return &Queue{
		jobs:     make(map[string]*Job),
		handlers: make(map[string]Handler),
		mutex:    sync.RWMutex{},
	}
}

// RegisterHandler registers a handler for a job type
func (q *Queue) RegisterHandler(jobType string, handler Handler) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.handlers[jobType] = handler
}

// Submit adds a new job to the queue
func (q *Queue) Submit(jobType string, payload map[string]interface{}) (string, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Check if we have a handler for this job type
	_, exists := q.handlers[jobType]
	if !exists {
		return "", fmt.Errorf("no handler registered for job type: %s", jobType)
	}

	// Create a new job
	jobID := fmt.Sprintf("%d", time.Now().UnixNano())
	job := &Job{
		ID:        jobID,
		Type:      jobType,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	// Add it to the map
	q.jobs[jobID] = job

	// Start processing the job in the background
	go q.processJob(job)

	return jobID, nil
}

// GetJob returns a job by ID
func (q *Queue) GetJob(jobID string) (*Job, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	job, exists := q.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

// processJob processes a job in the background
func (q *Queue) processJob(job *Job) {
	// Get the handler
	q.mutex.RLock()
	handler, exists := q.handlers[job.Type]
	q.mutex.RUnlock()

	if !exists {
		// This shouldn't happen since we check before submitting
		q.mutex.Lock()
		job.Status = StatusFailed
		job.Error = fmt.Errorf("no handler registered for job type: %s", job.Type)
		q.mutex.Unlock()
		return
	}

	// Update status to running
	q.mutex.Lock()
	job.Status = StatusRunning
	now := time.Now()
	job.StartedAt = &now
	q.mutex.Unlock()

	// Process the job
	result, err := handler(context.Background(), job.Payload)

	// Update job with result
	q.mutex.Lock()
	defer q.mutex.Unlock()

	completeTime := time.Now()
	job.CompletedAt = &completeTime

	if err != nil {
		job.Status = StatusFailed
		job.Error = err
	} else {
		job.Status = StatusCompleted
		job.Result = result
	}
}