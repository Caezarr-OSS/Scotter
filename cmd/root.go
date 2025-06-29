// Package cmd implements command-line interface for Scotter
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scotter",
	Short: "Scotter is a scaffolding tool for Go projects",
	Long: `Scotter is a scaffolding tool that allows rapid generation of
project structures with integrated CI/CD workflows.

It supports multiple project types and CI providers.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands will be registered here
}
