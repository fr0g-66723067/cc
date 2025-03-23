package job_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fr0g-66723067/cc/internal/job"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	queue := job.NewQueue()
	assert.NotNil(t, queue)
}

func TestRegisterHandler(t *testing.T) {
	queue := job.NewQueue()
	
	// Register a handler
	queue.RegisterHandler("test", func(ctx context.Context, payload map[string]interface{}) (interface{}, error) {
		return "test result", nil
	})
	
	// No way to directly test registration worked, but we'll test it indirectly by submitting a job
}

func TestSubmitJob(t *testing.T) {
	queue := job.NewQueue()
	
	// Register a handler
	handlerCalled := false
	queue.RegisterHandler("test", func(ctx context.Context, payload map[string]interface{}) (interface{}, error) {
		handlerCalled = true
		return "test result", nil
	})
	
	// Submit a job
	jobID, err := queue.Submit("test", nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, jobID)
	
	// Give the handler time to run
	time.Sleep(100 * time.Millisecond)
	
	// Check that the handler was called
	assert.True(t, handlerCalled)
	
	// Get the job
	j, err := queue.GetJob(jobID)
	assert.NoError(t, err)
	assert.NotNil(t, j)
	
	// Check job properties
	assert.Equal(t, jobID, j.ID)
	assert.Equal(t, "test", j.Type)
	assert.Equal(t, job.StatusCompleted, j.Status)
	assert.Equal(t, "test result", j.Result)
	assert.Nil(t, j.Error)
	assert.NotNil(t, j.StartedAt)
	assert.NotNil(t, j.CompletedAt)
}

func TestSubmitJobWithPayload(t *testing.T) {
	queue := job.NewQueue()
	
	// Register a handler
	queue.RegisterHandler("test", func(ctx context.Context, payload map[string]interface{}) (interface{}, error) {
		// Check that payload was passed correctly
		if name, ok := payload["name"].(string); ok {
			return "Hello, " + name, nil
		}
		return nil, errors.New("invalid payload")
	})
	
	// Submit a job with payload
	payload := map[string]interface{}{
		"name": "Test",
	}
	jobID, err := queue.Submit("test", payload)
	assert.NoError(t, err)
	
	// Give the handler time to run
	time.Sleep(100 * time.Millisecond)
	
	// Get the job
	j, err := queue.GetJob(jobID)
	assert.NoError(t, err)
	assert.Equal(t, job.StatusCompleted, j.Status)
	assert.Equal(t, "Hello, Test", j.Result)
}

func TestSubmitJobError(t *testing.T) {
	queue := job.NewQueue()
	
	// Register a handler that returns an error
	expectedErr := errors.New("test error")
	queue.RegisterHandler("error", func(ctx context.Context, payload map[string]interface{}) (interface{}, error) {
		return nil, expectedErr
	})
	
	// Submit a job
	jobID, err := queue.Submit("error", nil)
	assert.NoError(t, err)
	
	// Give the handler time to run
	time.Sleep(100 * time.Millisecond)
	
	// Get the job
	j, err := queue.GetJob(jobID)
	assert.NoError(t, err)
	assert.Equal(t, job.StatusFailed, j.Status)
	assert.Error(t, j.Error)
	assert.Equal(t, expectedErr.Error(), j.Error.Error())
}

func TestSubmitUnknownJobType(t *testing.T) {
	queue := job.NewQueue()
	
	// Submit a job with an unknown type
	_, err := queue.Submit("unknown", nil)
	assert.Error(t, err)
}

func TestGetNonexistentJob(t *testing.T) {
	queue := job.NewQueue()
	
	// Get a non-existent job
	_, err := queue.GetJob("nonexistent")
	assert.Error(t, err)
}