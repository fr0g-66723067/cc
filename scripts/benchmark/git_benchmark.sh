#\!/bin/bash
set -e

# Configuration 
TEMP_DIR="/tmp/cc-git-benchmark"
NUM_BRANCHES=20  # Number of branches to create
NUM_FILES=5      # Number of files per branch

echo "=== Git Operations Benchmark ==="
echo "Testing Git operations performance with $NUM_BRANCHES branches and $NUM_FILES files per branch"

# Setup
echo "Setting up test environment..."
rm -rf $TEMP_DIR
mkdir -p $TEMP_DIR
cd $TEMP_DIR

# Initialize git repo
echo "Initializing git repository..."
git init
git config user.name "Benchmark Test"
git config user.email "benchmark@example.com"

# Create and commit initial file
echo "# Git Benchmark Test" > README.md
git add README.md
git commit -m "Initial commit"

# Get the name of the default branch
MAIN_BRANCH=$(git branch --show-current)
echo "Main branch is: $MAIN_BRANCH"

# Record start time
START_TIME=$(date +%s)

# Create branches with files
echo "Creating $NUM_BRANCHES branches with $NUM_FILES files each..."
for ((i=1; i<=$NUM_BRANCHES; i++))
do
    # Create branch
    BRANCH_NAME="impl-test-$i"
    git checkout -b $BRANCH_NAME $MAIN_BRANCH
    
    # Create files
    for ((j=1; j<=$NUM_FILES; j++))
    do
        echo "This is test file $j on branch $BRANCH_NAME" > "file-$j.txt"
        git add "file-$j.txt"
    done
    
    # Commit changes
    git commit -m "Implementation $i with $NUM_FILES files"
    
    echo "Created branch $BRANCH_NAME with $NUM_FILES files"
done

# Measure branch creation time
BRANCH_END_TIME=$(date +%s)
BRANCH_DURATION=$((BRANCH_END_TIME - START_TIME))
echo "Branch creation completed in $BRANCH_DURATION seconds"

# Test feature branch creation
echo "Testing feature branch creation..."
FEATURE_START=$(date +%s)

# Select a random implementation branch
RANDOM_IMPL=$((1 + RANDOM % NUM_BRANCHES))
IMPL_BRANCH="impl-test-$RANDOM_IMPL"
git checkout $IMPL_BRANCH

# Create 5 feature branches
for ((i=1; i<=5; i++))
do
    # Create feature branch
    FEATURE_BRANCH="feat-test-$i-from-$RANDOM_IMPL"
    git checkout -b $FEATURE_BRANCH $IMPL_BRANCH
    
    # Create and modify files
    echo "Feature $i on top of implementation $RANDOM_IMPL" > "feature-$i.txt"
    echo "Modified by feature $i" >> "file-1.txt"
    
    # Commit changes
    git add "feature-$i.txt" "file-1.txt"
    git commit -m "Feature $i added to implementation $RANDOM_IMPL"
    
    echo "Created feature branch $FEATURE_BRANCH"
    
    # Go back to implementation branch
    git checkout $IMPL_BRANCH
done

# Measure feature branch time
FEATURE_END_TIME=$(date +%s)
FEATURE_DURATION=$((FEATURE_END_TIME - FEATURE_START))
echo "Feature branch creation completed in $FEATURE_DURATION seconds"

# Test branch switching
echo "Testing branch switching performance..."
SWITCH_START=$(date +%s)

# Get all branches
BRANCHES=$(git branch | cut -c 3-)

# Switch to each branch and back
for BRANCH in $BRANCHES
do
    git checkout $BRANCH >/dev/null 2>&1
    git checkout $MAIN_BRANCH >/dev/null 2>&1
done

# Measure switching time
SWITCH_END_TIME=$(date +%s)
SWITCH_DURATION=$((SWITCH_END_TIME - SWITCH_START))
echo "Branch switching completed in $SWITCH_DURATION seconds"

# Test diffing
echo "Testing branch comparison performance..."
DIFF_START=$(date +%s)

# Compare each implementation with main branch
for ((i=1; i<=$NUM_BRANCHES; i++))
do
    BRANCH_NAME="impl-test-$i"
    git diff --name-only $MAIN_BRANCH..$BRANCH_NAME >/dev/null
done

# Compare each feature with its base implementation
for ((i=1; i<=5; i++))
do
    git diff --name-only $IMPL_BRANCH.."feat-test-$i-from-$RANDOM_IMPL" >/dev/null
done

# Measure diff time
DIFF_END_TIME=$(date +%s)
DIFF_DURATION=$((DIFF_END_TIME - DIFF_START))
echo "Branch comparison completed in $DIFF_DURATION seconds"

# Calculate total time
END_TIME=$(date +%s)
TOTAL_DURATION=$((END_TIME - START_TIME))

# Show repository stats
BRANCH_COUNT=$(git branch | wc -l)
COMMIT_COUNT=$(git rev-list --all --count)
REPO_SIZE=$(du -sh . | cut -f1)

# Print results
echo ""
echo "=== Benchmark Results ==="
echo "Total branches created: $BRANCH_COUNT"
echo "Total commits: $COMMIT_COUNT"
echo "Repository size: $REPO_SIZE"
echo "Implementation branch creation time: $BRANCH_DURATION seconds"
echo "Feature branch creation time: $FEATURE_DURATION seconds"
echo "Branch switching time: $SWITCH_DURATION seconds"
echo "Branch comparison time: $DIFF_DURATION seconds"
echo "Total test time: $TOTAL_DURATION seconds"

echo ""
echo "Operations per second:"
echo "Branch creation rate: $(echo "scale=2; $NUM_BRANCHES / $BRANCH_DURATION" | bc -l) branches/s"
echo "Branch switching rate: $(echo "scale=2; $(($BRANCH_COUNT * 2)) / $SWITCH_DURATION" | bc -l) switches/s"
echo "Diff comparison rate: $(echo "scale=2; $(($NUM_BRANCHES + 5)) / $DIFF_DURATION" | bc -l) diffs/s"

# Ask if user wants to clean up
read -p "Do you want to clean up the benchmark artifacts? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Cleaning up..."
    cd /tmp
    rm -rf $TEMP_DIR
    echo "Cleanup complete"
else
    echo "Artifacts kept at $TEMP_DIR"
fi

echo "Benchmark complete\!"
