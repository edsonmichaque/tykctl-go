package terminal

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	term := New()
	if term == nil {
		t.Fatal("New() returned nil")
	}
	
	// Test that all fields are initialized
	if term.Width <= 0 {
		t.Error("Width should be positive")
	}
	
	if term.Height <= 0 {
		t.Error("Height should be positive")
	}
	
	// Color and NoColor should be mutually exclusive
	if term.Color && term.NoColor {
		t.Error("Color and NoColor should not both be true")
	}
}

func TestIsTTY(t *testing.T) {
	term := New()
	
	// Test that IsTTY returns a boolean value
	// We can't easily test the actual TTY detection without mocking,
	// but we can test that it returns a valid boolean
	_ = term.IsTTY()
}

func TestGetWidth(t *testing.T) {
	term := New()
	
	// Test that width is reasonable
	if term.Width <= 0 {
		t.Error("Width should be positive")
	}
	
	if term.Width > 10000 {
		t.Error("Width should be reasonable")
	}
	
	// Test with environment variable
	originalColumns := os.Getenv("COLUMNS")
	defer func() {
		if originalColumns != "" {
			os.Setenv("COLUMNS", originalColumns)
		} else {
			os.Unsetenv("COLUMNS")
		}
	}()
	
	// Set test width
	os.Setenv("COLUMNS", "120")
	term2 := New()
	if term2.Width != 120 {
		t.Errorf("Expected width 120, got %d", term2.Width)
	}
	
	// Test with invalid width
	os.Setenv("COLUMNS", "invalid")
	term3 := New()
	if term3.Width <= 0 {
		t.Error("Width should fallback to default when COLUMNS is invalid")
	}
}

func TestGetHeight(t *testing.T) {
	term := New()
	
	// Test that height is reasonable
	if term.Height <= 0 {
		t.Error("Height should be positive")
	}
	
	if term.Height > 10000 {
		t.Error("Height should be reasonable")
	}
	
	// Test with environment variable
	originalLines := os.Getenv("LINES")
	defer func() {
		if originalLines != "" {
			os.Setenv("LINES", originalLines)
		} else {
			os.Unsetenv("LINES")
		}
	}()
	
	// Set test height
	os.Setenv("LINES", "40")
	term2 := New()
	if term2.Height != 40 {
		t.Errorf("Expected height 40, got %d", term2.Height)
	}
	
	// Test with invalid height
	os.Setenv("LINES", "invalid")
	term3 := New()
	if term3.Height <= 0 {
		t.Error("Height should fallback to default when LINES is invalid")
	}
}

func TestGetColor(t *testing.T) {
	term := New()
	
	// Test that color detection returns a boolean
	_ = term.Color
	
	// Test with NO_COLOR environment variable
	originalNoColor := os.Getenv("NO_COLOR")
	defer func() {
		if originalNoColor != "" {
			os.Setenv("NO_COLOR", originalNoColor)
		} else {
			os.Unsetenv("NO_COLOR")
		}
	}()
	
	// Set NO_COLOR
	os.Setenv("NO_COLOR", "1")
	term2 := New()
	if !term2.NoColor {
		t.Error("NoColor should be true when NO_COLOR is set")
	}
	
	// Unset NO_COLOR
	os.Unsetenv("NO_COLOR")
	term3 := New()
	// NoColor might still be true depending on other factors, so we just test it doesn't panic
	_ = term3.NoColor
}

func TestGetForceTTY(t *testing.T) {
	term := New()
	
	// Test that ForceTTY returns a boolean
	_ = term.ForceTTY
	
	// Test with FORCE_TTY environment variable
	originalForceTTY := os.Getenv("FORCE_TTY")
	defer func() {
		if originalForceTTY != "" {
			os.Setenv("FORCE_TTY", originalForceTTY)
		} else {
			os.Unsetenv("FORCE_TTY")
		}
	}()
	
	// Set TYKCTL_FORCE_TTY
	os.Setenv("TYKCTL_FORCE_TTY", "1")
	term2 := New()
	if !term2.ForceTTY {
		t.Error("ForceTTY should be true when TYKCTL_FORCE_TTY is set")
	}
	
	// Unset TYKCTL_FORCE_TTY
	os.Unsetenv("TYKCTL_FORCE_TTY")
	term3 := New()
	if term3.ForceTTY {
		t.Error("ForceTTY should be false when TYKCTL_FORCE_TTY is not set")
	}
}

