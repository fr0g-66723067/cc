package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI code generation services
type Provider interface {
	// Initialize sets up the AI provider with necessary configuration
	Initialize(ctx context.Context, config map[string]string) error

	// GenerateProject generates a project structure based on description
	GenerateProject(ctx context.Context, description string) (string, error)

	// GenerateImplementation generates code with a specific framework
	GenerateImplementation(ctx context.Context, description string, framework string) (string, error)

	// AddFeature adds a feature to existing code
	AddFeature(ctx context.Context, codeDir string, description string) (string, error)

	// AnalyzeCode analyzes existing code and provides feedback
	AnalyzeCode(ctx context.Context, codeDir string) (string, error)

	// Name returns the provider's name
	Name() string

	// SupportedFrameworks returns the frameworks this provider can work with
	SupportedFrameworks() []string

	// Cleanup performs necessary cleanup operations
	Cleanup(ctx context.Context) error
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
		return nil, fmt.Errorf("unknown AI provider: %s", name)
	}
	return factory(config)
}