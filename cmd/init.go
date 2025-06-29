package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caezarr-oss/scotter/pkg/config"
	"github.com/caezarr-oss/scotter/pkg/plugin"
	"github.com/spf13/cobra"
)

var (
	projectType string
	language    string
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project",
	Long:  `Initialize a new project with the specified structure and settings`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectPath, err := filepath.Abs(projectName)
		if err != nil {
			return fmt.Errorf("unable to resolve project path: %w", err)
		}

		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("unable to create project directory: %w", err)
		}

		// Initialize configuration
		configManager := config.NewManager(projectPath)
		configManager.Config.ProjectName = projectName
		configManager.Config.ProjectType = projectType
		configManager.Config.Language = language
		
		// Get language provider
		pluginLoader := plugin.NewPluginLoader()
		registerPlugins(pluginLoader)
		
		langProvider, err := pluginLoader.GetLanguageProvider(language)
		if err != nil {
			return fmt.Errorf("language provider not available: %w", err)
		}
		
		// Add default platforms
		if err := configManager.AddPlatform("linux", langProvider); err != nil {
			// Just log the error and continue, as this is initial setup
			fmt.Printf("Warning: Failed to add platform 'linux': %s\n", err)
		}
		if err := configManager.AddPlatform("darwin", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add platform 'darwin': %s\n", err)
		}
		if err := configManager.AddPlatform("windows", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add platform 'windows': %s\n", err)
		}
		
		// Add default architectures
		if err := configManager.AddArchitecture("amd64", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add architecture 'amd64': %s\n", err)
		}
		if err := configManager.AddArchitecture("arm64", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add architecture 'arm64': %s\n", err)
		}
		
		// Add default release assets
		if err := configManager.AddReleaseAsset("checksum", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add release asset 'checksum': %s\n", err)
		}
		if err := configManager.AddReleaseAsset("sbom", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add release asset 'sbom': %s\n", err)
		}
		if err := configManager.AddReleaseAsset("archive", langProvider); err != nil {
			fmt.Printf("Warning: Failed to add release asset 'archive': %s\n", err)
		}
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}

		// Language provider already obtained above

		// Initialize project
		if err := langProvider.Initialize(projectName, projectType, configManager.Config.ExtraConfig); err != nil {
			return fmt.Errorf("failed to initialize project: %w", err)
		}

		fmt.Printf("Project '%s' successfully initialized with type '%s' using language '%s'\n", 
			projectName, projectType, language)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	
	// Define flags
	initCmd.Flags().StringVar(&projectType, "type", "default", "Project type (cli, api, library, default)")
	initCmd.Flags().StringVar(&language, "language", "go", "Programming language")
}
