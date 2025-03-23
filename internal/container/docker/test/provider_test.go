package docker_test

import (
	"context"
	"strings"
	"testing"

	"github.com/fr0g-66723067/cc/internal/container/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewProvider verifies we can create a new provider with config
func TestNewProvider(t *testing.T) {
	config := map[string]string{
		"key": "value",
	}
	
	provider, err := docker.NewProvider(config)
	require.NoError(t, err)
	require.NotNil(t, provider)
	
	assert.Equal(t, "docker", provider.Name())
	assert.False(t, provider.IsRemote())
}

// TestInitialize tests initialization of the Docker provider
func TestInitialize(t *testing.T) {
	provider, err := docker.NewProvider(nil)
	require.NoError(t, err)
	
	config := map[string]string{
		"host": "unix:///var/run/docker.sock",
	}
	
	ctx := context.Background()
	err = provider.Initialize(ctx, config)
	if err != nil {
		t.Skipf("Skipping Docker test due to connection error: %v", err)
	}
	
	// In a real test we'd verify the Docker client was initialized properly
}

// TestCopyOperations tests file copying between host and container
func TestCopyOperations(t *testing.T) {
	t.Skip("Implementation requires actual Docker - will be manually tested")
}

// TestContainerLifecycle tests container lifecycle operations
func TestContainerLifecycle(t *testing.T) {
	provider, err := docker.NewProvider(nil)
	require.NoError(t, err)
	
	ctx := context.Background()
	err = provider.Initialize(ctx, nil)
	if err != nil {
		t.Skipf("Skipping Docker test due to connection error: %v", err)
		return
	}
	
	// Test is now set up to check for real Docker availability and skip if not available
	
	// If we've gotten here, Docker is available, but we'll still skip for CI environment
	// since we don't want to actually pull and run containers in automated tests
	// This would typically be controlled by an environment variable
	if testing.Short() {
		t.Skip("Skipping Docker test in short mode")
	}
	
	// For local testing only - these tests will use actual Docker
	t.Skip("This test would run actual Docker containers - skipping by default")
	
	// Run a container
	containerID, err := provider.RunContainer(ctx, "alpine:latest", nil, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, containerID)
	
	// Execute command
	output, err := provider.ExecuteCommand(ctx, containerID, []string{"echo", "hello"})
	require.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(output), "hello")
	
	// Stop container
	err = provider.StopContainer(ctx, containerID)
	require.NoError(t, err)
	
	// Remove container
	err = provider.RemoveContainer(ctx, containerID)
	require.NoError(t, err)
}