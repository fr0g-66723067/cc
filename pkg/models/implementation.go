package models

import "time"

// Implementation represents a single implementation version of a project
type Implementation struct {
	// Framework used for this implementation
	Framework string `json:"framework"`

	// Name of the branch storing this implementation
	BranchName string `json:"branchName"`

	// Description of what this implementation provides
	Description string `json:"description"`

	// Creation timestamp
	CreatedAt time.Time `json:"createdAt"`

	// AI provider that generated this implementation
	Provider string `json:"provider"`

	// Tags for categorizing implementations
	Tags []string `json:"tags"`

	// Metrics for evaluating implementation
	Metrics map[string]float64 `json:"metrics"`

	// Evaluation score (0-100)
	Score int `json:"score"`

	// Features implemented in this branch
	Features []Feature `json:"features"`
}

// Feature represents a feature added to an implementation
type Feature struct {
	// Name of the feature
	Name string `json:"name"`

	// Name of the branch storing this feature
	BranchName string `json:"branchName"`

	// Description of what this feature provides
	Description string `json:"description"`

	// Creation timestamp
	CreatedAt time.Time `json:"createdAt"`

	// Name of branch this feature is based on
	BaseBranch string `json:"baseBranch"`

	// AI provider that generated this feature
	Provider string `json:"provider"`

	// Status of the feature (e.g. "completed", "in-progress", "failed")
	Status string `json:"status"`

	// Tags for categorizing features
	Tags []string `json:"tags"`
}