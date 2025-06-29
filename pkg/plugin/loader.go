package plugin

import (
	"fmt"
	"sync"
)

// DefaultPluginLoader is the default implementation of PluginLoader
type DefaultPluginLoader struct {
	languageProviders map[string]LanguageProvider
	ciProviders       map[string]CIProvider
	mu                sync.RWMutex
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader() *DefaultPluginLoader {
	return &DefaultPluginLoader{
		languageProviders: make(map[string]LanguageProvider),
		ciProviders:       make(map[string]CIProvider),
	}
}

// RegisterLanguageProvider registers a new language provider
func (l *DefaultPluginLoader) RegisterLanguageProvider(provider LanguageProvider) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.languageProviders[provider.Name()] = provider
}

// RegisterCIProvider registers a new CI provider
func (l *DefaultPluginLoader) RegisterCIProvider(provider CIProvider) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ciProviders[provider.Name()] = provider
}

// GetLanguageProvider retrieves a language provider by name
func (l *DefaultPluginLoader) GetLanguageProvider(name string) (LanguageProvider, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	provider, ok := l.languageProviders[name]
	if !ok {
		return nil, fmt.Errorf("language provider '%s' not found", name)
	}
	return provider, nil
}

// GetCIProvider retrieves a CI provider by name
func (l *DefaultPluginLoader) GetCIProvider(name string) (CIProvider, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	provider, ok := l.ciProviders[name]
	if !ok {
		return nil, fmt.Errorf("CI provider '%s' not found", name)
	}
	return provider, nil
}

// GetLanguageProviders returns all registered language providers
func (l *DefaultPluginLoader) GetLanguageProviders() []LanguageProvider {
	l.mu.RLock()
	defer l.mu.RUnlock()
	providers := make([]LanguageProvider, 0, len(l.languageProviders))
	for _, provider := range l.languageProviders {
		providers = append(providers, provider)
	}
	return providers
}

// GetCIProviders returns all registered CI providers
func (l *DefaultPluginLoader) GetCIProviders() []CIProvider {
	l.mu.RLock()
	defer l.mu.RUnlock()
	providers := make([]CIProvider, 0, len(l.ciProviders))
	for _, provider := range l.ciProviders {
		providers = append(providers, provider)
	}
	return providers
}
