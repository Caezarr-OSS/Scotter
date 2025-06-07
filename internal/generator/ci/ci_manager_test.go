package ci

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

func TestCIManagerFactory(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "ci-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Créer des sous-répertoires pour les templates
	templatesDir := filepath.Join(tmpDir, "templates")
	githubDir := filepath.Join(templatesDir, "github")
	gitlabDir := filepath.Join(templatesDir, "gitlab")
	circleciDir := filepath.Join(templatesDir, ".circleci")
	travisDir := filepath.Join(templatesDir, "travis")

	dirs := []string{githubDir, gitlabDir, circleciDir, travisDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create template directory %s: %v", dir, err)
		}
	}

	// Créer des templates de test
	testTemplates := map[string]string{
		filepath.Join(githubDir, "ci.yml.tmpl"):     "name: {{ .ProjectName }} CI",
		filepath.Join(gitlabDir, ".gitlab-ci.yml.tmpl"): "# GitLab CI for {{ .ProjectName }}",
		filepath.Join(circleciDir, "config.yml.tmpl"): "# CircleCI config for {{ .ProjectName }}",
		filepath.Join(travisDir, ".travis.yml.tmpl"): "# Travis CI for {{ .ProjectName }}",
	}

	for path, content := range testTemplates {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", path, err)
		}
	}

	// Test de la factory
	t.Run("FactoryCreateManager", func(t *testing.T) {
		factory := NewCIManagerFactory(templatesDir)
		
		testCases := []struct{
			name       string
			ciType     model.CIType
			expectType model.CIType
			expectNil  bool
		}{
			{"GithubActions", model.GithubActionsCI, model.GithubActionsCI, false},
			{"GitlabCI", model.GitlabCI, model.GitlabCI, false},
			{"CircleCI", model.CircleCI, model.CircleCI, false},
			{"TravisCI", model.TravisCI, model.TravisCI, false},
			{"None", model.NoneCI, "", true},
			{"Empty", "", "", true},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cfg := &model.Config{
					ProjectName: "TestProject",
					Language:    "go",
					Pipeline: model.PipelineConfig{
						CIType: tc.ciType,
					},
				}
				
				manager, err := factory.CreateManager(cfg)
				if err != nil {
					t.Fatalf("Failed to create manager: %v", err)
				}
				
				if tc.expectNil {
					if manager != nil {
						t.Errorf("Expected nil manager, got %T", manager)
					}
					return
				}
				
				if manager == nil {
					t.Fatal("Expected non-nil manager")
				}
				
				if manager.GetType() != tc.expectType {
					t.Errorf("Expected manager type %s, got %s", tc.expectType, manager.GetType())
				}
			})
		}
	})
}

func TestGitHubActionsManager(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "github-actions-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Créer des sous-répertoires pour les templates et la sortie
	templatesDir := filepath.Join(tmpDir, "templates")
	githubDir := filepath.Join(templatesDir, "github")
	outputDir := filepath.Join(tmpDir, "output")
	workflowsDir := filepath.Join(outputDir, ".github", "workflows")

	for _, dir := range []string{githubDir, workflowsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Créer un template de test GitHub Actions
	ciTemplate := `name: "[[ .ProjectName ]] CI"
on:
  push:
    branches: ["main"]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.18'
`
	ciPath := filepath.Join(githubDir, "ci.yml.tmpl")
	if err := os.WriteFile(ciPath, []byte(ciTemplate), 0644); err != nil {
		t.Fatalf("Failed to write CI template: %v", err)
	}

	// Créer une configuration
	cfg := &model.Config{
		ProjectName: "TestGitHubActions",
		Language:    model.GoLang,
		Pipeline: model.PipelineConfig{
			CIType: model.GithubActionsCI,
		},
		Go: model.GoConfig{
			ModulePath: "github.com/example/test",
		},
	}

	// Créer un template manager
	templateMgr := common.NewTemplateManager(cfg, templatesDir, common.LangGo)
	templateMgr = templateMgr.WithDelimiters(common.CustomDelimiters())

	// Créer le manager GitHub Actions
	manager := &GitHubActionsManager{
		config:      cfg,
		templateMgr: templateMgr,
	}

	// Tester la génération d'un workflow
	t.Run("GenerateWorkflow", func(t *testing.T) {
		// Configurer le chemin de sortie
		workflowPath := filepath.Join(workflowsDir, "ci.yml")
		
		// Générer le workflow
		err := manager.GenerateWorkflow("ci", workflowPath)
		if err != nil {
			t.Fatalf("Failed to generate workflow: %v", err)
		}
		
		// Vérifier que le fichier existe
		if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
			t.Fatalf("Workflow file was not created")
		}
		
		// Lire le contenu généré
		content, err := os.ReadFile(workflowPath)
		if err != nil {
			t.Fatalf("Failed to read generated workflow: %v", err)
		}
		
		// Vérifier que le contenu a été correctement généré avec les délimiteurs personnalisés
		expectedContent := `name: "TestGitHubActions CI"
on:
  push:
    branches: ["main"]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.18'
`
		if string(content) != expectedContent {
			t.Errorf("Generated workflow doesn't match.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
		}
	})
}
