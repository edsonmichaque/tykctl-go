package spinner

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if s.spinner == nil {
		t.Fatal("spinner is nil")
	}
}

func TestNewWithCharSet(t *testing.T) {
	s := NewWithCharSet(14, 200*time.Millisecond)
	if s == nil {
		t.Fatal("NewWithCharSet() returned nil")
	}
	if s.spinner == nil {
		t.Fatal("spinner is nil")
	}
}

func TestSpinnerOperations(t *testing.T) {
	s := New()

	// Test basic operations
	s.Start("Testing")
	time.Sleep(100 * time.Millisecond)
	s.Update("Updated message")
	time.Sleep(100 * time.Millisecond)
	s.Stop()
}

func TestWithContext(t *testing.T) {
	s := New()

	// Test successful operation
	err := s.WithContext(context.Background(), "Testing", func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err != nil {
		t.Errorf("WithContext returned error: %v", err)
	}
}

func TestWithContextCancellation(t *testing.T) {
	s := New()

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := s.WithContext(ctx, "Testing", func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
}

func TestWithTimeout(t *testing.T) {
	s := New()

	// Test timeout
	err := s.WithTimeout(50*time.Millisecond, "Testing", func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