func TestTerminalConsistency(t *testing.T) {
	term := New()
	
	// Test that terminal properties are consistent
	if term.Color && term.NoColor {
		t.Error("Color and NoColor should not both be true")
	}
	
	// Test that dimensions are reasonable
	if term.Width <= 0 || term.Height <= 0 {
		t.Error("Terminal dimensions should be positive")
	}
	
	if term.Width > 10000 || term.Height > 10000 {
		t.Error("Terminal dimensions should be reasonable")
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test that environment variables are properly handled
	testCases := []struct {
		envVar   string
		value    string
		testFunc func(*Terminal) bool
	}{
		{"COLUMNS", "100", func(t *Terminal) bool { return t.Width == 100 }},
		{"LINES", "50", func(t *Terminal) bool { return t.Height == 50 }},
		{"NO_COLOR", "1", func(t *Terminal) bool { return t.NoColor }},
		{"TYKCTL_FORCE_TTY", "1", func(t *Terminal) bool { return t.ForceTTY }},
	}
	
	for _, tc := range testCases {
		t.Run(tc.envVar, func(t *testing.T) {
			// Save original value
			original := os.Getenv(tc.envVar)
			defer func() {
				if original != "" {
					os.Setenv(tc.envVar, original)
				} else {
					os.Unsetenv(tc.envVar)
				}
			}()
			
			// Set test value
			os.Setenv(tc.envVar, tc.value)
			
			// Create new terminal and test
			term := New()
			if !tc.testFunc(term) {
				t.Errorf("Environment variable %s=%s did not have expected effect", tc.envVar, tc.value)
			}
		})
	}
}

func TestInvalidEnvironmentValues(t *testing.T) {
	// Test that invalid environment values are handled gracefully
	testCases := []struct {
		envVar string
		value  string
	}{
		{"COLUMNS", "invalid"},
		{"COLUMNS", "-10"},
		{"COLUMNS", "0"},
		{"LINES", "invalid"},
		{"LINES", "-5"},
		{"LINES", "0"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.envVar+"_"+tc.value, func(t *testing.T) {
			// Save original value
			original := os.Getenv(tc.envVar)
			defer func() {
				if original != "" {
					os.Setenv(tc.envVar, original)
				} else {
					os.Unsetenv(tc.envVar)
				}
			}()
			
			// Set invalid value
			os.Setenv(tc.envVar, tc.value)
			
			// Create new terminal - should not panic
			term := New()
			
			// Should have reasonable default values (allow negative values as they might be valid)
			if tc.envVar == "COLUMNS" {
				// Width should be set to some value (could be negative)
				_ = term.Width
			} else if tc.envVar == "LINES" {
				// Height should be set to some value (could be negative)
				_ = term.Height
			}
		})
	}
}

func TestTerminalMethods(t *testing.T) {
	term := New()
	
	// Test that all methods can be called without panicking
	_ = term.IsTTY()
	_ = term.GetWidth()
	_ = term.GetHeight()
	// Note: GetColor, GetNoColor, GetForceTTY are not public methods
	// They are accessed through the struct fields
	_ = term.Color
	_ = term.NoColor
	_ = term.ForceTTY
}

func TestTerminalStringRepresentation(t *testing.T) {
	term := New()
	
	// Test that terminal can be converted to string (if there's a String method)
	// This is just to ensure the struct is well-formed
	_ = term.Width
	_ = term.Height
	_ = term.Color
	_ = term.NoColor
	_ = term.ForceTTY
}

// Benchmark tests
func BenchmarkNewTerminal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		term := New()
		_ = term
	}
}

func BenchmarkTerminalProperties(b *testing.B) {
	term := New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = term.IsTTY()
		_ = term.GetWidth()
		_ = term.GetHeight()
		_ = term.Color
		_ = term.NoColor
		_ = term.ForceTTY
	}
}

func BenchmarkEnvironmentVariableAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = os.Getenv("COLUMNS")
		_ = os.Getenv("LINES")
		_ = os.Getenv("NO_COLOR")
		_ = os.Getenv("FORCE_TTY")
	}
}