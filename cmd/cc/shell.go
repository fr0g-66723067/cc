package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fr0g-66723067/cc/pkg/config"
)

// Shell represents the interactive shell
type Shell struct {
	cfg        *config.Config
	configPath string
}

// NewShell creates a new interactive shell
func NewShell(cfg *config.Config, configPath string) *Shell {
	return &Shell{
		cfg:        cfg,
		configPath: configPath,
	}
}

// Start starts the interactive shell
func (s *Shell) Start() {
	fmt.Println("Welcome to Code Controller Interactive Shell")
	fmt.Println("Type 'help' for a list of commands, 'exit' to quit")
	
	p := prompt.New(
		s.executor,
		s.completer,
		prompt.OptionTitle("Code Controller"),
		prompt.OptionPrefix(s.getPrompt()),
		prompt.OptionLivePrefix(s.getLivePrefix),
	)
	p.Run()
}

// getLivePrefix returns the current prompt prefix based on context
func (s *Shell) getLivePrefix() (string, bool) {
	return s.getPrompt(), true
}

// getPrompt returns the current prompt string based on context
func (s *Shell) getPrompt() string {
	path := s.cfg.GetContextPath()
	return fmt.Sprintf("cc%s > ", path)
}

// GetPromptForTest exposes getPrompt for testing
func (s *Shell) GetPromptForTest() string {
	return s.getPrompt()
}

