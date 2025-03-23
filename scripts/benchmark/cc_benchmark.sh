#\!/bin/bash
set -e

# Configuration 
TEMP_DIR="/tmp/cc-benchmark-project"
CONFIG_FILE="/tmp/cc-benchmark-config.json"
PROJECT_NAME="bench-project"
NUM_IMPLEMENTATIONS=5
FRAMEWORKS="react,vue,angular,svelte,lit"

# Setup timing function
timestamp() {
  date +%s.%N
}

echo "=== CC Git Performance Benchmark ==="
echo "Testing Git operations in Code Controller with $NUM_IMPLEMENTATIONS implementations"

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
cd /home/user/src/cc
if [ \! -f "./cc" ]; then
    echo "Building CC..."
    go build -o cc ./cmd/cc
fi

# Function to run a command and measure time
run_timed() {
    local cmd="$1"
    local desc="$2"
    
    echo -n "Running $desc... "
    
    local start=$(timestamp)
    eval "$cmd"
    local end=$(timestamp)
    
    # Calculate time difference
    local duration=$(echo "$end - $start" | bc)
    printf "Completed in %.2f seconds\n" $duration
    
    # Return the duration
    echo $duration
}

# Step 1: Initialize project
echo -e "\n1. Initializing project"
init_time=$(run_timed "./cc init $PROJECT_NAME -d \"Benchmark project for testing Git operations\" --config $CONFIG_FILE" "project initialization")

# Step 2: Generate implementations
echo -e "\n2. Generating implementations"
gen_start=$(timestamp)
./cc generate "Create a simple todo application with basic CRUD operations" --frameworks $FRAMEWORKS --count $NUM_IMPLEMENTATIONS --config $CONFIG_FILE
gen_end=$(timestamp)
gen_time=$(echo "$gen_end - $gen_start" | bc)
printf "Generation completed in %.2f seconds (%.2f implementations per second)\n" $gen_time $(echo "$NUM_IMPLEMENTATIONS / $gen_time" | bc -l)

# Step 3: List implementations
echo -e "\n3. Listing implementations"
implementations=$(./cc list implementations --config $CONFIG_FILE)
echo "$implementations"

# Get the first implementation branch
impl_branch=$(echo "$implementations" | grep -o "impl-[^ ]*" | head -1)
echo "Selected implementation: $impl_branch"

# Step 4: Select implementation
echo -e "\n4. Selecting implementation"
select_time=$(run_timed "./cc select $impl_branch --config $CONFIG_FILE" "implementation selection")

# Step 5: Add features
echo -e "\n5. Adding features"
features=("Add dark mode support" "Add user authentication" "Add filtering capabilities")
feature_times=()

for feature in "${features[@]}"; do
    feat_time=$(run_timed "./cc feature \"$feature\" --config $CONFIG_FILE" "feature addition: $feature")
    feature_times+=($feat_time)
done

# Step 6: List features
echo -e "\n6. Listing features"
features_list=$(./cc list features --config $CONFIG_FILE)
echo "$features_list"

# Get a feature branch
feature_branch=$(echo "$features_list" | grep -o "feat-[^ ]*" | head -1)
echo "Selected feature: $feature_branch"

# Step 7: Compare branches
echo -e "\n7. Comparing branches"
compare_time=$(run_timed "./cc compare $impl_branch $feature_branch --config $CONFIG_FILE > /dev/null" "branch comparison")

# Step 8: Check status
echo -e "\n8. Checking status"
status_time=$(run_timed "./cc status --config $CONFIG_FILE > /dev/null" "status check")

# Check git repository size
repo_size=$(du -sh "$TEMP_DIR/$PROJECT_NAME" | cut -f1)

# Summarize results
echo -e "\n=== Benchmark Results ==="
echo "Project directory: $TEMP_DIR/$PROJECT_NAME"
echo "Repository size: $repo_size"
echo "Project initialization time: $init_time seconds"
echo "Implementation generation time: $gen_time seconds for $NUM_IMPLEMENTATIONS implementations"
echo "Implementation rate: $(echo "$NUM_IMPLEMENTATIONS / $gen_time" | bc -l) implementations per second"
echo "Implementation selection time: $select_time seconds"

# Calculate average feature time
total_feature_time=0
for t in "${feature_times[@]}"; do
    total_feature_time=$(echo "$total_feature_time + $t" | bc)
done
avg_feature_time=$(echo "$total_feature_time / ${#feature_times[@]}" | bc -l)

echo "Average feature addition time: $avg_feature_time seconds"
echo "Branch comparison time: $compare_time seconds"
echo "Status check time: $status_time seconds"

# Ask if user wants to clean up
read -p "Do you want to clean up the benchmark artifacts? (y/n) " -n 1 -r
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

echo "Benchmark complete\!"
