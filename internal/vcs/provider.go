package vcs

import (
	"fmt"
)

// Provider defines the interface for version control systems
type Provider interface {
	// Initialize initializes a repository
	Initialize(path string) error

	// CreateBranch creates a new branch
	CreateBranch(name string, baseBranch string) error

	// SwitchBranch switches to a branch
	SwitchBranch(name string) error

	// GetCurrentBranch returns the current branch name
	GetCurrentBranch() (string, error)

	// ListBranches lists all branches
	ListBranches() ([]string, error)

	// AddFiles adds files to be committed
	AddFiles(paths []string) error

	// CommitChanges commits changes with a message
	CommitChanges(message string) error

	// HasChanges returns whether there are uncommitted changes
	HasChanges() (bool, error)

	// GetBranchMetadata gets metadata for a branch
	GetBranchMetadata(branch string) (map[string]string, error)

	// SetBranchMetadata sets metadata for a branch
	SetBranchMetadata(branch string, metadata map[string]string) error

	// ExportDiff exports a diff between branches
	ExportDiff(fromBranch, toBranch string) (string, error)

	// Name returns the provider's name
	Name() string
}

// Factory creates a provider based on name
type Factory func(config map[string]string) (Provider, error)

var providers = make(map[string]Factory)

// Register registers a provider factory
func Register(name string, factory Factory) {
	providers[name] = factory
}

// Create creates a provider with the given name
func Create(name string, config map[string]string) (Provider, error) {
	factory, exists := providers[name]
	if !exists {
		return nil, fmt.Errorf("unknown VCS provider: %s", name)
	}
	return factory(config)
}