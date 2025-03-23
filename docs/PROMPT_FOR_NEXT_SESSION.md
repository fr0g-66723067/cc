# Code Controller (CC) - Detailed Prompt for Next Session

## Project Context
I'm developing "Code Controller" (CC), a CLI tool built in Go that serves as an interface to Claude Code. CC enables rapid prototyping by generating multiple implementations of the same project using different frameworks and managing them through Git branching.

## Current State
- **Architecture**: Interface-based providers for AI (Claude), Container (Docker), and VCS (Git)
- **CLI Commands**: Implemented using Cobra with all core commands (init, generate, select, feature, list, compare, status)
- **Status**: The CLI is fully functional with actual Git operations, but lacks actual Docker integration and real Claude Code generation
- **Git Integration**: Git operations now fully implemented and working as expected
- **Testing**: Unit tests with mock providers are implemented, integration tests are in progress

## Focus for Next Session
I'd like to focus on implementing the actual Docker container provider to run Claude Code CLI within Docker. Currently, we're using a mock Docker provider, but we need to integrate with the real Docker SDK.

## Files to Work On
1. `/internal/container/docker/provider.go` - Contains the Docker provider implementation
2. `/internal/container/provider.go` - Contains the container provider interface
3. `/internal/ai/claude/provider.go` - Contains the Claude AI provider that uses the container

## Key Tasks

### 1. Complete Docker Provider Implementation
- Fix issues with Docker SDK imports and types
- Implement container creation, starting, stopping
- Add file copying between host and container
- Implement proper error handling for Docker operations
- Add configuration options for Docker settings

### 2. Integrate Docker Provider with AI Provider
- Update Claude AI provider to use the Docker container
- Initialize the container in the AI provider
- Pass prompts and receive responses through the container
- Implement proper cleanup when done

### 3. Update Tests
- Update the Docker provider tests to test actual Docker operations
- Create test fixtures that use real Docker containers
- Add integration tests for the full workflow

## Technical Requirements
- Use the Docker SDK for Go
- Follow TDD principles - update tests first, then implementation
- Maintain the current error handling patterns
- Ensure all tests pass after changes
- Document any changes to behavior

## Design Principles to Follow
1. **Separation of Concerns**: Keep Docker operations in the Docker provider
2. **Error Handling**: Provide meaningful error messages that help users understand issues
3. **Testability**: Ensure all code can be tested efficiently
4. **Idiomatic Go**: Follow Go best practices and conventions
5. **Modularity**: Keep components loosely coupled for easy replacement or extension

## Special Considerations
- Docker operations should be safe and properly handle permissions
- Consider edge cases like missing Docker daemon or failed container operations
- Ensure commands are idempotent when possible
- Add detailed logging for Docker operations for troubleshooting
- Consider users without Docker installed

## Implementation Guidelines

### Code Style
- Use meaningful variable and function names
- Add comments for complex logic
- Follow Go formatting standards (use gofmt)
- Implement proper error wrapping to maintain context
- Limit function complexity by breaking down operations

### Testing Approach
- Test happy paths and error cases
- Use test fixtures for Docker operations
- Mock external dependencies when appropriate
- Test edge cases specifically

### Commits
- Make small, focused commits with clear messages
- Group related changes together
- Ensure tests pass before committing

## Evaluation Criteria
- Correctness: Does the code correctly implement Docker operations?
- Error handling: Are errors handled appropriately?
- Test coverage: Are all code paths tested?
- Maintainability: Is the code easy to understand and modify?
- Performance: Are operations reasonably efficient?

## Potential Challenges
- Docker SDK versioning and compatibility
- Error handling for Docker operations
- Testing Docker operations in CI environments
- Handling Docker authentication
- Cross-platform issues (Windows vs. Linux/Mac)

## Resources
Project files:
- Architecture in `/ARCHITECTURE.md`
- Roadmap in `/DETAILED_ROADMAP.md`
- Example commands in `/USAGE.md`
- Current status in `/SUMMARY.md`