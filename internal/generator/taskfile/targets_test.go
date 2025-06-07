package taskfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

func TestBuildTargetManager(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tempDir := filepath.Join(os.TempDir(), "scotter_test_targets")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Erreur lors de la création du répertoire de test: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Créer un taskfile de test
	taskfileContent := `version: '3'

vars:
  BINARY_NAME: test-app
  BUILD_DIR: dist

tasks:
  build:
    desc: Build for the current platform
    cmds:
      - go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}} .

  build-linux-amd64:
    desc: Build for linux/amd64
    cmds:
      - go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}}-linux-amd64 .
    env:
      GOOS: linux
      GOARCH: amd64

  build-windows-amd64:
    desc: Build for windows/amd64
    cmds:
      - go build -o {{.BUILD_DIR}}/{{.BINARY_NAME}}-windows-amd64.exe .
    env:
      GOOS: windows
      GOARCH: amd64

  release:
    desc: Build release binaries
    deps: [build-linux-amd64, build-windows-amd64]
    cmds:
      - echo "Release binaries created"
`

	taskfilePath := filepath.Join(tempDir, "Taskfile.yml")
	err = os.WriteFile(taskfilePath, []byte(taskfileContent), 0644)
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier Taskfile de test: %v", err)
	}

	// Créer un gestionnaire de cibles de build
	manager := NewBuildTargetManager(tempDir)

	// Test 1: Lister les cibles de build
	t.Run("ListBuildTargets", func(t *testing.T) {
		targets, err := manager.ListBuildTargets()
		if err != nil {
			t.Fatalf("ListBuildTargets a échoué: %v", err)
		}

		if len(targets) != 2 {
			t.Errorf("Nombre incorrect de cibles: %d, attendu: 2", len(targets))
		}

		// Vérifier les cibles attendues
		foundLinuxAmd64 := false
		foundWindowsAmd64 := false

		for _, target := range targets {
			if target.OS == "linux" && target.Arch == "amd64" {
				foundLinuxAmd64 = true
			}
			if target.OS == "windows" && target.Arch == "amd64" {
				foundWindowsAmd64 = true
			}
		}

		if !foundLinuxAmd64 {
			t.Error("Cible linux/amd64 non trouvée")
		}
		if !foundWindowsAmd64 {
			t.Error("Cible windows/amd64 non trouvée")
		}
	})

	// Test 2: Ajouter une cible de build
	t.Run("AddBuildTarget", func(t *testing.T) {
		newTarget := model.BuildTarget{
			OS:   "darwin",
			Arch: "arm64",
		}

		err := manager.AddBuildTarget(newTarget)
		if err != nil {
			t.Fatalf("AddBuildTarget a échoué: %v", err)
		}

		// Vérifier que la cible a été ajoutée
		content, err := os.ReadFile(taskfilePath)
		if err != nil {
			t.Fatalf("Impossible de lire le Taskfile: %v", err)
		}

		if !strings.Contains(string(content), "build-darwin-arm64") {
			t.Error("La cible darwin/arm64 n'a pas été ajoutée au Taskfile")
		}

		targets, _ := manager.ListBuildTargets()
		found := false
		for _, target := range targets {
			if target.OS == "darwin" && target.Arch == "arm64" {
				found = true
				break
			}
		}
		if !found {
			t.Error("La cible darwin/arm64 n'est pas listée après ajout")
		}
	})

	// Test 3: Supprimer une cible de build
	t.Run("RemoveBuildTarget", func(t *testing.T) {
		targetToRemove := model.BuildTarget{
			OS:   "windows",
			Arch: "amd64",
		}

		err := manager.RemoveBuildTarget(targetToRemove)
		if err != nil {
			t.Fatalf("RemoveBuildTarget a échoué: %v", err)
		}

		// Vérifier que la cible a été supprimée
		content, err := os.ReadFile(taskfilePath)
		if err != nil {
			t.Fatalf("Impossible de lire le Taskfile: %v", err)
		}

		fileContent := string(content)
		t.Logf("Contenu du Taskfile après suppression:\n%s", fileContent)

		if strings.Contains(fileContent, "build-windows-amd64:") {
			t.Error("La cible windows/amd64 n'a pas été supprimée du Taskfile")
		}

		targets, _ := manager.ListBuildTargets()
		for _, target := range targets {
			if target.OS == "windows" && target.Arch == "amd64" {
				t.Error("La cible windows/amd64 est encore listée après suppression")
			}
		}
	})
}
