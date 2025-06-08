package ci

import (
	"fmt"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// GitHubActionsManager est l'implémentation pour GitHub Actions
type GitHubActionsManager struct {
	config      *model.Config
	templateMgr *common.TemplateManager
}

// NewGitHubActionsManager crée un nouveau manager GitHub Actions
func NewGitHubActionsManager(config *model.Config, templateMgr *common.TemplateManager) *GitHubActionsManager {
	// Si les délimiteurs personnalisés ne sont pas déjà définis, les configurer
	if templateMgr != nil {
		tmplMgr := templateMgr.WithDelimiters(common.CustomDelimiters())
		return &GitHubActionsManager{
			config:      config,
			templateMgr: tmplMgr,
		}
	}
	
	return &GitHubActionsManager{
		config:      config,
		templateMgr: nil,
	}
}

// Generate génère tous les fichiers de workflow GitHub Actions
func (m *GitHubActionsManager) Generate() error {
	fmt.Println("Generating GitHub Actions workflows...")
	
	// Vérifier quelles fonctionnalités sont activées
	hasCI := false
	hasCommitLint := false
	hasRelease := false
	hasDependabot := false
	hasChangelog := false
	
	for _, feature := range m.config.Pipeline.SelectedFeatures {
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
	
	// Générer les workflows configurés
	if hasCI {
		if err := m.GenerateWorkflow("ci", filepath.Join(".github", "workflows", "ci.yml")); err != nil {
			return fmt.Errorf("failed to generate CI workflow: %w", err)
		}
	}
	
	if hasCommitLint {
		if err := m.GenerateWorkflow("commitlint", filepath.Join(".github", "workflows", "commitlint.yml")); err != nil {
			return fmt.Errorf("failed to generate commit-lint workflow: %w", err)
		}
		
		// Générer également la configuration de commitlint
		if err := m.generateCommitlintConfig(); err != nil {
			return fmt.Errorf("failed to generate commitlint config: %w", err)
		}
	}
	
	if hasChangelog {
		if err := m.GenerateWorkflow("changelog", filepath.Join(".github", "workflows", "changelog.yml")); err != nil {
			return fmt.Errorf("failed to generate changelog workflow: %w", err)
		}
	}
	
	if hasRelease {
		// Pour les bibliothèques locales, on traite différemment
		skipReleaseWorkflow := m.config.Language == model.GoLang && m.config.Go.ProjectType == model.LocalLibraryGoType
		// Pour les bibliothèques Go distribuées, utiliser le workflow spécifique
		useGoLibraryWorkflow := m.config.Language == model.GoLang && 
			(m.config.Go.ProjectType == model.DistributedLibraryGoType || m.config.Go.ProjectType == model.LibraryGoType)
		
		// Générer la config GoReleaser uniquement pour les projets non-bibliothèques Go
		if m.config.Language == model.GoLang && !useGoLibraryWorkflow {
			if err := m.generateGoReleaserConfig(); err != nil {
				return fmt.Errorf("failed to generate GoReleaser config: %w", err)
			}
		}
		
		// Ne générer aucun workflow release pour les bibliothèques locales
		if skipReleaseWorkflow {
			fmt.Println("Skipping release workflow for local library...")
		} else if useGoLibraryWorkflow {
			// Pour les bibliothèques Go, utiliser notre nouveau workflow spécifique
			fmt.Println("Generating Go library specific release workflow...")
			if err := m.GenerateWorkflow("go-library-release", filepath.Join(".github", "workflows", "release.yml")); err != nil {
				return fmt.Errorf("failed to generate Go library release workflow: %w", err)
			}
		} else {
			// Pour les autres types de projets, utiliser le workflow standard
			if err := m.GenerateWorkflow("release", filepath.Join(".github", "workflows", "release.yml")); err != nil {
				return fmt.Errorf("failed to generate release workflow: %w", err)
			}
		}
		
		// Note: Même pour les bibliothèques locales, le fichier .goreleaser.yml est généré
		// mais avec une configuration qui ne publie pas sur les dépôts Go
	}
	
	if hasDependabot {
		if err := m.generateDependabot(); err != nil {
			return fmt.Errorf("failed to generate Dependabot configuration: %w", err)
		}
	}
	
	return nil
}

// GetType retourne le type de système CI
func (m *GitHubActionsManager) GetType() model.CIType {
	return model.GithubActionsCI
}

// GenerateWorkflow génère un workflow GitHub Actions spécifique
func (m *GitHubActionsManager) GenerateWorkflow(workflowName, outputPath string) error {
	// Extensions de fichiers acceptées pour les templates
	extensions := []string{".yml.tmpl", ".yaml.tmpl"}
	
	// Le chemin des templates suit une hiérarchie bien définie:
	// 1. templates/github/{language}/{workflowName}.yml.tmpl - Spécifique au langage et à GitHub
	// 2. templates/github/{workflowName}.yml.tmpl - Générique pour GitHub
	// 3. templates/{workflowName}.yml.tmpl - Fallback générique
	return m.templateMgr.GenerateFileFromTemplate(
		filepath.Join("github", workflowName), 
		outputPath, 
		extensions,
	)
}

// generateCommitlintConfig génère la configuration commitlint
func (m *GitHubActionsManager) generateCommitlintConfig() error {
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
	
	return m.templateMgr.GenerateFileFromContent(".commitlintrc.js", content)
}

// generateGoReleaserConfig génère la configuration GoReleaser
func (m *GitHubActionsManager) generateGoReleaserConfig() error {
	// Déterminer le template à utiliser en fonction du type de projet
	
	if m.config.Language == model.GoLang {
		// Vérifier si c'est une bibliothèque Go (qui n'a pas besoin de GoReleaser)
		isGoLibrary := m.config.Go.ProjectType == model.LibraryGoType || 
			m.config.Go.ProjectType == model.DistributedLibraryGoType || 
			m.config.Go.ProjectType == model.LocalLibraryGoType
		
		// Les bibliothèques Go n'ont pas besoin de GoReleaser (pas de binaires à construire)
		if isGoLibrary {
			fmt.Println("Skipping GoReleaser configuration for Go library... (not needed for libraries)")
			return nil
		}
		
		// Pour tous les autres types de projets Go (api, cli...)
		fmt.Println("Generating GoReleaser configuration...")
		return m.templateMgr.GenerateFileFromTemplate(
			"goreleaser",
			".goreleaser.yml",
			[]string{".yml.tmpl", ".yaml.tmpl"},
		)
	}
	
	// Pour les autres langages, générer la configuration standard
	return m.templateMgr.GenerateFileFromTemplate(
		"goreleaser", 
		".goreleaser.yml",
		[]string{".yml.tmpl", ".yaml.tmpl"},
	)
}

// generateDependabot génère la configuration Dependabot
func (m *GitHubActionsManager) generateDependabot() error {
	return m.templateMgr.GenerateFileFromTemplate(
		filepath.Join("github", "dependabot"), 
		filepath.Join(".github", "dependabot.yml"),
		[]string{".yml.tmpl", ".yaml.tmpl"},
	)
}
