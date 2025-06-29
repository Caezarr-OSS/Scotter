package golang

// SupportedGoPlatforms contains all supported Go platforms (GOOS values)
var SupportedGoPlatforms = []string{
	"aix",
	"android",
	"darwin",
	"dragonfly",
	"freebsd",
	"illumos",
	"ios",
	"js",
	"linux",
	"netbsd",
	"openbsd",
	"plan9",
	"solaris",
	"wasip1", // WebAssembly System Interface (WASI) Preview 1
	"windows",
}

// SupportedGoArchitectures contains all supported Go architectures (GOARCH values)
var SupportedGoArchitectures = []string{
	"386",
	"amd64",
	"arm",
	"arm64",
	"loong64",
	"mips",
	"mips64",
	"mips64le",
	"mipsle",
	"ppc64",
	"ppc64le",
	"riscv64",
	"s390x",
	"wasm",
}

// SupportedGoReleaseAssets contains all supported release asset types for Go projects
var SupportedGoReleaseAssets = []string{
	"checksum", // SHA-256 checksums for binaries
	"sbom",     // Software Bill of Materials
	"archive",  // Compressed archives (tar.gz, zip)
	// Note: binary and source assets aren't currently implemented in AddReleaseAsset
}

// Note: The contains function is already defined in provider.go
