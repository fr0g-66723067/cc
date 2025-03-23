package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/fr0g-66723067/cc/pkg/config"
)

var configPath string
var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "cc",
	Short: "Code Controller - AI-powered code generation manager",
	Long: `Code Controller (CC) is an interface on top of Claude Code CLI
that automates code generation and manages multiple implementation versions
through git branching.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Set up flags
	home, _ := os.UserHomeDir()
	defaultConfigPath := filepath.Join(home, ".cc", "config.json")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "config file path")

	// Initialize commands
	initCmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			
			// Get project description
			description, _ := cmd.Flags().GetString("description")
			if description == "" {
				description = projectName
			}
			
			fmt.Printf("Initializing project: %s\n", projectName)
			
			err := executeInitCommand(configPath, projectName, description)
			if err != nil {
				fmt.Printf("Error initializing project: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Project %s initialized successfully\n", projectName)
		},
	}
	
	// Add flags to init command
	initCmd.Flags().StringP("description", "d", "", "Project description")

	generateCmd := &cobra.Command{
		Use:   "generate [description]",
		Short: "Generate implementation versions from description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			
			// Get flags
			frameworks, _ := cmd.Flags().GetStringSlice("frameworks")
			count, _ := cmd.Flags().GetInt("count")
			parallel, _ := cmd.Flags().GetBool("parallel")
			
			fmt.Printf("Generating implementations for: %s\n", description)
			fmt.Printf("Frameworks: %v\n", frameworks)
			fmt.Printf("Count: %d\n", count)
			fmt.Printf("Parallel: %v\n", parallel)
			
			err := executeGenerateCommand(configPath, description, frameworks, count, parallel)
			if err != nil {
				fmt.Printf("Error generating implementations: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Println("Implementations generated successfully")
		},
	}

	// Add flags to generate command
	generateCmd.Flags().StringSlice("frameworks", []string{}, "Frameworks to generate (comma-separated)")
	generateCmd.Flags().Int("count", 3, "Number of implementations to generate")
	generateCmd.Flags().Bool("parallel", true, "Generate implementations in parallel")

	selectCmd := &cobra.Command{
		Use:   "select [branch]",
		Short: "Select an implementation branch as base",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			branch := args[0]
			fmt.Printf("Selecting implementation: %s\n", branch)
			
			err := executeSelectCommand(configPath, branch)
			if err != nil {
				fmt.Printf("Error selecting implementation: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Implementation %s selected successfully\n", branch)
		},
	}

	featureCmd := &cobra.Command{
		Use:   "feature [description]",
		Short: "Add a new feature to the current implementation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			fmt.Printf("Adding feature: %s\n", description)
			
			err := executeFeatureCommand(configPath, description)
			if err != nil {
				fmt.Printf("Error adding feature: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Println("Feature added successfully")
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [resource]",
		Short: "List resources (projects, implementations, features)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resource := args[0]
			
			if resource != "projects" && resource != "implementations" && resource != "features" {
				fmt.Printf("Unknown resource type: %s\n", resource)
				fmt.Println("Valid resources: projects, implementations, features")
				os.Exit(1)
			}
			
			fmt.Printf("Listing %s\n", resource)
			
			items, err := executeListCommand(configPath, resource)
			if err != nil {
				fmt.Printf("Error listing %s: %s\n", resource, err)
				os.Exit(1)
			}
			
			if len(items) == 0 {
				fmt.Printf("No %s found\n", resource)
				return
			}
			
			fmt.Printf("%s:\n", resource)
			for i, item := range items {
				fmt.Printf("  %d. %s\n", i+1, item)
			}
		},
	}

	compareCmd := &cobra.Command{
		Use:   "compare [branch1] [branch2]",
		Short: "Compare two implementations or features",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			branch1 := args[0]
			branch2 := args[1]
			fmt.Printf("Comparing %s and %s\n", branch1, branch2)
			
			diff, err := executeCompareCommand(configPath, branch1, branch2)
			if err != nil {
				fmt.Printf("Error comparing branches: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Println("\nDiff:")
			fmt.Println(diff)
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show the current project status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Current project status")
			
			status, err := executeStatusCommand(configPath)
			if err != nil {
				fmt.Printf("Error getting status: %s\n", err)
				os.Exit(1)
			}
			
			fmt.Println(status)
		},
	}

	// Add commands to root command
	rootCmd.AddCommand(
		initCmd,
		generateCmd,
		selectCmd,
		featureCmd,
		listCmd,
		compareCmd,
		statusCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}