package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/caezarr-oss/scotter/pkg/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove features from a project",
	Long:  `Remove features such as platforms, architectures, or release assets from a project`,
}

var removePlatformCmd = &cobra.Command{
	Use:   "platform [platform]",
	Short: "Remove support for a platform",
	Long:  `Remove support for a specific platform (linux, darwin, windows)`,
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
		
		// Remove platform from configuration
		if err := configManager.RemovePlatform(platformName); err != nil {
			return fmt.Errorf("unable to remove platform: %w", err)
		}
		
		// Platform changes only affect the configuration file
		// No need to update language provider
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Platform '%s' successfully removed from the project\n", platformName)
		return nil
	},
}

var removeReleaseAssetCmd = &cobra.Command{
	Use:   "release-asset [type]",
	Short: "Remove a release asset type",
	Long:  `Remove support for a release asset type (checksum, sbom, archive)`,
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
		
		// Remove asset type from configuration
		if err := configManager.RemoveReleaseAsset(assetType); err != nil {
			return fmt.Errorf("unable to remove release asset type: %w", err)
		}
		
		// Asset type changes only affect the configuration file
		// No need to update language provider
		
		// Save configuration
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("unable to save configuration: %w", err)
		}
		
		fmt.Printf("Release asset type '%s' successfully removed from the project\n", assetType)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.AddCommand(removePlatformCmd)
	removeCmd.AddCommand(removeReleaseAssetCmd)
}
