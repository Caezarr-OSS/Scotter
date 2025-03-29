package structure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of project structure
type Generator struct {
	Config *model.Config
}

// NewGenerator creates a new structure generator
func NewGenerator(cfg *model.Config) *Generator {
	return &Generator{
		Config: cfg,
	}
}

// Generate creates the directory structure for a project
func (g *Generator) Generate() error {
	fmt.Println("Generating project structure...")

	// Base directories needed for all project types
	dirs := []string{
		".",                                        // Root directory
		filepath.Join(".github"),                   // GitHub configuration directory
		filepath.Join(".github", "workflows"),      // GitHub Actions workflows
	}

	// Add type-specific directories
	switch g.Config.ProjectType {
	case model.DefaultType:
		// Minimal structure, nothing to add
	case model.LibraryType:
		dirs = append(dirs, "pkg")
	case model.CLIType:
		dirs = append(dirs, 
			"cmd",
			filepath.Join("cmd", g.Config.ProjectName),
			"internal",
			filepath.Join("internal", "config"),
		)
	case model.APIType:
		dirs = append(dirs, 
			"cmd",
			filepath.Join("cmd", g.Config.ProjectName),
			"internal",
			filepath.Join("internal", "api"),
			filepath.Join("internal", "config"),
			filepath.Join("internal", "middleware"),
			filepath.Join("internal", "handler"),
		)
	case model.CompleteType:
		dirs = append(dirs, 
			"cmd",
			filepath.Join("cmd", g.Config.ProjectName),
			"internal",
			filepath.Join("internal", "api"),
			filepath.Join("internal", "config"),
			filepath.Join("internal", "middleware"),
			filepath.Join("internal", "handler"),
			"pkg",
			"docs",
			"scripts",
		)
	}

	// Create all directories
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	fmt.Println("Project structure created successfully!")
	return nil
}

// GenerateGitIgnore creates a .gitignore file
func (g *Generator) GenerateGitIgnore() error {
	content := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Build artifacts
dist/
bin/

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# OS specific files
.DS_Store
Thumbs.db
`

	return os.WriteFile(".gitignore", []byte(content), 0644)
}
