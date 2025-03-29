package taskfile

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of Taskfile
type Generator struct {
	Config *model.Config
	// Base directory for templates
	TemplateDir string
}

// NewGenerator creates a new Taskfile generator
func NewGenerator(cfg *model.Config, templateDir string) *Generator {
	return &Generator{
		Config:      cfg,
		TemplateDir: templateDir,
	}
}

// Generate creates the Taskfile.yml
func (g *Generator) Generate() error {
	if !g.Config.Features.UseTaskFile {
		fmt.Println("Skipping Taskfile generation...")
		return nil
	}

	fmt.Println("Generating Taskfile...")

	templatePath := filepath.Join(g.TemplateDir, "taskfile.yml.tmpl")
	outputPath := "Taskfile.yml"
	
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
	
	fmt.Println("Taskfile generated successfully!")
	return nil
}
