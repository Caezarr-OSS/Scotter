package config

import (
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// GetDefaultConfig returns a default configuration
// This function exists for backward compatibility
func GetDefaultConfig() *model.Config {
	return model.NewConfig()
}

// ValidateConfig checks if the configuration is valid
// This function exists for backward compatibility
func ValidateConfig(config *model.Config) error {
	return model.ValidateConfig(config)
}
