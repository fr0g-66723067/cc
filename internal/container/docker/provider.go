package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

// Provider implements the container provider interface for Docker
type Provider struct {
	client *client.Client
	config map[string]string
}

// NewProvider creates a new Docker provider
func NewProvider(config map[string]string) (*Provider, error) {
	return &Provider{
		config: config,
	}, nil
}

// Initialize sets up the Docker provider
func (p *Provider) Initialize(ctx context.Context, config map[string]string) error {
	// Merge configs
	for k, v := range config {
		p.config[k] = v
	}

	// Create Docker client
	var err error
	p.client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Verify Docker connection
	_, err = p.client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	return nil
}

// RunContainer starts a container with the given image and returns its ID
func (p *Provider) RunContainer(ctx context.Context, image string, volumeMounts map[string]string, env map[string]string) (string, error) {
	// Pull image if needed
	_, _, err := p.client.ImageInspectWithRaw(ctx, image)
	if client.IsErrNotFound(err) {
		// Image doesn't exist, pull it
		reader, err := p.client.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to pull image: %w", err)
		}
		defer reader.Close()

		// Wait for the pull to complete
		io.Copy(io.Discard, reader)
	} else if err != nil {
		return "", fmt.Errorf("failed to inspect image: %w", err)
	}

	// Prepare volume mounts
	var mounts []mount.Mount
	for hostPath, containerPath := range volumeMounts {
		// Create host directory if it doesn't exist
		if err := os.MkdirAll(hostPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create host directory: %w", err)
		}

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostPath,
			Target: containerPath,
		})
	}

	// Prepare environment variables
	var envVars []string
	for k, v := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	// Create container
	resp, err := p.client.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"tail", "-f", "/dev/null"}, // Keep container running
		Env:   envVars,
		Tty:   true,
	}, &container.HostConfig{
		Mounts: mounts,
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := p.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// ExecuteCommand executes a command in the container
func (p *Provider) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	// Create exec configuration
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	}

	// Create exec instance
	exec, err := p.client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	// Start exec instance
	response, err := p.client.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("failed to start exec: %w", err)
	}
	defer response.Close()

	// Read output
	var stdout, stderr strings.Builder
	if _, err := stdcopy.StdCopy(&stdout, &stderr, response.Reader); err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	// Check if command failed
	inspect, err := p.client.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect exec: %w", err)
	}

	if inspect.ExitCode != 0 {
		return stdout.String(), fmt.Errorf("command failed with exit code %d: %s", inspect.ExitCode, stderr.String())
	}

	return stdout.String(), nil
}

// CopyFilesToContainer copies files from local to container
func (p *Provider) CopyFilesToContainer(ctx context.Context, containerID string, localPath string, containerPath string) error {
	// Create tar archive of source directory
	srcInfo, err := archive.CopyInfoSourcePath(localPath, false)
	if err != nil {
		return fmt.Errorf("failed to get source info: %w", err)
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return fmt.Errorf("failed to create tar archive: %w", err)
	}
	defer srcArchive.Close()

	// Copy tar archive to container
	if err := p.client.CopyToContainer(ctx, containerID, containerPath, srcArchive, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                false,
	}); err != nil {
		return fmt.Errorf("failed to copy to container: %w", err)
	}

	return nil
}

// CopyFilesFromContainer copies files from container to local
func (p *Provider) CopyFilesFromContainer(ctx context.Context, containerID string, containerPath string, localPath string) error {
	// Get tar archive from container
	reader, _, err := p.client.CopyFromContainer(ctx, containerID, containerPath)
	if err != nil {
		return fmt.Errorf("failed to copy from container: %w", err)
	}
	defer reader.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract tar archive to destination directory
	if err := archive.Untar(reader, localPath, &archive.TarOptions{
		NoLchown: true,
	}); err != nil {
		return fmt.Errorf("failed to extract tar archive: %w", err)
	}

	return nil
}

// StopContainer stops a running container
func (p *Provider) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10 // 10 seconds
	if err := p.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// RemoveContainer removes a container
func (p *Provider) RemoveContainer(ctx context.Context, containerID string) error {
	if err := p.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
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
	container.Register("docker", func(config map[string]string) (container.Provider, error) {
		return NewProvider(config)
	})
}