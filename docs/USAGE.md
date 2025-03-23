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

## Command Reference

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

## Advanced Usage

### Custom Configuration

You can specify a custom configuration file:

```bash
cc --config /path/to/config.json [command]
```

### Working with Multiple Projects

CC supports managing multiple projects. Use `cc list projects` to see all projects, and switch between them by initializing or selecting implementations in different projects.

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