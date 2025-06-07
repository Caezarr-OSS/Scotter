package github

import (
	"fmt"
	"os"
	"path/filepath"

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
// This is a legacy method for backward compatibility
func (g *Generator) Generate() error {
	// Check if GitHub Actions is enabled in the pipeline
	if !g.Config.Pipeline.UseGitHubActions {
		fmt.Println("Skipping GitHub workflows generation...")
		return nil
	}

	fmt.Println("Generating GitHub configuration...")

	// Check which features are enabled
	hasCI := false
	hasCommitLint := false
	hasRelease := false
	hasDependabot := false
	hasChangelog := false

	for _, feature := range g.Config.Pipeline.SelectedFeatures {
		switch feature {
		case "ci":
			hasCI = true
		case "commit-lint":
			hasCommitLint = true
		case "release":
			hasRelease = true
		case "dependabot":
			hasDependabot = true
		case "changelog":
			hasChangelog = true
		}
	}

	// Create CI workflow
	if hasCI {
		if err := g.GenerateCIWorkflow(); err != nil {
			return err
		}
	}

	// Create commitlint workflow if enabled
	if hasCommitLint {
		if err := g.GenerateCommitLintWorkflow(); err != nil {
			return err
		}
		
		// Generate commitlint configuration file
		if err := g.generateCommitlintConfig(); err != nil {
			return fmt.Errorf("failed to generate commitlint config: %w", err)
		}
	}

	// Create release workflow if enabled
	if hasRelease {
		if err := g.GenerateReleaseWorkflow(); err != nil {
			return err
		}
	}

	// Create Dependabot configuration if enabled
	if hasDependabot {
		if err := g.GenerateDependabotConfig(); err != nil {
			return err
		}
	}
	
	// Create Changelog workflow if enabled
	if hasChangelog {
		if err := g.GenerateChangelogWorkflow(); err != nil {
			return err
		}
	}

	fmt.Println("GitHub configuration generated successfully!")
	return nil
}

// GenerateCIWorkflow generates a CI workflow file
func (g *Generator) GenerateCIWorkflow() error {
	fmt.Println("Generating CI workflow...")
	return g.generateWorkflow("ci.yml")
}

// GenerateCommitLintWorkflow generates a commit lint workflow file
func (g *Generator) GenerateCommitLintWorkflow() error {
	fmt.Println("Generating commit lint workflow...")
	return g.generateWorkflow("commitlint.yml")
}

// GenerateChangelogWorkflow generates a changelog workflow file
func (g *Generator) GenerateChangelogWorkflow() error {
	fmt.Println("Generating changelog workflow...")
	return g.generateWorkflow("changelog.yml")
}

// GenerateReleaseWorkflow generates a release workflow file
func (g *Generator) GenerateReleaseWorkflow() error {
	fmt.Println("Generating release workflow...")
	
	// Generate the workflow file
	if err := g.generateWorkflow("release.yml"); err != nil {
		return err
	}
	
	// For Go projects, also generate GoReleaser configuration
	if g.Config.Language == model.GoLang {
		if err := g.generateGoReleaserConfig(); err != nil {
			return err
		}
	}
	
	return nil
}

// GenerateDependabotConfig generates Dependabot configuration
func (g *Generator) GenerateDependabotConfig() error {
	fmt.Println("Generating Dependabot configuration...")
	return g.generateDependabot()
}

// generateWorkflow generates a GitHub workflow file from template
func (g *Generator) generateWorkflow(workflowName string) error {
	// Liste des chemins possibles pour trouver le template
	possiblePaths := []string{
		// Chemin relatif standard
		filepath.Join("templates", "github", workflowName+".tmpl"),
		// Chemin basé sur TemplateDir (pour les tests)
		filepath.Join(g.TemplateDir, workflowName),
		// Chemin basé sur TemplateDir avec extension .tmpl
		filepath.Join(g.TemplateDir, workflowName+".tmpl"),
		// Chemin basé sur TemplateDir avec sous-dossier github
		filepath.Join(g.TemplateDir, "github", workflowName),
		// Chemin basé sur TemplateDir avec sous-dossier github et extension .tmpl
		filepath.Join(g.TemplateDir, "github", workflowName+".tmpl"),
		// Essayer avec l'extension .yml (pour les tests)
		filepath.Join(g.TemplateDir, workflowName+".yml"),
		// Essayer avec l'extension .yml dans le sous-dossier github
		filepath.Join(g.TemplateDir, "github", workflowName+".yml"),
	}
	
	// Chercher le template dans tous les chemins possibles
	var templatePath string
	var templateFound bool
	
	for _, path := range possiblePaths {
		_, err := os.Stat(path)
		if err == nil {
			templateFound = true
			templatePath = path
			break
		}
	}
	
	if !templateFound {
		return fmt.Errorf("template not found: %s (tried %v)", workflowName, possiblePaths)
	}
	
	// Créer le répertoire .github/workflows s'il n'existe pas
	workflowsDir := filepath.Join(".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflows directory: %v", err)
	}
	
	// Vérifier si le répertoire a bien été créé
	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		return fmt.Errorf("workflows directory does not exist after creation: %v", err)
	}

	outputPath := filepath.Join(workflowsDir, workflowName)
	
	// Utiliser les délimiteurs personnalisés pour les templates GitHub Actions
	// Cela permet d'éviter les conflits avec la syntaxe GitHub Actions ${{ ... }}
	delimiters := GitHubWorkflowDelimiters()
	
	// Exécuter le template avec les délimiteurs personnalisés
	result, err := ExecuteTemplateWithDelimiters(templatePath, g.Config, delimiters)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	
	// Écrire le résultat dans le fichier de sortie
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(result); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
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
	
	// Utiliser les délimiteurs personnalisés pour les templates YAML
	delimiters := GitHubWorkflowDelimiters()
	
	// Exécuter le template avec les délimiteurs personnalisés
	result, err := ExecuteTemplateWithDelimiters(templatePath, g.Config, delimiters)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	
	// Écrire le résultat dans le fichier de sortie
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(result); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
	}
	
	fmt.Println("GoReleaser configuration generated successfully!")
	return nil
}

// generateDependabot generates Dependabot configuration
func (g *Generator) generateDependabot() error {
	// Liste des chemins possibles pour trouver le template
	possiblePaths := []string{
		filepath.Join("templates", "github", "dependabot.yml.tmpl"),
		filepath.Join(g.TemplateDir, "github", "dependabot.yml.tmpl"),
		filepath.Join(g.TemplateDir, "dependabot.yml.tmpl"),
	}
	
	// Chercher le template dans tous les chemins possibles
	var templatePath string
	var templateFound bool
	
	for _, path := range possiblePaths {
		_, err := os.Stat(path)
		if err == nil {
			templateFound = true
			templatePath = path
			break
		}
	}
	
	if !templateFound {
		return fmt.Errorf("template not found: dependabot.yml.tmpl (tried %v)", possiblePaths)
	}

	// Créer le répertoire .github s'il n'existe pas
	dependabotDir := filepath.Join(".github")
	if err := os.MkdirAll(dependabotDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dependabotDir, err)
	}

	outputPath := filepath.Join(dependabotDir, "dependabot.yml")
	
	// Utiliser les délimiteurs personnalisés pour les templates YAML
	delimiters := GitHubWorkflowDelimiters()
	
	// Exécuter le template avec les délimiteurs personnalisés
	result, err := ExecuteTemplateWithDelimiters(templatePath, g.Config, delimiters)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	
	// Écrire le résultat dans le fichier de sortie
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(result); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", outputPath, err)
	}
	
	return nil
}