// completer provides tab completion suggestions
func (s *Shell) completer(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest
	
	// Get the current word being typed
	word := d.GetWordBeforeCursor()
	
	// Root level commands
	rootCommands := []prompt.Suggest{
		{Text: "help", Description: "Show help"},
		{Text: "exit", Description: "Exit the shell"},
		{Text: "projects", Description: "Projects commands"},
		{Text: "cd", Description: "Change to a context level"},
		{Text: "pwd", Description: "Show current context path"},
		{Text: "use", Description: "Use a specific resource"},
		{Text: "status", Description: "Show current status"},
	}
	
	// Projects level commands
	projectCommands := []prompt.Suggest{
		{Text: "list", Description: "List all projects"},
		{Text: "create", Description: "Create a new project"},
		{Text: "remove", Description: "Remove a project"},
		{Text: "rename", Description: "Rename a project"},
	}
	
	// Implementation level commands
	implCommands := []prompt.Suggest{
		{Text: "list", Description: "List implementations"},
		{Text: "generate", Description: "Generate new implementations"},
		{Text: "select", Description: "Select an implementation"},
		{Text: "remove", Description: "Remove an implementation"},
		{Text: "rename", Description: "Rename an implementation"},
		{Text: "compare", Description: "Compare implementations"},
	}
	
	// Feature level commands
	featureCommands := []prompt.Suggest{
		{Text: "list", Description: "List features"},
		{Text: "add", Description: "Add a new feature"},
		{Text: "remove", Description: "Remove a feature"},
		{Text: "rename", Description: "Rename a feature"},
	}

	// Decide which suggestions to show based on context and input
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) <= 1 {
		// At the root level
		suggestions = rootCommands
	} else {
		switch args[0] {
		case "projects":
			suggestions = projectCommands
			
			// If we have "projects list" already, suggest project names
			if len(args) >= 2 && args[1] == "list" {
				// Empty completion for list
			} else if len(args) >= 2 && (args[1] == "remove" || args[1] == "rename" || args[1] == "cd") {
				// Show project names for removal/rename
				for name := range s.cfg.Projects {
					suggestions = append(suggestions, prompt.Suggest{Text: name, Description: "Project"})
				}
			}
		case "implementations":
			suggestions = implCommands
			
			// If we have context for a project
			if s.cfg.Context.Level == "project" || s.cfg.Context.Level == "implementation" {
				project := s.cfg.GetProject(s.cfg.Context.ProjectName)
				if project != nil && len(args) >= 2 && (args[1] == "remove" || args[1] == "rename" || args[1] == "select") {
					// Show implementation names
					for _, impl := range project.Implementations {
						suggestions = append(suggestions, prompt.Suggest{Text: impl.BranchName, Description: impl.Framework})
					}
				}
			}
		case "features":
			suggestions = featureCommands
			
			// If we have context for an implementation
			if s.cfg.Context.Level == "implementation" || s.cfg.Context.Level == "feature" {
				project := s.cfg.GetProject(s.cfg.Context.ProjectName)
				if project != nil {
					impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
					if impl != nil && len(args) >= 2 && (args[1] == "remove" || args[1] == "rename") {
						// Show feature names
						for _, feature := range impl.Features {
							suggestions = append(suggestions, prompt.Suggest{Text: feature.BranchName, Description: feature.Description})
						}
					}
				}
			}
		case "cd":
			// Suggest context navigation paths
			suggestions = []prompt.Suggest{
				{Text: "/", Description: "Root level"},
				{Text: "projects", Description: "Projects level"},
			}
			
			// Add project paths if at root
			if s.cfg.Context.Level == "root" {
				for name := range s.cfg.Projects {
					suggestions = append(suggestions, prompt.Suggest{Text: "projects/" + name, Description: "Project"})
				}
			}
			
			// Add implementation paths if at project level
			if s.cfg.Context.Level == "project" {
				project := s.cfg.GetProject(s.cfg.Context.ProjectName)
				if project != nil {
					for _, impl := range project.Implementations {
						suggestions = append(suggestions, prompt.Suggest{Text: "implementations/" + impl.BranchName, Description: impl.Framework})
					}
				}
			}
			
			// Add feature paths if at implementation level
			if s.cfg.Context.Level == "implementation" {
				project := s.cfg.GetProject(s.cfg.Context.ProjectName)
				if project != nil {
					impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
					if impl != nil {
						for _, feature := range impl.Features {
							suggestions = append(suggestions, prompt.Suggest{Text: "features/" + feature.BranchName, Description: feature.Description})
						}
					}
				}
			}
		case "use":
			// Suggest resources to use
			suggestions = []prompt.Suggest{
				{Text: "project", Description: "Use a project"},
				{Text: "implementation", Description: "Use an implementation"},
				{Text: "feature", Description: "Use a feature"},
			}
			
			if len(args) >= 2 {
				switch args[1] {
				case "project":
					// Show project names
					for name := range s.cfg.Projects {
						suggestions = append(suggestions, prompt.Suggest{Text: name, Description: "Project"})
					}
				case "implementation":
					// Show implementation names if we have a project context
					if s.cfg.Context.Level != "root" {
						project := s.cfg.GetProject(s.cfg.Context.ProjectName)
						if project != nil {
							for _, impl := range project.Implementations {
								suggestions = append(suggestions, prompt.Suggest{Text: impl.BranchName, Description: impl.Framework})
							}
						}
					}
				case "feature":
					// Show feature names if we have an implementation context
					if s.cfg.Context.Level == "implementation" || s.cfg.Context.Level == "feature" {
						project := s.cfg.GetProject(s.cfg.Context.ProjectName)
						if project != nil {
							impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
							if impl != nil {
								for _, feature := range impl.Features {
									suggestions = append(suggestions, prompt.Suggest{Text: feature.BranchName, Description: feature.Description})
								}
							}
						}
					}
				}
			}
		}
	}
	
	// Filter suggestions based on current word
	if word == "" {
		return suggestions
	}
	
	filtered := []prompt.Suggest{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Text), strings.ToLower(word)) {
			filtered = append(filtered, s)
		}
	}
	
	return filtered
}

// executor processes the entered command
func (s *Shell) executor(cmd string) {
	cmd = strings.TrimSpace(cmd)
	
	if cmd == "" {
		return
	}
	
	if cmd == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
		return
	}
	
	args := strings.Split(cmd, " ")
	
	switch args[0] {
	case "help":
		s.showHelp()
	case "pwd":
		fmt.Println(s.cfg.GetContextPath())
	case "status":
		s.showStatus()
	case "cd":
		s.changeDirectory(args)
	case "use":
		s.useResource(args)
	case "projects":
		s.handleProjectsCommands(args)
	case "implementations":
		s.handleImplementationsCommands(args)
	case "features":
		s.handleFeaturesCommands(args)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Type 'help' for a list of commands")
	}
}

