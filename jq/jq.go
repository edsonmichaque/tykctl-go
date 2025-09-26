package jq

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// JQ represents a JQ processor
type JQ struct {
	program string
}

// New creates a new JQ processor
func New() *JQ {
	return &JQ{}
}

// NewWithProgram creates a new JQ processor with a program
func NewWithProgram(program string) *JQ {
	return &JQ{program: program}
}

// Process processes JSON data with JQ
func (j *JQ) Process(data []byte) ([]byte, error) {
	if j.program == "" {
		return data, nil
	}
	
	cmd := exec.Command("jq", j.program)
	cmd.Stdin = strings.NewReader(string(data))
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("jq processing failed: %w", err)
	}
	
	return output, nil
}

// ProcessString processes a JSON string with JQ
func (j *JQ) ProcessString(data string) (string, error) {
	result, err := j.Process([]byte(data))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// ProcessObject processes a JSON object with JQ
func (j *JQ) ProcessObject(data interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	
	result, err := j.Process(jsonData)
	if err != nil {
		return nil, err
	}
	
	var output interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	return output, nil
}

// SetProgram sets the JQ program
func (j *JQ) SetProgram(program string) {
	j.program = program
}

// GetProgram returns the JQ program
func (j *JQ) GetProgram() string {
	return j.program
}

// IsAvailable checks if JQ is available
func IsAvailable() bool {
	_, err := exec.LookPath("jq")
	return err == nil
}

// Process processes JSON data with a JQ program
func Process(data []byte, program string) ([]byte, error) {
	jq := NewWithProgram(program)
	return jq.Process(data)
}

// ProcessString processes a JSON string with a JQ program
func ProcessString(data string, program string) (string, error) {
	jq := NewWithProgram(program)
	return jq.ProcessString(data)
}

// ProcessObject processes a JSON object with a JQ program
func ProcessObject(data interface{}, program string) (interface{}, error) {
	jq := NewWithProgram(program)
	return jq.ProcessObject(data)
}

// Filter filters JSON data with a JQ filter
func Filter(data []byte, filter string) ([]byte, error) {
	return Process(data, filter)
}

// Select selects fields from JSON data
func Select(data []byte, fields ...string) ([]byte, error) {
	if len(fields) == 0 {
		return data, nil
	}
	
	selector := strings.Join(fields, ", ")
	return Process(data, selector)
}

// Format formats JSON data
func Format(data []byte) ([]byte, error) {
	return Process(data, ".")
}

// Compact compacts JSON data
func Compact(data []byte) ([]byte, error) {
	return Process(data, "-c")
}

// Pretty prints JSON data in a pretty format
func Pretty(data []byte) ([]byte, error) {
	return Process(data, ".")
}

// Keys extracts keys from JSON data
func Keys(data []byte) ([]byte, error) {
	return Process(data, "keys")
}

// Values extracts values from JSON data
func Values(data []byte) ([]byte, error) {
	return Process(data, "values")
}

// Length gets the length of JSON data
func Length(data []byte) ([]byte, error) {
	return Process(data, "length")
}

// Type gets the type of JSON data
func Type(data []byte) ([]byte, error) {
	return Process(data, "type")
}

// Has checks if JSON data has a field
func Has(data []byte, field string) ([]byte, error) {
	return Process(data, fmt.Sprintf("has(\"%s\")", field))
}

// Get gets a field from JSON data
func Get(data []byte, field string) ([]byte, error) {
	return Process(data, fmt.Sprintf(".%s", field))
}

// Set sets a field in JSON data
func Set(data []byte, field, value string) ([]byte, error) {
	return Process(data, fmt.Sprintf(".%s = \"%s\"", field, value))
}

// Delete deletes a field from JSON data
func Delete(data []byte, field string) ([]byte, error) {
	return Process(data, fmt.Sprintf("del(.%s)", field))
}

// Map maps over JSON data
func Map(data []byte, expression string) ([]byte, error) {
	return Process(data, fmt.Sprintf("map(%s)", expression))
}

// FilterArray filters an array
func FilterArray(data []byte, expression string) ([]byte, error) {
	return Process(data, fmt.Sprintf("map(select(%s))", expression))
}

// Sort sorts JSON data
func Sort(data []byte) ([]byte, error) {
	return Process(data, "sort")
}

// SortBy sorts JSON data by a field
func SortBy(data []byte, field string) ([]byte, error) {
	return Process(data, fmt.Sprintf("sort_by(.%s)", field))
}

// GroupBy groups JSON data by a field
func GroupBy(data []byte, field string) ([]byte, error) {
	return Process(data, fmt.Sprintf("group_by(.%s)", field))
}

// Unique gets unique values from JSON data
func Unique(data []byte) ([]byte, error) {
	return Process(data, "unique")
}

// Reverse reverses JSON data
func Reverse(data []byte) ([]byte, error) {
	return Process(data, "reverse")
}

// Slice slices JSON data
func Slice(data []byte, start, end int) ([]byte, error) {
	return Process(data, fmt.Sprintf("[%d:%d]", start, end))
}

// Head gets the first N items
func Head(data []byte, n int) ([]byte, error) {
	return Process(data, fmt.Sprintf(".[:%d]", n))
}

// Tail gets the last N items
func Tail(data []byte, n int) ([]byte, error) {
	return Process(data, fmt.Sprintf(".[-%d:]", n))
}
