package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TempDir creates a temporary directory for testing and returns the path.
// The directory will be automatically removed after the test.
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cc-test-*")
	require.NoError(t, err, "failed to create temp directory")
	
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	
	return dir
}

// CreateTestFile creates a test file with the given content.
func CreateTestFile(t *testing.T, path string, content string) {
	t.Helper()
	
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err, "failed to create directory for test file")
	
	// Write file
	err = os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err, "failed to write test file")
}

// ReadTestFile reads the content of a test file.
func ReadTestFile(t *testing.T, path string) string {
	t.Helper()
	
	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read test file")
	
	return string(data)
}