// showHelp displays available commands
func (s *Shell) showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                    - Show this help message")
	fmt.Println("  exit                    - Exit the shell")
	fmt.Println("  pwd                     - Show current context path")
	fmt.Println("  cd <path>               - Change context level")
	fmt.Println("  status                  - Show current status")
	fmt.Println("  use <resource> <name>   - Use a specific resource")
	fmt.Println()
	fmt.Println("Projects commands:")
	fmt.Println("  projects list                      - List all projects")
	fmt.Println("  projects create <name> [<desc>]    - Create a new project")
	fmt.Println("  projects remove <name>             - Remove a project")
	fmt.Println("  projects rename <old> <new>        - Rename a project")
	fmt.Println()
	fmt.Println("Implementations commands:")
	fmt.Println("  implementations list                         - List implementations")
	fmt.Println("  implementations generate <desc> [--frameworks] - Generate implementations")
	fmt.Println("  implementations select <branch>              - Select an implementation")
	fmt.Println("  implementations remove <branch>              - Remove an implementation")
	fmt.Println("  implementations rename <old> <new>           - Rename an implementation")
	fmt.Println("  implementations compare <branch1> <branch2>  - Compare implementations")
	fmt.Println()
	fmt.Println("Features commands:")
	fmt.Println("  features list                - List features")
	fmt.Println("  features add <description>   - Add a new feature")
	fmt.Println("  features remove <branch>     - Remove a feature")
	fmt.Println("  features rename <old> <new>  - Rename a feature")
}

// showStatus displays the current context status
func (s *Shell) showStatus() {
	status, err := executeStatusCommand(s.configPath)
	if err != nil {
		fmt.Printf("Error getting status: %s\n", err)
		return
	}
	fmt.Println(status)
}

// ChangeDirectoryForTest exposes changeDirectory for testing
func (s *Shell) ChangeDirectoryForTest(args []string) {
	s.changeDirectory(args)
}

// changeDirectory changes the context level
func (s *Shell) changeDirectory(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cd <path>")
		return
	}
	
	path := args[1]
	
	// Handle root path
	if path == "/" {
		s.cfg.SetContextLevel("root")
		s.cfg.Context.ProjectName = ""
		s.cfg.Context.ImplementationBranch = ""
		s.cfg.Context.FeatureBranch = ""
		config.SaveConfig(s.cfg, s.configPath)
		return
	}
	
	// Handle relative or absolute paths
	segments := strings.Split(path, "/")
	
	// Remove empty segments
	var cleanSegments []string
	for _, seg := range segments {
		if seg != "" {
			cleanSegments = append(cleanSegments, seg)
		}
	}
	
	// Handle absolute paths (starts with /)
	if strings.HasPrefix(path, "/") {
		// Reset to root first
		s.cfg.SetContextLevel("root")
		s.cfg.Context.ProjectName = ""
		s.cfg.Context.ImplementationBranch = ""
		s.cfg.Context.FeatureBranch = ""
	}
	
	// Parse path segments
	for i, segment := range cleanSegments {
		switch i {
		case 0:
			if segment == "projects" {
				s.cfg.SetContextLevel("root") // Just at /projects, not in a specific project yet
			} else {
				fmt.Printf("Invalid path segment: %s\n", segment)
				return
			}
		case 1:
			// Project name
			project := s.cfg.GetProject(segment)
			if project == nil {
				fmt.Printf("Project not found: %s\n", segment)
				return
			}
			s.cfg.SetContextProject(segment)
		case 2:
			if segment != "implementations" {
				fmt.Printf("Invalid path segment: %s\n", segment)
				return
			}
		case 3:
			// Implementation branch
			project := s.cfg.GetProject(s.cfg.Context.ProjectName)
			if project == nil {
				fmt.Printf("Project context lost\n")
				return
			}
			
			impl := project.GetImplementation(segment)
			if impl == nil {
				fmt.Printf("Implementation not found: %s\n", segment)
				return
			}
			
			s.cfg.SetContextImplementation(segment)
		case 4:
			if segment != "features" {
				fmt.Printf("Invalid path segment: %s\n", segment)
				return
			}
		case 5:
			// Feature branch
			project := s.cfg.GetProject(s.cfg.Context.ProjectName)
			if project == nil {
				fmt.Printf("Project context lost\n")
				return
			}
			
			impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
			if impl == nil {
				fmt.Printf("Implementation context lost\n")
				return
			}
			
			// Find feature
			featureFound := false
			for _, feature := range impl.Features {
				if feature.BranchName == segment {
					featureFound = true
					break
				}
			}
			
			if !featureFound {
				fmt.Printf("Feature not found: %s\n", segment)
				return
			}
			
			s.cfg.SetContextFeature(segment)
		}
	}
	
	// Save updated context
	config.SaveConfig(s.cfg, s.configPath)
}

