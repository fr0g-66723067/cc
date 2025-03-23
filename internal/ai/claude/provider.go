package claude

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fr0g-66723067/cc/internal/ai"
	"github.com/fr0g-66723067/cc/internal/container"
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

	// Create a workspace directory in the container
	workspacePath := "/workspace"
	createDirCmd := []string{"mkdir", "-p", workspacePath}
	_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, createDirCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Execute command in container with proper Claude Code CLI arguments
	// The Claude Code CLI syntax is typically:
	// claude code generate [--output DIR] "prompt"
	cmd := []string{"claude", "code", "generate", "--output", workspacePath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Log the generation output
	fmt.Printf("Claude generation complete. Output summary:\n%s\n", 
		truncateString(output, 500))
	
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
			"Focus on best practices, performance, and maintainability.\n\n"+
			"Make sure to adhere to the following guidelines:\n"+
			"1. Use modern %s patterns and libraries\n"+
			"2. Include proper error handling\n"+
			"3. Add comments explaining key logic\n"+
			"4. Include necessary dependencies and configuration\n"+
			"5. Implement a modular, maintainable architecture",
		framework, description, framework, framework,
	)

	// Create a clean workspace directory in the container
	workspacePath := "/workspace"
	cleanCmd := []string{"rm", "-rf", workspacePath + "/*"}
	_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cleanCmd)
	if err != nil {
		return "", fmt.Errorf("failed to clean workspace directory: %w", err)
	}

	// Create workspace directory
	createDirCmd := []string{"mkdir", "-p", workspacePath}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, createDirCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Execute command in container with proper Claude Code CLI arguments
	fmt.Printf("Generating %s implementation for: %s\n", framework, description)
	cmd := []string{"claude", "code", "generate", "--output", workspacePath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Log the generation output
	fmt.Printf("Claude generation complete. Output summary:\n%s\n", 
		truncateString(output, 500))

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

	// Clean workspace directory in container
	cleanCmd := []string{"rm", "-rf", containerPath + "/*"}
	_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cleanCmd)
	if err != nil {
		return "", fmt.Errorf("failed to clean workspace directory: %w", err)
	}

	// Create workspace directory
	createDirCmd := []string{"mkdir", "-p", containerPath}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, createDirCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Copy files to container
	fmt.Printf("Copying project files from %s to container...\n", localPath)
	if err := p.containerProvider.CopyFilesToContainer(ctx, p.containerID, localPath, containerPath); err != nil {
		return "", fmt.Errorf("failed to copy files to container: %w", err)
	}

	// Verify files were copied correctly
	lsCmd := []string{"ls", "-la", containerPath}
	lsOutput, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, lsCmd)
	if err != nil {
		return "", fmt.Errorf("failed to list files in container: %w", err)
	}
	fmt.Printf("Files in container workspace:\n%s\n", lsOutput)

	// Create prompt for Claude
	prompt := fmt.Sprintf(
		"Add a new feature to the existing code: %s\n\n"+
			"Please analyze the existing codebase in %s and add a new feature that: %s\n\n"+
			"Make all necessary changes to implement this feature completely while maintaining the following:\n"+
			"1. Keep the existing code style and architecture\n"+
			"2. Follow the same patterns as existing code\n"+
			"3. Add appropriate error handling\n"+
			"4. Include unit tests for new functionality\n"+
			"5. Update documentation as needed\n"+
			"6. Ensure the feature is fully integrated with the existing functionality\n\n"+
			"Describe your changes in detail and explain your implementation choices.",
		description, containerPath, description,
	)

	// Execute command in container
	fmt.Printf("Asking Claude to add feature: %s\n", description)
	cmd := []string{"claude", "code", "modify", "--dir", containerPath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Log the changes that Claude made
	fmt.Printf("Claude made the following changes:\n%s\n", 
		truncateString(output, 500))

	// Copy files back from container
	fmt.Printf("Copying modified files back from container to %s...\n", localPath)
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

	// Clean workspace directory in container
	cleanCmd := []string{"rm", "-rf", containerPath + "/*"}
	_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cleanCmd)
	if err != nil {
		return "", fmt.Errorf("failed to clean workspace directory: %w", err)
	}

	// Create workspace directory
	createDirCmd := []string{"mkdir", "-p", containerPath}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, createDirCmd)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Copy files to container
	fmt.Printf("Copying project files from %s to container for analysis...\n", localPath)
	if err := p.containerProvider.CopyFilesToContainer(ctx, p.containerID, localPath, containerPath); err != nil {
		return "", fmt.Errorf("failed to copy files to container: %w", err)
	}

	// Verify files were copied correctly
	lsCmd := []string{"find", containerPath, "-type", "f", "-not", "-path", "*/\\.*", "|", "sort"}
	lsOutput, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, lsCmd)
	if err != nil {
		return "", fmt.Errorf("failed to list files in container: %w", err)
	}
	fmt.Printf("Files in container for analysis:\n%s\n", 
		truncateString(lsOutput, 200))

	// Create prompt for Claude
	prompt := "Analyze the codebase in /workspace\n\n" +
		"Please provide a comprehensive analysis of the codebase including:\n" +
		"1. Overall architecture and design patterns used\n" +
		"2. Code quality assessment with specific examples\n" +
		"3. Potential issues or improvements, prioritized by impact\n" +
		"4. Security concerns and recommendations\n" +
		"5. Performance considerations and possible optimizations\n" +
		"6. Maintainability assessment and technical debt\n" +
		"7. Documentation quality and completeness\n" +
		"8. Test coverage and quality assessment\n\n" +
		"Format your analysis as a structured report with clear sections and bullet points."

	// Execute command in container
	fmt.Printf("Starting Claude code analysis...\n")
	cmd := []string{"claude", "code", "analyze", "--dir", containerPath, prompt}
	output, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	fmt.Printf("Code analysis complete. Generated %d characters of analysis.\n", len(output))

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
		// Check if container is still running
		pingCmd := []string{"echo", "ping"}
		_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, pingCmd)
		if err == nil {
			// Container is still running
			return nil
		}
		// If we get here, the container is no longer available
		p.containerID = ""
	}

	// Get container image
	image := p.config["claude_image"]
	if image == "" {
		// Use local mock image if available, otherwise fall back to official image
		mockImage := "claude-code-mock:latest"
		// Check if mock image exists
		checkImageCmd := exec.Command("docker", "images", "-q", mockImage)
		output, err := checkImageCmd.Output()
		if err == nil && len(output) > 0 {
			image = mockImage
			fmt.Printf("Using local mock Claude Code image: %s\n", mockImage)
		} else {
			image = "anthropic/claude-code:latest"
			fmt.Printf("Local mock image not found, trying to use: %s\n", image)
		}
	}

	// Create a persistent temporary directory for workspace
	tmpDir := filepath.Join("/tmp", fmt.Sprintf("claude-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Set up volume mounts
	volumeMounts := make(map[string]string)
	volumeMounts[tmpDir] = "/workspace"

	// Set up environment variables
	env := make(map[string]string)
	// Add Claude API key if provided
	if apiKey, ok := p.config["claude_api_key"]; ok && apiKey != "" {
		env["CLAUDE_API_KEY"] = apiKey
	}
	
	// Add any other environment variables from config
	for k, v := range p.config {
		if strings.HasPrefix(k, "env_") {
			env[strings.TrimPrefix(k, "env_")] = v
		}
	}

	// Ensure environment has minimum required variables
	env["CLAUDE_CLI_LOG_LEVEL"] = "info"  // Set logging level

	// Run container with restart policy to ensure it stays running
	fmt.Printf("Starting Claude container with image: %s\n", image)
	containerID, err := p.containerProvider.RunContainer(ctx, image, volumeMounts, env)
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	p.containerID = containerID
	fmt.Printf("Claude container started with ID: %s\n", containerID)
	
	// Wait a moment for container to initialize
	time.Sleep(2 * time.Second)
	
	return nil
}

// truncateString truncates a string to a maximum length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Register registers this provider factory
func init() {
	ai.Register("claude", func(config map[string]string) (ai.Provider, error) {
		return NewProvider(config)
	})
}