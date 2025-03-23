// Docker provider implementation
package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	cccontainer "github.com/fr0g-66723067/cc/internal/container"
)

// Provider implements the container provider interface for Docker
type Provider struct {
	config map[string]string
}

// NewProvider creates a new Docker provider
func NewProvider(config map[string]string) (*Provider, error) {
	if config == nil {
		config = make(map[string]string)
	}
	
	return &Provider{
		config: config,
	}, nil
}

// Initialize sets up the Docker provider
func (p *Provider) Initialize(ctx context.Context, config map[string]string) error {
	// Merge configs
	if config != nil {
		for k, v := range config {
			p.config[k] = v
		}
	}

	// Check if sudo should be used
	useSudo := false
	if val, ok := p.config["use_sudo"]; ok && val == "true" {
		useSudo = true
	}
	
	// Set config value for future use
	p.config["use_sudo"] = fmt.Sprintf("%t", useSudo)

	// Check if Docker is installed and running
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.CommandContext(ctx, "sudo", "docker", "info")
	} else {
		cmd = exec.CommandContext(ctx, "docker", "info")
	}
	
	if err := cmd.Run(); err != nil {
		// If failed without sudo, try with sudo
		if !useSudo {
			sudoCmd := exec.CommandContext(ctx, "sudo", "docker", "info")
			if sudoErr := sudoCmd.Run(); sudoErr == nil {
				// Docker works with sudo, update config
				p.config["use_sudo"] = "true"
				fmt.Println("Docker requires sudo privileges, enabling sudo for Docker commands")
				return nil
			}
		}
		return fmt.Errorf("Docker is not running or not installed: %w", err)
	}

	return nil
}

// RunContainer starts a container with the given image and returns its ID
func (p *Provider) RunContainer(ctx context.Context, image string, volumeMounts map[string]string, env map[string]string) (string, error) {
	// Build docker run command
	args := []string{"run", "-d"}
	
	// Add environment variables
	for k, v := range env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	
	// Add volume mounts
	for host, container := range volumeMounts {
		args = append(args, "-v", fmt.Sprintf("%s:%s", host, container))
	}
	
	// Add the image and command to keep container running
	args = append(args, image, "tail", "-f", "/dev/null")
	
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Run the container
	var cmd *exec.Cmd
	if useSudo {
		sudoArgs := append([]string{"docker"}, args...)
		cmd = exec.CommandContext(ctx, "sudo", sudoArgs...)
	} else {
		cmd = exec.CommandContext(ctx, "docker", args...)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the image doesn't exist, try to pull it
		if strings.Contains(string(output), "No such image") || strings.Contains(string(output), "not found") {
			fmt.Printf("Image %s not found, attempting to pull...\n", image)
			var pullCmd *exec.Cmd
			if useSudo {
				pullCmd = exec.CommandContext(ctx, "sudo", "docker", "pull", image)
			} else {
				pullCmd = exec.CommandContext(ctx, "docker", "pull", image)
			}
			
			pullOutput, pullErr := pullCmd.CombinedOutput()
			if pullErr != nil {
				return "", fmt.Errorf("failed to pull image %s: %w\n%s", image, pullErr, pullOutput)
			}
			
			// Try to run the container again
			if useSudo {
				sudoArgs := append([]string{"docker"}, args...)
				cmd = exec.CommandContext(ctx, "sudo", sudoArgs...)
			} else {
				cmd = exec.CommandContext(ctx, "docker", args...)
			}
			
			output, err = cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to run container after pulling image: %w\n%s", err, output)
			}
		} else {
			return "", fmt.Errorf("failed to run container: %w\n%s", err, output)
		}
	}
	
	// Get container ID from output
	containerID := strings.TrimSpace(string(output))
	fmt.Printf("Container started with ID: %s\n", containerID)
	
	return containerID, nil
}

// ExecuteCommand executes a command in the container
func (p *Provider) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	// Build docker exec command
	args := []string{"exec", containerID}
	args = append(args, command...)
	
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Execute the command
	var cmd *exec.Cmd
	if useSudo {
		sudoArgs := append([]string{"docker"}, args...)
		cmd = exec.CommandContext(ctx, "sudo", sudoArgs...)
	} else {
		cmd = exec.CommandContext(ctx, "docker", args...)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command in container: %w\n%s", err, output)
	}
	
	return string(output), nil
}

// CopyFilesToContainer copies files from local to container
func (p *Provider) CopyFilesToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error {
	// Verify local path exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("local path does not exist: %s", localPath)
	}
	
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Use docker cp to copy files
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.CommandContext(ctx, "sudo", "docker", "cp", localPath, fmt.Sprintf("%s:%s", containerID, containerPath))
	} else {
		cmd = exec.CommandContext(ctx, "docker", "cp", localPath, fmt.Sprintf("%s:%s", containerID, containerPath))
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy files to container: %w\n%s", err, output)
	}
	
	return nil
}

// CopyFilesFromContainer copies files from container to local
func (p *Provider) CopyFilesFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Use docker cp to copy files
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.CommandContext(ctx, "sudo", "docker", "cp", fmt.Sprintf("%s:%s", containerID, containerPath), localPath)
	} else {
		cmd = exec.CommandContext(ctx, "docker", "cp", fmt.Sprintf("%s:%s", containerID, containerPath), localPath)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy files from container: %w\n%s", err, output)
	}
	
	// Fix permissions if using sudo
	if useSudo {
		userInfo, err := user.Current()
		if err == nil {
			// Change ownership of the copied files to current user
			chownCmd := exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s:%s", userInfo.Uid, userInfo.Gid), localPath)
			if chownErr := chownCmd.Run(); chownErr != nil {
				fmt.Printf("Warning: Failed to fix permissions on copied files: %v\n", chownErr)
			}
		}
	}
	
	return nil
}

// StopContainer stops a running container
func (p *Provider) StopContainer(ctx context.Context, containerID string) error {
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Stop the container
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.CommandContext(ctx, "sudo", "docker", "stop", containerID)
	} else {
		cmd = exec.CommandContext(ctx, "docker", "stop", containerID)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %w\n%s", err, output)
	}
	
	return nil
}

// RemoveContainer removes a container
func (p *Provider) RemoveContainer(ctx context.Context, containerID string) error {
	// Check if sudo should be used
	useSudo := p.config["use_sudo"] == "true"
	
	// Remove the container with force
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.CommandContext(ctx, "sudo", "docker", "rm", "-f", containerID)
	} else {
		cmd = exec.CommandContext(ctx, "docker", "rm", "-f", containerID)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove container: %w\n%s", err, output)
	}
	
	return nil
}

// Name returns the provider's name
func (p *Provider) Name() string {
	return "docker"
}

// IsRemote returns whether the provider is running containers remotely
func (p *Provider) IsRemote() bool {
	return false
}

// Register registers this provider factory
func init() {
	cccontainer.Register("docker", func(config map[string]string) (cccontainer.Provider, error) {
		return NewProvider(config)
	})
}