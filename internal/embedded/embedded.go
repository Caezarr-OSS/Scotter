// Package embedded provides access to embedded templates and resources
package embedded

import (
	"embed"
	"io/fs"
	"text/template"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/caezarr-oss/scotter/pkg/plugin"
)

//go:embed templates/placeholder.txt
var templateFS embed.FS

// TemplateManager is the implementation of plugin.TemplateManager for embedded templates
type TemplateManager struct{}

// NewTemplateManager creates a new embedded template manager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{}
}

// Filesystem returns the embedded filesystem
func (m *TemplateManager) Filesystem() fs.FS {
	return templateFS
}

// RenderToFile renders a template to a file
func (m *TemplateManager) RenderToFile(templatePath, targetPath string, data interface{}) error {
	// Ensure the target directory exists
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Read template from embedded filesystem
	templateContent, err := fs.ReadFile(templateFS, templatePath)
	if err != nil {
		return err
	}

	// Parse and execute the template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return err
	}

	// Create target file
	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template to the file
	return tmpl.Execute(file, data)
}

// RenderToString renders a template as a string
func (m *TemplateManager) RenderToString(templatePath string, data interface{}) (string, error) {
	// Read template from embedded filesystem
	templateContent, err := fs.ReadFile(templateFS, templatePath)
	if err != nil {
		return "", err
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return "", err
	}

	// Execute the template to a string
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}

// Ensure TemplateManager implements plugin.TemplateManager
var _ plugin.TemplateManager = (*TemplateManager)(nil)
