package models

import (
	"time"
)

// Project represents a Code Controller project
type Project struct {
	// Project name
	Name string `json:"name"`

	// Path to the project directory
	Path string `json:"path"`

	// Current active branch
	ActiveBranch string `json:"activeBranch"`

	// Original project description
	Description string `json:"description"`

	// Creation timestamp
	CreatedAt time.Time `json:"createdAt"`

	// Last modification timestamp
	UpdatedAt time.Time `json:"updatedAt"`

	// Container configuration
	ContainerConfig map[string]string `json:"containerConfig"`

	// AI provider configuration
	AIConfig map[string]string `json:"aiConfig"`

	// Version control configuration
	VCSConfig map[string]string `json:"vcsConfig"`

	// Implementations generated for this project
	Implementations []Implementation `json:"implementations"`

	// Selected implementation (if any)
	SelectedImplementation string `json:"selectedImplementation"`

	// Project settings
	Settings map[string]interface{} `json:"settings"`

	// Current project status
	Status string `json:"status"`

	// Project tags
	Tags []string `json:"tags"`
}

// NewProject creates a new project
func NewProject(name string, path string, description string) *Project {
	now := time.Now()

	return &Project{
		Name:                  name,
		Path:                  path,
		Description:           description,
		ActiveBranch:          "main",
		CreatedAt:             now,
		UpdatedAt:             now,
		ContainerConfig:       make(map[string]string),
		AIConfig:              make(map[string]string),
		VCSConfig:             make(map[string]string),
		Implementations:       []Implementation{},
		SelectedImplementation: "",
		Settings:              make(map[string]interface{}),
		Status:                "initialized",
		Tags:                  []string{},
	}
}

// AddImplementation adds an implementation to the project
func (p *Project) AddImplementation(impl Implementation) {
	p.Implementations = append(p.Implementations, impl)
	p.UpdatedAt = time.Now()
}

// GetImplementation returns an implementation by branch name
func (p *Project) GetImplementation(branchName string) *Implementation {
	for i, impl := range p.Implementations {
		if impl.BranchName == branchName {
			return &p.Implementations[i]
		}
	}
	return nil
}

// SetSelectedImplementation sets the selected implementation
func (p *Project) SetSelectedImplementation(branchName string) {
	p.SelectedImplementation = branchName
	p.UpdatedAt = time.Now()
}

// GetSelectedImplementation returns the selected implementation
func (p *Project) GetSelectedImplementation() *Implementation {
	if p.SelectedImplementation == "" {
		return nil
	}
	return p.GetImplementation(p.SelectedImplementation)
}