package fs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFS_MkdirAll(t *testing.T) {
	fs := New()
	ctx := context.Background()

	// Test creating a new directory
	dir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(dir)

	err := fs.MkdirAll(ctx, dir, 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// Test idempotent behavior - should succeed if directory already exists
	err = fs.MkdirAll(ctx, dir, 0755)
	if err != nil {
		t.Fatalf("MkdirAll should be idempotent: %v", err)
	}
}

func TestFS_WriteFile(t *testing.T) {
	fs := New()
	ctx := context.Background()

	// Test writing a new file
	file := filepath.Join(os.TempDir(), "tykctl-test.txt")
	defer os.Remove(file)

	data := []byte("test content")
	err := fs.WriteFile(ctx, file, data, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Test idempotent behavior - should succeed if content is same
	err = fs.WriteFile(ctx, file, data, 0644)
	if err != nil {
		t.Fatalf("WriteFile should be idempotent: %v", err)
	}

	// Test writing different content
	newData := []byte("different content")
	err = fs.WriteFile(ctx, file, newData, 0644)
	if err != nil {
		t.Fatalf("WriteFile with different content failed: %v", err)
	}
}

func TestFS_RemoveIfExists(t *testing.T) {
	fs := New()
	ctx := context.Background()

	// Test removing non-existent file - should succeed
	file := filepath.Join(os.TempDir(), "tykctl-nonexistent.txt")
	err := fs.RemoveIfExists(ctx, file)
	if err != nil {
		t.Fatalf("RemoveIfExists should succeed for non-existent file: %v", err)
	}

	// Test removing existing file
	file = filepath.Join(os.TempDir(), "tykctl-existing.txt")
	err = fs.WriteFile(ctx, file, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	defer os.Remove(file)

	err = fs.RemoveIfExists(ctx, file)
	if err != nil {
		t.Fatalf("RemoveIfExists failed: %v", err)
	}
}

func TestFS_Exists(t *testing.T) {
	fs := New()
	ctx := context.Background()

	// Test non-existent file
	file := filepath.Join(os.TempDir(), "tykctl-nonexistent.txt")
	exists, err := fs.Exists(ctx, file)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Fatal("File should not exist")
	}

	// Test existing file
	file = filepath.Join(os.TempDir(), "tykctl-existing.txt")
	err = fs.WriteFile(ctx, file, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	defer os.Remove(file)

	exists, err = fs.Exists(ctx, file)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Fatal("File should exist")
	}
}

func TestFS_ContextCancellation(t *testing.T) {
	fs := New()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context immediately
	cancel()

	// Test that operations respect context cancellation
	dir := filepath.Join(os.TempDir(), "tykctl-cancel-test")
	err := fs.MkdirAll(ctx, dir, 0755)
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, got: %v", err)
	}
}

func TestFS_Timeout(t *testing.T) {
	fs := New()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(2 * time.Millisecond)

	// Test that operations respect timeout
	dir := filepath.Join(os.TempDir(), "tykctl-timeout-test")
	err := fs.MkdirAll(ctx, dir, 0755)
	if err != context.DeadlineExceeded {
		t.Fatalf("Expected context.DeadlineExceeded, got: %v", err)
	}
}

func TestMemFS(t *testing.T) {
	fs := NewMem()
	ctx := context.Background()

	// Test in-memory filesystem
	dir := "/test/dir"
	err := fs.MkdirAll(ctx, dir, 0755)
	if err != nil {
		t.Fatalf("MemFS MkdirAll failed: %v", err)
	}

	file := "/test/dir/file.txt"
	data := []byte("test content")
	err = fs.WriteFile(ctx, file, data, 0644)
	if err != nil {
		t.Fatalf("MemFS WriteFile failed: %v", err)
	}

	readData, err := fs.ReadFile(ctx, file)
	if err != nil {
		t.Fatalf("MemFS ReadFile failed: %v", err)
	}

	if string(readData) != string(data) {
		t.Fatalf("Expected %s, got %s", string(data), string(readData))
	}
}
