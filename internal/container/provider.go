package container

import (
	"context"
	"fmt"
)

// Provider defines the interface for container systems
type Provider interface {
	// Initialize sets up the container environment
	Initialize(ctx context.Context, config map[string]string) error

	// RunContainer starts a container with the given image and returns its ID
	RunContainer(ctx context.Context, image string, volumeMounts map[string]string, env map[string]string) (string, error)

	// ExecuteCommand executes a command in the container
	ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error)

	// CopyFilesToContainer copies files from local to container
	CopyFilesToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error

	// CopyFilesFromContainer copies files from container to local
	CopyFilesFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error

	// StopContainer stops a running container
	StopContainer(ctx context.Context, containerID string) error

	// RemoveContainer removes a container
	RemoveContainer(ctx context.Context, containerID string) error

	// Name returns the provider's name
	Name() string

	// IsRemote returns whether the provider is running containers remotely
	IsRemote() bool
}

// Factory creates a provider based on name
type Factory func(config map[string]string) (Provider, error)

var providers = make(map[string]Factory)

// Register registers a provider factory
func Register(name string, factory Factory) {
	providers[name] = factory
}

// Create creates a provider with the given name
func Create(name string, config map[string]string) (Provider, error) {
	factory, exists := providers[name]
	if !exists {
		return nil, fmt.Errorf("unknown container provider: %s", name)
	}
	return factory(config)
}