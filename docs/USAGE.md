# Code Controller (CC) Usage Guide

Code Controller (CC) is a powerful command-line tool that interfaces with Claude Code CLI to automate code generation and manage multiple implementation approaches through Git branching.

## Getting Started

### Installation

```bash
# Build from source
git clone https://github.com/fr0g-66723067/cc.git
cd cc
go build -o cc ./cmd/cc

# Move to a location in your PATH
sudo mv cc /usr/local/bin/
```

### First-time Setup

Code Controller automatically creates its configuration directory at `~/.cc/` the first time you run a command. The configuration file is stored at `~/.cc/config.json`.

## Basic Workflow

The typical workflow with CC follows these steps:

1. **Initialize a project**: Create a new project with a description
2. **Generate implementations**: Generate multiple implementations using different frameworks
3. **Select an implementation**: Choose your preferred implementation
4. **Add features**: Incrementally add features to your selected implementation
5. **Compare implementations**: Compare different approaches as needed

## Interaction Methods

Code Controller offers three ways to interact with your projects:

1. **Traditional CLI Commands** - The original command structure (backward compatible)
2. **Hierarchical Command Structure** - Commands organized by resource type (projects, implementations, features)
3. **Interactive Shell Mode** - A shell-like interface with tab completion and context navigation

## Traditional Command Reference

### Initialize a Project

```bash
cc init my-project -d "A task management application"
```

This command:
- Creates a new project named "my-project"
- Sets the description to "A task management application"
- Initializes a Git repository in the project directory
- Sets this project as the active project

Options:
- `-d, --description`: Project description (optional)

### Generate Implementations

```bash
cc generate "Create a web app with user authentication and task management"
```

This command:
- Takes your description and generates multiple implementations
- Creates a separate Git branch for each implementation
- Uses different frameworks based on your options

Options:
- `--frameworks`: Comma-separated list of frameworks to use (e.g., `--frameworks react,vue,svelte`)
- `--count`: Number of implementations to generate (default: 3)
- `--parallel`: Generate in parallel (default: true)

### List Resources

```bash
cc list projects      # List all projects
cc list implementations  # List implementations for the active project
cc list features      # List features for the selected implementation
```

### Select an Implementation

```bash
cc select impl-react-1234567890
```

This command:
- Sets the specified implementation as the selected implementation
- Switches the Git branch to the implementation branch

### Add Features

```bash
cc feature "Add dark mode support"
```

This command:
- Creates a new branch based on the selected implementation
- Adds the requested feature to the codebase
- Commits the changes

### Compare Implementations

```bash
cc compare impl-react-1234567890 impl-vue-0987654321
```

This command:
- Shows a diff between two implementation branches

### Show Status

```bash
cc status
```

This command:
- Displays information about the active project
- Shows all implementations and their features
- Indicates which implementation is currently selected
- Shows your current context path in the hierarchy

### Remove Resources

```bash
cc remove project project-name
cc remove implementation impl-branch-name
cc remove feature feat-branch-name
```

This command:
- Removes the specified project, implementation, or feature
- For projects: removes the project from the configuration
- For implementations: removes the implementation from the project
- For features: removes the feature from the implementation

### Rename Resources

```bash
cc rename project old-name new-name
cc rename implementation old-branch-name new-branch-name
cc rename feature old-feature-branch new-feature-branch
```

This command:
- Renames the specified project, implementation, or feature
- For projects: updates the project name and directory
- For implementations: creates a new Git branch with the new name
- For features: creates a new Git branch with the new name

## Hierarchical Command Structure

The hierarchical command structure organizes commands by resource type:

```bash
# Projects Commands
cc projects list
cc projects create my-project "Project description"
cc projects remove my-project
cc projects rename old-name new-name

# Implementations Commands
cc implementations list
cc implementations generate "Create a web app" --frameworks=react,vue
cc implementations select branch-name
cc implementations remove branch-name
cc implementations rename old-branch new-branch
cc implementations compare branch1 branch2

# Features Commands
cc features list
cc features add "Add dark mode toggle"
cc features remove feature-branch
cc features rename old-feature-branch new-feature-branch
```

### Context Navigation

CC maintains a navigation context that remembers which project, implementation, and feature you're working with:

```bash
# Switch between projects, implementations, and features
cc use project my-project
cc use implementation impl-react-123456
cc use feature feat-darkmode-123456

# Show current context and status
cc status
```

## Interactive Shell Mode

For a more interactive experience, use the shell mode:

```bash
cc shell
```

The shell provides:
- Tab completion for commands and resources
- Context-aware prompts showing your location
- Command history
- Path-based navigation similar to a filesystem

Shell commands include:

```
# Context navigation
pwd                                   # Show current context path
cd /projects/my-project               # Navigate to a project
cd implementations/impl-react-123456  # Navigate to an implementation
cd /                                  # Return to root level

# Resource management
projects list
projects create my-new-project "Description"
implementations generate "Create a login page" --frameworks=react,vue
features add "Add password reset"

# Using resources
use project my-project
use implementation impl-react-123456

# Show status
status

# Exit the shell
exit
```

## Advanced Usage

### Custom Configuration

You can specify a custom configuration file:

```bash
cc --config /path/to/config.json [command]
```

### Working with Multiple Projects

CC supports managing multiple projects. Use `cc list projects` to see all projects, and switch between them by using the `use project` command or navigating through the hierarchical structure.

## Integrating with Git

CC fully integrates with Git for version control, creating and managing branches for each implementation and feature. All Git operations are handled automatically by CC, but you can also interact with the repository directly using standard Git commands:

```bash
cd ~/cc-projects/my-project
git branch             # View all branches
git log                # View commit history
git checkout <branch>  # Manually switch branches
git diff branch1 branch2  # Compare branches (similar to cc compare)
```

CC performs these Git operations for you:
- Creates new branches for each implementation when you run `cc generate`
- Switches branches when you run `cc select`
- Creates feature branches based on the selected implementation when you run `cc feature`
- Commits changes automatically for implementations and features
- Shows the current Git branch status when you run `cc status`

## Troubleshooting

If you encounter issues, check the following:

1. Make sure Docker is running if you're using Docker as the container provider
2. Verify that you have the required permissions for the project directory
3. Check if Claude Code CLI is properly installed and accessible
4. Try reloading configuration with `cc use project <your-project>` if projects don't appear correctly
5. If resources don't update properly, try exiting and restarting the shell
6. Check your current context with `pwd` in shell mode or `cc status` to ensure you're in the right context

## Running Tests

The Code Controller comes with a comprehensive test suite covering all major functionality:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./cmd/cc/test -v   # Test shell and command functionality
go test ./pkg/config -v    # Test context navigation
go test ./pkg/models -v    # Test model functionality

# Run tests with coverage report
go test ./... -cover

# Run specific tests by name
go test -v ./... -run TestShellContextNavigation
```

## Development and Extension

The hierarchical command structure and shell mode can be easily extended with new commands. The shell handles tab completion and command parsing automatically, making it easy to add new features.