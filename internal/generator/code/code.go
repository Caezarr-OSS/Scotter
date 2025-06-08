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
		// Pour la compatibilité avec les versions antérieures
		return g.generateLibraryProject()
	case model.LocalLibraryGoType:
		return g.generateLocalLibraryProject()
	case model.DistributedLibraryGoType:
		return g.generateDistributedLibraryProject()
	case model.CLIGoType:
		return g.generateCLIProject()
	case model.APIGoType:
		return g.generateAPIProject()
	case model.CompleteGoType:
		return g.generateCompleteProject()
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

// ExampleFunction is an example function that can be called directly
// without creating a client instance
func ExampleFunction(input string) string {
	return "Example output for: " + input
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

func TestExampleFunction(t *testing.T) {
	got := ExampleFunction("test")
	want := "Example output for: test"
	if got != want {
		t.Errorf("ExampleFunction() = %%q, want %%q", got, want)
	}
}
`, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s_test.go", g.Config.ProjectName)), []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	// Create API documentation
	docDir := "docs"
	if err := os.MkdirAll(docDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", docDir, err)
	}
	
	// Generate API.md from template
	apiTmplPath := filepath.Join(g.TemplateDir, "api.md.tmpl")
	apiTmpl, err := template.ParseFiles(apiTmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse API template: %w", err)
	}
	
	apiFile, err := os.Create(filepath.Join(docDir, "API.md"))
	if err != nil {
		return fmt.Errorf("failed to create API documentation file: %w", err)
	}
	defer apiFile.Close()
	
	data := struct {
		ProjectName string
		ModulePath  string
	}{
		ProjectName: g.Config.ProjectName,
		ModulePath:  g.Config.Go.ModulePath,
	}
	
	if err := apiTmpl.Execute(apiFile, data); err != nil {
		return fmt.Errorf("failed to execute API template: %w", err)
	}
	
	// Create examples directory and example file
	examplesDir := "examples"
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", examplesDir, err)
	}
	
	// Generate example.go from template
	exampleTmplPath := filepath.Join(g.TemplateDir, "example.go.tmpl")
	exampleTmpl, err := template.ParseFiles(exampleTmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse example template: %w", err)
	}
	
	exampleFile, err := os.Create(filepath.Join(examplesDir, "basic_usage.go"))
	if err != nil {
		return fmt.Errorf("failed to create example file: %w", err)
	}
	defer exampleFile.Close()
	
	if err := exampleTmpl.Execute(exampleFile, data); err != nil {
		return fmt.Errorf("failed to execute example template: %w", err)
	}
	
	// Create a README.md in the examples directory
	title := "# " + g.Config.ProjectName + " Examples"
	description := "This directory contains examples showing how to use the " + g.Config.ProjectName + " library."
	usage := "## Running examples\n\nTo run an example:\n\n```bash\ngo run examples/basic_usage.go\n```"
	exampleReadmeContent := title + "\n\n" + description + "\n\n" + usage
	
	if err := os.WriteFile(filepath.Join(examplesDir, "README.md"), []byte(exampleReadmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write examples README file: %w", err)
	}
	
	return nil
}

// generateCLIProject creates a CLI project structure
// generateLocalLibraryProject crée une structure de projet bibliothèque locale simplifiée
func (g *Generator) generateLocalLibraryProject() error {
	// Pour les bibliothèques locales, on crée une structure minimale sans documentation exhaustive
	// Create package main file in pkg directory
	pkgDir := filepath.Join("pkg", g.Config.ProjectName)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pkgDir, err)
	}
	
	// Create a package file - simplified for local use
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

// ExampleFunction is an example function that can be called directly
// without creating a client instance
func ExampleFunction(input string) string {
	return "Example output for: " + input
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
}
`, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s_test.go", g.Config.ProjectName)), []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	// Create a simple README to explain this is a local library
	projectName := g.Config.ProjectName
	modulePath := g.Config.Go.ModulePath
	readmeContent := "# " + projectName + "\n\n"
	readmeContent += "Local Go library for internal use.\n\n"
	readmeContent += "## Usage\n\n"
	readmeContent += "This library is intended for local development and is not configured for distribution.\n"
	readmeContent += "To use it in your projects, add it as a local dependency using Go's replace directive:\n\n"
	readmeContent += "```go\n"
	readmeContent += "// In your go.mod file\n"
	readmeContent += fmt.Sprintf("require %s v0.0.0-unpublished\n\n", modulePath)
	readmeContent += fmt.Sprintf("replace %s => ../path/to/%s\n", modulePath, projectName)
	readmeContent += "```"
	
	if err := os.WriteFile("README.md", []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	return nil
}

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
// generateDistributedLibraryProject crée une structure de projet bibliothèque destinée à la distribution
func (g *Generator) generateDistributedLibraryProject() error {
	// Pour les bibliothèques distribuées, on crée une structure complète avec documentation et exemples
	// Create package main file in pkg directory
	pkgDir := filepath.Join("pkg", g.Config.ProjectName)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pkgDir, err)
	}
	
	// Create a package file - optimisé pour la distribution
	content := fmt.Sprintf(`package %s

import (
	"time"
)

// Version is the current version of the library
const Version = "0.1.0"

// NewClient creates a new client instance
func NewClient(options ...Option) *Client {
	client := &Client{
		// Définir des valeurs par défaut
		options: defaultOptions(),
	}
	
	// Appliquer les options personnalisées
	for _, opt := range options {
		opt(client)
	}
	
	return client
}

// Option est un type de fonction pour configurer un Client
type Option func(*Client)

// ClientOptions contient les options configurables
type ClientOptions struct {
	Timeout time.Duration
	Retries int
}

// defaultOptions retourne les options par défaut
func defaultOptions() ClientOptions {
	return ClientOptions{
		Timeout: 30 * time.Second,
		Retries: 3,
	}
}

// WithTimeout définit le timeout du client
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.options.Timeout = timeout
	}
}

// WithRetries définit le nombre de tentatives
func WithRetries(retries int) Option {
	return func(c *Client) {
		c.options.Retries = retries
	}
}

// Client provides functionality for %s
type Client struct {
	// Options de configuration
	options ClientOptions
	// Ajouter d'autres champs ici
}

// Hello returns a friendly greeting
func (c *Client) Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}

// ExampleFunction is an example function that can be called directly
// without creating a client instance
func ExampleFunction(input string) string {
	return "Example output for: " + input
}
`, g.Config.ProjectName, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s.go", g.Config.ProjectName)), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	// Create a test file with des tests plus exhaustifs
	testContent := fmt.Sprintf(`package %s

import (
	"testing"
	"time"
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

func TestClientOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		client := NewClient()
		if client.options.Timeout != 30*time.Second {
			t.Errorf("Expected default timeout of 30s, got %%v", client.options.Timeout)
		}
		if client.options.Retries != 3 {
			t.Errorf("Expected default retries of 3, got %%d", client.options.Retries)
		}
	})
	
	t.Run("custom options", func(t *testing.T) {
		client := NewClient(
			WithTimeout(10*time.Second),
			WithRetries(5),
		)
		if client.options.Timeout != 10*time.Second {
			t.Errorf("Expected timeout of 10s, got %%v", client.options.Timeout)
		}
		if client.options.Retries != 5 {
			t.Errorf("Expected retries of 5, got %%d", client.options.Retries)
		}
	})
}

func TestExampleFunction(t *testing.T) {
	got := ExampleFunction("test")
	want := "Example output for: test"
	if got != want {
		t.Errorf("ExampleFunction() = %%q, want %%q", got, want)
	}
}
`, g.Config.ProjectName)

	if err := os.WriteFile(filepath.Join(pkgDir, fmt.Sprintf("%s_test.go", g.Config.ProjectName)), []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	// Create API documentation
	docDir := "docs"
	if err := os.MkdirAll(docDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", docDir, err)
	}
	
	// Generate API.md from template
	apiTmplPath := filepath.Join(g.TemplateDir, "api.md.tmpl")
	apiTmpl, err := template.ParseFiles(apiTmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse API template: %w", err)
	}
	
	apiFile, err := os.Create(filepath.Join(docDir, "API.md"))
	if err != nil {
		return fmt.Errorf("failed to create API documentation file: %w", err)
	}
	defer apiFile.Close()
	
	data := struct {
		ProjectName string
		ModulePath  string
	}{
		ProjectName: g.Config.ProjectName,
		ModulePath:  g.Config.Go.ModulePath,
	}
	
	if err := apiTmpl.Execute(apiFile, data); err != nil {
		return fmt.Errorf("failed to execute API template: %w", err)
	}
	
	// Create examples directory and example file
	examplesDir := "examples"
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", examplesDir, err)
	}
	
	// Generate example.go from template
	exampleTmplPath := filepath.Join(g.TemplateDir, "example.go.tmpl")
	exampleTmpl, err := template.ParseFiles(exampleTmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse example template: %w", err)
	}
	
	exampleFile, err := os.Create(filepath.Join(examplesDir, "basic_usage.go"))
	if err != nil {
		return fmt.Errorf("failed to create example file: %w", err)
	}
	defer exampleFile.Close()
	
	if err := exampleTmpl.Execute(exampleFile, data); err != nil {
		return fmt.Errorf("failed to execute example template: %w", err)
	}
	
	// Create a README.md in the examples directory
	title := "# " + g.Config.ProjectName + " Examples"
	description := "This directory contains examples showing how to use the " + g.Config.ProjectName + " library."
	usage := "## Running examples\n\nTo run an example:\n\n```bash\ngo run examples/basic_usage.go\n```"
	exampleReadmeContent := title + "\n\n" + description + "\n\n" + usage
	
	if err := os.WriteFile(filepath.Join(examplesDir, "README.md"), []byte(exampleReadmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write examples README file: %w", err)
	}
	
	// Create a complete README.md in the root directory
	rootReadmeContent := fmt.Sprintf("# %s\n\n", g.Config.ProjectName)
	rootReadmeContent += "A Go library designed for distribution via Go modules.\n\n"
	rootReadmeContent += "## Installation\n\n"
	rootReadmeContent += fmt.Sprintf("```bash\ngo get %s\n```\n\n", g.Config.Go.ModulePath)
	rootReadmeContent += "## Usage\n\n"
	rootReadmeContent += "```go\n"
	rootReadmeContent += fmt.Sprintf("import \"%s\"\n\n", g.Config.Go.ModulePath)
	rootReadmeContent += "// Create a client with default options\n"
	rootReadmeContent += fmt.Sprintf("client := %s.NewClient()\n\n", g.Config.ProjectName)
	rootReadmeContent += "// Or with custom options\n"
	rootReadmeContent += fmt.Sprintf("client := %s.NewClient(\n", g.Config.ProjectName)
	rootReadmeContent += fmt.Sprintf("\t%s.WithTimeout(10*time.Second),\n", g.Config.ProjectName)
	rootReadmeContent += fmt.Sprintf("\t%s.WithRetries(5),\n", g.Config.ProjectName)
	rootReadmeContent += ")\n"
	rootReadmeContent += "```\n\n"
	rootReadmeContent += "See [API documentation](docs/API.md) and [examples](examples/) for more details.\n"
	
	if err := os.WriteFile("README.md", []byte(rootReadmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}
	
	// Create additional files for OSS projects comme CONTRIBUTING.md, CODE_OF_CONDUCT.md, etc.
	contributingContent := fmt.Sprintf("# Contributing to %s\n\n", g.Config.ProjectName)
	contributingContent += "We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:\n\n"
	contributingContent += "- Reporting a bug\n"
	contributingContent += "- Discussing the current state of the code\n"
	contributingContent += "- Submitting a fix\n"
	contributingContent += "- Proposing new features\n"
	contributingContent += "- Becoming a maintainer\n\n"
	contributingContent += "## Development Process\n\n"
	contributingContent += "We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.\n"
	
	if err := os.WriteFile("CONTRIBUTING.md", []byte(contributingContent), 0644); err != nil {
		return fmt.Errorf("failed to write CONTRIBUTING file: %w", err)
	}
	
	return nil
}

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
