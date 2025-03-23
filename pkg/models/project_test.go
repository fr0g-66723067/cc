package models_test

import (
	"testing"
	"time"

	"github.com/fr0g-66723067/cc/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	// Test creating a new project
	projectName := "test-project"
	projectPath := "/path/to/project"
	projectDesc := "Test project description"

	project := models.NewProject(projectName, projectPath, projectDesc)

	// Verify project properties
	assert.Equal(t, projectName, project.Name)
	assert.Equal(t, projectPath, project.Path)
	assert.Equal(t, projectDesc, project.Description)
	assert.Equal(t, "main", project.ActiveBranch)
	assert.Equal(t, "initialized", project.Status)
	assert.Empty(t, project.Implementations)
	assert.Empty(t, project.SelectedImplementation)

	// Verify maps are initialized
	assert.NotNil(t, project.ContainerConfig)
	assert.NotNil(t, project.AIConfig)
	assert.NotNil(t, project.VCSConfig)
	assert.NotNil(t, project.Settings)

	// Verify timestamps
	assert.WithinDuration(t, time.Now(), project.CreatedAt, 1*time.Second)
	assert.WithinDuration(t, time.Now(), project.UpdatedAt, 1*time.Second)
}

func TestAddImplementation(t *testing.T) {
	// Create a test project
	project := models.NewProject("test", "/path/to/test", "Test project")
	
	// Create a test implementation
	impl := models.Implementation{
		Framework:   "react",
		BranchName:  "implementation/react",
		Description: "React implementation",
		CreatedAt:   time.Now(),
		Provider:    "claude",
	}
	
	// Add the implementation to the project
	project.AddImplementation(impl)
	
	// Verify the implementation was added
	assert.Len(t, project.Implementations, 1)
	assert.Equal(t, impl.Framework, project.Implementations[0].Framework)
	assert.Equal(t, impl.BranchName, project.Implementations[0].BranchName)
	
	// Verify the UpdatedAt timestamp was updated
	assert.WithinDuration(t, time.Now(), project.UpdatedAt, 1*time.Second)
}

func TestGetImplementation(t *testing.T) {
	// Create a test project
	project := models.NewProject("test", "/path/to/test", "Test project")
	
	// Create and add test implementations
	impl1 := models.Implementation{
		Framework:   "react",
		BranchName:  "implementation/react",
		Description: "React implementation",
		CreatedAt:   time.Now(),
	}
	
	impl2 := models.Implementation{
		Framework:   "vue",
		BranchName:  "implementation/vue",
		Description: "Vue implementation",
		CreatedAt:   time.Now(),
	}
	
	project.AddImplementation(impl1)
	project.AddImplementation(impl2)
	
	// Test getting an existing implementation
	result := project.GetImplementation("implementation/react")
	assert.NotNil(t, result)
	assert.Equal(t, "react", result.Framework)
	
	// Test getting another existing implementation
	result = project.GetImplementation("implementation/vue")
	assert.NotNil(t, result)
	assert.Equal(t, "vue", result.Framework)
	
	// Test getting a non-existent implementation
	result = project.GetImplementation("implementation/angular")
	assert.Nil(t, result)
}

func TestSelectedImplementation(t *testing.T) {
	// Create a test project
	project := models.NewProject("test", "/path/to/test", "Test project")
	
	// Add test implementations
	project.AddImplementation(models.Implementation{
		Framework:   "react",
		BranchName:  "implementation/react",
		Description: "React implementation",
		CreatedAt:   time.Now(),
	})
	
	project.AddImplementation(models.Implementation{
		Framework:   "vue",
		BranchName:  "implementation/vue",
		Description: "Vue implementation",
		CreatedAt:   time.Now(),
	})
	
	// Initially, no implementation is selected
	assert.Empty(t, project.SelectedImplementation)
	assert.Nil(t, project.GetSelectedImplementation())
	
	// Set a selected implementation
	project.SetSelectedImplementation("implementation/react")
	assert.Equal(t, "implementation/react", project.SelectedImplementation)
	
	// Get the selected implementation
	selected := project.GetSelectedImplementation()
	assert.NotNil(t, selected)
	assert.Equal(t, "react", selected.Framework)
	
	// Change the selected implementation
	project.SetSelectedImplementation("implementation/vue")
	selected = project.GetSelectedImplementation()
	assert.NotNil(t, selected)
	assert.Equal(t, "vue", selected.Framework)
	
	// Set a non-existent implementation
	project.SetSelectedImplementation("implementation/angular")
	selected = project.GetSelectedImplementation()
	assert.Nil(t, selected)
}