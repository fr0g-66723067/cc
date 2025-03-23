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

	// Set permissions separately
	chmodCmd := []string{"chmod", "777", workspacePath}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, chmodCmd)
	if err != nil {
		fmt.Printf("Warning: Failed to set workspace permissions: %v\n", err)
	}

	// Verify workspace directory was created and has correct permissions
	lsCmd := []string{"ls", "-la", "/"}
	lsOutput, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, lsCmd)
	if err != nil {
		fmt.Printf("Warning: Failed to list root directory: %v\n", err)
	} else {
		fmt.Printf("Root directory contents:\n%s\n", lsOutput)
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
	// Use absolute path for source
	absLocalPath, err := filepath.Abs(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for local directory: %w", err)
	}

	// Make sure the target directory exists and has correct permissions
	mkdirCmd := []string{"mkdir", "-p", containerPath, "&&", "chmod", "755", containerPath}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, mkdirCmd)
	if err != nil {
		fmt.Printf("Warning: Failed to prepare container directory: %v\n", err)
	}

	// Now copy files
	if err := p.containerProvider.CopyFilesToContainer(ctx, p.containerID, absLocalPath, containerPath); err != nil {
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
	// Make sure all files in the container have the right permissions
	chmodCmd := []string{"find", containerPath, "-type", "f", "-exec", "chmod", "644", "{}", ";"}
	_, err = p.containerProvider.ExecuteCommand(ctx, p.containerID, chmodCmd)
	if err != nil {
		fmt.Printf("Warning: Failed to set file permissions in container: %v\n", err)
	}

	// Create target directory if it doesn't exist
	absLocalPath, pathErr := filepath.Abs(localPath)
	if pathErr != nil {
		return "", fmt.Errorf("failed to get absolute path for local directory: %w", pathErr)
	}

	if err := os.MkdirAll(absLocalPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create local directory: %w", err)
	}

	// Copy files back
	if err := p.containerProvider.CopyFilesFromContainer(ctx, p.containerID, containerPath, absLocalPath); err != nil {
		return "", fmt.Errorf("failed to copy files from container: %w", err)
	}

	// List the files that were copied
	files, err := os.ReadDir(absLocalPath)
	if err != nil {
		fmt.Printf("Warning: Failed to list copied files: %v\n", err)
	} else {
		fmt.Printf("Copied %d files/directories from container\n", len(files))
		if len(files) < 10 {
			for _, file := range files {
				fmt.Printf(" - %s\n", file.Name())
			}
		}
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

// ensureContainer ensures a container is running with proper authentication
func (p *Provider) ensureContainer(ctx context.Context) error {
	if p.containerID != "" {
		// Check if container is still running
		pingCmd := []string{"echo", "ping"}
		_, err := p.containerProvider.ExecuteCommand(ctx, p.containerID, pingCmd)
		if err == nil {
			// Container is still running
			return nil
		}
		fmt.Printf("Container with ID %s is no longer available, recreating...\n", p.containerID)
		// If we get here, the container is no longer available
		p.containerID = ""
	}

	// Get container image
	image := p.config["claude_image"]
	if image == "" {
		// Check if image is set in environment
		envImage := os.Getenv("CLAUDE_CODE_IMAGE")
		if envImage != "" {
			image = envImage
			fmt.Printf("Using Claude Code image from environment: %s\n", image)
		} else {
			// Get absolute path to the Dockerfile in our project
			exePath, err := os.Executable()
			if err != nil {
				fmt.Printf("Warning: Failed to get executable path: %v\n", err)
				exePath = "."
			}

			exeDir := filepath.Dir(exePath)
			// Check for project-local image first
			localImage := "claude-code:latest"

			// Try to build the image if not present
			buildScript := filepath.Join(exeDir, "claude", "build.sh")
			if _, err := os.Stat(buildScript); err == nil {
				fmt.Println("Found build script, attempting to build Claude Code image...")
				buildCmd := exec.Command("/bin/bash", buildScript)
				buildOutput, buildErr := buildCmd.CombinedOutput()
				if buildErr != nil {
					fmt.Printf("Warning: Failed to build image: %v\n%s\n", buildErr, buildOutput)
				} else {
					fmt.Printf("Successfully built Claude Code image\n")
					image = localImage
				}
			} else {
				// No build script, try to check if our image exists
				checkImageCmd := exec.Command("docker", "images", "-q", localImage)
				output, err := checkImageCmd.Output()
				if err == nil && len(output) > 0 {
					image = localImage
					fmt.Printf("Using local Claude Code image: %s\n", localImage)
				} else {
					// Fall back to the official image
					image = "anthropic/claude-code:latest"
					fmt.Printf("Local image not found, trying to use: %s\n", image)
				}
			}
		}
	}

	// Create a persistent temporary directory for workspace
	// Use a more reliable path that works across different environments
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = "/tmp"
		fmt.Printf("Failed to get user home directory, using /tmp instead: %v\n", err)
	}

	tmpDir := filepath.Join(userHome, ".cc", "workspace", fmt.Sprintf("claude-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	fmt.Printf("Created workspace directory: %s\n", tmpDir)

	// Set up volume mounts - ensure absolute paths
	volumeMounts := make(map[string]string)
	absoluteTmpDir, err := filepath.Abs(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for workspace directory: %w", err)
	}
	volumeMounts[absoluteTmpDir] = "/workspace"

	// Set up environment variables
	env := make(map[string]string)

	// Get the API key but don't include it directly in environment variables
	// We'll use a .env file for security
	apiKey := ""
	apiKeyFound := false

	// First, check config
	if configKey, ok := p.config["claude_api_key"]; ok && configKey != "" {
		apiKey = configKey
		fmt.Println("Using Claude API key from config")
		apiKeyFound = true
	}

	// If not in config, try environment variables
	if !apiKeyFound {
		envKey := os.Getenv("CLAUDE_API_KEY")
		if envKey != "" {
			apiKey = envKey
			fmt.Println("Using Claude API key from environment variables")
			apiKeyFound = true
		}
	}

	// If no API key found by this point, return an error
	if !apiKeyFound {
		return fmt.Errorf("no Claude API key provided. Set the CLAUDE_API_KEY environment variable or add it to your configuration")
	}

	// Add any other environment variables from config
	for k, v := range p.config {
		if strings.HasPrefix(k, "env_") {
			env[strings.TrimPrefix(k, "env_")] = v
		}
	}

	// Ensure environment has minimum required variables
	env["CLAUDE_CLI_LOG_LEVEL"] = "info" // Set logging level
	env["HOME"] = "/home/node"           // Ensure HOME is set correctly for Claude CLI

	// Run container
	fmt.Printf("Starting Claude container with image: %s\n", image)
	containerID, err := p.containerProvider.RunContainer(ctx, image, volumeMounts, env)
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	p.containerID = containerID
	fmt.Printf("Claude container started with ID: %s\n", containerID)

	// Create a temporary .env file
	envFile := filepath.Join(tmpDir, ".env")
	envContent := fmt.Sprintf("CLAUDE_API_KEY=%s\n", apiKey)

	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		fmt.Printf("Warning: Failed to write .env file: %v\n", err)
	} else {
		fmt.Println("Created temporary .env file with API key")
	}

	// Copy the .env file to the container root
	// First try to use our helper script
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		helperScript := filepath.Join(exeDir, "claude", "env-handler.sh")

		if _, err := os.Stat(helperScript); err == nil {
			fmt.Println("Using env-handler.sh to copy .env to container")
			copyCmd := exec.Command("/bin/bash", helperScript, "copy-env", containerID, envFile)
			copyOutput, copyErr := copyCmd.CombinedOutput()
			if copyErr != nil {
				fmt.Printf("Warning: Failed to copy .env file using helper: %v\n%s\n", copyErr, copyOutput)
			} else {
				fmt.Printf("Successfully copied .env file to container\n")
			}
		}
	}

	// Wait a moment for container to initialize
	time.Sleep(2 * time.Second)

	// Verify that the container can load .env
	loadEnvCmd := []string{"/usr/local/bin/load-env", "&&", "env", "|", "grep", "CLAUDE"}
	loadEnvOutput, _ := p.containerProvider.ExecuteCommand(ctx, containerID, loadEnvCmd)
	if !strings.Contains(loadEnvOutput, "CLAUDE_API_KEY") {
		fmt.Printf("Warning: .env file not properly loaded in container\n")
		fmt.Println("Trying direct environment variable instead...")

		// Try setting environment directly in container
		setEnvCmd := []string{"bash", "-c", fmt.Sprintf("echo 'export CLAUDE_API_KEY=%s' >> ~/.bashrc && echo 'export CLAUDE_API_KEY=%s' >> ~/.profile", apiKey, apiKey)}
		_, setEnvErr := p.containerProvider.ExecuteCommand(ctx, containerID, setEnvCmd)
		if setEnvErr != nil {
			fmt.Printf("Warning: Failed to set environment variables in container: %v\n", setEnvErr)
		}
	}

	// Verify that Claude CLI is working with authentication
	testCmd := []string{"claude", "code", "--version"}
	testOutput, testErr := p.containerProvider.ExecuteCommand(ctx, containerID, testCmd)
	if testErr != nil {
		fmt.Printf("Warning: Claude CLI test command failed: %v\n", testErr)
		fmt.Printf("Output: %s\n", testOutput)

		// Try to get more detailed error information
		errorCmd := []string{"ls", "-la", "/usr/local/bin/claude"}
		errorOutput, _ := p.containerProvider.ExecuteCommand(ctx, containerID, errorCmd)
		fmt.Printf("claude executable details: %s\n", errorOutput)

		// Check if the error is due to API key
		if strings.Contains(testOutput, "API") && strings.Contains(testOutput, "key") {
			// Try one more approach - copy API key directly into .claude directory
			setupCmd := []string{"mkdir", "-p", "/home/node/.claude", "&&",
				"echo", fmt.Sprintf("'{\"api_key\":\"%s\"}'", apiKey), ">", "/home/node/.claude/config.json"}
			_, setupErr := p.containerProvider.ExecuteCommand(ctx, containerID, setupCmd)
			if setupErr != nil {
				fmt.Printf("Warning: Failed to set up config.json: %v\n", setupErr)
				return fmt.Errorf("authentication failed: invalid or missing API key. Please check your CLAUDE_API_KEY")
			}

			// Try testing again
			testOutput, testErr = p.containerProvider.ExecuteCommand(ctx, containerID, testCmd)
			if testErr != nil {
				return fmt.Errorf("authentication failed: invalid or missing API key. Please check your CLAUDE_API_KEY")
			}
		} else {
			fmt.Println("Error doesn't appear to be authentication related, continuing anyway, but issues may occur.")
		}
	}

	fmt.Printf("Claude CLI is operational: %s\n", strings.TrimSpace(testOutput))
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

// SetContainerProviderForTest allows setting a mock container provider for testing
// This method is only exported in test builds
func (p *Provider) SetContainerProviderForTest(containerProvider container.Provider) error {
	if containerProvider == nil {
		return fmt.Errorf("container provider cannot be nil")
	}
	p.containerProvider = containerProvider
	return nil
}
