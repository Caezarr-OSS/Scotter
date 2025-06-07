package ci

import (
	"fmt"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TravisCIManager est l'implémentation pour Travis CI
type TravisCIManager struct {
	config      *model.Config
	templateMgr *common.TemplateManager
}

// Generate génère le fichier .travis.yml
func (m *TravisCIManager) Generate() error {
	fmt.Println("Generating Travis CI configuration...")
	
	// Travis CI utilise un fichier .travis.yml à la racine
	if err := m.GenerateWorkflow("travis", ".travis.yml"); err != nil {
		return fmt.Errorf("failed to generate Travis CI configuration: %w", err)
	}
	
	// Générer d'autres fichiers spécifiques selon les fonctionnalités
	for _, feature := range m.config.Pipeline.SelectedFeatures {
		switch feature {
		case "commit-lint":
			if err := m.generateCommitlintConfig(); err != nil {
				return fmt.Errorf("failed to generate commitlint config: %w", err)
			}
		}
	}
	
	return nil
}

// GetType retourne le type de système CI
func (m *TravisCIManager) GetType() model.CIType {
	return model.TravisCI
}

// GenerateWorkflow génère un workflow Travis CI spécifique
func (m *TravisCIManager) GenerateWorkflow(workflowName, outputPath string) error {
	extensions := []string{".yml.tmpl", ".yaml.tmpl"}
	
	// Le chemin des templates suit une hiérarchie bien définie
	return m.templateMgr.GenerateFileFromTemplate(
		filepath.Join("travis", workflowName), 
		outputPath, 
		extensions,
	)
}

// generateCommitlintConfig génère la configuration commitlint
func (m *TravisCIManager) generateCommitlintConfig() error {
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
