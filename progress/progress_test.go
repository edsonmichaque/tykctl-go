package progress

import (
	"sync"
	"testing"
)

func TestNewSpinner(t *testing.T) {
	spinner := New()
	if spinner == nil {
		t.Fatal("New() returned nil")
	}
	
	if spinner.message != "Loading..." {
		t.Errorf("Expected default message 'Loading...', got '%s'", spinner.message)
	}
	
	if len(spinner.frames) == 0 {
		t.Error("Spinner should have frames")
	}
	
	if spinner.index != 0 {
		t.Error("Spinner index should start at 0")
	}
	
	if spinner.done {
		t.Error("Spinner should not be done initially")
	}
}

func TestNewBar(t *testing.T) {
	total := int64(100)
	bar := NewBar(total)
	if bar == nil {
		t.Fatal("NewBar() returned nil")
	}
	
	if bar.total != total {
		t.Errorf("Expected total %d, got %d", total, bar.total)
	}
	
	if bar.current != 0 {
		t.Error("Bar current should start at 0")
	}
	
	if bar.message != "Progress" {
		t.Errorf("Expected default message 'Progress', got '%s'", bar.message)
	}
	
	if bar.width != 50 {
		t.Errorf("Expected default width 50, got %d", bar.width)
	}
	
	if bar.fillChar != "█" {
		t.Errorf("Expected default fill char '█', got '%s'", bar.fillChar)
	}
	
	if bar.emptyChar != "░" {
		t.Errorf("Expected default empty char '░', got '%s'", bar.emptyChar)
	}
	
	if bar.done {
		t.Error("Bar should not be done initially")
	}
}

func TestSpinnerWithMessage(t *testing.T) {
	spinner := New()
	newMessage := "Custom message"
	result := spinner.WithMessage(newMessage)
	
	if result != spinner {
		t.Error("WithMessage should return the same spinner instance")
	}
	
	if spinner.message != newMessage {
		t.Errorf("Expected message '%s', got '%s'", newMessage, spinner.message)
	}
}

func TestBarWithMessage(t *testing.T) {
	bar := NewBar(100)
	newMessage := "Custom progress"
	result := bar.WithMessage(newMessage)
	
	if result != bar {
		t.Error("WithMessage should return the same bar instance")
	}
	
	if bar.message != newMessage {
		t.Errorf("Expected message '%s', got '%s'", newMessage, bar.message)
	}
}

func TestBarSetCurrent(t *testing.T) {
	bar := NewBar(100)
	
	// Test setting current progress
	bar.SetCurrent(50)
	if bar.current != 50 {
		t.Errorf("Expected current 50, got %d", bar.current)
	}
	
	// Test setting current beyond total
	bar.SetCurrent(150)
	if bar.current != 100 { // Should be capped at total
		t.Errorf("Expected current 100 (capped), got %d", bar.current)
	}
	
	// Test setting negative current
	bar.SetCurrent(-10)
	if bar.current != -10 {
		t.Errorf("Expected current -10, got %d", bar.current)
	}
}

func TestBarAdd(t *testing.T) {
	bar := NewBar(100)
	
	// Test adding progress
	bar.Add(30)
	if bar.current != 30 {
		t.Errorf("Expected current 30, got %d", bar.current)
	}
	
	// Test adding more progress
	bar.Add(50)
	if bar.current != 80 {
		t.Errorf("Expected current 80, got %d", bar.current)
	}
	
	// Test adding beyond total
	bar.Add(50)
	if bar.current != 100 { // Should be capped at total
		t.Errorf("Expected current 100 (capped), got %d", bar.current)
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
	
	// Should not panic and should be in a consistent state
	_ = spinner.message
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
	
	// Should not panic and should be in a consistent state
	if bar.current < 0 || bar.current > 1000 {
		t.Errorf("Bar current should be reasonable, got %d", bar.current)
	}
}

func TestSpinnerEdgeCases(t *testing.T) {
	spinner := New()
	
	// Test empty message
	spinner.WithMessage("")
	if spinner.message != "" {
		t.Error("Should be able to set empty message")
	}
}

func TestBarEdgeCases(t *testing.T) {
	// Test zero total
	bar := NewBar(0)
	if bar.total != 0 {
		t.Error("Should be able to create bar with zero total")
	}
	
	// Test negative total
	bar2 := NewBar(-100)
	if bar2.total != -100 {
		t.Error("Should be able to create bar with negative total")
	}
	
	// Test very large total
	bar3 := NewBar(1000000)
	if bar3.total != 1000000 {
		t.Error("Should be able to create bar with large total")
	}
	
	// Test empty message
	bar4 := NewBar(100)
	bar4.WithMessage("")
	if bar4.message != "" {
		t.Error("Should be able to set empty message")
	}
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