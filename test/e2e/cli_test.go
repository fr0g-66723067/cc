package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLICommands tests basic CLI commands
func TestCLICommands(t *testing.T) {
	t.Skip("End-to-end test requiring built binary - skipping for now")
	
	// Create a temporary directory for our test project
	projectDir := testutil.TempDir(t)
	
	// Create a configuration directory
	homeDir := testutil.TempDir(t)
	configDir := filepath.Join(homeDir, ".cc")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)
	
	// Set HOME to our temporary directory to control config location
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", homeDir)
	
	// Build the binary
	binaryPath := filepath.Join(homeDir, "cc")
	cmd := exec.Command("go", "build", "-o", binaryPath, "../../../cmd/cc")
	err = cmd.Run()
	require.NoError(t, err, "Failed to build binary")
	
	// Test the init command
	cmd = exec.Command(binaryPath, "init", "test-project")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Init command failed: %s", output)
	assert.Contains(t, string(output), "Initializing project: test-project")
	
	// Verify project directory was created
	projectPath := filepath.Join(projectDir, "test-project")
	_, err = os.Stat(projectPath)
	assert.NoError(t, err, "Project directory wasn't created")
	
	// Test the generate command
	cmd = exec.Command(binaryPath, "generate", "Create a CLI tool", 
		"--frameworks", "go", "--count", "1")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Generate command failed: %s", output)
	assert.Contains(t, string(output), "Generating implementations for: Create a CLI tool")
	
	// Test the list command
	cmd = exec.Command(binaryPath, "list", "implementations")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "List command failed: %s", output)
	assert.Contains(t, string(output), "Listing implementations")
	
	// Test the status command
	cmd = exec.Command(binaryPath, "status")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Status command failed: %s", output)
	assert.Contains(t, string(output), "Current project status")
}