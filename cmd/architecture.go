package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/caezarr-oss/scotter/pkg/config"
	"github.com/caezarr-oss/scotter/pkg/plugin"
	"github.com/spf13/cobra"
)

var addArchitectureCmd = &cobra.Command{
	Use:   "architecture [arch]",
	Short: "Add support for an architecture",
	Long:  `Add support for a specific architecture (amd64, arm64, etc.)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		archName := args[0]
		
		// Get current directory as project path
		projectPath, err := filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("unable to resolve project path: %w", err)
		}
		
		// Load configuration
		configManager := config.NewManager(projectPath)
		if err := configManager.Load(); err != nil {
			return fmt.Errorf("unable to load configuration: %w", err)
		}
		
		// Get language provider
		pluginLoader := plugin.NewPluginLoader()
		registerPlugins(pluginLoader)
		
		langProvider, err := pluginLoader.GetLanguageProvider(configManager.Config.Language)
		if err != nil {
			return fmt.Errorf("language provider not available: %w", err)
		}

		// Add architecture to configuration
		if err := configManager.AddArchitecture(archName, langProvider); err != nil {
			return fmt.Errorf("unable to add architecture: %w", err)
		}
		
		// Architecture changes only affect the configuration file
		// No need to update language provider
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Architecture '%s' successfully added to the project\n", archName)
		return nil
	},
}

var removeArchitectureCmd = &cobra.Command{
	Use:   "architecture [arch]",
	Short: "Remove support for an architecture",
	Long:  `Remove support for a specific architecture (amd64, arm64, etc.)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		archName := args[0]
		
		// Get current directory as project path
		projectPath, err := filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("unable to resolve project path: %w", err)
		}
		
		// Load configuration
		configManager := config.NewManager(projectPath)
		if err := configManager.Load(); err != nil {
			return fmt.Errorf("unable to load configuration: %w", err)
		}
		
		// Remove architecture from configuration
		if err := configManager.RemoveArchitecture(archName); err != nil {
			return fmt.Errorf("unable to remove architecture: %w", err)
		}
		
		// Architecture changes only affect the configuration file
		// No need to update language provider
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Architecture '%s' successfully removed from the project\n", archName)
		return nil
	},
}

func init() {
	addCmd.AddCommand(addArchitectureCmd)
	removeCmd.AddCommand(removeArchitectureCmd)
}
