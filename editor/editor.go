package editor

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// Editor represents a file editor
type Editor struct {
	editor  string
	args    []string
	timeout time.Duration
}

// New creates a new editor instance
func New() *Editor {
	return &Editor{
		editor:  getDefaultEditor(),
		timeout: 5 * time.Minute,
	}
}

// NewWithEditor creates a new editor with a specific editor
func NewWithEditor(editor string) *Editor {
	return &Editor{
		editor:  editor,
		timeout: 5 * time.Minute,
	}
}

// SetEditor sets the editor
func (e *Editor) SetEditor(editor string) {
	e.editor = editor
}

// SetArgs sets the editor arguments
func (e *Editor) SetArgs(args []string) {
	e.args = args
}

// SetTimeout sets the timeout
func (e *Editor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// EditFile edits a file
func (e *Editor) EditFile(ctx context.Context, filename string) error {
	cmd := exec.CommandContext(ctx, e.editor, append(e.args, filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// EditString edits a string and returns the result
func (e *Editor) EditString(ctx context.Context, content string) (string, error) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "tykctl-edit-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	// Write content to temp file
	if _, err := tmpfile.WriteString(content); err != nil {
		return "", err
	}
	tmpfile.Close()

	// Edit the file
	if err := e.EditFile(ctx, tmpfile.Name()); err != nil {
		return "", err
	}

	// Read the edited content
	edited, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return string(edited), nil
}

// getDefaultEditor gets the default editor from environment variables
func getDefaultEditor() string {
	if editor := os.Getenv("TYKCTL_EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "vim" // fallback
}
