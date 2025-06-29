// Package config provides functionality for managing Scotter project configuration
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caezarr-oss/scotter/pkg/plugin"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigFile is the default name of the Scotter configuration file
	DefaultConfigFile = ".scotter.yaml"
)

// Config represents the Scotter project configuration
type Config struct {
	ProjectName    string   `yaml:"project_name"`
	ProjectType    string   `yaml:"project_type"`
	Language       string   `yaml:"language"`
	Platforms      []string `yaml:"platforms"`
	Architectures  []string `yaml:"architectures"`
	ReleaseAssets  []string `yaml:"release_assets"`
	CIProvider     string   `yaml:"ci_provider"`
	ExtraConfig    map[string]interface{} `yaml:"extra_config,omitempty"`
}

// Manager handles configuration operations
type Manager struct {
	ConfigPath string
	Config     *Config
}

// NewManager creates a new configuration manager
func NewManager(projectPath string) *Manager {
	return &Manager{
		ConfigPath: filepath.Join(projectPath, DefaultConfigFile),
		Config: &Config{
			Platforms:     []string{},
			Architectures: []string{},
			ReleaseAssets: []string{},
			ExtraConfig:   make(map[string]interface{}),
		},
	}
}

// Load loads configuration from file
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, m.Config)
}

// Save saves configuration to file
func (m *Manager) Save() error {
	data, err := yaml.Marshal(m.Config)
	if err != nil {
		return err
	}

	return os.WriteFile(m.ConfigPath, data, 0644)
}

// AddPlatform adds a new platform if not already present
func (m *Manager) AddPlatform(platform string, langProvider plugin.LanguageProvider) error {
	// Validate that the platform is supported by the language provider
	if !langProvider.IsSupportedPlatform(platform) {
		return fmt.Errorf("platform '%s' is not supported by %s", platform, langProvider.Name())
	}

	// Check if platform already exists
	for _, p := range m.Config.Platforms {
		if p == platform {
			return fmt.Errorf("platform '%s' already exists", platform)
		}
	}
	
	m.Config.Platforms = append(m.Config.Platforms, platform)
	return nil
}

// RemovePlatform removes a platform
func (m *Manager) RemovePlatform(platform string) error {
	for i, p := range m.Config.Platforms {
		if p == platform {
			m.Config.Platforms = append(m.Config.Platforms[:i], m.Config.Platforms[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("platform '%s' not found", platform)
}

// AddArchitecture adds a new architecture if not already present
func (m *Manager) AddArchitecture(arch string, langProvider plugin.LanguageProvider) error {
	// Validate that the architecture is supported by the language provider
	if !langProvider.IsSupportedArchitecture(arch) {
		return fmt.Errorf("architecture '%s' is not supported by %s", arch, langProvider.Name())
	}

	// Check if architecture already exists
	for _, a := range m.Config.Architectures {
		if a == arch {
			return fmt.Errorf("architecture '%s' already exists", arch)
		}
	}
	
	m.Config.Architectures = append(m.Config.Architectures, arch)
	return nil
}

// RemoveArchitecture removes an architecture
func (m *Manager) RemoveArchitecture(arch string) error {
	for i, a := range m.Config.Architectures {
		if a == arch {
			m.Config.Architectures = append(m.Config.Architectures[:i], m.Config.Architectures[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("architecture '%s' not found", arch)
}

// AddReleaseAsset adds a new release asset type if not already present
func (m *Manager) AddReleaseAsset(assetType string, langProvider plugin.LanguageProvider) error {
	// Validate that the asset type is supported by the language provider
	if !langProvider.IsSupportedReleaseAsset(assetType) {
		return fmt.Errorf("release asset type '%s' is not supported by %s", assetType, langProvider.Name())
	}

	// Check if asset type already exists
	for _, a := range m.Config.ReleaseAssets {
		if a == assetType {
			return fmt.Errorf("release asset type '%s' already exists", assetType)
		}
	}
	
	m.Config.ReleaseAssets = append(m.Config.ReleaseAssets, assetType)
	return nil
}

// RemoveReleaseAsset removes a release asset type
func (m *Manager) RemoveReleaseAsset(assetType string) error {
	for i, a := range m.Config.ReleaseAssets {
		if a == assetType {
			m.Config.ReleaseAssets = append(m.Config.ReleaseAssets[:i], m.Config.ReleaseAssets[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("release asset '%s' not found", assetType)
}

// SetExtraConfig sets a value in the extra configuration
func (m *Manager) SetExtraConfig(key string, value interface{}) {
	m.Config.ExtraConfig[key] = value
}

// GetExtraConfig gets a value from the extra configuration
func (m *Manager) GetExtraConfig(key string) (interface{}, bool) {
	value, ok := m.Config.ExtraConfig[key]
	return value, ok
}
