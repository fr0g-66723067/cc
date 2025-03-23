package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/fr0g-66723067/cc/internal/ai/mocks"
	"github.com/fr0g-66723067/cc/pkg/config"
	"github.com/fr0g-66723067/cc/pkg/models"
	"github.com/fr0g-66723067/cc/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// This test simulates the complete workflow of:
// 1. Creating a project
// 2. Generating implementations
// 3. Selecting an implementation
// 4. Adding a feature
func TestProjectWorkflow(t *testing.T) {
	t.Skip("Integration test requiring mocks - skipping for now")
	
	// Create a controller for our mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Create a temporary directory for our test project
	projectDir := testutil.TempDir(t)
	
	// Create a configuration directory
	configDir := testutil.TempDir(t)
	configPath := filepath.Join(configDir, "config.json")
	
	// Load the default configuration
	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)
	
	// Set the projects directory to our temp directory
	cfg.ProjectsDir = projectDir
	
	// Create a mock AI provider
	mockAI := mocks.NewMockProvider(ctrl)
	
	// Set expectations for the AI provider
	mockAI.EXPECT().
		Name().
		Return("mock").
		AnyTimes()
	
	mockAI.EXPECT().
		Initialize(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	
	mockAI.EXPECT().
		SupportedFrameworks().
		Return([]string{"react", "vue"}).
		AnyTimes()
	
	// Expect GenerateProject to be called
	mockAI.EXPECT().
		GenerateProject(gomock.Any(), gomock.Eq("Create a todo app")).
		Return("Generated project structure", nil)
	
	// Expect GenerateImplementation to be called for each framework
	mockAI.EXPECT().
		GenerateImplementation(gomock.Any(), gomock.Eq("Create a todo app"), gomock.Eq("react")).
		Return("Generated React implementation", nil)
	
	mockAI.EXPECT().
		GenerateImplementation(gomock.Any(), gomock.Eq("Create a todo app"), gomock.Eq("vue")).
		Return("Generated Vue implementation", nil)
	
	// Expect AddFeature to be called
	mockAI.EXPECT().
		AddFeature(gomock.Any(), gomock.Any(), gomock.Eq("Add user authentication")).
		Return("Added authentication feature", nil)
	
	// Create a new project
	project := models.NewProject("todo-app", 
		filepath.Join(projectDir, "todo-app"), 
		"Create a todo app")
	
	// Add the project to the configuration
	cfg.AddProject(project)
	cfg.SetActiveProject("todo-app")
	
	// Create the project directory
	err = os.MkdirAll(project.Path, 0755)
	require.NoError(t, err)
	
	// Save the configuration
	err = config.SaveConfig(cfg, configPath)
	require.NoError(t, err)
	
	// Generate implementations
	ctx := context.Background()
	
	// Initialize the AI provider
	err = mockAI.Initialize(ctx, nil)
	require.NoError(t, err)
	
	// Generate the project structure
	_, err = mockAI.GenerateProject(ctx, project.Description)
	require.NoError(t, err)
	
	// Generate implementations for each framework
	frameworks := mockAI.SupportedFrameworks()
	for _, framework := range frameworks {
		// Generate the implementation
		_, err = mockAI.GenerateImplementation(ctx, project.Description, framework)
		require.NoError(t, err)
		
		// Create an implementation model
		impl := models.Implementation{
			Framework:   framework,
			BranchName:  "implementation/" + framework,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
			Provider:    mockAI.Name(),
		}
		
		// Add the implementation to the project
		project.AddImplementation(impl)
	}
	
	// Verify that implementations were added
	assert.Len(t, project.Implementations, 2)
	
	// Select the React implementation
	project.SetSelectedImplementation("implementation/react")
	assert.Equal(t, "implementation/react", project.SelectedImplementation)
	
	// Verify we can get the selected implementation
	selectedImpl := project.GetSelectedImplementation()
	assert.NotNil(t, selectedImpl)
	assert.Equal(t, "react", selectedImpl.Framework)
	
	// Add a feature
	_, err = mockAI.AddFeature(ctx, project.Path, "Add user authentication")
	require.NoError(t, err)
	
	// Create a feature model
	feature := models.Feature{
		Name:        "authentication",
		BranchName:  "feature/authentication",
		Description: "Add user authentication",
		CreatedAt:   project.CreatedAt,
		BaseBranch:  project.SelectedImplementation,
		Provider:    mockAI.Name(),
		Status:      "completed",
	}
	
	// Add the feature to the selected implementation
	selectedImpl.Features = append(selectedImpl.Features, feature)
	
	// Verify the feature was added
	assert.Len(t, selectedImpl.Features, 1)
	assert.Equal(t, "authentication", selectedImpl.Features[0].Name)
	
	// Update the project in the configuration
	cfg.AddProject(project)
	
	// Save the configuration again
	err = config.SaveConfig(cfg, configPath)
	require.NoError(t, err)
	
	// Load the configuration and verify everything was saved
	loadedCfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)
	
	// Get the project
	loadedProject := loadedCfg.GetProject("todo-app")
	assert.NotNil(t, loadedProject)
	
	// Verify project properties
	assert.Equal(t, "todo-app", loadedProject.Name)
	assert.Equal(t, "Create a todo app", loadedProject.Description)
	assert.Equal(t, "implementation/react", loadedProject.SelectedImplementation)
	
	// Verify implementations
	assert.Len(t, loadedProject.Implementations, 2)
	
	// Get the selected implementation
	loadedImpl := loadedProject.GetSelectedImplementation()
	assert.NotNil(t, loadedImpl)
	assert.Equal(t, "react", loadedImpl.Framework)
	
	// Verify features
	assert.Len(t, loadedImpl.Features, 1)
	assert.Equal(t, "authentication", loadedImpl.Features[0].Name)
}