# Code Controller Project Structure

## File Structure

```
cc/
├── cmd/
│   └── cc/
│       └── main.go           # CLI entry point
├── internal/
│   ├── ai/
│   │   ├── provider.go       # AI provider interface
│   │   └── claude/
│   │       └── provider.go   # Claude implementation
│   ├── container/
│   │   ├── provider.go       # Container abstraction
│   │   ├── docker/
│   │   │   └── provider.go   # Docker implementation
│   │   └── kubernetes/
│   │       └── provider.go   # Kubernetes implementation
│   ├── job/
│   │   └── queue.go          # Job queue for async operations
│   └── vcs/
│       ├── provider.go       # VCS abstraction
│       └── git/
│           └── provider.go   # Git implementation
├── pkg/
│   ├── config/
│   │   └── config.go         # Configuration handling
│   ├── models/
│   │   ├── project.go        # Project data model
│   │   └── implementation.go # Implementation data model
│   └── plugin/
│       └── manager.go        # Plugin system
├── plugins/
│   ├── frameworks/           # Framework-specific plugins
│   └── templates/            # Project templates
├── ARCHITECTURE.md           # Architecture documentation
├── CLAUDE.md                 # Guidelines for Claude
├── go.mod                    # Go module definition
├── PROJECT_STRUCTURE.md      # This file
└── README.md                 # Project overview
```

## Component Responsibilities

### cmd/cc
Main entry point for the CLI application. Defines commands and their handlers using Cobra.

### internal/ai
Abstraction layer for AI providers with Claude implementation.

### internal/container
Abstraction for container management, supporting Docker and Kubernetes.

### internal/job
Handles asynchronous job processing for long-running tasks.

### internal/vcs
Abstracts version control operations, with Git as primary implementation.

### pkg/config
Manages application configuration, including loading and saving settings.

### pkg/models
Defines data models used throughout the application.

### pkg/plugin
Pluggable extension system for custom functionality.

### plugins
Custom plugins for specific frameworks and project templates.

## Design Patterns

1. **Provider Pattern** - Abstractions with multiple implementations
2. **Strategy Pattern** - Interchangeable algorithms for different providers
3. **Factory Pattern** - Dynamic creation of appropriate provider implementations
4. **Observer Pattern** - Event-based notification for async operations
5. **Command Pattern** - Encapsulated operations for job queue

## Workflow Implementation

1. `cc init` - Creates a new project directory, initializes VCS
2. `cc generate` - Uses AI provider to generate multiple implementation versions
3. `cc select` - Switches to selected implementation branch
4. `cc feature` - Creates a new feature branch and uses AI to implement

## Development Plan

1. Implement core abstractions and interfaces
2. Implement Docker container integration
3. Implement Git VCS management
4. Implement Claude AI provider
5. Implement job queue system
6. Add plugin system
7. Add Kubernetes support
8. Extend with additional AI providers