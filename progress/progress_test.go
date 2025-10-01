package progress

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewSpinner(t *testing.T) {
	spinner := New()
	if spinner == nil {
		t.Fatal("New() returned nil")
	}
	
	// Test that WithMessage works
	result := spinner.WithMessage("Test message")
	if result != spinner {
		t.Error("WithMessage should return the same spinner instance")
	}
}

func TestNewBar(t *testing.T) {
	total := int64(100)
	bar := NewBar(total)
	if bar == nil {
		t.Fatal("NewBar() returned nil")
	}
	
	// Test that WithMessage works
	result := bar.WithMessage("Test progress")
	if result != bar {
		t.Error("WithMessage should return the same bar instance")
	}
}

func TestSpinnerWithMessage(t *testing.T) {
	spinner := New()
	newMessage := "Custom message"
	result := spinner.WithMessage(newMessage)
	
	if result != spinner {
		t.Error("WithMessage should return the same spinner instance")
	}
}

func TestBarWithMessage(t *testing.T) {
	bar := NewBar(100)
	newMessage := "Custom progress"
	result := bar.WithMessage(newMessage)
	
	if result != bar {
		t.Error("WithMessage should return the same bar instance")
	}
}

func TestBarSetCurrent(t *testing.T) {
	bar := NewBar(100)
	
	// Test setting current progress
	bar.SetCurrent(50)
	
	// Test setting current beyond total
	bar.SetCurrent(150)
	
	// Test setting negative current
	bar.SetCurrent(-10)
	
	// These operations should not panic
}

func TestBarAdd(t *testing.T) {
	bar := NewBar(100)
	
	// Test adding progress
	bar.Add(30)
	
	// Test adding more progress
	bar.Add(50)
	
	// Test adding beyond total
	bar.Add(50)
	
	// These operations should not panic
}

func TestSpinnerWithContext(t *testing.T) {
	spinner := New()
	
	// Test successful operation
	err := spinner.WithContext(context.Background(), "Testing...", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSpinnerWithContextError(t *testing.T) {
	spinner := New()
	
	// Test error operation
	expectedErr := context.DeadlineExceeded
	err := spinner.WithContext(context.Background(), "Testing...", func() error {
		return expectedErr
	})
	
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestBarWithContext(t *testing.T) {
	bar := NewBar(100)
	
	// Test successful operation
	err := bar.WithContext(context.Background(), "Testing...", 100, func(update func(int64)) error {
		for i := int64(0); i < 100; i += 10 {
			update(10)
			time.Sleep(1 * time.Millisecond)
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBarWithContextError(t *testing.T) {
	bar := NewBar(100)
	
	// Test error operation
	expectedErr := context.DeadlineExceeded
	err := bar.WithContext(context.Background(), "Testing...", 100, func(update func(int64)) error {
		update(50)
		return expectedErr
	})
	
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestSpinnerConcurrency(t *testing.T) {
	spinner := New()
	
	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 10
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			spinner.WithMessage("Concurrent test")
		}()
	}
	
	wg.Wait()
	
	// Should not panic
}

func TestBarConcurrency(t *testing.T) {
	bar := NewBar(1000)
	
	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 10
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(progress int64) {
			defer wg.Done()
			bar.SetCurrent(progress)
			bar.WithMessage("Concurrent test")
		}(int64(i * 100))
	}
	
	wg.Wait()
	
	// Should not panic
}

func TestSpinnerEdgeCases(t *testing.T) {
	spinner := New()
	
	// Test empty message
	spinner.WithMessage("")
	
	// Should not panic
}

func TestBarEdgeCases(t *testing.T) {
	// Test zero total
	bar := NewBar(0)
	bar.WithMessage("Zero total")
	
	// Test negative total
	bar2 := NewBar(-100)
	bar2.WithMessage("Negative total")
	
	// Test very large total
	bar3 := NewBar(1000000)
	bar3.WithMessage("Large total")
	
	// Test empty message
	bar4 := NewBar(100)
	bar4.WithMessage("")
	
	// Should not panic
}

// Benchmark tests
func BenchmarkNewSpinner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spinner := New()
		_ = spinner
	}
}

func BenchmarkNewBar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bar := NewBar(100)
		_ = bar
	}
}

func BenchmarkSpinnerOperations(b *testing.B) {
	spinner := New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spinner.WithMessage("Benchmark test")
	}
}

func BenchmarkBarOperations(b *testing.B) {
	bar := NewBar(1000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bar.SetCurrent(int64(i % 1000))
		bar.WithMessage("Benchmark test")
	}
}

func BenchmarkSpinnerConcurrent(b *testing.B) {
	spinner := New()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spinner.WithMessage("Concurrent benchmark")
		}
	})
}

func BenchmarkBarConcurrent(b *testing.B) {
	bar := NewBar(1000)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bar.SetCurrent(int64(b.N % 1000))
			bar.WithMessage("Concurrent benchmark")
		}
	})
}