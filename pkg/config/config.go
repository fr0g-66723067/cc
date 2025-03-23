package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/user/cc/pkg/models"
)

// Config stores the application configuration
type Config struct {
	// Container provider configuration
	Container struct {
		// Provider type (docker, kubernetes, etc.)
		Provider string `json:"provider"`

		// Provider-specific configuration
		Config map[string]string `json:"config"`

		// Claude Code image to use
		ClaudeImage string `json:"claudeImage"`
	} `json:"container"`

	// AI provider configuration
	AI struct {
		// Provider type (claude, etc.)
		Provider string `json:"provider"`

		// Provider-specific configuration
		Config map[string]string `json:"config"`
	} `json:"ai"`

	// Version control configuration
	VCS struct {
		// Provider type (git, etc.)
		Provider string `json:"provider"`

		// Provider-specific configuration
		Config map[string]string `json:"config"`
	} `json:"vcs"`

	// Projects directory
	ProjectsDir string `json:"projectsDir"`

	// Current active project
	ActiveProject string `json:"activeProject"`

	// Jobs configuration
	Jobs struct {
		// Max concurrent jobs
		MaxConcurrent int `json:"maxConcurrent"`

		// Job timeout in seconds
		Timeout int `json:"timeout"`
	} `json:"jobs"`

	// Plugin configuration
	Plugins struct {
		// Directory to load plugins from
		Dir string `json:"dir"`

		// Enabled plugins
		Enabled []string `json:"enabled"`
	} `json:"plugins"`

	// List of projects
	Projects map[string]*models.Project `json:"projects"`

	// Internal mutex for concurrent access
	mutex sync.RWMutex `json:"-"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	config := &Config{
		ProjectsDir: filepath.Join(homeDir, "cc-projects"),
		ActiveProject: "",
		Projects: make(map[string]*models.Project),
	}

	// Default container config
	config.Container.Provider = "docker"
	config.Container.Config = make(map[string]string)
	config.Container.ClaudeImage = "anthropic/claude-code:latest"

	// Default AI config
	config.AI.Provider = "claude"
	config.AI.Config = make(map[string]string)

	// Default VCS config
	config.VCS.Provider = "git"
	config.VCS.Config = make(map[string]string)

	// Default jobs config
	config.Jobs.MaxConcurrent = 4
	config.Jobs.Timeout = 3600 // 1 hour

	// Default plugins config
	config.Plugins.Dir = filepath.Join(homeDir, ".cc", "plugins")
	config.Plugins.Enabled = []string{}

	return config
}

// LoadConfig loads the configuration from disk
func LoadConfig(path string) (*Config, error) {
	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(config *Config, path string) error {
	config.mutex.RLock()
	defer config.mutex.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0644)
}

// GetProject gets a project by name
func (c *Config) GetProject(name string) *models.Project {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Projects[name]
}

// AddProject adds a project
func (c *Config) AddProject(project *models.Project) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Projects[project.Name] = project
}

// RemoveProject removes a project
func (c *Config) RemoveProject(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Projects, name)
	if c.ActiveProject == name {
		c.ActiveProject = ""
	}
}

// SetActiveProject sets the active project
func (c *Config) SetActiveProject(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ActiveProject = name
}

// GetActiveProject gets the active project
func (c *Config) GetActiveProject() *models.Project {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.ActiveProject == "" {
		return nil
	}
	return c.Projects[c.ActiveProject]
}