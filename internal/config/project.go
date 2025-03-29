package config

import (
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// GetDefaultConfig retourne une configuration par défaut
// Cette fonction existe pour rétrocompatibilité
func GetDefaultConfig() *model.Config {
	return model.NewDefaultConfig()
}

// ValidateConfig vérifie si la configuration est valide
// Cette fonction existe pour rétrocompatibilité
func ValidateConfig(config *model.Config) error {
	return model.ValidateConfig(config)
}
