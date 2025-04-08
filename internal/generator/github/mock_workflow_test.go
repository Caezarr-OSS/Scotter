package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TestMockGenerateWorkflows tests the workflow generation using mocks
func TestMockGenerateWorkflows(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create necessary subdirectories
	workflowsDir := filepath.Join(tempDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("Failed to create workflows dir: %v", err)
	}

	// Test cases
	testCases := []struct {
		name            string
		features        []string
		expectChangelog bool
	}{
		{
			name:            "With changelog feature",
			features:        []string{"changelog"},
			expectChangelog: true,
		},
		{
			name:            "Without changelog feature",
			features:        []string{"ci"},
			expectChangelog: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean the workflows directory
			files, _ := os.ReadDir(workflowsDir)
			for _, f := range files {
				os.Remove(filepath.Join(workflowsDir, f.Name()))
			}

			// Create a test config with the specified features
			cfg := &model.Config{
				ProjectName: "test-project",
				Pipeline: model.PipelineConfig{
					UseGitHubActions: true,
					SelectedFeatures: tc.features,
				},
			}

			// Create a mock generator
			generator := NewMockGenerator(cfg, workflowsDir)

			// Generate all workflows
			err = generator.Generate()
			if err != nil {
				t.Fatalf("Failed to generate workflows: %v", err)
			}

			// Check if the changelog workflow was generated
			if tc.expectChangelog {
				if !generator.ChangelogCalled {
					t.Errorf("Expected GenerateChangelogWorkflow to be called, but it wasn't")
				}

				changelogPath := filepath.Join(workflowsDir, "changelog.yml")
				if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
					t.Errorf("Expected changelog.yml to be created, but it wasn't")
				}
			} else {
				if generator.ChangelogCalled {
					t.Errorf("GenerateChangelogWorkflow was called when it shouldn't have been")
				}

				changelogPath := filepath.Join(workflowsDir, "changelog.yml")
				if _, err := os.Stat(changelogPath); err == nil {
					t.Errorf("changelog.yml was created when it shouldn't have been")
				}
			}
		})
	}
}

// TestMockChangelogWorkflow tests the changelog workflow generation specifically
func TestMockChangelogWorkflow(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config with changelog feature
	cfg := &model.Config{
		ProjectName: "test-project",
		Pipeline: model.PipelineConfig{
			UseGitHubActions: true,
			SelectedFeatures: []string{"changelog"},
		},
	}

	// Create a mock generator
	generator := NewMockGenerator(cfg, tempDir)

	// Generate the changelog workflow
	err = generator.GenerateChangelogWorkflow()
	if err != nil {
		t.Fatalf("Failed to generate changelog workflow: %v", err)
	}

	// Check if the changelog workflow was generated
	if !generator.ChangelogCalled {
		t.Errorf("Expected GenerateChangelogWorkflow to be called, but it wasn't")
	}

	changelogPath := filepath.Join(tempDir, "changelog.yml")
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		t.Errorf("Expected changelog.yml to be created, but it wasn't")
	}
}
