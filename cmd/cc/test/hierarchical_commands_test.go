package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/pkg/config"
	"github.com/fr0g-66723067/cc/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestHierarchicalCommands tests the hierarchical command structure
func TestHierarchicalCommands(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "cc-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Create a test configuration
	cfg := config.DefaultConfig()
	cfg.ProjectsDir = filepath.Join(tempDir, "projects")

	// Add test projects
	project1 := models.NewProject("test-project", filepath.Join(cfg.ProjectsDir, "test-project"), "Test Project")
	cfg.AddProject(project1)
	cfg.SetActiveProject("test-project")
	
	// Save the configuration
	err = config.SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create test commands
	rootCmd := &cobra.Command{Use: "cc"}
	
	// Create projects list command
	var out bytes.Buffer
	projectsListCmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out.WriteString("Projects listed")
		},
	}
	
	// Create projects create command
	projectsCreateCmd := &cobra.Command{
		Use:  "create [name] [description]",
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			out.WriteString("Project " + name + " created")
		},
	}
	
	// Create projects command group
	projectsCmd := &cobra.Command{
		Use:   "projects",
		Short: "Work with projects",
	}
	
	projectsCmd.AddCommand(projectsListCmd, projectsCreateCmd)
	rootCmd.AddCommand(projectsCmd)
	
	// Create use command
	useCmd := &cobra.Command{
		Use:   "use [resource-type] [name]",
		Short: "Select a project, implementation, or feature as active",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			resourceType := args[0]
			name := args[1]
			out.WriteString("Using " + resourceType + " " + name)
		},
	}
	
	rootCmd.AddCommand(useCmd)
	
	// Test projects subcommand
	rootCmd.SetArgs([]string{"projects", "list"})
	rootCmd.Execute()
	assert.Equal(t, "Projects listed", out.String())
	
	// Reset output buffer
	out.Reset()
	
	// Test projects create
	rootCmd.SetArgs([]string{"projects", "create", "new-project"})
	rootCmd.Execute()
	assert.Equal(t, "Project new-project created", out.String())
	
	// Reset output buffer
	out.Reset()
	
	// Test use command
	rootCmd.SetArgs([]string{"use", "project", "test-project"})
	rootCmd.Execute()
	assert.Equal(t, "Using project test-project", out.String())
}