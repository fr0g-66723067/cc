# Code Controller (CC) Guidelines

## Build & Development Commands
- `go build -o cc ./cmd/cc` - Build the CC binary
- `go run ./cmd/cc` - Run the application without building
- `go test ./...` - Run all tests
- `go test ./internal/docker` - Run tests for a specific package
- `go test -v -run TestName` - Run a specific test
- `go fmt ./...` - Format all code
- `go vet ./...` - Run static code analysis

## Code Style Guidelines
- **Language**: Go 1.21+ with modules
- **Formatting**: Use gofmt, 80-character soft line limit
- **Imports**: Group standard library, external, and internal imports
- **Naming**: CamelCase for exported items, camelCase for non-exported
- **Error handling**: Always check errors, use descriptive error messages
- **Documentation**: Document all exported functions, types, and packages
- **Testing**: Write tests for all functionality, aim for >80% coverage
- **Dependencies**: Minimize external dependencies
- **Git workflow**: Feature branches from main, descriptive commits
- **Docker**: Use multi-stage builds for smaller images

## Project Architecture
- CLI tool (Go) communicates with Claude Code in Docker
- Git operations managed internally for branch management
- Configuration stored in user's home directory
- Container executes Claude Code commands in isolated environment