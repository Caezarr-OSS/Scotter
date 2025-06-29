// Package plugin defines the core interfaces for Scotter's plugin system
package plugin

// LanguageProvider is the main interface for language plugins
type LanguageProvider interface {
	// Name returns the plugin name
	Name() string
	
	// SupportedProjectTypes returns project types supported by this language
	SupportedProjectTypes() []string
	
	// Initialize initializes a new project
	Initialize(projectName, projectType string, config map[string]interface{}) error
	
	// GenerateReleaseScript generates a release script
	GenerateReleaseScript(projectPath string, config map[string]interface{}) error
	
	// AddPlatform adds support for a new platform
	AddPlatform(projectPath, platform string) error
	
	// AddReleaseAsset adds support for a new release asset type
	AddReleaseAsset(projectPath, assetType string) error
	
	// IsSupportedPlatform checks if a platform is supported by this language
	IsSupportedPlatform(platform string) bool
	
	// IsSupportedArchitecture checks if an architecture is supported by this language
	IsSupportedArchitecture(arch string) bool
	
	// IsSupportedReleaseAsset checks if a release asset type is supported by this language
	IsSupportedReleaseAsset(assetType string) bool
	
	// GetSupportedPlatforms returns all supported platforms for this language
	GetSupportedPlatforms() []string
	
	// GetSupportedArchitectures returns all supported architectures for this language
	GetSupportedArchitectures() []string
	
	// GetSupportedReleaseAssets returns all supported release asset types for this language
	GetSupportedReleaseAssets() []string
}

// CIProvider is the interface for continuous integration plugins
type CIProvider interface {
	// Name returns the CI provider name
	Name() string
	
	// SupportedLanguages returns languages supported by this provider
	SupportedLanguages() []string
	
	// GenerateWorkflows generates CI workflows for a language and project type
	GenerateWorkflows(projectPath, language, projectType string, config map[string]interface{}) error
}

// PluginLoader handles the registration and management of plugins
type PluginLoader interface {
	// RegisterLanguageProvider registers a new language provider
	RegisterLanguageProvider(provider LanguageProvider)
	
	// RegisterCIProvider registers a new CI provider
	RegisterCIProvider(provider CIProvider)
	
	// GetLanguageProvider retrieves a language provider by name
	GetLanguageProvider(name string) (LanguageProvider, error)
	
	// GetCIProvider retrieves a CI provider by name
	GetCIProvider(name string) (CIProvider, error)
	
	// GetLanguageProviders returns all registered language providers
	GetLanguageProviders() []LanguageProvider
	
	// GetCIProviders returns all registered CI providers
	GetCIProviders() []CIProvider
}
