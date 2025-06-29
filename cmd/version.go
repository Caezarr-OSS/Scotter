package cmd

import (
	"fmt"

	"github.com/caezarr-oss/scotter/pkg/version"
	"github.com/spf13/cobra"
)

var (
	showDetailedVersion bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information of Scotter",
	Long:  `Display the version, commit, build date and other information about the Scotter binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showDetailedVersion {
			// Display detailed information
			buildInfo := version.BuildInfo()
			fmt.Println("Scotter - Detailed Build Information:")
			fmt.Println("===================================")
			for k, v := range buildInfo {
				fmt.Printf("%-12s: %s\n", k, v)
			}
		} else {
			// Afficher juste la version
			fmt.Printf("Scotter %s\n", version.Version())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	
	// Add a flag to display detailed information
	versionCmd.Flags().BoolVarP(&showDetailedVersion, "detailed", "d", false, "Show detailed version information")
}
