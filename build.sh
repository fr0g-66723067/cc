#!/bin/bash
# Main build script for Code Controller

set -e  # Exit on error

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "=== Code Controller Build Script ==="

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Using Go version: $GO_VERSION"

# Check for Docker to build the Claude container
if ! command -v docker &> /dev/null; then
    echo "Warning: Docker is not installed. Claude container will not be built."
    BUILD_CONTAINER=false
else
    BUILD_CONTAINER=true
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        echo "Warning: Docker is not running. Claude container will not be built."
        BUILD_CONTAINER=false
    fi
fi

# Check for Claude API key
if [ -z "$CLAUDE_API_KEY" ]; then
    echo "Warning: CLAUDE_API_KEY is not set. Some functionality may be limited."
fi

# Build the Claude Docker container if Docker is available
if [ "$BUILD_CONTAINER" = true ]; then
    echo "Building Claude Docker container..."
    if [ -f "$SCRIPT_DIR/claude/build.sh" ]; then
        chmod +x "$SCRIPT_DIR/claude/build.sh"
        "$SCRIPT_DIR/claude/build.sh"
    else
        echo "Warning: Claude Docker build script not found at $SCRIPT_DIR/claude/build.sh"
    fi
fi

# Build the CLI tool
echo "Building CC CLI tool..."
cd "$SCRIPT_DIR"
go build -o cc ./cmd/cc

if [ $? -eq 0 ]; then
    echo "Build successful! CC binary is at $SCRIPT_DIR/cc"
    
    # Make the binary executable
    chmod +x "$SCRIPT_DIR/cc"
    
    echo "You can now use the tool with: $SCRIPT_DIR/cc"
    echo "For help, run: $SCRIPT_DIR/cc --help"
    
    # Create the .cc directory if it doesn't exist
    mkdir -p "$HOME/.cc"
    
    echo "To install CC globally, run:"
    echo "sudo cp $SCRIPT_DIR/cc /usr/local/bin/"
else
    echo "Error: Build failed"
    exit 1
fi

# Run tests if requested
if [ "${1:-}" = "--with-tests" ]; then
    echo "Running tests..."
    go test ./...
fi