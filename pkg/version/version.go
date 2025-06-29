package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

var (
	// These variables will be set during compilation via ldflags
	// If they are not defined, default values will be used
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtBy   = "unknown"
	timestamp = time.Now().Format(time.RFC3339)
)

// Version returns the complete version information
func Version() string {
	if version == "dev" {
		// Local development version, use the timestamp
		return fmt.Sprintf("dev-%s", timestamp)
	}
	// Version officielle de release
	return fmt.Sprintf("%s (commit %s, built %s by %s)", version, commit, date, builtBy)
}

// Short retourne une version courte de la version
func Short() string {
	if version == "dev" {
		return fmt.Sprintf("dev-%s", timestamp)
	}
	return version
}

// BuildInfo returns detailed information about the build
func BuildInfo() map[string]string {
	info := map[string]string{
		"version":   version,
		"commit":    commit,
		"built_at":  date,
		"built_by":  builtBy,
		"timestamp": timestamp,
	}

	// Try to get additional information via the runtime/debug module
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" && commit == "none" {
				info["commit"] = setting.Value
			}
			if setting.Key == "vcs.time" && date == "unknown" {
				info["built_at"] = setting.Value
			}
		}
	}

	return info
}