// UseResourceForTest exposes useResource for testing
func (s *Shell) UseResourceForTest(args []string) {
	s.useResource(args)
}

// useResource sets a resource as the current context
func (s *Shell) useResource(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: use <resource_type> <name>")
		fmt.Println("Resource types: project, implementation, feature")
		return
	}
	
	resourceType := args[1]
	name := args[2]
	
	switch resourceType {
	case "project":
		// Reload config first to make sure we have the latest info
		updatedConfig, err := config.LoadConfig(s.configPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			return
		}
		s.cfg = updatedConfig
		
		project := s.cfg.GetProject(name)
		if project == nil {
			fmt.Printf("Project not found: %s\n", name)
			return
		}
		s.cfg.SetContextProject(name)
		s.cfg.SetActiveProject(name)
		fmt.Printf("Now using project: %s\n", name)
		
	case "implementation":
		// Must be in a project context
		if s.cfg.Context.Level == "root" {
			fmt.Println("You must select a project first")
			return
		}
		
		project := s.cfg.GetProject(s.cfg.Context.ProjectName)
		if project == nil {
			fmt.Println("Project context is invalid")
			return
		}
		
		impl := project.GetImplementation(name)
		if impl == nil {
			fmt.Printf("Implementation not found: %s\n", name)
			return
		}
		
		s.cfg.SetContextImplementation(name)
		project.SetSelectedImplementation(name)
		fmt.Printf("Now using implementation: %s\n", name)
		
	case "feature":
		// Must be in an implementation context
		if s.cfg.Context.Level == "root" || s.cfg.Context.Level == "project" {
			fmt.Println("You must select an implementation first")
			return
		}
		
		project := s.cfg.GetProject(s.cfg.Context.ProjectName)
		if project == nil {
			fmt.Println("Project context is invalid")
			return
		}
		
		impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
		if impl == nil {
			fmt.Println("Implementation context is invalid")
			return
		}
		
		// Find feature
		featureFound := false
		for _, feature := range impl.Features {
			if feature.BranchName == name {
				featureFound = true
				break
			}
		}
		
		if !featureFound {
			fmt.Printf("Feature not found: %s\n", name)
			return
		}
		
		s.cfg.SetContextFeature(name)
		fmt.Printf("Now using feature: %s\n", name)
		
	default:
		fmt.Printf("Unknown resource type: %s\n", resourceType)
		fmt.Println("Valid types: project, implementation, feature")
	}
	
	// Save updated context
	config.SaveConfig(s.cfg, s.configPath)
}

// HandleProjectsCommandsForTest exposes handleProjectsCommands for testing
func (s *Shell) HandleProjectsCommandsForTest(args []string) bool {
	if len(args) < 2 {
		return false
	}
	s.handleProjectsCommands(args)
	return true
}

