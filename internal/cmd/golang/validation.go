package golang

// IsSupportedPlatform checks if a platform is supported by Go
func (p *GoLanguageProvider) IsSupportedPlatform(platform string) bool {
	return contains(SupportedGoPlatforms, platform)
}

// IsSupportedArchitecture checks if an architecture is supported by Go
func (p *GoLanguageProvider) IsSupportedArchitecture(arch string) bool {
	return contains(SupportedGoArchitectures, arch)
}

// IsSupportedReleaseAsset checks if a release asset type is supported for Go projects
func (p *GoLanguageProvider) IsSupportedReleaseAsset(assetType string) bool {
	return contains(SupportedGoReleaseAssets, assetType)
}

// GetSupportedPlatforms returns all supported platforms for Go
func (p *GoLanguageProvider) GetSupportedPlatforms() []string {
	// Return a copy to prevent modification of the original slice
	platforms := make([]string, len(SupportedGoPlatforms))
	copy(platforms, SupportedGoPlatforms)
	return platforms
}

// GetSupportedArchitectures returns all supported architectures for Go
func (p *GoLanguageProvider) GetSupportedArchitectures() []string {
	// Return a copy to prevent modification of the original slice
	architectures := make([]string, len(SupportedGoArchitectures))
	copy(architectures, SupportedGoArchitectures)
	return architectures
}

// GetSupportedReleaseAssets returns all supported release asset types for Go projects
func (p *GoLanguageProvider) GetSupportedReleaseAssets() []string {
	// Return a copy to prevent modification of the original slice
	assets := make([]string, len(SupportedGoReleaseAssets))
	copy(assets, SupportedGoReleaseAssets)
	return assets
}
