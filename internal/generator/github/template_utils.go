package github

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateDelimiters defines custom delimiters for different template types
type TemplateDelimiters struct {
	Left  string
	Right string
}

// DefaultDelimiters returns the default Go template delimiters
func DefaultDelimiters() TemplateDelimiters {
	return TemplateDelimiters{
		Left:  "{{",
		Right: "}}",
	}
}

// GitHubWorkflowDelimiters returns the delimiters used for GitHub workflow files
// to avoid conflicts with GitHub Actions syntax
func GitHubWorkflowDelimiters() TemplateDelimiters {
	return TemplateDelimiters{
		Left:  "[[",
		Right: "]]",
	}
}

// ExecuteTemplateWithDelimiters renders a template with custom delimiters
func ExecuteTemplateWithDelimiters(
	templatePath string,
	data interface{},
	delimiters TemplateDelimiters,
) (string, error) {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Create a new template with custom delimiters
	name := filepath.Base(templatePath)
	tmpl := template.New(name).Delims(delimiters.Left, delimiters.Right)
	
	// Parse the template
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// ValidateTemplate checks if a template can be parsed with the given delimiters
func ValidateTemplate(templatePath string, delimiters TemplateDelimiters) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl := template.New(filepath.Base(templatePath)).Delims(delimiters.Left, delimiters.Right)
	_, err = tmpl.Parse(string(content))
	return err
}
