package cmd

import (
	"github.com/caezarr-oss/scotter/internal/ci/github"
	golangplugin "github.com/caezarr-oss/scotter/internal/cmd/golang"
	"github.com/caezarr-oss/scotter/pkg/plugin"
)

// registerPlugins registers all available plugins with the plugin loader
func registerPlugins(loader *plugin.DefaultPluginLoader) {
	// Register language providers
	loader.RegisterLanguageProvider(golangplugin.NewGoLanguageProvider())
	
	// Register CI providers
	loader.RegisterCIProvider(github.NewGitHubProvider())
}
