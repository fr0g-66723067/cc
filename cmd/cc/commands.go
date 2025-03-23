package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fr0g-66723067/cc/internal/ai"
	"github.com/fr0g-66723067/cc/internal/vcs"
	"github.com/fr0g-66723067/cc/pkg/config"
	"github.com/fr0g-66723067/cc/pkg/models"
)

// executeInitCommand initializes a new project
func executeInitCommand(configPath, projectName, description string) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if project already exists
	if cfg.GetProject(projectName) != nil {
		return fmt.Errorf("project %s already exists", projectName)
	}

	// Create project directory
	projectDir := filepath.Join(cfg.ProjectsDir, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize Git repository
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return fmt.Errorf("failed to create VCS provider: %w", err)
	}

	if err := vcsProvider.Initialize(projectDir); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Create project model
	project := models.NewProject(projectName, projectDir, description)

	// Set specific configurations if needed
	project.ContainerConfig = cfg.Container.Config
	project.AIConfig = cfg.AI.Config
	project.VCSConfig = cfg.VCS.Config

	// Create initial README file
	readmePath := filepath.Join(projectDir, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\n%s\n", projectName, description)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README file: %w", err)
	}

	// Add and commit the README file
	if err := vcsProvider.AddFiles([]string{readmePath}); err != nil {
		return fmt.Errorf("failed to add README file: %w", err)
	}

	if err := vcsProvider.CommitChanges("Initial commit"); err != nil {
		return fmt.Errorf("failed to commit initial changes: %w", err)
	}

	// Add project to config
	cfg.AddProject(project)
	cfg.SetActiveProject(projectName)

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// executeGenerateCommand generates implementations for a project
func executeGenerateCommand(configPath, description string, frameworks []string, count int, parallel bool) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get active project
	project := cfg.GetActiveProject()
	if project == nil {
		return fmt.Errorf("no active project")
	}

	// Create AI provider
	aiProvider, err := ai.Create(cfg.AI.Provider, cfg.AI.Config)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Initialize AI provider
	ctx := getContext()
	if err := aiProvider.Initialize(ctx, nil); err != nil {
		return fmt.Errorf("failed to initialize AI provider: %w", err)
	}
	defer aiProvider.Cleanup(ctx)

	// Create VCS provider
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return fmt.Errorf("failed to create VCS provider: %w", err)
	}

	// Initialize VCS
	if err := vcsProvider.Initialize(project.Path); err != nil {
		return fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// If no frameworks specified, use supported frameworks from AI provider
	if len(frameworks) == 0 {
		// Limit to the requested count
		allFrameworks := aiProvider.SupportedFrameworks()
		if count < len(allFrameworks) {
			frameworks = allFrameworks[:count]
		} else {
			frameworks = allFrameworks
		}
	}

	// Generate implementations (up to the count)
	if count > len(frameworks) {
		count = len(frameworks)
	}

	// Get current branch before we start creating new branches
	currentBranch, err := vcsProvider.GetCurrentBranch()
	if err != nil {
		// If we can't get the current branch, it might be because the repository
		// doesn't have any commits yet, so let's create an initial commit
		fmt.Printf("No current branch found, initializing with an initial commit\n")

		// Create a dummy file
		dummyPath := filepath.Join(project.Path, "init.txt")
		if err := os.WriteFile(dummyPath, []byte("Initial commit"), 0644); err != nil {
			return fmt.Errorf("failed to create init file: %w", err)
		}

		// Add and commit the file
		if err := vcsProvider.AddFiles([]string{dummyPath}); err != nil {
			return fmt.Errorf("failed to add init file: %w", err)
		}

		if err := vcsProvider.CommitChanges("Initial commit for tests"); err != nil {
			return fmt.Errorf("failed to make initial commit: %w", err)
		}

		// Now try to get the current branch again
		currentBranch, err = vcsProvider.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch after init: %w", err)
		}
	}
	
	fmt.Printf("Current branch is: %s\n", currentBranch)

	for i := 0; i < count; i++ {
		framework := frameworks[i]
		branchName := fmt.Sprintf("impl-%s-%d", framework, time.Now().Unix())

		// Create a new branch for this implementation
		fmt.Printf("Creating branch %s...\n", branchName)
		if err := vcsProvider.CreateBranch(branchName, currentBranch); err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}

		// Switch to the new branch
		if err := vcsProvider.SwitchBranch(branchName); err != nil {
			return fmt.Errorf("failed to switch to branch: %w", err)
		}

		// Generate code
		combinedDesc := fmt.Sprintf("%s using %s", description, framework)
		
		// Notify user
		fmt.Printf("Generating implementation for %s... This may take a while.\n", framework)
		
		// Generate code using Claude AI provider
		_, err := aiProvider.GenerateImplementation(ctx, description, framework)
		if err != nil {
			fmt.Printf("Warning: AI code generation failed: %v\n", err)
			fmt.Printf("Creating a placeholder implementation instead...\n")
			
			// Create a fallback file if generation fails
			readmePath := filepath.Join(project.Path, "README.md")
			content := fmt.Sprintf("# %s\n\n%s\n\nFramework: %s\n", project.Name, description, framework)
			if writeErr := os.WriteFile(readmePath, []byte(content), 0644); writeErr != nil {
				return fmt.Errorf("failed to write README.md: %w", writeErr)
			}
		} else {
			fmt.Printf("Successfully generated code for %s implementation.\n", framework)
			// Files have been generated in the project directory by the AI provider
		}
		
		// Add all changes and commit
		// Use the Git command to add all files in the project directory
		allFiles := []string{project.Path}
		if err := vcsProvider.AddFiles(allFiles); err != nil {
			return fmt.Errorf("failed to add files: %w", err)
		}
		
		commitMsg := fmt.Sprintf("Implementation: %s using %s", project.Name, framework)
		if err := vcsProvider.CommitChanges(commitMsg); err != nil {
			return fmt.Errorf("failed to commit changes: %w", err)
		}

		// Create implementation model
		impl := models.Implementation{
			Framework:   framework,
			BranchName:  branchName,
			Description: combinedDesc,
			CreatedAt:   time.Now(),
			Provider:    aiProvider.Name(),
			Tags:        []string{framework},
			Metrics:     make(map[string]float64),
			Score:       50, // Default score
			Features:    []models.Feature{},
		}

		// Add implementation to project
		project.AddImplementation(impl)
	}

	// Switch back to the original branch
	if err := vcsProvider.SwitchBranch(currentBranch); err != nil {
		return fmt.Errorf("failed to switch back to original branch: %w", err)
	}

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// executeSelectCommand selects an implementation
func executeSelectCommand(configPath, branchName string) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get active project
	project := cfg.GetActiveProject()
	if project == nil {
		return fmt.Errorf("no active project")
	}

	// Check if implementation exists
	impl := project.GetImplementation(branchName)
	if impl == nil {
		return fmt.Errorf("implementation %s not found", branchName)
	}

	// Create VCS provider
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return fmt.Errorf("failed to create VCS provider: %w", err)
	}

	// Initialize VCS
	if err := vcsProvider.Initialize(project.Path); err != nil {
		return fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Check if branch exists in git
	branches, err := vcsProvider.ListBranches()
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}
	
	branchExists := false
	for _, branch := range branches {
		if branch == branchName {
			branchExists = true
			break
		}
	}
	
	if !branchExists {
		return fmt.Errorf("branch %s does not exist in the repository", branchName)
	}

	// Switch to the branch
	if err := vcsProvider.SwitchBranch(branchName); err != nil {
		return fmt.Errorf("failed to switch to branch %s: %w", branchName, err)
	}

	// Update project model
	project.SetSelectedImplementation(branchName)
	project.ActiveBranch = branchName

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// executeFeatureCommand adds a feature to the current implementation
func executeFeatureCommand(configPath, description string) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get active project
	project := cfg.GetActiveProject()
	if project == nil {
		return fmt.Errorf("no active project")
	}

	// Check if an implementation is selected
	selectedImpl := project.GetSelectedImplementation()
	if selectedImpl == nil {
		return fmt.Errorf("no implementation selected")
	}

	// Create feature name and branch
	featureName := sanitizeForBranchName(description)
	featureBranch := fmt.Sprintf("feat-%s-%d", featureName, time.Now().Unix())

	// Create AI provider
	aiProvider, err := ai.Create(cfg.AI.Provider, cfg.AI.Config)
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Initialize AI provider
	ctx := getContext()
	if err := aiProvider.Initialize(ctx, nil); err != nil {
		return fmt.Errorf("failed to initialize AI provider: %w", err)
	}
	defer aiProvider.Cleanup(ctx)

	// Create VCS provider
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return fmt.Errorf("failed to create VCS provider: %w", err)
	}

	// Initialize VCS
	if err := vcsProvider.Initialize(project.Path); err != nil {
		return fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Get current branch to verify we're on the correct implementation branch
	currentBranch, err := vcsProvider.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Ensure we're on the selected implementation branch
	if currentBranch != selectedImpl.BranchName {
		return fmt.Errorf("current branch (%s) is not the selected implementation branch (%s), please run 'select' first", 
			currentBranch, selectedImpl.BranchName)
	}

	// Create a new branch for this feature based on the implementation branch
	fmt.Printf("Creating feature branch %s...\n", featureBranch)
	if err := vcsProvider.CreateBranch(featureBranch, selectedImpl.BranchName); err != nil {
		return fmt.Errorf("failed to create feature branch: %w", err)
	}

	// Switch to the feature branch
	if err := vcsProvider.SwitchBranch(featureBranch); err != nil {
		return fmt.Errorf("failed to switch to feature branch: %w", err)
	}

	// Use Claude AI provider to add the feature
	fmt.Printf("Adding feature: %s\n", description)
	output, err := aiProvider.AddFeature(ctx, project.Path, description)
	if err != nil {
		fmt.Printf("Warning: AI feature generation failed: %v\n", err)
		fmt.Printf("Creating a placeholder feature instead...\n")
		
		// Create a fallback feature file if AI fails
		featureFile := filepath.Join(project.Path, fmt.Sprintf("feature-%s.txt", featureName))
		content := fmt.Sprintf("Feature: %s\nImplementation: %s\n", description, selectedImpl.Framework)
		if writeErr := os.WriteFile(featureFile, []byte(content), 0644); writeErr != nil {
			return fmt.Errorf("failed to write feature file: %w", writeErr)
		}
	} else {
		fmt.Printf("Successfully added feature: %s\n", description)
		fmt.Printf("AI Output Summary: %s\n", truncateString(output, 200))
	}

	// Add and commit the changes
	// Use the Git command to add all files in the project directory
	allFiles := []string{project.Path}
	if err := vcsProvider.AddFiles(allFiles); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	commitMsg := fmt.Sprintf("Feature: %s", description)
	if err := vcsProvider.CommitChanges(commitMsg); err != nil {
		return fmt.Errorf("failed to commit feature changes: %w", err)
	}

	// Create feature model
	feature := models.Feature{
		Name:        featureName,
		BranchName:  featureBranch,
		Description: description,
		CreatedAt:   time.Now(),
		BaseBranch:  selectedImpl.BranchName,
		Provider:    aiProvider.Name(),
		Status:      "completed",
		Tags:        []string{},
	}

	// Add feature to implementation
	impl := project.GetImplementation(selectedImpl.BranchName)
	impl.Features = append(impl.Features, feature)

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Switch back to the implementation branch
	if err := vcsProvider.SwitchBranch(selectedImpl.BranchName); err != nil {
		fmt.Printf("Warning: Failed to switch back to implementation branch (%s): %v\n", selectedImpl.BranchName, err)
		fmt.Printf("You are currently on the feature branch: %s\n", featureBranch)
	}

	return nil
}

