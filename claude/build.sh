#!/bin/bash
# Build script for Claude Code Docker image

set -e  # Exit on error

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: docker command not found"
    exit 1
fi

# Check if Dockerfile exists
if [ ! -f "$SCRIPT_DIR/Dockerfile" ]; then
    echo "Error: Dockerfile not found at $SCRIPT_DIR/Dockerfile"
    exit 1
fi

# Check if env-handler.sh exists, and if not, create it
if [ ! -f "$SCRIPT_DIR/env-handler.sh" ]; then
    echo "Warning: env-handler.sh not found, creating it..."
    cat > "$SCRIPT_DIR/env-handler.sh" << 'EOF'
#!/bin/bash
# This script helps manage environment variables for the Claude container
# It can copy .env files to the container and check for the presence of required variables

set -e  # Exit on error

# Check if we have required commands
if ! command -v docker &> /dev/null; then
    echo "Error: docker command not found"
    exit 1
fi

# Function to check if a container exists
container_exists() {
    local container_id="$1"
    docker ps -a --format "{{.ID}}" | grep -q "$container_id"
    return $?
}

# Function to copy .env file to container
copy_env_to_container() {
    local container_id="$1"
    local env_file="$2"
    
    if [ -z "$container_id" ]; then
        echo "Error: No container ID provided"
        return 1
    fi
    
    if [ ! -f "$env_file" ]; then
        echo "Error: .env file not found at $env_file"
        return 1
    fi
    
    if ! container_exists "$container_id"; then
        echo "Error: Container $container_id not found"
        return 1
    fi
    
    echo "Copying $env_file to container $container_id"
    docker cp "$env_file" "$container_id:/.env"
    
    # Verify the file was copied
    if docker exec "$container_id" test -f "/.env"; then
        echo "Successfully copied .env file to container"
        return 0
    else
        echo "Failed to copy .env file to container"
        return 1
    fi
}

# Function to check for CLAUDE_API_KEY in container
check_api_key_in_container() {
    local container_id="$1"
    
    if [ -z "$container_id" ]; then
        echo "Error: No container ID provided"
        return 1
    fi
    
    if ! container_exists "$container_id"; then
        echo "Error: Container $container_id not found"
        return 1
    fi
    
    echo "Checking for CLAUDE_API_KEY in container..."
    if docker exec "$container_id" bash -c 'source /usr/local/bin/load-env && [ -n "$CLAUDE_API_KEY" ]' &>/dev/null; then
        echo "CLAUDE_API_KEY is available in the container"
        return 0
    else
        echo "WARNING: CLAUDE_API_KEY is NOT available in the container"
        return 1
    fi
}

# Parse command line arguments
case "${1:-}" in
    copy-env)
        if [ $# -lt 3 ]; then
            echo "Usage: $0 copy-env CONTAINER_ID ENV_FILE"
            exit 1
        fi
        copy_env_to_container "$2" "$3"
        ;;
    check-key)
        if [ $# -lt 2 ]; then
            echo "Usage: $0 check-key CONTAINER_ID"
            exit 1
        fi
        check_api_key_in_container "$2"
        ;;
    *)
        echo "Usage: $0 [command] [args...]"
        echo "Commands:"
        echo "  copy-env CONTAINER_ID ENV_FILE - Copy .env file to container"
        echo "  check-key CONTAINER_ID - Check if CLAUDE_API_KEY is available in container"
        exit 1
        ;;
esac
EOF
    chmod +x "$SCRIPT_DIR/env-handler.sh"
fi

# Check if user wants to build with cache
USE_CACHE="--no-cache"
if [ "${1:-}" == "--use-cache" ]; then
    USE_CACHE=""
    shift
fi

# Set image name
IMAGE_NAME="${1:-claude-code}"
TAG="${2:-latest}"

echo "Building Claude Code Docker image: $IMAGE_NAME:$TAG"
echo "Using Dockerfile at: $SCRIPT_DIR/Dockerfile"

# Build the Docker image
docker build $USE_CACHE -t "$IMAGE_NAME:$TAG" -f "$SCRIPT_DIR/Dockerfile" "$SCRIPT_DIR"

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Successfully built Claude Code Docker image: $IMAGE_NAME:$TAG"
    
    # Show image details
    docker images "$IMAGE_NAME:$TAG"
    
    echo ""
    echo "To use this image with CC, set the environment variable:"
    echo "export CLAUDE_CODE_IMAGE=$IMAGE_NAME:$TAG"
    echo ""
    echo "Or add it to your configuration file at ~/.cc/config.json in the container section:"
    echo '"claudeImage": "'$IMAGE_NAME:$TAG'"'
else
    echo "Error: Failed to build Claude Code Docker image"
    exit 1
fi