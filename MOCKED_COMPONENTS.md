# Mocked Components in Code Controller (CC)

This document lists all the mocked components in the Code Controller (CC) codebase that need to be replaced with real implementations. These were created for demo/testing purposes only.

## 1. Claude Code CLI Mock

**Status**: Currently mocked with a shell script
**Files**:
- `/home/user/src/cc/claude/Dockerfile` - Contains mock Claude CLI script
- `/home/user/src/cc/claude/simple_mock_claude.sh` - Simplified mock script

**What it does**:
- Creates a Docker container with a fake Claude CLI script
- Checks for an API key for authentication
- Has hard-coded handlers for code generation, modification, and analysis
- Generates basic Python Flask and JavaScript templates regardless of prompt

**What needs to be implemented**:
- Real integration with Claude AI API
- Proper authentication with Claude API keys
- Actual code generation from AI-based prompts
- Proper error handling and response parsing
- Support for various project types and frameworks

## 2. Docker Container Provider

**Status**: Partially implemented, with some functions mocked
**Files**:
- `/home/user/src/cc/internal/container/docker/provider.go`

**What it does**:
- Uses os/exec to run Docker commands
- Handles basic container operations (run, exec, copy, stop, remove)
- Contains sudo support functionality

**What needs to be implemented**:
- Better error handling for Docker operations
- Support for more complex Docker configurations
- Improved file handling between host and container
- Volume management and workspace handling
- Container resource management

## 3. Claude AI Provider

**Status**: Partially implemented with heavy mocking
**Files**:
- `/home/user/src/cc/internal/ai/claude/provider.go`

**What it does**:
- Manages interaction with Claude container
- Contains handlers for project generation, implementation, feature addition
- Has mock code for API key authentication
- Creates workspace directories and manages containers

**What needs to be implemented**:
- Real API integration with Claude
- Proper prompt formatting for code generation
- Handling different project types and frameworks
- Proper error handling for API responses
- Context management for multi-step operations

## 4. Test Scripts

**Status**: Mocked for basic functionality testing
**Files**:
- `/home/user/src/cc/test/docker_auth_test.sh`
- `/home/user/src/cc/test/auth_test.go`

**What they do**:
- Test Docker container creation and authentication
- Check Claude API key loading from environment
- Verify basic code generation functionality

**What needs to be implemented**:
- Comprehensive test suite for all functionality
- Integration tests with real API
- Unit tests for all provider functions
- Error condition testing

## 5. Command Implementations

**Status**: Basic workflow is implemented but relies on mocks
**Files**:
- `/home/user/src/cc/cmd/cc/commands.go`

**What it does**:
- Implements command handlers for CLI operations
- Creates projects, implementations, features
- Manages Git branches and operations
- Uses mocked AI and container providers

**What needs to be implemented**:
- Proper error handling
- Validation of input and output
- Progress reporting for long-running operations
- Support for real AI-generated code

## 6. Git VCS Provider

**Status**: Partially implemented, integration with mocks
**Files**:
- `/home/user/src/cc/internal/vcs/git/provider.go`

**What it does**:
- Manages Git operations using go-git library
- Creates repositories, branches, commits
- Handles file changes and commits

**What needs to be implemented**:
- Better error handling
- More robust branch management
- Conflict resolution
- Remote repository support
- Better file tracking

## Summary of Critical Issues

1. **Missing Real Claude API Integration**: The current implementation has no actual integration with Claude's API. A real implementation should use the official Claude APIs with proper authentication.

2. **Missing Real Code Generation**: All code generation is currently simulated with hard-coded templates. A real implementation should use Claude's AI capabilities to generate code from prompts.

3. **Project Management Issues**: Projects are created in config but not properly on the filesystem.

4. **File Permission Problems**: Issues with file permissions between host and container need proper fixing.

5. **Authentication Flow**: The API key authentication flow is mocked and needs real implementation.

6. **Workspace Management**: Proper workspace handling between host and container is not fully implemented.

7. **Testing**: Proper unit and integration tests with real API calls are needed.

## Next Steps

1. Implement real Claude API integration
2. Set up proper authentication flow
3. Implement real code generation and feature addition
4. Fix file permission and workspace management issues
5. Add comprehensive tests
6. Implement proper error handling and recovery