package github

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of GitHub-related files
type Generator struct {
	Config *model.Config
	// Base directory for templates
	TemplateDir string
}

// NewGenerator creates a new GitHub files generator
func NewGenerator(cfg *model.Config, templateDir string) *Generator {
	return &Generator{
		Config:      cfg,
		TemplateDir: templateDir,
	}
}

// Generate creates GitHub-related files
func (g *Generator) Generate() error {
	if !g.Config.Features.GitHub.UseWorkflows {
		fmt.Println("Skipping GitHub workflows generation...")
		return nil
	}

	fmt.Println("Generating GitHub configuration...")

	// Create CI workflow
	if err := g.generateWorkflow("ci.yml"); err != nil {
		return err
	}

	// Create commitlint workflow if enabled
	if g.Config.Features.GitHub.UseCommitLint {
		if err := g.generateWorkflow("commitlint.yml"); err != nil {
			return err
		}
		
		// Generate commitlint.config.js
		if err := g.generateCommitlintConfig(); err != nil {
			return err
		}
	}

	// Create release workflow if enabled
	if g.Config.Features.GitHub.UseReleaseWorkflow {
		if err := g.generateWorkflow("release.yml"); err != nil {
			return err
		}
		
		// Generate GoReleaser configuration
		if err := g.generateGoReleaserConfig(); err != nil {
			return err
		}
	}

	// Create Dependabot configuration if enabled
	if g.Config.Features.GitHub.UseDependabot {
		if err := g.generateDependabot(); err != nil {
			return err
		}
	}

	fmt.Println("GitHub configuration generated successfully!")
	return nil
}

// generateWorkflow generates a GitHub workflow file from template
func (g *Generator) generateWorkflow(workflowName string) error {
	templatePath := filepath.Join(g.TemplateDir, "github", workflowName+".tmpl")
	outputPath := filepath.Join(".github", "workflows", workflowName)
	
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

// generateCommitlintConfig generates the commitlint configuration file
func (g *Generator) generateCommitlintConfig() error {
	content := `module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'body-leading-blank': [1, 'always'],
    'body-max-line-length': [2, 'always', 100],
    'footer-leading-blank': [1, 'always'],
    'footer-max-line-length': [2, 'always', 100],
    'header-max-length': [2, 'always', 100],
    'subject-case': [
      2,
      'never',
      ['sentence-case', 'start-case', 'pascal-case', 'upper-case'],
    ],
    'subject-empty': [2, 'never'],
    'subject-full-stop': [2, 'never', '.'],
    'type-case': [2, 'always', 'lower-case'],
    'type-empty': [2, 'never'],
    'type-enum': [
      2,
      'always',
      [
        'build',
        'chore',
        'ci',
        'docs',
        'feat',
        'fix',
        'perf',
        'refactor',
        'revert',
        'style',
        'test',
      ],
    ],
  },
};`

	return os.WriteFile("commitlint.config.js", []byte(content), 0644)
}

// generateGoReleaserConfig generates GoReleaser configuration
func (g *Generator) generateGoReleaserConfig() error {
	templatePath := filepath.Join(g.TemplateDir, "goreleaser.yml.tmpl")
	outputPath := ".goreleaser.yml"
	
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
	
	fmt.Println("GoReleaser configuration generated successfully!")
	return nil
}

// generateDependabot generates Dependabot configuration
func (g *Generator) generateDependabot() error {
	content := `version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10

  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
`

	dependabotDir := filepath.Join(".github")
	if err := os.MkdirAll(dependabotDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dependabotDir, err)
	}

	return os.WriteFile(filepath.Join(dependabotDir, "dependabot.yml"), []byte(content), 0644)
}
