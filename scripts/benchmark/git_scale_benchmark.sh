#\!/bin/bash
set -e

# Configuration 
TEMP_DIR="/tmp/cc-git-scale-benchmark"
NUM_BRANCHES=50  # Increased number of branches
NUM_FILES=20     # Increased number of files per branch
FILE_SIZE=500    # Size of each file in lines
NUM_FEATURES=10  # Number of feature branches
NUM_COMMITS=5    # Commits per branch

echo "=== Git Operations Scale Benchmark ==="
echo "Testing Git scalability with:"
echo "- $NUM_BRANCHES implementation branches"
echo "- $NUM_FILES files per branch (each with $FILE_SIZE lines)"
echo "- $NUM_FEATURES feature branches"
echo "- $NUM_COMMITS commits per branch"

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

# Generate test file content
generate_content() {
    local branch=$1
    local file=$2
    local size=$3
    
    for ((i=1; i<=$size; i++)); do
        echo "Line $i of test file $file on branch $branch, with some additional padding text to make the file larger and more realistic in terms of content size, simulating actual code that might be generated in a real world scenario with variable length lines and different patterns of text."
    done
}

# Record start time
START_TIME=$(date +%s)

# Create branches with files
echo "Creating $NUM_BRANCHES branches with $NUM_FILES files and $NUM_COMMITS commits each..."
for ((i=1; i<=$NUM_BRANCHES; i++))
do
    # Create branch
    BRANCH_NAME="impl-test-$i"
    git checkout -b $BRANCH_NAME $MAIN_BRANCH
    
    # Create initial files
    for ((j=1; j<=$NUM_FILES; j++))
    do
        generate_content "$BRANCH_NAME" "$j" "$FILE_SIZE" > "file-$j.txt"
        git add "file-$j.txt"
    done
    
    # Initial commit
    git commit -m "Implementation $i - Initial commit with $NUM_FILES files"
    
    # Add additional commits with modifications
    for ((c=2; c<=$NUM_COMMITS; c++))
    do
        # Modify some random files
        NUM_MODS=$((1 + RANDOM % 5))
        for ((m=1; m<=$NUM_MODS; m++))
        do
            FILE_TO_MOD=$((1 + RANDOM % NUM_FILES))
            echo "Additional content from commit $c" >> "file-$FILE_TO_MOD.txt"
            git add "file-$FILE_TO_MOD.txt"
        done
        
        git commit -m "Implementation $i - Commit $c: Modified $NUM_MODS files"
    done
    
    echo "Created branch $BRANCH_NAME with $NUM_FILES files and $NUM_COMMITS commits"
    
    # Every 10 branches, report progress and time
    if [ $((i % 10)) -eq 0 ]; then
        CURRENT_TIME=$(date +%s)
        ELAPSED=$((CURRENT_TIME - START_TIME))
        echo "Progress: $i/$NUM_BRANCHES branches created in $ELAPSED seconds"
    fi
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

# Create feature branches
for ((i=1; i<=$NUM_FEATURES; i++))
do
    # Create feature branch
    FEATURE_BRANCH="feat-test-$i-from-$RANDOM_IMPL"
    git checkout -b $FEATURE_BRANCH $IMPL_BRANCH
    
    # Create a feature file and modify several existing files
    echo "Feature $i on top of implementation $RANDOM_IMPL" > "feature-$i.txt"
    git add "feature-$i.txt"
    
    # Modify some existing files
    NUM_MODS=$((2 + RANDOM % 5))
    for ((m=1; m<=$NUM_MODS; m++))
    do
        FILE_TO_MOD=$((1 + RANDOM % NUM_FILES))
        echo "Modified by feature $i" >> "file-$FILE_TO_MOD.txt"
        git add "file-$FILE_TO_MOD.txt"
    done
    
    # Commit changes
    git commit -m "Feature $i added to implementation $RANDOM_IMPL"
    
    # Make additional commits
    for ((c=2; c<=3; c++))
    do
        # Modify some random files
        NUM_MODS=$((1 + RANDOM % 3))
        for ((m=1; m<=$NUM_MODS; m++))
        do
            FILE_TO_MOD=$((1 + RANDOM % NUM_FILES))
            echo "Additional feature $i change from commit $c" >> "file-$FILE_TO_MOD.txt"
            git add "file-$FILE_TO_MOD.txt"
        done
        
        git commit -m "Feature $i - Commit $c: Modified $NUM_MODS files"
    done
    
    echo "Created feature branch $FEATURE_BRANCH with 3 commits"
    
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
BRANCHES=$(git branch | cut -c 3- | head -20)  # Limit to 20 branches for switching test

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

# Compare a sample of implementations with main branch
for ((i=1; i<=10; i++))
do
    SAMPLE_IDX=$(( i * (NUM_BRANCHES / 10) ))
    BRANCH_NAME="impl-test-$SAMPLE_IDX"
    git diff --stat $MAIN_BRANCH..$BRANCH_NAME >/dev/null
done

# Compare all features with their base implementation
for ((i=1; i<=$NUM_FEATURES; i++))
do
    git diff --stat $IMPL_BRANCH.."feat-test-$i-from-$RANDOM_IMPL" >/dev/null
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
echo "Total branches: $BRANCH_COUNT"
echo "Total commits: $COMMIT_COUNT"
echo "Repository size: $REPO_SIZE"
echo "Implementation branch creation time: $BRANCH_DURATION seconds ($(echo "scale=2; $NUM_BRANCHES / $BRANCH_DURATION" | bc -l) branches/s)"
echo "Feature branch creation time: $FEATURE_DURATION seconds ($(echo "scale=2; $NUM_FEATURES / $FEATURE_DURATION" | bc -l) branches/s)"
echo "Branch switching time: $SWITCH_DURATION seconds ($(echo "scale=2; 40 / $SWITCH_DURATION" | bc -l) switches/s)"
echo "Branch comparison time: $DIFF_DURATION seconds ($(echo "scale=2; 20 / $DIFF_DURATION" | bc -l) diffs/s)"
echo "Total test time: $TOTAL_DURATION seconds"

# Generate report
echo ""
echo "=== Performance Summary ==="
if [ $BRANCH_DURATION -gt 0 ]; then
    BRANCH_RATE=$(echo "scale=2; $NUM_BRANCHES / $BRANCH_DURATION" | bc -l)
    echo "- Can create $BRANCH_RATE implementation branches per second"
else
    echo "- Branch creation too fast to measure accurately"
fi

if [ $FEATURE_DURATION -gt 0 ]; then
    FEATURE_RATE=$(echo "scale=2; $NUM_FEATURES / $FEATURE_DURATION" | bc -l)
    echo "- Can create $FEATURE_RATE feature branches per second"
else
    echo "- Feature creation too fast to measure accurately"
fi

if [ $SWITCH_DURATION -gt 0 ]; then
    SWITCH_RATE=$(echo "scale=2; 40 / $SWITCH_DURATION" | bc -l)
    echo "- Can perform $SWITCH_RATE branch switches per second"
else
    echo "- Branch switching too fast to measure accurately"
fi

if [ $DIFF_DURATION -gt 0 ]; then
    DIFF_RATE=$(echo "scale=2; 20 / $DIFF_DURATION" | bc -l)
    echo "- Can compute $DIFF_RATE branch diffs per second"
else
    echo "- Branch diffing too fast to measure accurately"
fi

TOTAL_FILES=$((NUM_BRANCHES * NUM_FILES + NUM_FEATURES))
echo "- Repository contains approximately $TOTAL_FILES files"
echo "- Repository size is $REPO_SIZE"

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
