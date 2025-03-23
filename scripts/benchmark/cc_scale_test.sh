#\!/bin/bash
set -e

# Configuration 
TEMP_DIR="/tmp/cc-scale-test"
CONFIG_FILE="/tmp/cc-scale-config.json"
PROJECT_NAME="scale-test"

# Test size parameters
NUM_IMPLEMENTATIONS=20   # Number of implementations to generate
NUM_FEATURES=10         # Number of features to add to first implementation
FRAMEWORK_BASE="framework"  # Base name for frameworks

echo "=== CC Git Scalability Test ==="
echo "Testing CC Git operations with $NUM_IMPLEMENTATIONS implementations and $NUM_FEATURES features"

# Clean up any previous run
rm -rf $TEMP_DIR
rm -f $CONFIG_FILE

# Create config file
mkdir -p $(dirname $CONFIG_FILE)
cat > $CONFIG_FILE << CONF
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
      "user.name": "Scale Test",
      "user.email": "scale@example.com"
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
cd /home/user/src/cc
if [ \! -f "./cc" ]; then
    echo "Building CC..."
    go build -o cc ./cmd/cc
fi

START_TIME=$(date +%s)

# Step 1: Initialize project
echo "Initializing project..."
./cc init $PROJECT_NAME -d "Scalability test project" --config $CONFIG_FILE

# Step 2: Generate frameworks list
FRAMEWORKS=""
for ((i=1; i<=$NUM_IMPLEMENTATIONS; i++)); do
    if [ $i -gt 1 ]; then
        FRAMEWORKS+=","
    fi
    FRAMEWORKS+="${FRAMEWORK_BASE}-$i"
done

# Step 3: Generate implementations 
echo "Generating $NUM_IMPLEMENTATIONS implementations..."
GEN_START=$(date +%s)
./cc generate "Create a scalable application" --frameworks $FRAMEWORKS --count $NUM_IMPLEMENTATIONS --config $CONFIG_FILE
GEN_END=$(date +%s)
GEN_TIME=$((GEN_END - GEN_START))
echo "Generated $NUM_IMPLEMENTATIONS implementations in $GEN_TIME seconds"

# Step 4: List and get first implementation
IMPLEMENTATIONS=$(./cc list implementations --config $CONFIG_FILE)
FIRST_IMPL=$(echo "$IMPLEMENTATIONS" | grep "impl-" | head -1 | awk '{print $2}')
echo "First implementation: $FIRST_IMPL"

# Step 5: Select the first implementation
echo "Selecting implementation..."
./cc select $FIRST_IMPL --config $CONFIG_FILE

# Step 6: Add features to the first implementation
echo "Adding $NUM_FEATURES features..."
FEATURE_START=$(date +%s)
for ((i=1; i<=$NUM_FEATURES; i++)); do
    echo "Adding feature $i/$NUM_FEATURES..."
    ./cc feature "Add feature $i functionality" --config $CONFIG_FILE
done
FEATURE_END=$(date +%s)
FEATURE_TIME=$((FEATURE_END - FEATURE_START))
echo "Added $NUM_FEATURES features in $FEATURE_TIME seconds"

# Step 7: Check git repository statistics
cd $TEMP_DIR/$PROJECT_NAME
BRANCH_COUNT=$(git branch | wc -l)
COMMIT_COUNT=$(git rev-list --all --count)
REPO_SIZE=$(du -sh . | cut -f1)

# Get branches grouped by type
IMPL_BRANCHES=$(git branch | grep "impl-" | wc -l)
FEAT_BRANCHES=$(git branch | grep "feat-" | wc -l)

# Calculate total time
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

# Print results
echo ""
echo "=== Scalability Test Results ==="
echo "Total time: $TOTAL_TIME seconds"
echo "Repository statistics:"
echo "- Repository size: $REPO_SIZE"
echo "- Total branches: $BRANCH_COUNT"
echo "- Implementation branches: $IMPL_BRANCHES"
echo "- Feature branches: $FEAT_BRANCHES"
echo "- Total commits: $COMMIT_COUNT"
echo ""
echo "Performance metrics:"
echo "- Implementation generation: $NUM_IMPLEMENTATIONS implementations in $GEN_TIME seconds"
if [ $GEN_TIME -gt 0 ]; then
    echo "  - Rate: $(echo "scale=2; $NUM_IMPLEMENTATIONS / $GEN_TIME" | bc) implementations/second"
else
    echo "  - Rate: Too fast to measure accurately"
fi
echo "- Feature addition: $NUM_FEATURES features in $FEATURE_TIME seconds"
if [ $FEATURE_TIME -gt 0 ]; then
    echo "  - Rate: $(echo "scale=2; $NUM_FEATURES / $FEATURE_TIME" | bc) features/second"
else
    echo "  - Rate: Too fast to measure accurately"
fi
echo "- Average commit creation: $(echo "scale=2; $COMMIT_COUNT / $TOTAL_TIME" | bc) commits/second"
echo "- Average branch creation: $(echo "scale=2; $BRANCH_COUNT / $TOTAL_TIME" | bc) branches/second"

# Ask if user wants to clean up
read -p "Do you want to clean up the test artifacts? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Cleaning up..."
    rm -rf $TEMP_DIR
    rm -f $CONFIG_FILE
    echo "Cleanup complete"
else
    echo "Artifacts kept at $TEMP_DIR and $CONFIG_FILE"
fi

echo "Scalability test complete\!"
