package claude

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/cc/internal/container"
)

// Provider implements the AI provider interface for Claude
type Provider struct {
	containerProvider container.Provider
	containerID       string
	config            map[string]string
	frameworks        []string
}

// NewProvider creates a new Claude provider
func NewProvider(config map[string]string) (*Provider, error) {
	// Set default frameworks if not specified
	frameworks := []string{
		"react", "vue", "svelte", "angular", 
		"nextjs", "nuxt", "express", "fastify",
		"django", "flask", "spring", "rails",
	}

	if customFrameworks, ok := config["frameworks"]; ok && customFrameworks != "" {
		frameworks = strings.Split(customFrameworks, ",")
	}

	return &Provider{
		config:     config,
		frameworks: frameworks,
	}, nil
}

// Initialize sets up the Claude provider
func (p *Provider) Initialize(ctx context.Context, config map[string]string) error {
	// Merge configs
	for k, v := range config {
		p.config[k] = v
	}

	// Create container provider
	containerType := p.config["container_provider"]
	if containerType == "" {
		containerType = "docker"
	}

	containerConfig := make(map[string]string)
	for k, v := range p.config {
		if strings.HasPrefix(k, "container_") {
			containerConfig[strings.TrimPrefix(k, "container_")] = v
		}
	}

	var err error
	p.containerProvider, err = container.Create(containerType, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container provider: %w", err)
	}

	// Initialize container provider
	if err := p.containerProvider.Initialize(ctx, containerConfig); err != nil {
		return fmt.Errorf("failed to initialize container provider: %w", err)
	}

	return nil
}

// GenerateProject generates a project structure based on description
func (p *Provider) GenerateProject(ctx context.Context, description string) (string, error) {
	// Make sure we have a container
	if err := p.ensureContainer(ctx); err != nil {
		return "", err
	}

	// Create prompt for Claude
	prompt := fmt.Sprintf(
		"Create a project structure for: %s\n\n"+
			"Please create the necessary directory structure, configuration files, and basic scaffolding "+
			"for a new project based on this description. Focus on setting up a solid foundation "+
			"that can be used for multiple implementation approaches.",
		description,
	)

	// Execute command in container
	cmd := []string{"claude", "code", "generate", prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}

// GenerateImplementation generates code with a specific framework
func (p *Provider) GenerateImplementation(ctx context.Context, description string, framework string) (string, error) {
	// Make sure we have a container
	if err := p.ensureContainer(ctx); err != nil {
		return "", err
	}

	// Create prompt for Claude
	prompt := fmt.Sprintf(
		"Create a %s implementation for: %s\n\n"+
			"Please implement a complete, working application using the %s framework based on this description. "+
			"Include all necessary files, configuration, and code to make the application fully functional. "+
			"Focus on best practices, performance, and maintainability.",
		framework, description, framework,
	)

	// Execute command in container
	cmd := []string{"claude", "code", "generate", prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}

// AddFeature adds a feature to existing code
func (p *Provider) AddFeature(ctx context.Context, codeDir string, description string) (string, error) {
	// Make sure we have a container
	if err := p.ensureContainer(ctx); err != nil {
		return "", err
	}

	// Create relative path for code directory
	containerPath := "/workspace"
	localPath := codeDir

	// Copy files to container
	if err := p.containerProvider.CopyFilesToContainer(ctx, p.containerID, localPath, containerPath); err != nil {
		return "", fmt.Errorf("failed to copy files to container: %w", err)
	}

	// Create prompt for Claude
	prompt := fmt.Sprintf(
		"Add a new feature to the existing code: %s\n\n"+
			"Please analyze the existing codebase in /workspace and add a new feature that: %s\n"+
			"Make all necessary changes to implement this feature completely while maintaining the existing code style "+
			"and architecture. Ensure the feature is fully integrated with the existing functionality.",
		description, description,
	)

	// Execute command in container
	cmd := []string{"claude", "code", "modify", "--dir", containerPath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Copy files back from container
	if err := p.containerProvider.CopyFilesFromContainer(ctx, p.containerID, containerPath, localPath); err != nil {
		return "", fmt.Errorf("failed to copy files from container: %w", err)
	}

	return output, nil
}

// AnalyzeCode analyzes existing code and provides feedback
func (p *Provider) AnalyzeCode(ctx context.Context, codeDir string) (string, error) {
	// Make sure we have a container
	if err := p.ensureContainer(ctx); err != nil {
		return "", err
	}

	// Create relative path for code directory
	containerPath := "/workspace"
	localPath := codeDir

	// Copy files to container
	if err := p.containerProvider.CopyFilesToContainer(ctx, p.containerID, localPath, containerPath); err != nil {
		return "", fmt.Errorf("failed to copy files to container: %w", err)
	}

	// Create prompt for Claude
	prompt := "Analyze the codebase in /workspace\n\n" +
		"Please provide a comprehensive analysis of the codebase including:\n" +
		"1. Overall architecture and design patterns\n" +
		"2. Code quality assessment\n" +
		"3. Potential issues or improvements\n" +
		"4. Security concerns\n" +
		"5. Performance considerations\n" +
		"6. Maintainability assessment"

	// Execute command in container
	cmd := []string{"claude", "code", "analyze", "--dir", containerPath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}

// Name returns the provider's name
func (p *Provider) Name() string {
	return "claude"
}

// SupportedFrameworks returns the frameworks this provider can work with
func (p *Provider) SupportedFrameworks() []string {
	return p.frameworks
}

// Cleanup performs necessary cleanup operations
func (p *Provider) Cleanup(ctx context.Context) error {
	if p.containerID != "" {
		// Stop container
		if err := p.containerProvider.StopContainer(ctx, p.containerID); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}

		// Remove container
		if err := p.containerProvider.RemoveContainer(ctx, p.containerID); err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}

		p.containerID = ""
	}

	return nil
}

// ensureContainer ensures a container is running
func (p *Provider) ensureContainer(ctx context.Context) error {
	if p.containerID != "" {
		// Container already running
		return nil
	}

	// Get container image
	image := p.config["claude_image"]
	if image == "" {
		image = "anthropic/claude-code:latest"
	}

	// Set up volume mounts
	volumeMounts := make(map[string]string)
	tmpDir := filepath.Join("/tmp", fmt.Sprintf("claude-%d", time.Now().UnixNano()))
	volumeMounts[tmpDir] = "/workspace"

	// Set up environment variables
	env := make(map[string]string)
	for k, v := range p.config {
		if strings.HasPrefix(k, "env_") {
			env[strings.TrimPrefix(k, "env_")] = v
		}
	}

	// Run container
	containerID, err := p.containerProvider.RunContainer(ctx, image, volumeMounts, env)
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	p.containerID = containerID
	return nil
}

// Register registers this provider factory
func init() {
	ai.Register("claude", func(config map[string]string) (ai.Provider, error) {
		return NewProvider(config)
	})
}