package jq

import (
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	if !IsAvailable() {
		t.Skip("jq not available")
	}

	jsonData := []byte(`{"name": "John", "age": 30, "city": "NYC"}`)
	
	// Test basic field access
	result, err := Process(jsonData, ".name")
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	
	expected := `"John"`
	actual := strings.TrimSpace(string(result))
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestProcessString(t *testing.T) {
	if !IsAvailable() {
		t.Skip("jq not available")
	}

	jsonData := `{"users": [{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]}`
	
	result, err := ProcessString(jsonData, ".users[0].name")
	if err != nil {
		t.Fatalf("ProcessString failed: %v", err)
	}
	
	expected := `"John"`
	actual := strings.TrimSpace(result)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestProcessObject(t *testing.T) {
	if !IsAvailable() {
		t.Skip("jq not available")
	}

	data := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}
	
	result, err := ProcessObject(data, ".name")
	if err != nil {
		t.Fatalf("ProcessObject failed: %v", err)
	}
	
	expected := "John"
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIsAvailable(t *testing.T) {
	// This test will pass regardless of whether jq is available
	// It just tests that the function doesn't panic
	available := IsAvailable()
	_ = available // Use the variable to avoid unused variable warning
}