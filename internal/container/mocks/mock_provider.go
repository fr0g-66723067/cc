package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockProvider is a mock implementation of the container.Provider interface
type MockProvider struct {
	mock.Mock
}

// Initialize sets up the container environment
func (m *MockProvider) Initialize(ctx context.Context, config map[string]string) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// RunContainer starts a container with the given image and returns its ID
func (m *MockProvider) RunContainer(ctx context.Context, image string, volumeMounts map[string]string, env map[string]string) (string, error) {
	args := m.Called(ctx, image, volumeMounts, env)
	return args.String(0), args.Error(1)
}

// ExecuteCommand executes a command in the container
func (m *MockProvider) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	args := m.Called(ctx, containerID, command)
	return args.String(0), args.Error(1)
}

// CopyFilesToContainer copies files from local to container
func (m *MockProvider) CopyFilesToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error {
	args := m.Called(ctx, containerID, localPath, containerPath)
	return args.Error(0)
}

// CopyFilesFromContainer copies files from container to local
func (m *MockProvider) CopyFilesFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error {
	args := m.Called(ctx, containerID, containerPath, localPath)
	return args.Error(0)
}

// StopContainer stops a running container
func (m *MockProvider) StopContainer(ctx context.Context, containerID string) error {
	args := m.Called(ctx, containerID)
	return args.Error(0)
}

// RemoveContainer removes a container
func (m *MockProvider) RemoveContainer(ctx context.Context, containerID string) error {
	args := m.Called(ctx, containerID)
	return args.Error(0)
}

// Name returns the provider's name
func (m *MockProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

// IsRemote returns whether the provider is running containers remotely
func (m *MockProvider) IsRemote() bool {
	args := m.Called()
	return args.Bool(0)
}
