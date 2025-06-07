package ci

import (
	"fmt"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// CIManager définit l'interface commune pour tous les générateurs de CI/CD
type CIManager interface {
	// Generate génère tous les fichiers de pipeline pour un système CI donné
	Generate() error
	
	// GetType retourne le type de système CI
	GetType() model.CIType
	
	// GenerateWorkflow génère un workflow spécifique
	GenerateWorkflow(workflowName string, outputPath string) error
}

// CIManagerFactory est responsable de la création du bon manager CI
type CIManagerFactory struct {
	TemplateRootDir string
}

// NewCIManagerFactory crée une nouvelle fabrique de managers CI
func NewCIManagerFactory(templateDir string) *CIManagerFactory {
	return &CIManagerFactory{
		TemplateRootDir: templateDir,
	}
}

// CreateManager crée un manager CI basé sur la configuration
func (f *CIManagerFactory) CreateManager(config *model.Config) (CIManager, error) {
	// Pour assurer la compatibilité avec l'ancienne configuration
	ciType := config.Pipeline.CIType
	if ciType == "" && config.Pipeline.UseGitHubActions {
		ciType = model.GithubActionsCI
	}
	
	// Si aucun CI n'est configuré
	if ciType == "" || ciType == model.NoneCI {
		return nil, nil // Pas d'erreur, mais pas de manager
	}
	
	// Créer le TemplateManager adapté au langage du projet
	var language common.TemplateLanguage
	switch config.Language {
	case model.GoLang:
		language = common.LangGo
	case "python":
		language = common.LangPython
	case "rust":
		language = common.LangRust
	case "typescript":
		language = common.LangTypeScript
	default:
		language = ""
	}
	
	templateMgr := common.NewTemplateManager(
		config, 
		f.TemplateRootDir,
		language,
	)
	
	// Créer le manager approprié selon le type de CI
	switch ciType {
	case model.GithubActionsCI:
		return &GitHubActionsManager{
			config:      config,
			templateMgr: templateMgr.WithDelimiters(common.CustomDelimiters()),
		}, nil
	case model.GitlabCI:
		return &GitLabCIManager{
			config:      config,
			templateMgr: templateMgr,
		}, nil
	case model.CircleCI:
		return &CircleCIManager{
			config:      config,
			templateMgr: templateMgr,
		}, nil
	case model.TravisCI:
		return &TravisCIManager{
			config:      config,
			templateMgr: templateMgr,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported CI type: %s", ciType)
	}
}

// GitHubActionsManager, GitLabCIManager, CircleCIManager et TravisCIManager
// sont implémentés dans leurs fichiers respectifs github_actions.go, gitlab_ci.go, etc.
