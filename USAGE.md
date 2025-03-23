# Code Controller (CC) Usage Guide

This document provides detailed instructions on how to use the Code Controller (CC) tool with Claude Code.

## Prerequisites

Before using CC, ensure you have:
- Go 1.21+ installed
- Docker installed and running
- Git installed
- A Claude API key from Anthropic

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/fr0g-66723067/cc.git
   cd cc
   ```

2. Run the setup script:
   ```bash
   ./claude/setup.sh
   ```
   
   This script will:
   - Check for and prompt for your Claude API key if needed
   - Build the Claude Docker image
   - Build the CC tool

3. Alternatively, you can set up manually:
   ```bash
   # Set your API key
   export CLAUDE_API_KEY=your-api-key
   
   # Build the Docker image
   cd claude
   ./build.sh
   
   # Build the CC tool
   cd ..
   go build -o cc ./cmd/cc
   ```

4. (Optional) Move the binary to your PATH:
   ```bash
   sudo mv cc /usr/local/bin/
   ```

## Basic Usage

### Initialize a New Project

Create a new project with CC:

```bash
cc init my-project
```

This creates a new directory structure and initializes a Git repository.

You can add a description:

```bash
cc init my-project --description "A web application for task management"
```

### Generate Implementations

Once you've initialized a project, you can generate different implementations:

```bash
cc generate "Create a web app for tracking daily tasks"
```

By default, this creates implementations for the top 3 supported frameworks. You can specify frameworks:

```bash
cc generate "Create a web app for tracking daily tasks" --frameworks react,vue,svelte
```

Or change the number of implementations:

```bash
cc generate "Create a web app for tracking daily tasks" --count 5
```

### List Generated Implementations

To see what implementations have been generated:

```bash
cc list implementations
```

### Select an Implementation

After reviewing the generated implementations, select one as your base:

```bash
cc select branch-name
```

Replace `branch-name` with the implementation branch name from the list.

### Add a Feature

Once you've selected an implementation, you can add features:

```bash
cc feature "Add a dark mode toggle"
```

This creates a new feature branch based on your selected implementation.

### Compare Implementations

You can compare different implementations or features:

```bash
cc compare branch1 branch2
```

### Show Project Status

To see the current status of your project:

```bash
cc status
```

## Working with the Claude Docker Container

CC manages the Claude Docker container for you, but if you need to interact with it directly:

### Building the Container Manually

```bash
cd claude
./build.sh
```

### Checking Container Status

```bash
docker ps | grep claude-code
```

### Executing Commands in the Container

```bash
docker exec -it [container-id] bash
```

Once inside the container, you can run Claude commands:

```bash
claude code --version
```

### Environment Variables

CC automatically handles passing your Claude API key to the container. If you need to add other environment variables, you can:

1. Add them to your `.env` file:
   ```
   CLAUDE_API_KEY=your-api-key
   OTHER_VAR=other-value
   ```

2. Or set them in your configuration file (`~/.cc/config.json`):
   ```json
   {
     "ai": {
       "config": {
         "env_OTHER_VAR": "other-value"
       }
     }
   }
   ```

## Advanced Configuration

CC stores its configuration in `~/.cc/config.json`. You can manually edit this file to change default behavior.

### Common Configuration Options

- Change default Docker image:
  ```json
  {
    "container": {
      "claudeImage": "your-custom-image:tag"
    }
  }
  ```

- Change projects directory:
  ```json
  {
    "projectsDir": "/path/to/projects"
  }
  ```

- Set default frameworks:
  ```json
  {
    "ai": {
      "config": {
        "frameworks": "react,vue,angular"
      }
    }
  }
  ```

## Troubleshooting

### API Key Issues

If you encounter API key authentication errors:

1. Check that your API key is correctly set in the environment:
   ```bash
   echo $CLAUDE_API_KEY
   ```

2. Try adding it directly to the command:
   ```bash
   cc --api-key=your-api-key status
   ```

3. Verify the key is being passed to the container:
   ```bash
   # Get container ID
   docker ps
   
   # Check environment in container
   docker exec [container-id] /usr/local/bin/load-env
   ```

### Docker Issues

If Docker fails to start or run:

1. Check Docker status:
   ```bash
   docker info
   ```

2. If permissions are an issue, add your user to the docker group:
   ```bash
   sudo usermod -aG docker $USER
   ```
   Then log out and back in.

3. If sudo is required, you can tell CC to use sudo:
   ```json
   {
     "container": {
       "config": {
         "use_sudo": "true"
       }
     }
   }
   ```

### Container Issues

If you encounter issues with the Claude container:

1. Rebuild the image:
   ```bash
   cd claude
   ./build.sh --no-cache
   ```

2. Check container logs:
   ```bash
   docker logs [container-id]
   ```

## Example Workflow

Here's a complete workflow example:

```bash
# Set your API key
export CLAUDE_API_KEY=your-api-key

# Initialize a project
cc init task-app --description "A web application for task management"

# Generate different implementations
cc generate "Create a web app that tracks daily tasks with categories, due dates, and priority levels"

# List implementations
cc list implementations

# Select preferred implementation
cc select impl-react-1234567890

# Add features
cc feature "Add a dark mode toggle"
cc feature "Add task filtering by category"
cc feature "Add a statistics dashboard"

# Check status
cc status
```

This will create a complete project with multiple implementations and features, all managed through git branches for easy comparison and selection.