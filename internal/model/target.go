package model

// BuildTarget represents a build target with OS and architecture
type BuildTarget struct {
	// OS is the target operating system (linux, darwin, windows)
	OS string
	// Arch is the target architecture (amd64, arm64, etc.)
	Arch string
}

// DefaultBuildTargets returns a set of common build targets
func DefaultBuildTargets() []BuildTarget {
	return []BuildTarget{
		{OS: "linux", Arch: "amd64"},
		{OS: "linux", Arch: "arm64"},
		{OS: "darwin", Arch: "amd64"},
		{OS: "darwin", Arch: "arm64"},
		{OS: "windows", Arch: "amd64"},
	}
}

// ValidOS returns whether the provided OS is valid
func ValidOS(os string) bool {
	validOS := []string{"linux", "darwin", "windows"}
	for _, valid := range validOS {
		if os == valid {
			return true
		}
	}
	return false
}

// ValidArch returns whether the provided architecture is valid
func ValidArch(arch string) bool {
	validArch := []string{"amd64", "arm64", "386"}
	for _, valid := range validArch {
		if arch == valid {
			return true
		}
	}
	return false
}

// TargetExists checks if a target exists in a list of targets
func TargetExists(target BuildTarget, targets []BuildTarget) bool {
	for _, t := range targets {
		if t.OS == target.OS && t.Arch == target.Arch {
			return true
		}
	}
	return false
}
