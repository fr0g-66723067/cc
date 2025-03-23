# Code Controller (CC)

## Overview
Code Controller (CC) is an advanced interface built on top of Claude Code CLI that automates and streamlines AI-powered code generation. It enables rapid prototyping and exploration of different implementation approaches through automated git branching.

**Current Status**: Fully functional CLI with actual Git operations, Docker integration, and Claude Code integration. The core functionality is implemented with real Git branch management, Docker container support (with sudo support), and Claude Code AI for code generation. Successfully tested end-to-end with a mock Claude Code implementation.

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

### Installation from Source
```bash
git clone https://github.com/fr0g-66723067/cc.git
cd cc
go build -o cc ./cmd/cc
```

Move the binary to a location in your PATH to use it from anywhere:
```bash
sudo mv cc /usr/local/bin/
```

## Usage

### Initialize a New Project
```bash
cc init my-project
```

### Generate Implementations
```bash
cc generate "Create a web app to track daily tasks"
```

With specific frameworks:
```bash
cc generate "Create a web app to track daily tasks" --frameworks react,vue,svelte
```

### List Generated Implementations
```bash
cc list implementations
```

### Select an Implementation
```bash
cc select implementation-branch-name
```

### Add a Feature
```bash
cc feature "Add a dark mode toggle"
```

### Compare Implementations
```bash
cc compare branch1 branch2
```

### Show Project Status
```bash
cc status
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

### Current Limitations
- Advanced Git operations (merging, rebasing) not yet implemented
- Comprehensive logging not yet implemented
- No authentication support for private Docker registries
- Limited plugin system functionality
- No support for custom project templates

### Next Steps
- Add comprehensive error handling and logging
- Implement retries and recovery strategies
- Add progress reporting and better user feedback
- Expand the plugin system for custom frameworks
- Implement advanced Git operations like merging branches
- Add improved metrics and diagnostics
- Add support for private Docker registries
- Implement project templates and framework-specific configurations
- Create a configuration UI for easier setup

### Contributing
Contributions are welcome! Please feel free to submit a Pull Request.