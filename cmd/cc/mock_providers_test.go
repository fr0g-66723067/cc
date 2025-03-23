package main

import (
	"context"

	"github.com/fr0g-66723067/cc/internal/ai"
	"github.com/fr0g-66723067/cc/internal/vcs"
)

// Initialize mock providers for testing
func init() {
	// Register mock VCS provider
	vcs.Register("git", func(config map[string]string) (vcs.Provider, error) {
		return &mockVCSProvider{
			config: config,
		}, nil
	})

	// Register mock AI provider
	ai.Register("claude", func(config map[string]string) (ai.Provider, error) {
		return &mockAIProvider{
			config: config,
		}, nil
	})
}

// mockVCSProvider implements a mock VCS provider for testing
type mockVCSProvider struct {
	config map[string]string
	path   string
}

// Initialize initializes the mock VCS
func (m *mockVCSProvider) Initialize(path string) error {
	m.path = path
	return nil
}

// CreateBranch creates a branch in the mock VCS
func (m *mockVCSProvider) CreateBranch(name string, baseBranch string) error {
	return nil
}

// SwitchBranch switches to a branch in the mock VCS
func (m *mockVCSProvider) SwitchBranch(name string) error {
	return nil
}

// GetCurrentBranch returns the current branch in the mock VCS
func (m *mockVCSProvider) GetCurrentBranch() (string, error) {
	return "main", nil
}

// ListBranches lists all branches in the mock VCS
func (m *mockVCSProvider) ListBranches() ([]string, error) {
	return []string{"main"}, nil
}

// AddFiles adds files to the mock VCS
func (m *mockVCSProvider) AddFiles(paths []string) error {
	return nil
}

// CommitChanges commits changes to the mock VCS
func (m *mockVCSProvider) CommitChanges(message string) error {
	return nil
}

// HasChanges returns whether the mock VCS has changes
func (m *mockVCSProvider) HasChanges() (bool, error) {
	return false, nil
}

// GetBranchMetadata gets metadata for a branch in the mock VCS
func (m *mockVCSProvider) GetBranchMetadata(branch string) (map[string]string, error) {
	return make(map[string]string), nil
}

// SetBranchMetadata sets metadata for a branch in the mock VCS
func (m *mockVCSProvider) SetBranchMetadata(branch string, metadata map[string]string) error {
	return nil
}

// ExportDiff exports a diff between branches in the mock VCS
func (m *mockVCSProvider) ExportDiff(fromBranch, toBranch string) (string, error) {
	return "Mock diff between " + fromBranch + " and " + toBranch, nil
}

// Name returns the name of the mock VCS provider
func (m *mockVCSProvider) Name() string {
	return "git"
}

// mockAIProvider implements a mock AI provider for testing
type mockAIProvider struct {
	config     map[string]string
	frameworks []string
}

// Initialize initializes the mock AI provider
func (m *mockAIProvider) Initialize(ctx context.Context, config map[string]string) error {
	m.frameworks = []string{"react", "vue", "angular"}
	return nil
}

// GenerateProject generates a project in the mock AI provider
func (m *mockAIProvider) GenerateProject(ctx context.Context, description string) (string, error) {
	return "Mock project for " + description, nil
}

// GenerateImplementation generates an implementation in the mock AI provider
func (m *mockAIProvider) GenerateImplementation(ctx context.Context, description string, framework string) (string, error) {
	return "Mock implementation for " + description + " using " + framework, nil
}

// AddFeature adds a feature in the mock AI provider
func (m *mockAIProvider) AddFeature(ctx context.Context, codeDir string, description string) (string, error) {
	return "Mock feature for " + description, nil
}

// AnalyzeCode analyzes code in the mock AI provider
func (m *mockAIProvider) AnalyzeCode(ctx context.Context, codeDir string) (string, error) {
	return "Mock analysis for " + codeDir, nil
}

// Name returns the name of the mock AI provider
func (m *mockAIProvider) Name() string {
	return "claude"
}

// SupportedFrameworks returns the frameworks the mock AI provider supports
func (m *mockAIProvider) SupportedFrameworks() []string {
	return m.frameworks
}

// Cleanup cleans up the mock AI provider
func (m *mockAIProvider) Cleanup(ctx context.Context) error {
	return nil
}