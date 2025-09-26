package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

// String returns the version string
func String() string {
	return fmt.Sprintf("v%s", Version)
}

// Full returns the full version information
func Full() string {
	return fmt.Sprintf(`version %s
Git commit: %s
Build date: %s
Go version: %s`, Version, GitCommit, BuildDate, GoVersion)
}

// Short returns the short version string
func Short() string {
	return Version
}

// Info returns version info for a specific extension
func Info(extensionName string) string {
	return fmt.Sprintf("%s %s", extensionName, String())
}

// InfoFull returns full version info for a specific extension
func InfoFull(extensionName string) string {
	return fmt.Sprintf(`%s %s
Git commit: %s
Build date: %s
Go version: %s`, extensionName, String(), GitCommit, BuildDate, GoVersion)
}
