package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextNavigation(t *testing.T) {
	// Create a test configuration
	cfg := DefaultConfig()

	// Test initial state
	assert.Equal(t, "root", cfg.Context.Level)
	assert.Equal(t, "", cfg.Context.ProjectName)
	assert.Equal(t, "", cfg.Context.ImplementationBranch)
	assert.Equal(t, "", cfg.Context.FeatureBranch)
	assert.Equal(t, "/", cfg.GetContextPath())

	// Test setting project context
	cfg.SetContextProject("test-project")
	assert.Equal(t, "project", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)
	assert.Equal(t, "", cfg.Context.ImplementationBranch)
	assert.Equal(t, "", cfg.Context.FeatureBranch)
	assert.Equal(t, "/projects/test-project", cfg.GetContextPath())

	// Test setting implementation context
	cfg.SetContextImplementation("impl-test-123")
	assert.Equal(t, "implementation", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)
	assert.Equal(t, "impl-test-123", cfg.Context.ImplementationBranch)
	assert.Equal(t, "", cfg.Context.FeatureBranch)
	assert.Equal(t, "/projects/test-project/implementations/impl-test-123", cfg.GetContextPath())

	// Test setting feature context
	cfg.SetContextFeature("feat-test-123")
	assert.Equal(t, "feature", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)
	assert.Equal(t, "impl-test-123", cfg.Context.ImplementationBranch)
	assert.Equal(t, "feat-test-123", cfg.Context.FeatureBranch)
	assert.Equal(t, "/projects/test-project/implementations/impl-test-123/features/feat-test-123", cfg.GetContextPath())

	// Test resetting to root
	cfg.SetContextLevel("root")
	assert.Equal(t, "root", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName) // These don't change automatically
	assert.Equal(t, "impl-test-123", cfg.Context.ImplementationBranch)
	assert.Equal(t, "feat-test-123", cfg.Context.FeatureBranch)
	assert.Equal(t, "/", cfg.GetContextPath()) // Path is based on level, not values
}

func TestContextParentChildRelationships(t *testing.T) {
	// Create a test configuration
	cfg := DefaultConfig()

	// Setting project
	cfg.SetContextProject("test-project")
	assert.Equal(t, "project", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)

	// Setting implementation clears feature
	cfg.SetContextFeature("should-be-cleared")
	cfg.SetContextImplementation("impl-test-123")
	assert.Equal(t, "implementation", cfg.Context.Level)
	assert.Equal(t, "impl-test-123", cfg.Context.ImplementationBranch)
	assert.Equal(t, "", cfg.Context.FeatureBranch) // Feature should be cleared

	// Setting project clears both implementation and feature
	cfg.SetContextImplementation("impl-test-456")
	cfg.SetContextFeature("feat-test-123")
	cfg.SetContextProject("another-project")
	assert.Equal(t, "project", cfg.Context.Level)
	assert.Equal(t, "another-project", cfg.Context.ProjectName)
	assert.Equal(t, "", cfg.Context.ImplementationBranch) // Implementation should be cleared
	assert.Equal(t, "", cfg.Context.FeatureBranch)        // Feature should be cleared
}