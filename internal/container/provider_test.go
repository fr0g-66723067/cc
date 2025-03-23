package container_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fr0g-66723067/cc/internal/container"
	"github.com/stretchr/testify/assert"
)

// TestProviderFactory tests the provider factory function
func TestProviderFactory(t *testing.T) {
	// Create a mock factory function
	mockFactory := func(config map[string]string) (container.Provider, error) {
		return &mockProvider{name: "mock"}, nil
	}

	// Register the mock factory
	container.Register("mock", mockFactory)

	// Test creating a provider
	provider, err := container.Create("mock", nil)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "mock", provider.Name())

	// Test creating a non-existent provider
	provider, err = container.Create("nonexistent", nil)
	assert.Error(t, err)
	assert.Nil(t, provider)
}

// TestErrorPropagation tests error propagation from the factory
func TestErrorPropagation(t *testing.T) {
	// Create a mock factory function that returns an error
	expectedErr := fmt.Errorf("mock error")
	mockFactory := func(config map[string]string) (container.Provider, error) {
		return nil, expectedErr
	}

	// Register the mock factory
	container.Register("error", mockFactory)

	// Test error propagation
	provider, err := container.Create("error", nil)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, provider)
}

// mockProvider is a simple implementation of the Provider interface for testing
type mockProvider struct {
	name string
}

func (p *mockProvider) Initialize(ctx context.Context, config map[string]string) error {
	return nil
}

func (p *mockProvider) RunContainer(ctx context.Context, image string, volumeMounts map[string]string, env map[string]string) (string, error) {
	return "container-id", nil
}

func (p *mockProvider) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	return "command output", nil
}

func (p *mockProvider) CopyFilesToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error {
	return nil
}

func (p *mockProvider) CopyFilesFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error {
	return nil
}

func (p *mockProvider) StopContainer(ctx context.Context, containerID string) error {
	return nil
}

func (p *mockProvider) RemoveContainer(ctx context.Context, containerID string) error {
	return nil
}

func (p *mockProvider) Name() string {
	return p.name
}

func (p *mockProvider) IsRemote() bool {
	return false
}