// handleProjectsCommands processes projects-related commands
func (s *Shell) handleProjectsCommands(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: projects <command> [args...]")
		fmt.Println("Commands: list, create, remove, rename")
		return
	}
	
	command := args[1]
	
	switch command {
	case "list":
		// List all projects
		items, err := executeListCommand(s.configPath, "projects")
		if err != nil {
			fmt.Printf("Error listing projects: %s\n", err)
			return
		}
		
		if len(items) == 0 {
			fmt.Println("No projects found")
			return
		}
		
		fmt.Println("Projects:")
		for i, item := range items {
			fmt.Printf("  %d. %s\n", i+1, item)
		}
		
	case "create":
		if len(args) < 3 {
			fmt.Println("Usage: projects create <name> [description]")
			return
		}
		
		name := args[2]
		description := name
		if len(args) > 3 {
			description = args[3]
		}
		
		err := executeInitCommand(s.configPath, name, description)
		if err != nil {
			fmt.Printf("Error creating project: %s\n", err)
			return
		}
		
		// Reload config to get the new project
		updatedConfig, err := config.LoadConfig(s.configPath)
		if err != nil {
			fmt.Printf("Error reloading config: %s\n", err)
			return
		}
		
		// Update the shell's config object
		s.cfg = updatedConfig
		
		// Automatically switch to the newly created project
		s.cfg.SetContextProject(name)
		s.cfg.SetActiveProject(name)
		
		// Save the updated context
		config.SaveConfig(s.cfg, s.configPath)
		
		fmt.Printf("Project %s created successfully\n", name)
		fmt.Printf("Now using project: %s\n", name)
		
	case "remove":
		if len(args) < 3 {
			fmt.Println("Usage: projects remove <name>")
			return
		}
		
		name := args[2]
		
		err := executeRemoveCommand(s.configPath, "project", name)
		if err != nil {
			fmt.Printf("Error removing project: %s\n", err)
			return
		}
		
		// Update context if we removed the current project
		if s.cfg.Context.ProjectName == name {
			s.cfg.SetContextLevel("root")
		}
		
		fmt.Printf("Project %s removed successfully\n", name)
		
	case "rename":
		if len(args) < 4 {
			fmt.Println("Usage: projects rename <old-name> <new-name>")
			return
		}
		
		oldName := args[2]
		newName := args[3]
		
		err := executeRenameCommand(s.configPath, "project", oldName, newName)
		if err != nil {
			fmt.Printf("Error renaming project: %s\n", err)
			return
		}
		
		// Update context if we renamed the current project
		if s.cfg.Context.ProjectName == oldName {
			s.cfg.SetContextProject(newName)
		}
		
		fmt.Printf("Project renamed from %s to %s successfully\n", oldName, newName)
		
	default:
		fmt.Printf("Unknown projects command: %s\n", command)
		fmt.Println("Valid commands: list, create, remove, rename")
	}
}

// HandleImplementationsCommandsForTest exposes handleImplementationsCommands for testing
func (s *Shell) HandleImplementationsCommandsForTest(args []string) bool {
	// Must be in a project context
	if s.cfg.Context.Level == "root" {
		fmt.Println("You must select a project first")
		fmt.Println("Use 'use project <n>' or 'cd /projects/<n>'")
		return false
	}
	
	if len(args) < 2 {
		return false
	}
	
	s.handleImplementationsCommands(args)
	return true
}