// executeListCommand lists projects, implementations, or features
func executeListCommand(configPath, resourceType string) ([]string, error) {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Handle different resource types
	switch resourceType {
	case "projects":
		// List projects
		var projectNames []string
		for name := range cfg.Projects {
			projectNames = append(projectNames, name)
		}
		return projectNames, nil

	case "implementations":
		// List implementations for active project
		project := cfg.GetActiveProject()
		if project == nil {
			return nil, fmt.Errorf("no active project")
		}

		var implNames []string
		for _, impl := range project.Implementations {
			implNames = append(implNames, impl.BranchName)
		}
		return implNames, nil

	case "features":
		// List features for selected implementation
		project := cfg.GetActiveProject()
		if project == nil {
			return nil, fmt.Errorf("no active project")
		}

		selectedImpl := project.GetSelectedImplementation()
		if selectedImpl == nil {
			return nil, fmt.Errorf("no implementation selected")
		}

		var featureNames []string
		for _, feature := range selectedImpl.Features {
			featureNames = append(featureNames, feature.BranchName)
		}
		return featureNames, nil

	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

// executeStatusCommand shows the current project status
func executeStatusCommand(configPath string) (string, error) {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Get active project
	project := cfg.GetActiveProject()
	if project == nil {
		return "", fmt.Errorf("no active project")
	}

	// Create VCS provider to get real git information
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return "", fmt.Errorf("failed to create VCS provider: %w", err)
	}

	// Initialize VCS
	if err := vcsProvider.Initialize(project.Path); err != nil {
		return "", fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Get current git branch
	currentBranch, err := vcsProvider.GetCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	// Check for uncommitted changes
	hasChanges, err := vcsProvider.HasChanges()
	if err != nil {
		return "", fmt.Errorf("failed to check for changes: %w", err)
	}

	// Build status string
	var status strings.Builder
	status.WriteString(fmt.Sprintf("Project: %s\n", project.Name))
	status.WriteString(fmt.Sprintf("Description: %s\n", project.Description))
	status.WriteString(fmt.Sprintf("Path: %s\n", project.Path))
	status.WriteString(fmt.Sprintf("Current Git Branch: %s\n", currentBranch))
	status.WriteString(fmt.Sprintf("Selected Implementation: %s\n", project.SelectedImplementation))
	status.WriteString(fmt.Sprintf("Uncommitted Changes: %t\n", hasChanges))
	status.WriteString(fmt.Sprintf("Status: %s\n", project.Status))
	status.WriteString(fmt.Sprintf("Created: %s\n", project.CreatedAt.Format(time.RFC3339)))
	status.WriteString(fmt.Sprintf("Updated: %s\n", project.UpdatedAt.Format(time.RFC3339)))
	
	// Add implementations information
	status.WriteString(fmt.Sprintf("\nImplementations (%d):\n", len(project.Implementations)))
	for i, impl := range project.Implementations {
		status.WriteString(fmt.Sprintf("  %d. %s (%s)\n", i+1, impl.Framework, impl.BranchName))
		
		// Mark selected implementation
		if project.SelectedImplementation == impl.BranchName {
			status.WriteString("     * SELECTED *\n")
		}
		
		// Mark current branch
		if currentBranch == impl.BranchName {
			status.WriteString("     * CURRENT BRANCH *\n")
		}
		
		// List features for this implementation
		if len(impl.Features) > 0 {
			status.WriteString("     Features:\n")
			for j, feature := range impl.Features {
				status.WriteString(fmt.Sprintf("       %d. %s (%s)\n", j+1, feature.Description, feature.BranchName))
				
				// Mark current branch if we're on a feature branch
				if currentBranch == feature.BranchName {
					status.WriteString("         * CURRENT BRANCH *\n")
				}
			}
		}
	}

	return status.String(), nil
}

// executeCompareCommand compares two branches
func executeCompareCommand(configPath, branch1, branch2 string) (string, error) {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Get active project
	project := cfg.GetActiveProject()
	if project == nil {
		return "", fmt.Errorf("no active project")
	}

	// Create VCS provider
	vcsProvider, err := vcs.Create(cfg.VCS.Provider, cfg.VCS.Config)
	if err != nil {
		return "", fmt.Errorf("failed to create VCS provider: %w", err)
	}

	// Initialize VCS
	if err := vcsProvider.Initialize(project.Path); err != nil {
		return "", fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Check if branches exist
	branches, err := vcsProvider.ListBranches()
	if err != nil {
		return "", fmt.Errorf("failed to list branches: %w", err)
	}

	// Verify branch1 exists
	branch1Exists := false
	for _, branch := range branches {
		if branch == branch1 {
			branch1Exists = true
			break
		}
	}
	if !branch1Exists {
		return "", fmt.Errorf("branch %s does not exist", branch1)
	}

	// Verify branch2 exists
	branch2Exists := false
	for _, branch := range branches {
		if branch == branch2 {
			branch2Exists = true
			break
		}
	}
	if !branch2Exists {
		return "", fmt.Errorf("branch %s does not exist", branch2)
	}

	// Generate the diff using the VCS provider
	diff, err := vcsProvider.ExportDiff(branch1, branch2)
	if err != nil {
		return "", fmt.Errorf("failed to export diff: %w", err)
	}

	// Add header to the diff output
	header := fmt.Sprintf("Comparison between %s and %s:\n\n", branch1, branch2)
	return header + diff, nil
}

// Helper functions

// sanitizeForBranchName converts a description to a valid branch name
func sanitizeForBranchName(description string) string {
	// Convert to lowercase
	name := strings.ToLower(description)
	
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	
	// Remove special characters
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, name)
	
	// Truncate to a reasonable length
	if len(name) > 30 {
		name = name[:30]
	}
	
	return name
}

// getContext returns a context for AI operations
func getContext() context.Context {
	return context.Background()
}

// truncateString truncates a string to a maximum length and adds "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}