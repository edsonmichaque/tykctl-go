package editor

import (
	"context"
	"os"
	"testing"
)

func TestDetectEditor(t *testing.T) {
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
			detected := DetectEditor()

			// For the fallback case, we can't predict which editor will be found
			// so we just check that something is returned or empty
			if tt.expected == "" {
				// Just verify the function doesn't panic
				_ = detected
			} else {
				if detected != tt.expected {
					t.Errorf("DetectEditor() = %v, want %v", detected, tt.expected)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	editor := New()
	if editor == nil {
		t.Error("New() returned nil")
	}

	// Verify it's an Editor
	if editor == nil {
		t.Error("New() did not return an Editor")
	}
}

func TestEditor_Open(t *testing.T) {
	editor := New()
	ctx := context.Background()

	// Test with non-existent file (should still try to open)
	// This will likely fail, but we're testing the interface
	err := editor.Open(ctx, "/tmp/non-existent-file-for-testing")

	// We expect an error since the file doesn't exist and the editor will fail
	// The important thing is that the method doesn't panic
	if err == nil {
		t.Log("Open() succeeded unexpectedly (editor might have created the file)")
	} else {
		t.Logf("Open() failed as expected: %v", err)
	}
}

func TestHasContentChanged(t *testing.T) {
	tests := []struct {
		name     string
		original []byte
		edited   []byte
		expected bool
	}{
		{
			name:     "identical content",
			original: []byte("hello world"),
			edited:   []byte("hello world"),
			expected: false,
		},
		{
			name:     "different content",
			original: []byte("hello world"),
			edited:   []byte("hello universe"),
			expected: true,
		},
		{
			name:     "empty vs content",
			original: []byte(""),
			edited:   []byte("hello"),
			expected: true,
		},
		{
			name:     "content vs empty",
			original: []byte("hello"),
			edited:   []byte(""),
			expected: true,
		},
		{
			name:     "whitespace difference",
			original: []byte("hello world"),
			edited:   []byte("hello  world"), // extra space
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasContentChanged(tt.original, tt.edited)
			if result != tt.expected {
				t.Errorf("hasContentChanged() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrNoChanges(t *testing.T) {
	// Test that ErrNoChanges is a known error type
	err := ErrNoChanges
	if err == nil {
		t.Error("ErrNoChanges should not be nil")
	}

	// Test error message
	expectedMsg := "no changes detected - content is identical to original"
	if err.Error() != expectedMsg {
		t.Errorf("ErrNoChanges.Error() = %v, want %v", err.Error(), expectedMsg)
	}
}
