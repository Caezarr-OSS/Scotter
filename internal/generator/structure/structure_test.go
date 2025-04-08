package structure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TestNewGenerator tests the constructor
func TestNewGenerator(t *testing.T) {
	cfg := &model.Config{
		ProjectName: "testproject",
	}
	generator := NewGenerator(cfg)
	
	if generator == nil {
		t.Fatal("expected generator to be created, got nil")
	}
	
	if generator.Config != cfg {
		t.Errorf("expected Config to be %v, got %v", cfg, generator.Config)
	}
}

// TestGenerateMinimalStructure tests creating a minimal project structure
func TestGenerateMinimalStructure(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore original directory: %v", err)
		}
	}()
	
	// Create a config with Go language and default type
	cfg := &model.Config{
		ProjectName: "testproject",
		Language: model.GoLang,
		Go: model.GoConfig{
			ProjectType: model.DefaultGoType,
		},
		Pipeline: model.PipelineConfig{
			UseGitHubActions: true,
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	
	// Check that the required directories were created
	requiredDirs := []string{
		".github",
		filepath.Join(".github", "workflows"),
	}
	
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("expected directory %s to exist", dir)
		}
	}
}

// TestGenerateAPIStructure tests creating an API project structure
func TestGenerateAPIStructure(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore original directory: %v", err)
		}
	}()
	
	// Create a config with Go language and API type
	cfg := &model.Config{
		ProjectName: "testapi",
		Language: model.GoLang,
		Go: model.GoConfig{
			ProjectType: model.APIGoType,
		},
		Pipeline: model.PipelineConfig{
			UseGitHubActions: true,
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	
	// Check that the required directories were created
	requiredDirs := []string{
		".github",
		filepath.Join(".github", "workflows"),
		"cmd",
		filepath.Join("cmd", "testapi"),
		"internal",
		filepath.Join("internal", "api"),
		filepath.Join("internal", "config"),
		filepath.Join("internal", "middleware"),
		filepath.Join("internal", "handler"),
	}
	
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("expected directory %s to exist", dir)
		}
	}
}

// TestGenerateGitIgnore tests creating a .gitignore file
func TestGenerateGitIgnore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore original directory: %v", err)
		}
	}()
	
	// Create a config with Go language
	cfg := &model.Config{
		ProjectName: "testproject",
		Language: model.GoLang,
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.GenerateGitIgnore(); err != nil {
		t.Fatalf("GenerateGitIgnore() failed: %v", err)
	}
	
	// Check that .gitignore was created
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Error("expected .gitignore to exist")
	}
	
	// Check contents
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	
	// Check for key patterns
	expectedPatterns := []string{
		"*.exe",
		"*.test",
		"vendor/",
		".idea/",
		".DS_Store",
	}
	
	for _, pattern := range expectedPatterns {
		if !containsString(string(content), pattern) {
			t.Errorf("expected .gitignore to contain %q", pattern)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
