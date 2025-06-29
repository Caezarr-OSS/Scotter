package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/caezarr-oss/scotter/pkg/config"
	"github.com/caezarr-oss/scotter/pkg/plugin"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add features to a project",
	Long:  `Add features such as CI providers, platforms, or release assets to a project`,
}

var addCICmd = &cobra.Command{
	Use:   "ci [provider]",
	Short: "Add CI workflows to a project",
	Long:  `Add CI workflows from a specific provider to a project`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName := args[0]
		
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
		
		// Get CI provider
		pluginLoader := plugin.NewPluginLoader()
		registerPlugins(pluginLoader)
		
		ciProvider, err := pluginLoader.GetCIProvider(providerName)
		if err != nil {
			return fmt.Errorf("CI provider not available: %w", err)
		}
		
		// Check if language is supported by the provider
		language := configManager.Config.Language
		projectType := configManager.Config.ProjectType
		
		supported := false
		for _, l := range ciProvider.SupportedLanguages() {
			if l == language {
				supported = true
				break
			}
		}
		
		if !supported {
			return fmt.Errorf("language '%s' is not supported by CI provider '%s'", 
				language, providerName)
		}
		
		// Generate workflows
		if err := ciProvider.GenerateWorkflows(projectPath, language, projectType, 
			configManager.Config.ExtraConfig); err != nil {
			return fmt.Errorf("failed to generate workflows: %w", err)
		}
		
		// Get the language provider to generate release script (GoReleaser config)
		langProvider, err := pluginLoader.GetLanguageProvider(language)
		if err != nil {
			return fmt.Errorf("language provider not available: %w", err)
		}
		
		// Generate release script if applicable for this language
		config := map[string]interface{}{
			"project_type": projectType,
		}
		
		// Try to generate the release script, but don't fail if it already exists
		err = langProvider.GenerateReleaseScript(projectPath, config)
		if err != nil && err.Error() != ".goreleaser.yaml already exists" {
			fmt.Printf("Warning: Could not generate release script: %v\n", err)
		}
		
		// Update configuration
		configManager.Config.CIProvider = providerName
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("CI workflows for provider '%s' successfully added to the project\n", 
			providerName)
		return nil
	},
}

var addPlatformCmd = &cobra.Command{
	Use:   "platform [platform]",
	Short: "Add support for a platform",
	Long:  `Add support for a specific platform (linux, darwin, windows)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		platformName := args[0]
		
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

		// Add platform to configuration
		if err := configManager.AddPlatform(platformName, langProvider); err != nil {
			return fmt.Errorf("unable to add platform: %w", err)
		}
		
		// Add platform to project
		if err := langProvider.AddPlatform(projectPath, platformName); err != nil {
			return fmt.Errorf("failed to add platform: %w", err)
		}
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Platform '%s' successfully added to the project\n", platformName)
		return nil
	},
}

var addReleaseAssetCmd = &cobra.Command{
	Use:   "release-asset [type]",
	Short: "Add a release asset type",
	Long:  `Add support for a release asset type (checksum, sbom, archive)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assetType := args[0]
		
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

		// Add asset type to configuration with validation
		if err := configManager.AddReleaseAsset(assetType, langProvider); err != nil {
			return fmt.Errorf("unable to add release asset: %w", err)
		}
		
		// Add release asset to project
		if err := langProvider.AddReleaseAsset(projectPath, assetType); err != nil {
			return fmt.Errorf("failed to add release asset: %w", err)
		}
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Release asset type '%s' successfully added to the project\n", assetType)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addCICmd)
	addCmd.AddCommand(addPlatformCmd)
	addCmd.AddCommand(addReleaseAssetCmd)
}
