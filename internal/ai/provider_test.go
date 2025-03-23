package ai_test

import (
	"context"
	"testing"

	"github.com/fr0g-66723067/cc/internal/ai"
	"github.com/stretchr/testify/assert"
)

// TestProviderFactory tests the provider factory function
func TestProviderFactory(t *testing.T) {
	// Create a mock factory function
	mockFactory := func(config map[string]string) (ai.Provider, error) {
		return &mockProvider{name: "mock"}, nil
	}

	// Register the mock factory
	ai.Register("mock", mockFactory)

	// Test creating a provider
	provider, err := ai.Create("mock", nil)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "mock", provider.Name())

	// Test creating a non-existent provider
	provider, err = ai.Create("nonexistent", nil)
	assert.Error(t, err)
	assert.Nil(t, provider)
}

// mockProvider is a simple implementation of the Provider interface for testing
type mockProvider struct {
	name string
}

func (p *mockProvider) Initialize(ctx context.Context, config map[string]string) error {
	return nil
}

func (p *mockProvider) GenerateProject(ctx context.Context, description string) (string, error) {
	return "generated project", nil
}

func (p *mockProvider) GenerateImplementation(ctx context.Context, description string, framework string) (string, error) {
	return "generated implementation", nil
}

func (p *mockProvider) AddFeature(ctx context.Context, codeDir string, description string) (string, error) {
	return "added feature", nil
}

func (p *mockProvider) AnalyzeCode(ctx context.Context, codeDir string) (string, error) {
	return "code analysis", nil
}

func (p *mockProvider) Name() string {
	return p.name
}

func (p *mockProvider) SupportedFrameworks() []string {
	return []string{"react", "vue"}
}

func (p *mockProvider) Cleanup(ctx context.Context) error {
	return nil
}