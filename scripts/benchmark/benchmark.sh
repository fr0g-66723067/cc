#\!/bin/bash
set -e

# Configuration 
CONFIG_PATH="/tmp/cc-benchmark-config.json"
PROJECT_NAME="benchmark-test"
TEMP_DIR="/tmp/cc-benchmark"
NUM_IMPLEMENTATIONS=10  # Number of implementations to generate
FRAMEWORKS="nodejs,python,go,java,ruby,php,rust,cpp,typescript,kotlin"

echo "=== CC Benchmark Test ==="
echo "Testing Git operations performance with $NUM_IMPLEMENTATIONS implementations"

# Create start time
START_TIME=$(date +%s)

# Setup
echo "Setting up test environment..."
mkdir -p $TEMP_DIR
mkdir -p $(dirname $CONFIG_PATH)

# Create temporary config file
cat > $CONFIG_PATH << CONF
{
  "container": {
    "provider": "docker",
    "config": {},
    "claudeImage": "anthropic/claude-code:latest"
  },
  "ai": {
    "provider": "claude",
    "config": {}
  },
  "vcs": {
    "provider": "git",
    "config": {
      "user.name": "Benchmark Test",
      "user.email": "benchmark@example.com"
    }
  },
  "projectsDir": "$TEMP_DIR",
  "activeProject": "",
  "jobs": {
    "maxConcurrent": 4,
    "timeout": 3600
  },
  "plugins": {
    "dir": "/home/user/.cc/plugins",
    "enabled": []
  },
  "projects": {}
}
CONF

# Build CC if needed
if [ \! -f "/home/user/src/cc/cc" ]; then
    echo "Building CC..."
    cd /home/user/src/cc
    go build -o cc ./cmd/cc
fi

# Initialize project
echo "Initializing project..."
cd /home/user/src/cc
./cc init $PROJECT_NAME -d "Benchmark test project" --config $CONFIG_PATH

# Generate implementations
echo "Generating $NUM_IMPLEMENTATIONS implementations (this will take some time)..."
./cc generate "Create a simple Hello World application" --frameworks $FRAMEWORKS --count $NUM_IMPLEMENTATIONS --config $CONFIG_PATH

# Measure generation time
GEN_END_TIME=$(date +%s)
GEN_DURATION=$((GEN_END_TIME - START_TIME))
echo "Generation completed in $GEN_DURATION seconds"

# Get number of branches
cd $TEMP_DIR/$PROJECT_NAME
BRANCH_COUNT=$(git branch | wc -l)
echo "Number of branches created: $BRANCH_COUNT"

# Test branch switching speed
echo "Testing branch switching performance..."
SWITCH_START=$(date +%s)

# Get all implementation branches
IMPL_BRANCHES=$(./cc list implementations --config $CONFIG_PATH | grep "impl-" | sed 's/^[[:space:]]*[0-9]\+\. //')

# Select first implementation
FIRST_IMPL=$(echo "$IMPL_BRANCHES" | head -n 1)
cd /home/user/src/cc
./cc select $FIRST_IMPL --config $CONFIG_PATH

# Add a feature
echo "Adding a feature..."
./cc feature "Add a simple feature" --config $CONFIG_PATH

# Measure feature time
FEATURE_END_TIME=$(date +%s)
FEATURE_DURATION=$((FEATURE_END_TIME - SWITCH_START))
echo "Feature addition completed in $FEATURE_DURATION seconds"

# Get feature branch
FEATURE_BRANCH=$(cd /home/user/src/cc && ./cc list features --config $CONFIG_PATH | grep "feat-" | sed 's/^[[:space:]]*[0-9]\+\. //')

# Compare branches
echo "Testing branch comparison..."
COMPARE_START=$(date +%s)
cd /home/user/src/cc
./cc compare $FIRST_IMPL $FEATURE_BRANCH --config $CONFIG_PATH > /dev/null

# Measure comparison time
COMPARE_END_TIME=$(date +%s)
COMPARE_DURATION=$((COMPARE_END_TIME - COMPARE_START))
echo "Branch comparison completed in $COMPARE_DURATION seconds"

# Calculate total time
END_TIME=$(date +%s)
TOTAL_DURATION=$((END_TIME - START_TIME))

# Print results
echo ""
echo "=== Benchmark Results ==="
echo "Total branches created: $BRANCH_COUNT"
echo "Implementation generation time: $GEN_DURATION seconds"
echo "Feature addition time: $FEATURE_DURATION seconds"
echo "Branch comparison time: $COMPARE_DURATION seconds"
echo "Total test time: $TOTAL_DURATION seconds"

# Ask if user wants to clean up
read -p "Do you want to clean up the benchmark artifacts? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Cleaning up..."
    rm -rf $TEMP_DIR
    rm -f $CONFIG_PATH
    echo "Cleanup complete"
else
    echo "Artifacts kept at $TEMP_DIR and $CONFIG_PATH"
fi

echo "Benchmark complete\!"
