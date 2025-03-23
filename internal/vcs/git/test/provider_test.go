package git_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/internal/vcs/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewProvider tests creating a new Git provider
func TestNewProvider(t *testing.T) {
	provider, err := git.NewProvider(nil)
	require.NoError(t, err)
	require.NotNil(t, provider)
	
	assert.Equal(t, "git", provider.Name())
}

// TestInitialize tests initializing a Git repository
func TestInitialize(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := ioutil.TempDir("", "git-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a Git provider
	provider, err := git.NewProvider(nil)
	require.NoError(t, err)
	
	// Initialize the repository
	err = provider.Initialize(tempDir)
	require.NoError(t, err)
	
	// Verify that the repository was created
	_, err = os.Stat(filepath.Join(tempDir, ".git"))
	assert.False(t, os.IsNotExist(err), "Repository .git directory should exist")
	
	// Verify that the README file was created
	_, err = os.Stat(filepath.Join(tempDir, "README.md"))
	assert.False(t, os.IsNotExist(err), "README.md should exist")
}

// TestBranchOperations tests creating and switching branches
func TestBranchOperations(t *testing.T) {
	t.Skip("Skip branch operations test - requires actual Git client")
}

// TestMetadataOperations tests getting and setting branch metadata
func TestMetadataOperations(t *testing.T) {
	t.Skip("Skip metadata operations test - requires actual Git client")
}

// TestAddCommitFiles tests adding and committing files
func TestAddCommitFiles(t *testing.T) {
	t.Skip("Skip add/commit test - requires actual Git client")
}

// TestExportDiff tests exporting diffs between branches
func TestExportDiff(t *testing.T) {
	t.Skip("Skip export diff test - requires actual Git client")
}