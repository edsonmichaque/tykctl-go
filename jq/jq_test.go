package jq

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
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

func TestProcessArray(t *testing.T) {
	jsonData := []byte(`[1, 2, 3, 4, 5]`)
	
	result, err := Process(jsonData, ".[] | select(. > 2)")
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	
	// Should return array of values > 2
	expected := `[3,4,5]`
	actual := strings.TrimSpace(string(result))
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestProcessComplex(t *testing.T) {
	jsonData := []byte(`{
		"users": [
			{"name": "John", "age": 30, "active": true},
			{"name": "Jane", "age": 25, "active": false},
			{"name": "Bob", "age": 35, "active": true}
		]
	}`)
	
	result, err := Process(jsonData, ".users[] | select(.active) | .name")
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	
	// Should return array of names of active users
	expected := `["John","Bob"]`
	actual := strings.TrimSpace(string(result))
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestEmptyProgram(t *testing.T) {
	jsonData := []byte(`{"name": "John"}`)
	
	result, err := Process(jsonData, "")
	if err != nil {
		t.Fatalf("Process with empty program failed: %v", err)
	}
	
	// Parse both JSONs to compare content, not formatting
	var expected, actual interface{}
	if err := json.Unmarshal(jsonData, &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected: %v", err)
	}
	if err := json.Unmarshal(result, &actual); err != nil {
		t.Fatalf("Failed to unmarshal actual: %v", err)
	}
	
	// Compare the parsed objects
	expectedStr := fmt.Sprintf("%v", expected)
	actualStr := fmt.Sprintf("%v", actual)
	if expectedStr != actualStr {
		t.Errorf("Expected %s, got %s", expectedStr, actualStr)
	}
}

func TestInvalidProgram(t *testing.T) {
	jsonData := []byte(`{"name": "John"}`)
	
	_, err := Process(jsonData, "invalid syntax")
	if err == nil {
		t.Error("Expected error for invalid jq program")
	}
}

func TestInvalidJSON(t *testing.T) {
	jsonData := []byte(`invalid json`)
	
	_, err := Process(jsonData, ".name")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}