package code

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of Go code files
type Generator struct {
	Config *model.Config
	// Base directory for templates
	TemplateDir string
}

// NewGenerator creates a new code generator
func NewGenerator(cfg *model.Config, templateDir string) *Generator {
	return &Generator{
		Config:      cfg,
		TemplateDir: templateDir,
	}
}

// Generate creates the main Go files based on project type
func (g *Generator) Generate() error {
	// Only generate Go code for Go projects
	if g.Config.Language != model.GoLang {
		fmt.Println("Skipping Go code generation for non-Go project...")
		return nil
	}

	fmt.Println("Generating Go code files...")

	// Initialize go.mod
	if err := g.generateGoMod(); err != nil {
		return err
	}

	// Generate README
	if err := g.generateReadme(); err != nil {
		return err
	}

	// Generate code based on project type
	switch g.Config.Go.ProjectType {
	case model.DefaultGoType:
		return g.generateDefaultProject()
	case model.LibraryGoType:
		return g.generateLibraryProject()
	case model.CLIGoType:
		return g.generateCLIProject()
	case model.APIGoType:
		return g.generateAPIProject()
	default:
		return fmt.Errorf("unknown project type: %s", g.Config.Go.ProjectType)
	}
}

// generateGoMod initializes the go.mod file
func (g *Generator) generateGoMod() error {
	cmdStr := fmt.Sprintf("go mod init %s", g.Config.Go.ModulePath)
	
	// Execute the command
	fmt.Println("Initializing Go module:", cmdStr)
	cmd := exec.Command("go", "mod", "init", g.Config.Go.ModulePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}
	
	return nil
}

// generateReadme creates the README.md file
func (g *Generator) generateReadme() error {
	templatePath := filepath.Join(g.TemplateDir, "readme.md.tmpl")
	outputPath := "README.md"
	
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}
	
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()
	
	// Create template data with additional fields
	data := struct {
		*model.Config
		Description string
	}{
		Config:      g.Config,
		Description: fmt.Sprintf("A %s Go project", g.Config.Go.ProjectType),
	}
	
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	
	return nil
}

// generateDefaultProject creates a minimal Go project
func (g *Generator) generateDefaultProject() error {
	// Create a simple main.go in the root
	templatePath := filepath.Join(g.TemplateDir, "default_main.go.tmpl")
	outputPath := "main.go"
	
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}
	
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, g.Config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	
	return nil
}

// generateLibraryProject creates a library project structure
func (g *Generator) generateLibraryProject() error {
	// Create package main file in pkg directory
	pkgDir := filepath.Join("pkg", g.Config.ProjectName)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pkgDir, err)
	}
	
	// Create a package file
	content := fmt.Sprintf(`package %s

// Version is the current version of the library
const Version = "0.1.0"

// NewClient creates a new client instance
func NewClient() *Client {
	return &Client{}
}

// Client provides functionality for %s
type Client struct {
	// Add your fields here
}

// Hello returns a friendly greeting
func (c *Client) Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}
`, g.Config.ProjectName, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s.go", g.Config.ProjectName)), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	// Create a test file
	testContent := fmt.Sprintf(`package %s

import (
	"testing"
)

func TestHello(t *testing.T) {
	client := NewClient()
	
	t.Run("with name", func(t *testing.T) {
		got := client.Hello("Test")
		want := "Hello, Test!"
		if got != want {
			t.Errorf("Hello() = %%q, want %%q", got, want)
		}
	})
	
	t.Run("without name", func(t *testing.T) {
		got := client.Hello("")
		want := "Hello, World!"
		if got != want {
			t.Errorf("Hello() = %%q, want %%q", got, want)
		}
	})
}
`, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s_test.go", g.Config.ProjectName)), []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	return nil
}

// generateCLIProject creates a CLI project structure
func (g *Generator) generateCLIProject() error {
	// Create main.go in cmd/projectname
	mainPath := filepath.Join("cmd", g.Config.ProjectName)
	if err := os.MkdirAll(mainPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", mainPath, err)
	}
	
	// Create main.go
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

// Version information (will be set during build)
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("%s v%%s (%%s) built on %%s\n", Version, CommitSHA, BuildDate)
		return
	}
	
	fmt.Println("Hello from %s!")
	// TODO: Add your CLI logic here
}
`, g.Config.ProjectName, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(mainPath, "main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}
	
	// Create a basic config.go in internal/config
	configDir := filepath.Join("internal", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", configDir, err)
	}
	
	configContent := `package config

// Config holds the application configuration
type Config struct {
	// Add your configuration fields here
	Debug bool
	LogLevel string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Debug: false,
		LogLevel: "info",
	}
}
`

	if err := os.WriteFile(filepath.Join(configDir, "config.go"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.go: %w", err)
	}
	
	return nil
}

// generateAPIProject creates an API project structure
func (g *Generator) generateAPIProject() error {
	// Implement API project structure
	// This is a placeholder for a more sophisticated implementation
	fmt.Println("Generating API project - basic structure only for now")
	
	// Create main.go in cmd/projectname
	mainPath := filepath.Join("cmd", g.Config.ProjectName)
	if err := os.MkdirAll(mainPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", mainPath, err)
	}
	
	// Create API server entry point
	projectName := g.Config.ProjectName
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Version information (will be set during build)
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Setup HTTP routes
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	
	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"version\":\"%s\",\"commit\":\"%s\",\"buildDate\":\"%s\"}", Version, CommitSHA, BuildDate)
	})
	
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting ` + projectName + ` API server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}`

	if err := os.WriteFile(filepath.Join(mainPath, "main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}
	
	return nil
}

// generateCompleteProject creates a complete project with all features
func (g *Generator) generateCompleteProject() error {
	// For now, we'll just use the API project structure as a starting point
	if err := g.generateAPIProject(); err != nil {
		return err
	}
	
	// Add additional directories and files here
	docsDir := "docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", docsDir, err)
	}
	
	// Create a README in the docs directory
	docsContent := fmt.Sprintf(`# %s Documentation

## Overview

This documentation covers how to use and contribute to %s.

## API Reference

### Endpoints

- GET /api/health - Health check endpoint
- GET /api/version - Version information

## Development

### Building from Source

See the main README.md file for build instructions.
`, g.Config.ProjectName, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(docsDir, "README.md"), []byte(docsContent), 0644); err != nil {
		return fmt.Errorf("failed to write docs README: %w", err)
	}
	
	return nil
}
