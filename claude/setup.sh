#!/bin/bash
# Setup script for Claude Code Controller

set -e  # Exit on error

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." &> /dev/null && pwd )"

echo "=== Claude Code Controller Setup ==="
echo "This script will help you set up the Claude Code Controller."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Check for Claude API key
if [ -z "$CLAUDE_API_KEY" ]; then
    # Check if .env file exists in project root
    if [ -f "$PROJECT_ROOT/.env" ]; then
        echo "Found .env file in project root, checking for CLAUDE_API_KEY..."
        if grep -q "CLAUDE_API_KEY" "$PROJECT_ROOT/.env"; then
            echo "CLAUDE_API_KEY found in .env file."
            source "$PROJECT_ROOT/.env"
        else
            echo "CLAUDE_API_KEY not found in .env file."
            ENV_FILE="$PROJECT_ROOT/.env"
            CREATE_ENV=true
        fi
    else
        # Check if .env file exists in user home directory
        USER_ENV="$HOME/.env"
        if [ -f "$USER_ENV" ] && grep -q "CLAUDE_API_KEY" "$USER_ENV"; then
            echo "CLAUDE_API_KEY found in $USER_ENV file."
            source "$USER_ENV"
        else
            echo "CLAUDE_API_KEY not found in any .env file."
            ENV_FILE="$PROJECT_ROOT/.env"
            CREATE_ENV=true
        fi
    fi
    
    # If still no API key, prompt user
    if [ -z "$CLAUDE_API_KEY" ]; then
        echo ""
        echo "No Claude API key found in environment or .env files."
        echo "Please enter your Claude API key (starts with 'sk-'):"
        read -p "> " API_KEY
        
        if [ -z "$API_KEY" ]; then
            echo "No API key provided. Exiting."
            exit 1
        fi
        
        export CLAUDE_API_KEY="$API_KEY"
        
        # Create or update .env file
        if [ "$CREATE_ENV" = true ]; then
            echo "Creating .env file at $ENV_FILE..."
            echo "CLAUDE_API_KEY=$API_KEY" > "$ENV_FILE"
            echo ".env file created with API key."
        else
            echo "Please add your API key to your environment or .env file manually."
        fi
    fi
fi

# Build the Claude Code Docker image
echo ""
echo "Building Claude Code Docker image..."
if [ -f "$SCRIPT_DIR/build.sh" ]; then
    chmod +x "$SCRIPT_DIR/build.sh"
    "$SCRIPT_DIR/build.sh"
else
    echo "Error: build.sh not found at $SCRIPT_DIR/build.sh"
    exit 1
fi

# Build the CC tool
echo ""
echo "Building CC tool..."
cd "$PROJECT_ROOT"
go build -o cc ./cmd/cc
if [ $? -ne 0 ]; then
    echo "Error: Failed to build CC tool."
    exit 1
fi

echo ""
echo "Setup complete! You can now use the CC tool."
echo "Example command: ./cc init my-project"
echo ""
echo "For more information, see the README.md file."