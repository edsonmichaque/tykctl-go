package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	result := String()
	expected := "v" + Version
	
	if result != expected {
		t.Errorf("String() = %s, expected %s", result, expected)
	}
	
	// Should start with 'v'
	if !strings.HasPrefix(result, "v") {
		t.Error("String() should start with 'v'")
	}
}

func TestShort(t *testing.T) {
	result := Short()
	
	if result != Version {
		t.Errorf("Short() = %s, expected %s", result, Version)
	}
}

func TestFull(t *testing.T) {
	result := Full()
	
	// Should contain all version components
	if !strings.Contains(result, Version) {
		t.Error("Full() should contain Version")
	}
	
	if !strings.Contains(result, GitCommit) {
		t.Error("Full() should contain GitCommit")
	}
	
	if !strings.Contains(result, BuildDate) {
		t.Error("Full() should contain BuildDate")
	}
	
	if !strings.Contains(result, GoVersion) {
		t.Error("Full() should contain GoVersion")
	}
	
	// Should contain labels
	if !strings.Contains(result, "version") {
		t.Error("Full() should contain 'version' label")
	}
	
	if !strings.Contains(result, "Git commit:") {
		t.Error("Full() should contain 'Git commit:' label")
	}
	
	if !strings.Contains(result, "Build date:") {
		t.Error("Full() should contain 'Build date:' label")
	}
	
	if !strings.Contains(result, "Go version:") {
		t.Error("Full() should contain 'Go version:' label")
	}
}

func TestInfo(t *testing.T) {
	extensionName := "test-extension"
	result := Info(extensionName)
	expected := extensionName + " " + String()
	
	if result != expected {
		t.Errorf("Info() = %s, expected %s", result, expected)
	}
	
	// Should contain extension name
	if !strings.Contains(result, extensionName) {
		t.Error("Info() should contain extension name")
	}
	
	// Should contain version string
	if !strings.Contains(result, String()) {
		t.Error("Info() should contain version string")
	}
}

func TestInfoFull(t *testing.T) {
	extensionName := "test-extension"
	result := InfoFull(extensionName)
	
	// Should contain extension name
	if !strings.Contains(result, extensionName) {
		t.Error("InfoFull() should contain extension name")
	}
	
	// Should contain version string
	if !strings.Contains(result, String()) {
		t.Error("InfoFull() should contain version string")
	}
	
	// Should contain all version components
	if !strings.Contains(result, Version) {
		t.Error("InfoFull() should contain Version")
	}
	
	if !strings.Contains(result, GitCommit) {
		t.Error("InfoFull() should contain GitCommit")
	}
	
	if !strings.Contains(result, BuildDate) {
		t.Error("InfoFull() should contain BuildDate")
	}
	
	if !strings.Contains(result, GoVersion) {
		t.Error("InfoFull() should contain GoVersion")
	}
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are set
	if Version == "" {
		t.Error("Version should not be empty")
	}
	
	if GitCommit == "" {
		t.Error("GitCommit should not be empty")
	}
	
	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}
	
	if GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that version follows expected format
	if !strings.Contains(Version, ".") {
		t.Error("Version should contain dots (e.g., 1.0.0)")
	}
	
	// Test that String() follows expected format
	versionStr := String()
	if !strings.HasPrefix(versionStr, "v") {
		t.Error("String() should start with 'v'")
	}
}

func TestInfoWithEmptyExtensionName(t *testing.T) {
	result := Info("")
	expected := " " + String() // Should still work with empty name
	
	if result != expected {
		t.Errorf("Info(\"\") = %s, expected %s", result, expected)
	}
}

func TestInfoFullWithEmptyExtensionName(t *testing.T) {
	result := InfoFull("")
	
	// Should still contain version information
	if !strings.Contains(result, String()) {
		t.Error("InfoFull(\"\") should contain version string")
	}
	
	if !strings.Contains(result, Version) {
		t.Error("InfoFull(\"\") should contain Version")
	}
}

func TestConsistency(t *testing.T) {
	// Test that all functions return consistent results
	versionStr := String()
	shortVersion := Short()
	fullVersion := Full()
	infoVersion := Info("test")
	infoFullVersion := InfoFull("test")
	
	// All should contain the base version
	if !strings.Contains(versionStr, Version) {
		t.Error("String() should contain Version")
	}
	
	if !strings.Contains(shortVersion, Version) {
		t.Error("Short() should contain Version")
	}
	
	if !strings.Contains(fullVersion, Version) {
		t.Error("Full() should contain Version")
	}
	
	if !strings.Contains(infoVersion, Version) {
		t.Error("Info() should contain Version")
	}
	
	if !strings.Contains(infoFullVersion, Version) {
		t.Error("InfoFull() should contain Version")
	}
}

// Benchmark tests
func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = String()
	}
}

func BenchmarkShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Short()
	}
}

func BenchmarkFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Full()
	}
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Info("test-extension")
	}
}

func BenchmarkInfoFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = InfoFull("test-extension")
	}
}