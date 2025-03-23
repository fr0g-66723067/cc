package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/pkg/config"
	"github.com/fr0g-66723067/cc/pkg/models"
	"github.com/fr0g-66723067/cc/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	// Test the default configuration
	cfg := config.DefaultConfig()
	
	// Check default values
	assert.Equal(t, "docker", cfg.Container.Provider)
	assert.Equal(t, "claude", cfg.AI.Provider)
	assert.Equal(t, "git", cfg.VCS.Provider)
	assert.Equal(t, "anthropic/claude-code:latest", cfg.Container.ClaudeImage)
	assert.Equal(t, 4, cfg.Jobs.MaxConcurrent)
	assert.Greater(t, cfg.Jobs.Timeout, 0)
	assert.NotEmpty(t, cfg.ProjectsDir)
	assert.Empty(t, cfg.ActiveProject)
	assert.Empty(t, cfg.Projects)
	
	// Check that maps are initialized
	assert.NotNil(t, cfg.Container.Config)
	assert.NotNil(t, cfg.AI.Config)
	assert.NotNil(t, cfg.VCS.Config)
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)
	configPath := filepath.Join(tempDir, "config.json")
	
	// Test loading non-existent file
	cfg, err := config.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	
	// Create a test config file
	testConfig := `{
		"container": {
			"provider": "test-container",
			"config": {
				"key1": "value1"
			},
			"claudeImage": "test/claude:latest"
		},
		"ai": {
			"provider": "test-ai",
			"config": {
				"key2": "value2"
			}
		},
		"vcs": {
			"provider": "test-vcs",
			"config": {
				"key3": "value3"
			}
		},
		"projectsDir": "/test/projects",
		"activeProject": "test-project"
	}`
	err = os.WriteFile(configPath, []byte(testConfig), 0644)
	require.NoError(t, err)
	
	// Test loading existing file
	cfg, err = config.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	
	// Check loaded values
	assert.Equal(t, "test-container", cfg.Container.Provider)
	assert.Equal(t, "test-ai", cfg.AI.Provider)
	assert.Equal(t, "test-vcs", cfg.VCS.Provider)
	assert.Equal(t, "test/claude:latest", cfg.Container.ClaudeImage)
	assert.Equal(t, "/test/projects", cfg.ProjectsDir)
	assert.Equal(t, "test-project", cfg.ActiveProject)
	
	// Check map values
	assert.Equal(t, "value1", cfg.Container.Config["key1"])
	assert.Equal(t, "value2", cfg.AI.Config["key2"])
	assert.Equal(t, "value3", cfg.VCS.Config["key3"])
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)
	configPath := filepath.Join(tempDir, "config.json")
	
	// Create a test config
	cfg := config.DefaultConfig()
	cfg.Container.Provider = "test-container"
	cfg.AI.Provider = "test-ai"
	cfg.VCS.Provider = "test-vcs"
	cfg.Container.ClaudeImage = "test/claude:latest"
	cfg.ProjectsDir = "/test/projects"
	cfg.ActiveProject = "test-project"
	
	// Add some map values
	cfg.Container.Config["key1"] = "value1"
	cfg.AI.Config["key2"] = "value2"
	cfg.VCS.Config["key3"] = "value3"
	
	// Save the config
	err := config.SaveConfig(cfg, configPath)
	assert.NoError(t, err)
	
	// Check that the file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)
	
	// Load the config back and verify its contents
	loadedCfg, err := config.LoadConfig(configPath)
	assert.NoError(t, err)
	
	// Check values
	assert.Equal(t, "test-container", loadedCfg.Container.Provider)
	assert.Equal(t, "test-ai", loadedCfg.AI.Provider)
	assert.Equal(t, "test-vcs", loadedCfg.VCS.Provider)
	assert.Equal(t, "test/claude:latest", loadedCfg.Container.ClaudeImage)
	assert.Equal(t, "/test/projects", loadedCfg.ProjectsDir)
	assert.Equal(t, "test-project", loadedCfg.ActiveProject)
	
	// Check map values
	assert.Equal(t, "value1", loadedCfg.Container.Config["key1"])
	assert.Equal(t, "value2", loadedCfg.AI.Config["key2"])
	assert.Equal(t, "value3", loadedCfg.VCS.Config["key3"])
}

func TestProjectManagement(t *testing.T) {
	// Create a config
	cfg := config.DefaultConfig()
	
	// Create a test project
	project := models.NewProject("test-project", "/path/to/project", "Test project")
	
	// Add the project
	cfg.AddProject(project)
	
	// Check that the project was added
	assert.Len(t, cfg.Projects, 1)
	assert.Contains(t, cfg.Projects, "test-project")
	
	// Get the project
	retrievedProject := cfg.GetProject("test-project")
	assert.NotNil(t, retrievedProject)
	assert.Equal(t, "test-project", retrievedProject.Name)
	assert.Equal(t, "/path/to/project", retrievedProject.Path)
	
	// Set the active project
	cfg.SetActiveProject("test-project")
	assert.Equal(t, "test-project", cfg.ActiveProject)
	
	// Get the active project
	activeProject := cfg.GetActiveProject()
	assert.NotNil(t, activeProject)
	assert.Equal(t, "test-project", activeProject.Name)
	
	// Remove the project
	cfg.RemoveProject("test-project")
	assert.Empty(t, cfg.Projects)
	assert.Empty(t, cfg.ActiveProject)
	assert.Nil(t, cfg.GetActiveProject())
}