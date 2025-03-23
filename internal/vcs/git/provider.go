package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Provider implements the VCS provider interface for Git
type Provider struct {
	repo     *git.Repository
	repoPath string
	config   map[string]string
}

// NewProvider creates a new Git provider
func NewProvider(config map[string]string) (*Provider, error) {
	return &Provider{
		config: config,
	}, nil
}

// Initialize initializes a Git repository
func (p *Provider) Initialize(path string) error {
	p.repoPath = path

	// Check if repo exists
	_, err := os.Stat(filepath.Join(path, ".git"))
	if err == nil {
		// Repository exists, open it
		repo, err := git.PlainOpen(path)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}
		p.repo = repo
		return nil
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Initialize new repository
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	p.repo = repo

	// Create initial commit
	wt, err := p.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create README file
	readmePath := filepath.Join(path, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Project\n\nInitialized by Code Controller\n"), 0644); err != nil {
		return fmt.Errorf("failed to create README file: %w", err)
	}

	// Add README file
	if _, err := wt.Add("README.md"); err != nil {
		return fmt.Errorf("failed to add README file: %w", err)
	}

	// Commit changes
	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Code Controller",
			Email: "cc@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

// CreateBranch creates a new branch
func (p *Provider) CreateBranch(name string, baseBranch string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	// Get the base branch reference
	var baseRef *plumbing.Reference
	if baseBranch == "" {
		// Use HEAD as base
		ref, err := p.repo.Head()
		if err != nil {
			return fmt.Errorf("failed to get HEAD: %w", err)
		}
		baseRef = ref
	} else {
		// Use specified base branch
		ref, err := p.repo.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
		if err != nil {
			return fmt.Errorf("failed to get base branch reference: %w", err)
		}
		baseRef = ref
	}

	// Create branch reference
	refName := plumbing.NewBranchReferenceName(name)
	ref := plumbing.NewHashReference(refName, baseRef.Hash())

	// Create branch
	if err := p.repo.Storer.SetReference(ref); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// SwitchBranch switches to a branch
func (p *Provider) SwitchBranch(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	// Get worktree
	wt, err := p.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Check if branch exists
	_, err = p.repo.Reference(plumbing.NewBranchReferenceName(name), true)
	if err != nil {
		return fmt.Errorf("branch does not exist: %w", err)
	}

	// Switch branch
	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Force:  false,
	}
	if err := wt.Checkout(opts); err != nil {
		return fmt.Errorf("failed to switch branch: %w", err)
	}

	return nil
}

// GetCurrentBranch returns the current branch name
func (p *Provider) GetCurrentBranch() (string, error) {
	// Get HEAD reference
	ref, err := p.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Check if HEAD is a branch
	if ref.Name().IsBranch() {
		return ref.Name().Short(), nil
	}

	return "", fmt.Errorf("HEAD is not a branch")
}

// ListBranches lists all branches
func (p *Provider) ListBranches() ([]string, error) {
	// Get branches iterator
	branchesIter, err := p.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// Collect branch names
	var branches []string
	if err := branchesIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return branches, nil
}

// AddFiles adds files to be committed
func (p *Provider) AddFiles(paths []string) error {
	// Get worktree
	wt, err := p.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add each file
	for _, path := range paths {
		// Get relative path to repository root
		relPath, err := filepath.Rel(p.repoPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Add file
		if _, err := wt.Add(relPath); err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}

	return nil
}

// CommitChanges commits changes with a message
func (p *Provider) CommitChanges(message string) error {
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	// Get worktree
	wt, err := p.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Commit changes
	_, err = wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Code Controller",
			Email: "cc@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

// HasChanges returns whether there are uncommitted changes
func (p *Provider) HasChanges() (bool, error) {
	// Get worktree
	wt, err := p.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get status
	status, err := wt.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get status: %w", err)
	}

	return !status.IsClean(), nil
}

// GetBranchMetadata gets metadata for a branch
func (p *Provider) GetBranchMetadata(branch string) (map[string]string, error) {
	metadataPath := filepath.Join(p.repoPath, ".git", "cc", "metadata", branch+".json")

	// Check if metadata file exists
	_, err := os.Stat(metadataPath)
	if os.IsNotExist(err) {
		// No metadata yet
		return make(map[string]string), nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check metadata file: %w", err)
	}

	// Read metadata file
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Parse metadata
	var metadata map[string]string
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return metadata, nil
}

// SetBranchMetadata sets metadata for a branch
func (p *Provider) SetBranchMetadata(branch string, metadata map[string]string) error {
	// Create metadata directory if it doesn't exist
	metadataDir := filepath.Join(p.repoPath, ".git", "cc", "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Marshal metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write metadata file
	metadataPath := filepath.Join(metadataDir, branch+".json")
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// ExportDiff exports a diff between branches
func (p *Provider) ExportDiff(fromBranch, toBranch string) (string, error) {
	// Get repository
	if p.repo == nil {
		return "", fmt.Errorf("repository not initialized")
	}

	// Get from branch reference
	fromRef, err := p.repo.Reference(plumbing.NewBranchReferenceName(fromBranch), true)
	if err != nil {
		return "", fmt.Errorf("failed to get from branch reference: %w", err)
	}

	// Get to branch reference
	toRef, err := p.repo.Reference(plumbing.NewBranchReferenceName(toBranch), true)
	if err != nil {
		return "", fmt.Errorf("failed to get to branch reference: %w", err)
	}

	// Get from commit
	fromCommit, err := p.repo.CommitObject(fromRef.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get from commit: %w", err)
	}

	// Get to commit
	toCommit, err := p.repo.CommitObject(toRef.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get to commit: %w", err)
	}

	// Get diff
	patch, err := fromCommit.Patch(toCommit)
	if err != nil {
		return "", fmt.Errorf("failed to get patch: %w", err)
	}

	return patch.String(), nil
}

// Name returns the provider's name
func (p *Provider) Name() string {
	return "git"
}

// Register registers this provider factory
func init() {
	vcs.Register("git", func(config map[string]string) (vcs.Provider, error) {
		return NewProvider(config)
	})
}