// handleImplementationsCommands processes implementations-related commands
func (s *Shell) handleImplementationsCommands(args []string) {
	// Must be in a project context
	if s.cfg.Context.Level == "root" {
		fmt.Println("You must select a project first")
		fmt.Println("Use 'use project <name>' or 'cd /projects/<name>'")
		return
	}
	
	project := s.cfg.GetProject(s.cfg.Context.ProjectName)
	if project == nil {
		fmt.Println("Project context is invalid")
		return
	}
	
	if len(args) < 2 {
		fmt.Println("Usage: implementations <command> [args...]")
		fmt.Println("Commands: list, generate, select, remove, rename, compare")
		return
	}
	
	command := args[1]
	
	switch command {
	case "list":
		// List implementations for current project
		items, err := executeListCommand(s.configPath, "implementations")
		if err != nil {
			fmt.Printf("Error listing implementations: %s\n", err)
			return
		}
		
		if len(items) == 0 {
			fmt.Println("No implementations found")
			return
		}
		
		fmt.Println("Implementations:")
		for i, item := range items {
			impl := project.GetImplementation(item)
			if impl != nil {
				fmt.Printf("  %d. %s (%s)\n", i+1, item, impl.Framework)
			} else {
				fmt.Printf("  %d. %s\n", i+1, item)
			}
		}
		
	case "generate":
		if len(args) < 3 {
			fmt.Println("Usage: implementations generate <description> [--frameworks=react,vue] [--count=3] [--parallel=true]")
			return
		}
		
		description := args[2]
		
		// Parse flags
		frameworks := []string{}
		count := 3
		parallel := true
		
		for i := 3; i < len(args); i++ {
			if strings.HasPrefix(args[i], "--frameworks=") {
				fwList := strings.TrimPrefix(args[i], "--frameworks=")
				frameworks = strings.Split(fwList, ",")
			} else if strings.HasPrefix(args[i], "--count=") {
				countStr := strings.TrimPrefix(args[i], "--count=")
				fmt.Sscanf(countStr, "%d", &count)
			} else if strings.HasPrefix(args[i], "--parallel=") {
				parallelStr := strings.TrimPrefix(args[i], "--parallel=")
				parallel = parallelStr == "true"
			}
		}
		
		err := executeGenerateCommand(s.configPath, description, frameworks, count, parallel)
		if err != nil {
			fmt.Printf("Error generating implementations: %s\n", err)
			return
		}
		
		fmt.Println("Implementations generated successfully")
		
	case "select":
		if len(args) < 3 {
			fmt.Println("Usage: implementations select <branch>")
			return
		}
		
		branch := args[2]
		
		err := executeSelectCommand(s.configPath, branch)
		if err != nil {
			fmt.Printf("Error selecting implementation: %s\n", err)
			return
		}
		
		// Update context to the selected implementation
		s.cfg.SetContextImplementation(branch)
		
		fmt.Printf("Implementation %s selected successfully\n", branch)
		
	case "remove":
		if len(args) < 3 {
			fmt.Println("Usage: implementations remove <branch>")
			return
		}
		
		branch := args[2]
		
		err := executeRemoveCommand(s.configPath, "implementation", branch)
		if err != nil {
			fmt.Printf("Error removing implementation: %s\n", err)
			return
		}
		
		// Update context if we removed the current implementation
		if s.cfg.Context.ImplementationBranch == branch {
			s.cfg.SetContextLevel("project")
		}
		
		fmt.Printf("Implementation %s removed successfully\n", branch)
		
	case "rename":
		if len(args) < 4 {
			fmt.Println("Usage: implementations rename <old-branch> <new-branch>")
			return
		}
		
		oldBranch := args[2]
		newBranch := args[3]
		
		err := executeRenameCommand(s.configPath, "implementation", oldBranch, newBranch)
		if err != nil {
			fmt.Printf("Error renaming implementation: %s\n", err)
			return
		}
		
		// Update context if we renamed the current implementation
		if s.cfg.Context.ImplementationBranch == oldBranch {
			s.cfg.SetContextImplementation(newBranch)
		}
		
		fmt.Printf("Implementation renamed from %s to %s successfully\n", oldBranch, newBranch)
		
	case "compare":
		if len(args) < 4 {
			fmt.Println("Usage: implementations compare <branch1> <branch2>")
			return
		}
		
		branch1 := args[2]
		branch2 := args[3]
		
		diff, err := executeCompareCommand(s.configPath, branch1, branch2)
		if err != nil {
			fmt.Printf("Error comparing implementations: %s\n", err)
			return
		}
		
		fmt.Println("\nComparison:")
		fmt.Println(diff)
		
	default:
		fmt.Printf("Unknown implementations command: %s\n", command)
		fmt.Println("Valid commands: list, generate, select, remove, rename, compare")
	}
}

