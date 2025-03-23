package claude_test

import (
	"context"
	"os"
	"testing"

	"github.com/fr0g-66723067/cc/internal/ai/claude"
	"github.com/fr0g-66723067/cc/internal/container/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestNewProvider tests creating a new Claude provider
func TestNewProvider(t *testing.T) {
	// Test with default config
	provider, err := claude.NewProvider(nil)
	require.NoError(t, err)
	require.NotNil(t, provider)

	assert.Equal(t, "claude", provider.Name())

	// Frameworks should have default values
	frameworks := provider.SupportedFrameworks()
	assert.NotEmpty(t, frameworks)
	assert.Contains(t, frameworks, "react")
	assert.Contains(t, frameworks, "vue")

	// Test with custom frameworks
	customConfig := map[string]string{
		"frameworks": "angular,svelte",
	}
	provider, err = claude.NewProvider(customConfig)
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Frameworks should match custom config
	frameworks = provider.SupportedFrameworks()
	assert.Len(t, frameworks, 2)
	assert.Contains(t, frameworks, "angular")
	assert.Contains(t, frameworks, "svelte")
}

// setupMockContainerProvider creates a mock container provider for testing
func setupMockContainerProvider(t *testing.T) *mocks.MockProvider {
	// Create mock container provider
	mockProvider := new(mocks.MockProvider)

	// Set up mock methods
	mockProvider.On("Initialize", mock.Anything, mock.Anything).Return(nil)
	mockProvider.On("RunContainer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("test-container-id", nil)
	mockProvider.On("ExecuteCommand", mock.Anything, "test-container-id", []string{"echo", "ping"}).Return("ping", nil)
	mockProvider.On("ExecuteCommand", mock.Anything, "test-container-id", []string{"claude", "code", "--version"}).Return("Claude Code CLI v1.0.0", nil)
	mockProvider.On("ExecuteCommand", mock.Anything, "test-container-id", mock.Anything).Return("Command executed successfully", nil)
	mockProvider.On("CopyFilesToContainer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockProvider.On("CopyFilesFromContainer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockProvider.On("StopContainer", mock.Anything, mock.Anything).Return(nil)
	mockProvider.On("RemoveContainer", mock.Anything, mock.Anything).Return(nil)

	return mockProvider
}

// TestInitializeWithMock tests initializing a Claude provider with a mock container provider
func TestInitializeWithMock(t *testing.T) {
	// Create mock container provider
	mockProvider := setupMockContainerProvider(t)

	// Create Claude provider with test config
	config := map[string]string{
		"container_provider": "mock",
		"claude_api_key":     "test-api-key",
	}
	provider, err := claude.NewProvider(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Set mock container provider
	err = provider.SetContainerProviderForTest(mockProvider)
	require.NoError(t, err)

	// Initialize provider
	ctx := context.Background()
	err = provider.Initialize(ctx, nil)
	require.NoError(t, err)

	// Verify mock was called
	mockProvider.AssertCalled(t, "Initialize", mock.Anything, mock.Anything)
}

// TestGenerateProjectWithMock tests generating a project with a mock container provider
func TestGenerateProjectWithMock(t *testing.T) {
	// Skip the test if we don't want to run integration tests
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Create mock container provider
	mockProvider := setupMockContainerProvider(t)

	// Create Claude provider with test config
	config := map[string]string{
		"container_provider": "mock",
		"claude_api_key":     "test-api-key",
	}
	provider, err := claude.NewProvider(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Set mock container provider
	err = provider.SetContainerProviderForTest(mockProvider)
	require.NoError(t, err)

	// Initialize provider
	ctx := context.Background()
	err = provider.Initialize(ctx, nil)
	require.NoError(t, err)

	// Test generate project
	description := "Create a test project with React"
	output, err := provider.GenerateProject(ctx, description)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Verify mock was called with expected commands
	mockProvider.AssertCalled(t, "ExecuteCommand", mock.Anything, "test-container-id", mock.Anything)
}

// TestAddFeatureWithMock tests adding a feature with a mock container provider
func TestAddFeatureWithMock(t *testing.T) {
	// Skip the test if we don't want to run integration tests
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Create mock container provider
	mockProvider := setupMockContainerProvider(t)

	// Create Claude provider with test config
	config := map[string]string{
		"container_provider": "mock",
		"claude_api_key":     "test-api-key",
	}
	provider, err := claude.NewProvider(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Set mock container provider
	err = provider.SetContainerProviderForTest(mockProvider)
	require.NoError(t, err)

	// Initialize provider
	ctx := context.Background()
	err = provider.Initialize(ctx, nil)
	require.NoError(t, err)

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "claude-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Write a sample file
	sampleFile := tempDir + "/sample.txt"
	err = os.WriteFile(sampleFile, []byte("Sample content"), 0644)
	require.NoError(t, err)

	// Test add feature
	description := "Add a dark mode feature"
	output, err := provider.AddFeature(ctx, tempDir, description)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Verify mocks were called with expected commands
	mockProvider.AssertCalled(t, "CopyFilesToContainer", mock.Anything, "test-container-id", mock.Anything, mock.Anything)
	mockProvider.AssertCalled(t, "ExecuteCommand", mock.Anything, "test-container-id", mock.Anything)
	mockProvider.AssertCalled(t, "CopyFilesFromContainer", mock.Anything, "test-container-id", mock.Anything, mock.Anything)
}

// TestCleanupWithMock tests cleanup operations with a mock container provider
func TestCleanupWithMock(t *testing.T) {
	// Create mock container provider
	mockProvider := setupMockContainerProvider(t)

	// Create Claude provider with test config
	config := map[string]string{
		"container_provider": "mock",
		"claude_api_key":     "test-api-key",
	}
	provider, err := claude.NewProvider(config)
	require.NoError(t, err)
	require.NotNil(t, provider)

	// Set mock container provider
	err = provider.SetContainerProviderForTest(mockProvider)
	require.NoError(t, err)

	// Initialize provider
	ctx := context.Background()
	err = provider.Initialize(ctx, nil)
	require.NoError(t, err)

	// Call a method that uses container to ensure containerID is set
	_, err = provider.GenerateProject(ctx, "Test project")
	require.NoError(t, err)

	// Test cleanup
	err = provider.Cleanup(ctx)
	require.NoError(t, err)

	// Verify mocks were called with expected commands
	mockProvider.AssertCalled(t, "StopContainer", mock.Anything, "test-container-id")
	mockProvider.AssertCalled(t, "RemoveContainer", mock.Anything, "test-container-id")
}
