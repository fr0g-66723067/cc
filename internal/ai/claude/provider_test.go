package claude_test

import (
	"testing"

	"github.com/fr0g-66723067/cc/internal/ai/claude"
	"github.com/stretchr/testify/assert"
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

// TestInitialize tests initializing a Claude provider
func TestInitialize(t *testing.T) {
	t.Skip("Requires mocking container.Create which is difficult without refactoring")
	
	// This test would verify that:
	// 1. Container provider is created with the right config
	// 2. Container provider is initialized correctly
	// 3. All config values are properly set
}

// TestGenerateProject tests generating a project
func TestGenerateProject(t *testing.T) {
	t.Skip("Requires mocking or interface exposure")
	
	// This test would verify that:
	// 1. Container is created if needed
	// 2. Command is executed with correct parameters
	// 3. Output is returned correctly
}

// TestAddFeature tests adding a feature to existing code
func TestAddFeature(t *testing.T) {
	t.Skip("Requires mocking or interface exposure")
	
	// This test would verify that:
	// 1. Container is created if needed
	// 2. Files are copied to container
	// 3. Command is executed with correct parameters
	// 4. Files are copied back from container
	// 5. Output is returned correctly
}

// TestAnalyzeCode tests analyzing code
func TestAnalyzeCode(t *testing.T) {
	t.Skip("Requires mocking or interface exposure")
	
	// This test would verify that:
	// 1. Container is created if needed
	// 2. Files are copied to container
	// 3. Command is executed with correct parameters
	// 4. Output is returned correctly
}

// TestCleanup tests cleanup operations
func TestCleanup(t *testing.T) {
	t.Skip("Requires mocking or interface exposure")
	
	// This test would verify that:
	// 1. Container is stopped and removed
}