// HandleFeaturesCommandsForTest exposes handleFeaturesCommands for testing
func (s *Shell) HandleFeaturesCommandsForTest(args []string) bool {
	// Must be in an implementation context
	if s.cfg.Context.Level == "root" || s.cfg.Context.Level == "project" {
		fmt.Println("You must select an implementation first")
		fmt.Println("Use 'use implementation <branch>' or navigate to the implementation")
		return false
	}
	
	if len(args) < 2 {
		return false
	}
	
	s.handleFeaturesCommands(args)
	return true
}

// handleFeaturesCommands processes features-related commands
func (s *Shell) handleFeaturesCommands(args []string) {
	// Must be in an implementation context
	if s.cfg.Context.Level == "root" || s.cfg.Context.Level == "project" {
		fmt.Println("You must select an implementation first")
		fmt.Println("Use 'use implementation <branch>' or navigate to the implementation")
		return
	}
	
	project := s.cfg.GetProject(s.cfg.Context.ProjectName)
	if project == nil {
		fmt.Println("Project context is invalid")
		return
	}
	
	impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
	if impl == nil {
		fmt.Println("Implementation context is invalid")
		return
	}
	
	if len(args) < 2 {
		fmt.Println("Usage: features <command> [args...]")
		fmt.Println("Commands: list, add, remove, rename")
		return
	}
	
	command := args[1]
	
	switch command {
	case "list":
		// List features for current implementation
		items, err := executeListCommand(s.configPath, "features")
		if err != nil {
			fmt.Printf("Error listing features: %s\n", err)
			return
		}
		
		if len(items) == 0 {
			fmt.Println("No features found")
			return
		}
		
		fmt.Println("Features:")
		for i, item := range items {
			// Find feature description
			description := ""
			for _, feat := range impl.Features {
				if feat.BranchName == item {
					description = feat.Description
					break
				}
			}
			
			if description != "" {
				fmt.Printf("  %d. %s (%s)\n", i+1, item, description)
			} else {
				fmt.Printf("  %d. %s\n", i+1, item)
			}
		}
		
	case "add":
		if len(args) < 3 {
			fmt.Println("Usage: features add <description>")
			return
		}
		
		description := args[2]
		
		err := executeFeatureCommand(s.configPath, description)
		if err != nil {
			fmt.Printf("Error adding feature: %s\n", err)
			return
		}
		
		fmt.Println("Feature added successfully")
		
		// Reload project data to get the new feature
		s.cfg, _ = config.LoadConfig(s.configPath)
		
		// Find the new feature branch
		project := s.cfg.GetProject(s.cfg.Context.ProjectName)
		if project != nil {
			impl := project.GetImplementation(s.cfg.Context.ImplementationBranch)
			if impl != nil && len(impl.Features) > 0 {
				// Use the last feature as the new one
				newFeature := impl.Features[len(impl.Features)-1]
				s.cfg.SetContextFeature(newFeature.BranchName)
			}
		}
		
	case "remove":
		if len(args) < 3 {
			fmt.Println("Usage: features remove <branch>")
			return
		}
		
		branch := args[2]
		
		err := executeRemoveCommand(s.configPath, "feature", branch)
		if err != nil {
			fmt.Printf("Error removing feature: %s\n", err)
			return
		}
		
		// Update context if we removed the current feature
		if s.cfg.Context.FeatureBranch == branch {
			s.cfg.SetContextLevel("implementation")
		}
		
		fmt.Printf("Feature %s removed successfully\n", branch)
		
	case "rename":
		if len(args) < 4 {
			fmt.Println("Usage: features rename <old-branch> <new-branch>")
			return
		}
		
		oldBranch := args[2]
		newBranch := args[3]
		
		err := executeRenameCommand(s.configPath, "feature", oldBranch, newBranch)
		if err != nil {
			fmt.Printf("Error renaming feature: %s\n", err)
			return
		}
		
		// Update context if we renamed the current feature
		if s.cfg.Context.FeatureBranch == oldBranch {
			s.cfg.SetContextFeature(newBranch)
		}
		
		fmt.Printf("Feature renamed from %s to %s successfully\n", oldBranch, newBranch)
		
	default:
		fmt.Printf("Unknown features command: %s\n", command)
		fmt.Println("Valid commands: list, add, remove, rename")
	}
}