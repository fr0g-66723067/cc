package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockProvider is a mock implementation of the ai.Provider interface
type MockProvider struct {
	mock.Mock
}

// Initialize sets up the AI provider with necessary configuration
func (m *MockProvider) Initialize(ctx context.Context, config map[string]string) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// GenerateProject generates a project structure based on description
func (m *MockProvider) GenerateProject(ctx context.Context, description string) (string, error) {
	args := m.Called(ctx, description)
	return args.String(0), args.Error(1)
}

// GenerateImplementation generates code with a specific framework
func (m *MockProvider) GenerateImplementation(ctx context.Context, description string, framework string) (string, error) {
	args := m.Called(ctx, description, framework)
	return args.String(0), args.Error(1)
}

// AddFeature adds a feature to existing code
func (m *MockProvider) AddFeature(ctx context.Context, codeDir string, description string) (string, error) {
	args := m.Called(ctx, codeDir, description)
	return args.String(0), args.Error(1)
}

// AnalyzeCode analyzes existing code and provides feedback
func (m *MockProvider) AnalyzeCode(ctx context.Context, codeDir string) (string, error) {
	args := m.Called(ctx, codeDir)
	return args.String(0), args.Error(1)
}

// Name returns the provider's name
func (m *MockProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

// SupportedFrameworks returns the frameworks this provider can work with
func (m *MockProvider) SupportedFrameworks() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// Cleanup performs necessary cleanup operations
func (m *MockProvider) Cleanup(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
