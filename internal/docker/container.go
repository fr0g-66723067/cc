package docker

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fr0g-66723067/cc/internal/container"
	"github.com/fr0g-66723067/cc/pkg/config"
)

// Manager handles Docker container operations
type Manager struct {
	provider container.Provider
	config   *config.Config
}

// NewManager creates a new Docker manager
func NewManager(cfg *config.Config) (*Manager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Create the Docker provider
	provider, err := container.Create(cfg.Container.Provider, cfg.Container.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create container provider: %w", err)
	}

	// Initialize the provider
	ctx := context.Background()
	if err := provider.Initialize(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to initialize container provider: %w", err)
	}

	return &Manager{
		provider: provider,
		config:   cfg,
	}, nil
}

// RunClaudeContainer starts a Claude Code container
func (m *Manager) RunClaudeContainer(ctx context.Context, projectPath string) (string, error) {
	// Get the Claude Code image from config
	image := m.config.Container.ClaudeImage
	if image == "" {
		image = "anthropic/claude-code:latest" // Default image
	}

	// Set up volume mounts
	volumeMounts := map[string]string{
		projectPath: "/workspace", // Mount project directory to /workspace in container
	}

	// Set up environment variables
	env := map[string]string{
		"CC_PROJECT_PATH": "/workspace",
	}

	// Add any additional environment variables from config
	for k, v := range m.config.Container.Config {
		if k != "host" && k != "socketPath" { // Skip Docker connection settings
			env[k] = v
		}
	}

	// Run the container
	containerID, err := m.provider.RunContainer(ctx, image, volumeMounts, env)
	if err != nil {
		return "", fmt.Errorf("failed to run Claude Code container: %w", err)
	}

	return containerID, nil
}

// ExecuteCommand runs a command in the Claude Code container
func (m *Manager) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	// Execute the command
	output, err := m.provider.ExecuteCommand(ctx, containerID, command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command in container: %w", err)
	}

	return output, nil
}

// CopyFileToContainer copies a file to the container
func (m *Manager) CopyFileToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error {
	// Resolve paths
	absLocalPath, err := filepath.Abs(localPath)
	if err != nil {
		return fmt.Errorf("failed to resolve local path: %w", err)
	}

	return m.provider.CopyFilesToContainer(ctx, containerID, absLocalPath, containerPath)
}

// CopyFileFromContainer copies a file from the container
func (m *Manager) CopyFileFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error {
	// Resolve paths
	absLocalPath, err := filepath.Abs(localPath)
	if err != nil {
		return fmt.Errorf("failed to resolve local path: %w", err)
	}

	return m.provider.CopyFilesFromContainer(ctx, containerID, containerPath, absLocalPath)
}

// StopContainer stops a running container
func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	return m.provider.StopContainer(ctx, containerID)
}

// CleanupContainer removes a container
func (m *Manager) CleanupContainer(ctx context.Context, containerID string) error {
	return m.provider.RemoveContainer(ctx, containerID)
}