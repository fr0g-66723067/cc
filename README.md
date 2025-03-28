# Code Controller (CC)

## Overview
Code Controller (CC) is an advanced interface built on top of Claude Code CLI that automates and streamlines AI-powered code generation. It enables rapid prototyping and exploration of different implementation approaches through automated git branching.

**Current Status**: Fully functional CLI with actual Git operations, Docker integration, and Claude Code integration. The core functionality is implemented with real Git branch management, Docker container support (with sudo support), and Claude Code AI for code generation. The Claude Code integration now uses the real Claude CLI in a Docker container with proper environment and API key management. Complete end-to-end workflow is supported from project generation to feature implementation.

## Key Features

### Multi-version Prototyping
- Generate multiple versions of the same application using different frameworks and approaches
- Each version is maintained in its own git branch for easy comparison and switching
- Example: "Create a web app to track daily tasks" will generate versions using React, Vue, Svelte, etc.

### Project Evolution
- Select your preferred version as the base implementation
- Incrementally enhance your chosen implementation by requesting new features
- Each feature request creates a new branch, allowing for easy feature comparison and rollback

### Workflow
1. Define your project requirements
2. CC generates multiple implementation versions using Claude
3. Review and select your preferred implementation approach
4. Incrementally add features through natural language requests
5. Navigate between different versions and feature branches easily

### Benefits
- Dramatically accelerates initial prototyping and exploration
- Enables testing different frameworks and architectural approaches without manual setup
- Provides a structured way to compare different implementation strategies
- Maintains a clean version history through automated git branching

## Installation

### Prerequisites
- Go 1.21+
- Docker (for running Claude Code)
- Git
- Claude API key (get one from [Anthropic](https://console.anthropic.com/))

### Installation from Source
```bash
git clone https://github.com/fr0g-66723067/cc.git
cd cc
# Run the setup script to build the Docker image and CC tool
./claude/setup.sh
```

This setup script will:
1. Check for your Claude API key in the environment or .env file
2. Build the Claude Docker image
3. Build the CC tool

If you prefer to do this manually:
```bash
# Build the Claude Docker image
cd claude
./build.sh

# Build the CC tool
cd ..
go build -o cc ./cmd/cc
```

Move the binary to a location in your PATH to use it from anywhere:
```bash
sudo mv cc /usr/local/bin/
```

### Setting up your API Key

You can provide your Claude API key in one of the following ways:
1. Environment variable: `export CLAUDE_API_KEY=your-api-key`
2. .env file in the project directory or your home directory: `CLAUDE_API_KEY=your-api-key`
3. Using the `--api-key` flag when running a command: `cc --api-key=your-api-key init my-project`

## Usage

Code Controller offers three ways to interact with your projects:

1. Traditional CLI commands (backward compatible)
2. Hierarchical command structure
3. Interactive shell mode

### Traditional CLI Commands (Legacy)

```bash
cc init my-project
cc generate "Create a web app to track daily tasks"
cc list implementations
cc select implementation-branch-name
cc feature "Add a dark mode toggle"
cc compare branch1 branch2
cc status
cc remove project project-name
cc rename implementation old-branch-name new-branch-name
```

### Hierarchical Command Structure

The hierarchical command structure organizes commands by resource type:

```bash
# Projects Commands
cc projects list
cc projects create my-project "Project description"
cc projects remove my-project
cc projects rename old-name new-name

# Implementations Commands
cc implementations list
cc implementations generate "Create a web app" --frameworks react,vue
cc implementations select branch-name
cc implementations remove branch-name
cc implementations rename old-branch new-branch
cc implementations compare branch1 branch2

# Features Commands
cc features list
cc features add "Add dark mode toggle"
cc features remove feature-branch
cc features rename old-feature-branch new-feature-branch
```

### Context Navigation

CC maintains a navigation context that remembers which project, implementation, and feature you're working with:

```bash
# Switch between projects, implementations, and features
cc use project my-project
cc use implementation impl-react-123456
cc use feature feat-darkmode-123456

# Show current context and status
cc status
```

### Interactive Shell Mode

For a more interactive experience, use the shell mode:

```bash
cc shell
```

The shell provides:
- Tab completion for commands and resources
- Context-aware prompts showing your location
- Command history
- Path-based navigation similar to a filesystem

Shell commands include:

```
# Context navigation
pwd                                   # Show current context path
cd /projects/my-project               # Navigate to a project
cd implementations/impl-react-123456  # Navigate to an implementation
cd /                                  # Return to root level

# Resource management
projects list
projects create my-new-project "Description"
implementations generate "Create a login page" --frameworks=react,vue
features add "Add password reset"

# Using resources
use project my-project
use implementation impl-react-123456

# Show status
status
```

## Project Structure
- `cmd/cc`: CLI implementation
- `internal/ai`: AI provider interfaces and implementations
- `internal/container`: Container provider interfaces and implementations
- `internal/vcs`: Version control interfaces and implementations
- `pkg/config`: Configuration management
- `pkg/models`: Data models
- `pkg/plugin`: Plugin system
- `docs`: Project documentation
- `scripts`: Utility scripts, including benchmarks
- `test`: Test files and utilities

## Configuration
CC uses a configuration file stored at `~/.cc/config.json`. You can specify a different configuration file with the `--config` flag.

## Documentation
Comprehensive documentation is available in the `docs` directory:

- [Usage Guide](docs/USAGE.md) - Detailed usage instructions
- [Project Summary](docs/SUMMARY.md) - Current project status
- [Architecture](docs/design/ARCHITECTURE.md) - System architecture
- [Project Structure](docs/design/PROJECT_STRUCTURE.md) - Detailed code organization
- [Development Roadmap](docs/design/DETAILED_ROADMAP.md) - Future development plans

## Development

### Building
```bash
go build -o cc ./cmd/cc
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./cmd/cc -v
go test ./internal/container -v
```

### Implementation Details
- The CLI is implemented using Cobra
- Providers use an interface-based design for easy mocking and testing
- Configuration is stored in JSON format
- Dependency injection is used for providers
- TDD approach is used for all features

### Implemented Features

- Real Claude Code CLI integration with Docker container
- Secure API key management with .env file support
- Multiple authentication methods (env vars, config files, CLI flags)
- Git branch management for implementations and features
- Docker container lifecycle management
- Testing infrastructure with mock providers
- Project, implementation, and feature management (creation, listing, removal, renaming)
- Hierarchical command structure with resource-specific subcommands
- Interactive shell mode with tab completion and context navigation
- Context-aware filesystem-like navigation for resources
- Comprehensive test suite for shell and navigation functionality

## Remaining Limitations
- Advanced Git operations (merging, rebasing) not yet implemented
- Comprehensive logging improvements needed
- No authentication support for private Docker registries
- Plugin system functionality needs expansion
- No support for custom project templates

### Next Steps
- Implement more robust error recovery strategies
- Add request retries for network or API failures
- Add progress reporting and terminal UI enhancements
- Implement container pooling to reduce startup times
- Add caching for repeated operations
- Optimize parallel implementation generation
- Expand the plugin system for custom frameworks
- Implement advanced Git operations like merging branches
- Add support for private Docker registries
- Implement project templates and framework-specific configurations
- Add shell history persistence across sessions
- Improve shell tab completion with more intelligent suggestions
- Add alias commands for common operations
- Enhance error messages with context-specific suggestions

### Contributing
Contributions are welcome! Please feel free to submit a Pull Request.