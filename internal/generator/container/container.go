package container

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of container-related files
type Generator struct {
	Config *model.Config
	// Base directory for templates
	TemplateDir string
}

// NewGenerator creates a new container files generator
func NewGenerator(cfg *model.Config, templateDir string) *Generator {
	return &Generator{
		Config:      cfg,
		TemplateDir: templateDir,
	}
}

// Generate creates container-related files
func (g *Generator) Generate() error {
	// Check if container feature is enabled in the pipeline
	hasContainer := false
	for _, feature := range g.Config.Pipeline.SelectedFeatures {
		if feature == "container" {
			hasContainer = true
			break
		}
	}

	if !hasContainer {
		fmt.Println("Skipping container files generation...")
		return nil
	}

	fmt.Println("Generating container configuration...")

	// Generate the container file based on the selected format
	if err := g.GenerateContainerFile(); err != nil {
		return err
	}

	// Generate GitHub workflow for container builds if GitHub Actions is enabled
	if g.Config.Pipeline.UseGitHubActions {
		if err := g.GenerateContainerWorkflow(); err != nil {
			return err
		}
	}

	fmt.Println("Container configuration generated successfully!")
	return nil
}

// GenerateContainerFile generates the appropriate container file (Dockerfile or Containerfile)
func (g *Generator) GenerateContainerFile() error {
	var templateName string
	var outputFileName string

	// Determine which template to use based on language
	switch g.Config.Language {
	case model.GoLang:
		templateName = "go.dockerfile.tmpl"
	case model.NoLang:
		templateName = "shell.dockerfile.tmpl"
	default:
		templateName = "default.dockerfile.tmpl"
	}

	// Determine output file name based on selected format
	switch g.Config.Pipeline.ContainerFormat {
	case model.DockerfileFormat:
		outputFileName = "Dockerfile"
	case model.ContainerfileFormat:
		outputFileName = "Containerfile"
	default:
		outputFileName = "Dockerfile"
	}

	fmt.Printf("Generating %s...\n", outputFileName)

	// Create the container file
	return g.generateContainerFile(templateName, outputFileName)
}

// GenerateContainerWorkflow generates a GitHub workflow for container builds
func (g *Generator) GenerateContainerWorkflow() error {
	fmt.Println("Generating container workflow...")
	return g.generateWorkflow("container.yml")
}

// generateContainerFile generates a container file from template
func (g *Generator) generateContainerFile(templateName, outputFileName string) error {
	// Get template path
	templatePath := filepath.Join(g.TemplateDir, "container", templateName)

	// Create template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create output file
	outputPath := filepath.Join(".", outputFileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	err = tmpl.Execute(outputFile, g.Config)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// generateWorkflow generates a GitHub workflow file from template
func (g *Generator) generateWorkflow(workflowName string) error {
	// Get template path
	templatePath := filepath.Join(g.TemplateDir, "github", workflowName)

	// Create output directory
	outputDir := filepath.Join(".", ".github", "workflows")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Create template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create output file
	outputPath := filepath.Join(outputDir, workflowName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	err = tmpl.Execute(outputFile, g.Config)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
