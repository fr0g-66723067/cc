package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	
	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	// We need to create a new command for testing to avoid side effects
	// from the global rootCmd in the real app
	cmd := &cobra.Command{
		Use:   "cc",
		Short: "Code Controller - AI-powered code generation manager",
		Long: `Code Controller (CC) is an interface on top of Claude Code CLI
that automates code generation and manages multiple implementation versions
through git branching.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Code Controller CLI")
		},
	}
	
	// Execute the root command
	output, err := executeCommand(cmd)
	assert.NoError(t, err)
	assert.Contains(t, output, "Code Controller CLI")
}

func TestInitCommand(t *testing.T) {
	// Create a test root command
	root := &cobra.Command{Use: "cc"}
	
	// Add the init command
	init := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			cmd.Printf("Initializing project: %s\n", projectName)
		},
	}
	root.AddCommand(init)
	
	// Execute the init command with a project name
	output, err := executeCommand(root, "init", "test-project")
	assert.NoError(t, err)
	assert.Contains(t, output, "Initializing project: test-project")
	
	// Execute the init command without a project name
	_, err = executeCommand(root, "init")
	assert.Error(t, err) // Should fail due to missing argument
}

func TestGenerateCommand(t *testing.T) {
	// Create a test root command
	root := &cobra.Command{Use: "cc"}
	
	// Add the generate command
	generate := &cobra.Command{
		Use:   "generate [description]",
		Short: "Generate implementation versions from description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			cmd.Printf("Generating implementations for: %s\n", description)
			
			// Get flags
			frameworks, _ := cmd.Flags().GetStringSlice("frameworks")
			count, _ := cmd.Flags().GetInt("count")
			parallel, _ := cmd.Flags().GetBool("parallel")
			
			cmd.Printf("Frameworks: %v\n", frameworks)
			cmd.Printf("Count: %d\n", count)
			cmd.Printf("Parallel: %v\n", parallel)
		},
	}
	
	// Add flags to generate command
	generate.Flags().StringSlice("frameworks", []string{}, "Frameworks to generate (comma-separated)")
	generate.Flags().Int("count", 3, "Number of implementations to generate")
	generate.Flags().Bool("parallel", true, "Generate implementations in parallel")
	
	root.AddCommand(generate)
	
	// Execute the generate command with a description
	output, err := executeCommand(root, "generate", "Create a web app")
	assert.NoError(t, err)
	assert.Contains(t, output, "Generating implementations for: Create a web app")
	assert.Contains(t, output, "Frameworks: []")
	assert.Contains(t, output, "Count: 3")
	assert.Contains(t, output, "Parallel: true")
	
	// Execute with custom flags
	output, err = executeCommand(root, "generate", "Create a web app", 
		"--frameworks", "react,vue", "--count", "5", "--parallel=false")
	assert.NoError(t, err)
	assert.Contains(t, output, "Frameworks: [react vue]")
	assert.Contains(t, output, "Count: 5")
	assert.Contains(t, output, "Parallel: false")
}