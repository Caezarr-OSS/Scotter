package ci

import (
	"fmt"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// GitLabCIManager est l'implémentation pour GitLab CI
type GitLabCIManager struct {
	config      *model.Config
	templateMgr *common.TemplateManager
}

// Generate génère le fichier .gitlab-ci.yml
func (m *GitLabCIManager) Generate() error {
	fmt.Println("Generating GitLab CI configuration...")
	
	// GitLab CI utilise un seul fichier principal de configuration
	if err := m.GenerateWorkflow("gitlab-ci", ".gitlab-ci.yml"); err != nil {
		return fmt.Errorf("failed to generate GitLab CI configuration: %w", err)
	}
	
	// Générer d'autres fichiers spécifiques à GitLab selon les fonctionnalités
	for _, feature := range m.config.Pipeline.SelectedFeatures {
		switch feature {
		case "dependabot":
			// GitLab utilise Renovate ou Dependabot Edge
			if err := m.GenerateWorkflow("renovate", "renovate.json"); err != nil {
				return fmt.Errorf("failed to generate Renovate configuration: %w", err)
			}
		case "commit-lint":
			if err := m.generateCommitlintConfig(); err != nil {
				return fmt.Errorf("failed to generate commitlint config: %w", err)
			}
		}
	}
	
	return nil
}

// GetType retourne le type de système CI
func (m *GitLabCIManager) GetType() model.CIType {
	return model.GitlabCI
}

// GenerateWorkflow génère un workflow GitLab CI spécifique
func (m *GitLabCIManager) GenerateWorkflow(workflowName, outputPath string) error {
	extensions := []string{".yml.tmpl", ".yaml.tmpl", ".json.tmpl"}
	
	// Le chemin des templates suit une hiérarchie bien définie:
	// 1. templates/gitlab/{language}/{workflowName}.yml.tmpl - Spécifique au langage
	// 2. templates/gitlab/{workflowName}.yml.tmpl - Générique pour GitLab
	// 3. templates/{workflowName}.yml.tmpl - Fallback générique
	return m.templateMgr.GenerateFileFromTemplate(
		filepath.Join("gitlab", workflowName), 
		outputPath, 
		extensions,
	)
}

// generateCommitlintConfig génère la configuration commitlint
func (m *GitLabCIManager) generateCommitlintConfig() error {
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
