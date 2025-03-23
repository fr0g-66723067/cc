FROM node:20

# Install basic tools
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    git \
    sudo \
    curl \
    jq \
    unzip

# Create workspace
RUN mkdir -p /workspace && chmod 777 /workspace
WORKDIR /workspace

# Install real Claude CLI
RUN curl -fsSL https://anthropic.github.io/claude-code-cli/install.sh | bash

# Create .env loader script
RUN echo '#!/bin/bash' > /usr/local/bin/load-env && \
    echo 'if [ -f "/.env" ]; then' >> /usr/local/bin/load-env && \
    echo '  echo "Loading environment variables from /.env"' >> /usr/local/bin/load-env && \
    echo '  set -a' >> /usr/local/bin/load-env && \
    echo '  source /.env' >> /usr/local/bin/load-env && \
    echo '  set +a' >> /usr/local/bin/load-env && \
    echo 'fi' >> /usr/local/bin/load-env && \
    echo 'if [ -f "/workspace/.env" ]; then' >> /usr/local/bin/load-env && \
    echo '  echo "Loading environment variables from /workspace/.env"' >> /usr/local/bin/load-env && \
    echo '  set -a' >> /usr/local/bin/load-env && \
    echo '  source /workspace/.env' >> /usr/local/bin/load-env && \
    echo '  set +a' >> /usr/local/bin/load-env && \
    echo 'fi' >> /usr/local/bin/load-env && \
    chmod +x /usr/local/bin/load-env

# Create claude wrapper script
RUN echo '#!/bin/bash' > /usr/local/bin/claude && \
    echo '# Load environment variables from .env files' >> /usr/local/bin/claude && \
    echo 'source /usr/local/bin/load-env' >> /usr/local/bin/claude && \
    echo '' >> /usr/local/bin/claude && \
    echo '# Check if Claude API key is available' >> /usr/local/bin/claude && \
    echo 'if [[ -z "$CLAUDE_API_KEY" ]]; then' >> /usr/local/bin/claude && \
    echo '  echo "Error: CLAUDE_API_KEY not set"' >> /usr/local/bin/claude && \
    echo '  echo "Please provide it in /.env or /workspace/.env file or as environment variable"' >> /usr/local/bin/claude && \
    echo '  exit 1' >> /usr/local/bin/claude && \
    echo 'fi' >> /usr/local/bin/claude && \
    echo '' >> /usr/local/bin/claude && \
    echo '# Check for Claude Code CLI' >> /usr/local/bin/claude && \
    echo 'if ! command -v claude-code &> /dev/null; then' >> /usr/local/bin/claude && \
    echo '  echo "Error: claude-code CLI not found"' >> /usr/local/bin/claude && \
    echo '  exit 1' >> /usr/local/bin/claude && \
    echo 'fi' >> /usr/local/bin/claude && \
    echo '' >> /usr/local/bin/claude && \
    echo '# Execute the real claude-code CLI' >> /usr/local/bin/claude && \
    echo 'if [[ "$1" == "code" ]]; then' >> /usr/local/bin/claude && \
    echo '  # Handle special version command directly' >> /usr/local/bin/claude && \
    echo '  if [[ "$2" == "--version" ]]; then' >> /usr/local/bin/claude && \
    echo '    claude-code --version' >> /usr/local/bin/claude && \
    echo '    exit $?' >> /usr/local/bin/claude && \
    echo '  fi' >> /usr/local/bin/claude && \
    echo '' >> /usr/local/bin/claude && \
    echo '  # Pass all arguments to claude-code' >> /usr/local/bin/claude && \
    echo '  shift  # Remove "code" from arguments' >> /usr/local/bin/claude && \
    echo '  exec claude-code "$@"' >> /usr/local/bin/claude && \
    echo 'else' >> /usr/local/bin/claude && \
    echo '  # For direct claude CLI commands' >> /usr/local/bin/claude && \
    echo '  echo "Running claude command: $@"' >> /usr/local/bin/claude && \
    echo '  exec claude-code "$@"' >> /usr/local/bin/claude && \
    echo 'fi' >> /usr/local/bin/claude && \
    chmod +x /usr/local/bin/claude && \
    ln -sf /usr/local/bin/claude /usr/local/bin/claude-code-wrapper

# Create entrypoint script that ensures environment variables are loaded
RUN echo '#!/bin/bash' > /usr/local/bin/entrypoint.sh && \
    echo 'source /usr/local/bin/load-env' >> /usr/local/bin/entrypoint.sh && \
    echo 'exec "$@"' >> /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

ENV HOME=/home/node
ENV PATH="/home/node/.local/bin:$PATH"
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["bash"]

USER node
