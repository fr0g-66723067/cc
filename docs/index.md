# Code Controller (CC) Documentation

Welcome to the Code Controller documentation. CC is an advanced interface built on top of Claude Code CLI that automates and streamlines AI-powered code generation through Git branching.

## Quick Links

### Getting Started
- [README](../README.md) - Project overview and installation
- [Usage Guide](USAGE.md) - Detailed usage instructions

### Project Information
- [Project Summary](SUMMARY.md) - Current project status and implementation details
- [Architecture](design/ARCHITECTURE.md) - System architecture and components
- [Project Structure](design/PROJECT_STRUCTURE.md) - Code organization and patterns
- [Development Roadmap](design/DETAILED_ROADMAP.md) - Future development plans

## Key Features

### Multi-version Prototyping
- Generate multiple versions of the same application using different frameworks and approaches
- Each version is maintained in its own git branch for easy comparison and switching
- Example: "Create a web app to track daily tasks" will generate versions using React, Vue, Svelte, etc.

### Project Evolution
- Select your preferred version as the base implementation
- Incrementally enhance your chosen implementation by requesting new features
- Each feature request creates a new branch, allowing for easy feature comparison and rollback

### Git Integration
CC fully integrates with Git for version control:
- Creates implementation branches when you run `cc generate`
- Switches branches when you run `cc select`
- Creates feature branches based on the selected implementation when you run `cc feature`
- Commits changes automatically for implementations and features
- Shows the current Git branch status when you run `cc status`

## Common Workflows

1. **Initialize a project**
   ```bash
   cc init my-project -d "A task management application"
   ```

2. **Generate implementations**
   ```bash
   cc generate "Create a web app with user authentication and task management" --frameworks react,vue,svelte
   ```

3. **Select an implementation**
   ```bash
   cc select impl-react-1234567890
   ```

4. **Add a feature**
   ```bash
   cc feature "Add dark mode support"
   ```

5. **Compare implementations**
   ```bash
   cc compare impl-react-1234567890 impl-vue-0987654321
   ```

6. **Show project status**
   ```bash
   cc status
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. See the README and documentation for development setup instructions.