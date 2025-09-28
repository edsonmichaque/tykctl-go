package editor

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	editor := New()
	if editor == nil {
		t.Error("New() returned nil")
	}
}

func TestNewWithEditor(t *testing.T) {
	editor := NewWithEditor("vim")
	if editor == nil {
		t.Error("NewWithEditor() returned nil")
	}
}

func TestSetEditor(t *testing.T) {
	editor := New()
	editor.SetEditor("nano")
	// We can't easily test the internal state, but we can verify it doesn't panic
}

func TestSetArgs(t *testing.T) {
	editor := New()
	args := []string{"-w", "-n"}
	editor.SetArgs(args)
	// We can't easily test the internal state, but we can verify it doesn't panic
}

func TestSetTimeout(t *testing.T) {
	editor := New()
	timeout := 10 * time.Minute
	editor.SetTimeout(timeout)
	// We can't easily test the internal state, but we can verify it doesn't panic
}

func TestEditFile(t *testing.T) {
	editor := New()
	ctx := context.Background()

	// Test with non-existent file (should still try to open)
	// This will likely fail, but we're testing the interface
	err := editor.EditFile(ctx, "/tmp/non-existent-file-for-testing")

	// We expect an error since the file doesn't exist and the editor will fail
	// The important thing is that the method doesn't panic
	if err == nil {
		t.Log("EditFile() succeeded unexpectedly (editor might have created the file)")
	} else {
		t.Logf("EditFile() failed as expected: %v", err)
	}
}

func TestEditString(t *testing.T) {
	editor := New()
	ctx := context.Background()

	// Test editing a string
	content := "Hello, World!"
	result, err := editor.EditString(ctx, content)

	// This will likely fail in a test environment, but we're testing the interface
	if err == nil {
		t.Logf("EditString() succeeded, result: %s", result)
	} else {
		t.Logf("EditString() failed as expected: %v", err)
	}
}

func TestGetDefaultEditor(t *testing.T) {
	// Save original environment
	originalEditor := os.Getenv("EDITOR")
	originalVisual := os.Getenv("VISUAL")

	// Clean up after test
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
		if originalVisual != "" {
			os.Setenv("VISUAL", originalVisual)
		} else {
			os.Unsetenv("VISUAL")
		}
	}()

	tests := []struct {
		name     string
		editor   string
		visual   string
		expected string
	}{
		{
			name:     "EDITOR takes precedence",
			editor:   "vim",
			visual:   "code",
			expected: "vim",
		},
		{
			name:     "VISUAL used when EDITOR not set",
			editor:   "",
			visual:   "code",
			expected: "code",
		},
		{
			name:     "fallback when neither set",
			editor:   "",
			visual:   "",
			expected: "", // Will be empty since we can't predict which fallback editor exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.editor != "" {
				os.Setenv("EDITOR", tt.editor)
			} else {
				os.Unsetenv("EDITOR")
			}

			if tt.visual != "" {
				os.Setenv("VISUAL", tt.visual)
			} else {
				os.Unsetenv("VISUAL")
			}

			// Test detection
			detected := getDefaultEditor()

			// For the fallback case, we can't predict which editor will be found
			// so we just check that something is returned or empty
			if tt.expected == "" {
				// Just verify the function doesn't panic
				_ = detected
			} else {
				if detected != tt.expected {
					t.Errorf("getDefaultEditor() = %v, want %v", detected, tt.expected)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		editor := New()
		_ = editor
	}
}

func BenchmarkNewWithEditor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		editor := NewWithEditor("vim")
		_ = editor
	}
}

func BenchmarkGetDefaultEditor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		editor := getDefaultEditor()
		_ = editor
	}
}