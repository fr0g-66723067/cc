package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/pkg/config"
	"github.com/fr0g-66723067/cc/pkg/models"
	"github.com/stretchr/testify/assert"

	main "github.com/fr0g-66723067/cc/cmd/cc"
)

// TestShellContextNavigation tests the shell's context navigation functionality
func TestShellContextNavigation(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "cc-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Create a test configuration
	cfg := config.DefaultConfig()

	// Add test projects
	project1 := models.NewProject("test-project", "/tmp/test-project", "Test Project")
	
	// Add implementations to project1
	impl1 := models.Implementation{
		Framework:   "react",
		BranchName:  "impl-react-123",
		Description: "React Implementation",
	}
	
	impl2 := models.Implementation{
		Framework:   "vue",
		BranchName:  "impl-vue-456",
		Description: "Vue Implementation",
	}
	
	// Add features to impl1
	feature1 := models.Feature{
		Name:        "login",
		BranchName:  "feat-login-123",
		Description: "Login Feature",
	}
	
	impl1.Features = append(impl1.Features, feature1)
	
	project1.AddImplementation(impl1)
	project1.AddImplementation(impl2)
	
	cfg.AddProject(project1)
	cfg.SetActiveProject("test-project")
	
	// Save the configuration
	err = config.SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create a new shell instance
	shell := main.NewShell(cfg, configPath)

	// Test getPrompt
	assert.Equal(t, "cc/ > ", shell.GetPromptForTest())

	// Test changeDirectory - move to projects/test-project
	shell.ChangeDirectoryForTest([]string{"cd", "/projects/test-project"})
	assert.Equal(t, "project", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)
	assert.Equal(t, "cc/projects/test-project > ", shell.GetPromptForTest())

	// Test changeDirectory - move to implementations/impl-react-123
	shell.ChangeDirectoryForTest([]string{"cd", "implementations/impl-react-123"})
	assert.Equal(t, "implementation", cfg.Context.Level)
	assert.Equal(t, "impl-react-123", cfg.Context.ImplementationBranch)
	assert.Equal(t, "cc/projects/test-project/implementations/impl-react-123 > ", shell.GetPromptForTest())

	// Test changeDirectory - move to features/feat-login-123
	shell.ChangeDirectoryForTest([]string{"cd", "features/feat-login-123"})
	assert.Equal(t, "feature", cfg.Context.Level)
	assert.Equal(t, "feat-login-123", cfg.Context.FeatureBranch)
	assert.Equal(t, "cc/projects/test-project/implementations/impl-react-123/features/feat-login-123 > ", shell.GetPromptForTest())

	// Test changeDirectory - move back to root
	shell.ChangeDirectoryForTest([]string{"cd", "/"})
	assert.Equal(t, "root", cfg.Context.Level)
	assert.Equal(t, "", cfg.Context.ProjectName)
	assert.Equal(t, "", cfg.Context.ImplementationBranch)
	assert.Equal(t, "", cfg.Context.FeatureBranch)
	assert.Equal(t, "cc/ > ", shell.GetPromptForTest())
}

// TestShellUseResource tests the shell's resource selection functionality
func TestShellUseResource(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "cc-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Create a test configuration
	cfg := config.DefaultConfig()

	// Add test projects
	project1 := models.NewProject("test-project", "/tmp/test-project", "Test Project")
	
	// Add implementations to project1
	impl1 := models.Implementation{
		Framework:   "react",
		BranchName:  "impl-react-123",
		Description: "React Implementation",
	}
	
	impl2 := models.Implementation{
		Framework:   "vue",
		BranchName:  "impl-vue-456",
		Description: "Vue Implementation",
	}
	
	// Add features to impl1
	feature1 := models.Feature{
		Name:        "login",
		BranchName:  "feat-login-123",
		Description: "Login Feature",
	}
	
	impl1.Features = append(impl1.Features, feature1)
	
	project1.AddImplementation(impl1)
	project1.AddImplementation(impl2)
	
	cfg.AddProject(project1)
	
	// Save the configuration
	err = config.SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create a new shell instance
	shell := main.NewShell(cfg, configPath)

	// Test useResource - use project
	shell.UseResourceForTest([]string{"use", "project", "test-project"})
	assert.Equal(t, "project", cfg.Context.Level)
	assert.Equal(t, "test-project", cfg.Context.ProjectName)
	assert.Equal(t, "test-project", cfg.ActiveProject)

	// Test useResource - use implementation
	shell.UseResourceForTest([]string{"use", "implementation", "impl-react-123"})
	assert.Equal(t, "implementation", cfg.Context.Level)
	assert.Equal(t, "impl-react-123", cfg.Context.ImplementationBranch)
	
	// Get the updated project to check if implementation is selected
	project := cfg.GetProject("test-project")
	assert.Equal(t, "impl-react-123", project.SelectedImplementation)

	// Test useResource - use feature
	shell.UseResourceForTest([]string{"use", "feature", "feat-login-123"})
	assert.Equal(t, "feature", cfg.Context.Level)
	assert.Equal(t, "feat-login-123", cfg.Context.FeatureBranch)
}

// TestShellCommandHandlers tests the shell's command handlers
func TestShellCommandHandlers(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "cc-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Create a test configuration
	cfg := config.DefaultConfig()

	// Add test projects
	project1 := models.NewProject("test-project", "/tmp/test-project", "Test Project")
	
	// Add implementations to project1
	impl1 := models.Implementation{
		Framework:   "react",
		BranchName:  "impl-react-123",
		Description: "React Implementation",
	}
	
	project1.AddImplementation(impl1)
	project1.SetSelectedImplementation("impl-react-123")
	
	cfg.AddProject(project1)
	cfg.SetActiveProject("test-project")
	cfg.SetContextProject("test-project")
	cfg.SetContextImplementation("impl-react-123")
	
	// Save the configuration
	err = config.SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create a new shell instance
	shell := main.NewShell(cfg, configPath)

	// Mock execution functions would be needed for thorough testing
	// but we'll test what we can without modifying the actual implementations
	
	// Test that the correct level checks are performed in the various handlers
	
	// Reset context to root
	cfg.SetContextLevel("root")
	
	// Test that implementations command requires project context
	executed := shell.HandleImplementationsCommandsForTest([]string{"implementations", "list"})
	assert.False(t, executed, "Should not execute implementation commands at root level")
	
	// Test that features command requires implementation context
	executed = shell.HandleFeaturesCommandsForTest([]string{"features", "list"})
	assert.False(t, executed, "Should not execute feature commands at root level")
	
	// Set to project level
	cfg.SetContextProject("test-project")
	
	// Now implementation commands should work
	executed = shell.HandleImplementationsCommandsForTest([]string{"implementations", "list"})
	assert.True(t, executed, "Should execute implementation commands at project level")
	
	// But feature commands still should not
	executed = shell.HandleFeaturesCommandsForTest([]string{"features", "list"})
	assert.False(t, executed, "Should not execute feature commands at project level")
	
	// Set to implementation level
	cfg.SetContextImplementation("impl-react-123")
	
	// Now feature commands should work
	executed = shell.HandleFeaturesCommandsForTest([]string{"features", "list"})
	assert.True(t, executed, "Should execute feature commands at implementation level")
}