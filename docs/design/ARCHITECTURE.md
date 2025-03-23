# Code Controller (CC) Architecture

## Overview
CC is a CLI tool written in Go that integrates with container systems to manage AI-generated code projects using various AI providers like Claude Code.

## Core Components

### CLI Tool (Go)
- **cmd/cc**: Main CLI entry point and commands
- **internal/container**: Container management abstraction layer
  - Docker implementation
  - Kubernetes implementation (future)
- **internal/vcs**: Version control system abstraction
  - Git implementation
  - Other VCS support (future)
- **internal/ai**: AI provider abstraction
  - Claude implementation
  - Other AI providers (future)
- **pkg/models**: Data models and interfaces
- **pkg/config**: Configuration handling
- **pkg/plugin**: Plugin system for extensibility

### Container Orchestration
- Abstracted container management
- Support for Docker and Kubernetes
- Local and remote execution environments

### AI Provider System
- Pluggable AI providers
- Claude Code as primary implementation
- Interface to support other AI code generators

## Workflow
1. User initiates a project with `cc init [project]`
2. CC prepares container environment with AI provider
3. CC manages version control for different implementation versions
4. User selects preferred implementation with `cc select [branch]`
5. User requests features with `cc feature "description"`
6. CC manages feature branches and updates

## Technical Decisions

### Language: Go
- Single binary distribution
- Excellent container integration
- High performance
- Cross-platform capability
- Strong concurrency support

### Abstraction Layers
- Provider abstraction for multiple AI systems
- Container abstraction for Docker/Kubernetes
- VCS abstraction for Git and alternatives
- Plugin system for extensibility

### Distributed Execution
- Support for parallel implementation generation
- Asynchronous job processing
- Remote execution capabilities

### State Management
- Decoupled from specific VCS
- Metadata storage for version comparison
- Project state tracking

### Extensibility
- Plugin system for custom templates
- Framework-specific generators
- Custom workflows