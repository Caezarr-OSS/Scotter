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
		".", // Root directory
	}

	// Add GitHub directories if GitHub Actions is enabled
	if g.Config.Pipeline.UseGitHubActions {
		dirs = append(dirs,
			filepath.Join(".github"),              // GitHub configuration directory
			filepath.Join(".github", "workflows"), // GitHub Actions workflows
		)
	}

	// Add language-specific directories
	if g.Config.Language == model.GoLang {
		// Add Go-specific directories based on project type
		switch g.Config.Go.ProjectType {
		case model.DefaultGoType:
			// Minimal structure, nothing to add
		case model.LibraryGoType:
			// Structure optimisée pour les bibliothèques Go (pour compatibilité)
			dirs = append(dirs, 
				"pkg",       // Code exportable et réutilisable
				"internal",  // Code interne non exporté
				"examples",  // Exemples d'utilisation de la bibliothèque
				"docs",      // Documentation détaillée
			)
		case model.DistributedLibraryGoType:
			// Structure optimisée pour les bibliothèques Go destinées à être distribuées
			dirs = append(dirs, 
				"pkg",       // Code exportable et réutilisable
				"internal",  // Code interne non exporté
				"examples",  // Exemples d'utilisation de la bibliothèque
				"docs",      // Documentation détaillée
			)
		case model.LocalLibraryGoType:
			// Structure simplifiée pour les bibliothèques Go locales
			dirs = append(dirs, 
				"pkg",       // Code exportable et réutilisable
			)
		case model.CLIGoType:
			dirs = append(dirs,
				"cmd",
				filepath.Join("cmd", g.Config.ProjectName),
				"internal",
				filepath.Join("internal", "config"),
			)
		case model.APIGoType:
			dirs = append(dirs,
				"cmd",
				filepath.Join("cmd", g.Config.ProjectName),
				"internal",
				filepath.Join("internal", "api"),
				filepath.Join("internal", "config"),
				filepath.Join("internal", "middleware"),
				filepath.Join("internal", "handler"),
			)
		}
	} else if g.Config.Language == model.NoLang {
		// For shell/script projects, create a minimal structure
		dirs = append(dirs,
			"scripts",
			"docs",
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
	// Base gitignore content
	content := "# OS specific files\n.DS_Store\nThumbs.db\n\n# IDE specific files\n.idea/\n.vscode/\n*.swp\n*.swo\n\n"

	// Add language-specific gitignore content
	if g.Config.Language == model.GoLang {
		content += `# Go specific ignores
# Binaries for programs and plugins
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
`
	} else {
		// Generic ignores for shell/script projects
		content += `# Build artifacts
dist/
bin/

# Logs
*.log

# Temporary files
*.tmp
*~
`
	}

	return os.WriteFile(".gitignore", []byte(content), 0644)
}
