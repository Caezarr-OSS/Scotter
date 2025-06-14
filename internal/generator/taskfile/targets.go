package taskfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

const (
	taskfilePath = "Taskfile.yml"
)

// TaskfileContent represents the structure of a Taskfile.yml
type TaskfileContent struct {
	Version    string                 `json:"version"`
	Tasks      map[string]TaskDef     `json:"tasks"`
	Vars       map[string]interface{} `json:"vars,omitempty"`
	Includes   map[string]string      `json:"includes,omitempty"`
	Output     string                 `json:"output,omitempty"`
	Silent     bool                   `json:"silent,omitempty"`
	IncludeVars map[string]bool       `json:"includes_vars,omitempty"`
}

// TaskDef represents a task definition in a Taskfile
type TaskDef struct {
	Desc       string                 `json:"desc,omitempty"`
	Cmds       []interface{}          `json:"cmds,omitempty"`
	Deps       []string               `json:"deps,omitempty"`
	Vars       map[string]string      `json:"vars,omitempty"`
	Sources    []string               `json:"sources,omitempty"`
	Generates  []string               `json:"generates,omitempty"`
	Status     []string               `json:"status,omitempty"`
	Dir        string                 `json:"dir,omitempty"`
	Silent     bool                   `json:"silent,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Platforms  []string               `json:"platforms,omitempty"`
}

// BuildTargetManager handles build target operations
type BuildTargetManager struct {
	ProjectPath string
}

// NewBuildTargetManager creates a new build target manager
func NewBuildTargetManager(projectPath string) *BuildTargetManager {
	return &BuildTargetManager{
		ProjectPath: projectPath,
	}
}

// ReadTaskfile reads the existing Taskfile.yml
func (m *BuildTargetManager) ReadTaskfile() (string, error) {
	taskPath := filepath.Join(m.ProjectPath, taskfilePath)
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return "", fmt.Errorf("failed to read Taskfile: %w", err)
	}
	return string(content), nil
}

// WriteTaskfile writes the updated Taskfile.yml
func (m *BuildTargetManager) WriteTaskfile(content string) error {
	taskPath := filepath.Join(m.ProjectPath, taskfilePath)
	err := os.WriteFile(taskPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Taskfile: %w", err)
	}
	return nil
}

// ListBuildTargets lists all build targets in the Taskfile
func (m *BuildTargetManager) ListBuildTargets() ([]model.BuildTarget, error) {
	content, err := m.ReadTaskfile()
	if err != nil {
		return nil, err
	}
	
	// Extract build targets from the Taskfile
	targets := []model.BuildTarget{}
	
	// Look for build tasks with patterns like:
	// task: build-linux-amd64, build-darwin-arm64, etc.
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Check if the line defines a build task
		if strings.HasPrefix(line, "build-") {
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				taskName := strings.TrimSpace(parts[0])
				
				// Extract OS and arch from task name (format: build-os-arch)
				taskParts := strings.Split(taskName, "-")
				if len(taskParts) >= 3 {
					os := taskParts[1]
					arch := taskParts[2]
					
					if model.ValidOS(os) && model.ValidArch(arch) {
						targets = append(targets, model.BuildTarget{
							OS:   os,
							Arch: arch,
						})
					}
				}
			}
		}
	}
	
	return targets, nil
}

// AddBuildTarget adds a new build target to the Taskfile
func (m *BuildTargetManager) AddBuildTarget(target model.BuildTarget) error {
	// Validate target
	if !model.ValidOS(target.OS) || !model.ValidArch(target.Arch) {
		return fmt.Errorf("invalid build target: OS=%s, Arch=%s", target.OS, target.Arch)
	}
	
	// Check if target already exists
	existingTargets, err := m.ListBuildTargets()
	if err != nil {
		return err
	}
	
	for _, existing := range existingTargets {
		if existing.OS == target.OS && existing.Arch == target.Arch {
			return fmt.Errorf("build target already exists: OS=%s, Arch=%s", target.OS, target.Arch)
		}
	}
	
	// Read existing Taskfile
	content, err := m.ReadTaskfile()
	if err != nil {
		return err
	}
	
	// Add new build task
	taskName := fmt.Sprintf("build-%s-%s", target.OS, target.Arch)
	newTaskContent := fmt.Sprintf(`
  %s:
    desc: Build for %s/%s
    cmds:
      - go build -o {{.BINARY_NAME}}-%s-%s {{.MAIN_PACKAGE}}
    env:
      GOOS: %s
      GOARCH: %s
    generates:
      - "{{.BINARY_NAME}}-%s-%s{{exeExt}}"

`, 
    taskName, target.OS, target.Arch, 
    target.OS, target.Arch, 
    target.OS, target.Arch,
    target.OS, target.Arch)
	
	// Find the position to insert the new task
	tasksSection := strings.Index(content, "tasks:")
	if tasksSection == -1 {
		return fmt.Errorf("could not find tasks section in Taskfile")
	}
	
	// Insert the new task after the last task
	lines := strings.Split(content, "\n")
	insertPosition := -1
	
	for i, line := range lines {
		if strings.TrimSpace(line) == "tasks:" {
			insertPosition = i
			break
		}
	}
	
	if insertPosition == -1 {
		return fmt.Errorf("could not find position to insert new task")
	}
	
	// Update the build-all task to include the new target
	buildAllTaskIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "build-all:" {
			buildAllTaskIndex = i
			break
		}
	}
	
	// Insert the new task
	// Ensure the new task is properly separated from existing tasks
	// by analyzing the last line before insertion
	hasNewlineBefore := false
	if insertPosition > 0 && insertPosition < len(lines) && strings.TrimSpace(lines[insertPosition]) == "" {
		hasNewlineBefore = true
	}
	
	// If the new task doesn't already have a newline before it, add one
	if !hasNewlineBefore {
		newTaskContent = "\n" + newTaskContent
	}
	
	updatedLines := append(lines[:insertPosition+1], append([]string{newTaskContent}, lines[insertPosition+1:]...)...)
	
	// Update the build-all task if it exists
	if buildAllTaskIndex != -1 {
		buildAllDepsLine := -1
		for i := buildAllTaskIndex; i < len(updatedLines) && i < buildAllTaskIndex+10; i++ {
			if strings.TrimSpace(updatedLines[i]) == "deps:" {
				buildAllDepsLine = i
				break
			}
		}
		
		if buildAllDepsLine != -1 {
			// Add the new target to the deps list
			updatedLines = append(updatedLines[:buildAllDepsLine+1], append([]string{fmt.Sprintf("      - %s", taskName)}, updatedLines[buildAllDepsLine+1:]...)...)
		}
	}
	
	// Write the updated content
	err = m.WriteTaskfile(strings.Join(updatedLines, "\n"))
	if err != nil {
		return err
	}
	
	return nil
}

// RemoveBuildTarget removes a build target from the Taskfile
func (m *BuildTargetManager) RemoveBuildTarget(target model.BuildTarget) error {
	// Validate target
	if !model.ValidOS(target.OS) || !model.ValidArch(target.Arch) {
		return fmt.Errorf("invalid build target: OS=%s, Arch=%s", target.OS, target.Arch)
	}
	
	// Check if target exists
	existingTargets, err := m.ListBuildTargets()
	if err != nil {
		return err
	}
	
	targetExists := false
	for _, existing := range existingTargets {
		if existing.OS == target.OS && existing.Arch == target.Arch {
			targetExists = true
			break
		}
	}
	
	if !targetExists {
		return fmt.Errorf("build target not found: OS=%s, Arch=%s", target.OS, target.Arch)
	}
	
	// Read existing Taskfile
	content, err := m.ReadTaskfile()
	if err != nil {
		return err
	}
	
	// Définir le nom de la tâche à supprimer
	taskName := fmt.Sprintf("build-%s-%s", target.OS, target.Arch)
	lines := strings.Split(content, "\n")
	
	// Nouvelle approche pour supprimer la tâche complète
	updatedLines := []string{}
	skipping := false
	indentation := ""
	
	// Parcourir chaque ligne pour construire le nouveau contenu
	for i, line := range lines {
		// Si la ligne correspond au début de la tâche à supprimer
		if strings.TrimSpace(line) == taskName+":" {
			skipping = true
			// Déterminer l'indentation attendue pour les sous-tâches
			if i > 0 && len(line) > 0 {
				// Calculer l'indentation de cette ligne
				for j := 0; j < len(line); j++ {
					if line[j] != ' ' {
						indentation = line[:j]
						break
					}
				}
			}
			continue
		}
		
		// Si on est en mode "saut" et qu'on rencontre une nouvelle tâche au même niveau ou supérieur
		if skipping {
			// Si la ligne n'est pas vide et n'est pas indentée plus que la tâche à supprimer
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" {
				// Si la ligne commence par une indentation inférieure ou égale,
				// c'est une nouvelle tâche ou une section au même niveau
				if !strings.HasPrefix(line, indentation+" ") {
					skipping = false
				}
			}
		}
		
		// Si on n'est plus en mode "saut", ajouter la ligne
		if !skipping {
			updatedLines = append(updatedLines, line)
		}
	}
	
	// Update the release task to remove the target from dependencies
	for i := 0; i < len(updatedLines); i++ {
		// Trouver la tâche release
		if strings.TrimSpace(updatedLines[i]) == "release:" {
			// Chercher sa liste de dépendances
			for j := i + 1; j < len(updatedLines) && j < i + 10; j++ {
				line := updatedLines[j]
				trimmedLine := strings.TrimSpace(line)
				
				// Si on trouve une ligne de dépendances qui contient notre cible
				if strings.HasPrefix(trimmedLine, "deps:") && strings.Contains(trimmedLine, taskName) {
					// Parser les dépendances
					start := strings.Index(trimmedLine, "[") + 1
					end := strings.Index(trimmedLine, "]")
					
					if start > 0 && end > start {
						depsStr := trimmedLine[start:end]
						deps := strings.Split(depsStr, ",")
						
						// Filtrer la dépendance à supprimer
						newDeps := []string{}
						for _, dep := range deps {
							dep = strings.TrimSpace(dep)
							if dep != taskName {
								newDeps = append(newDeps, dep)
							}
						}
						
						// Reconstruire la ligne
						newDepsStr := strings.Join(newDeps, ", ")
						newLine := trimmedLine[:start-1] + "[" + newDepsStr + "]" + trimmedLine[end+1:]
						
						// Recréer la ligne avec la même indentation
						indent := ""
						for k := 0; k < len(line); k++ {
							if line[k] != ' ' {
								indent = line[:k]
								break
							}
						}
						
						updatedLines[j] = indent + newLine
					}
					break
				}
			}
			break
		}
	}
	
	// Write the updated content
	err = m.WriteTaskfile(strings.Join(updatedLines, "\n"))
	if err != nil {
		return err
	}
	
	return nil
}
