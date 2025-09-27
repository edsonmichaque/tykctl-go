package jq

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Process processes JSON data with a JQ program
func Process(data []byte, program string) ([]byte, error) {
	if program == "" {
		return data, nil
	}
	
	cmd := exec.Command("jq", program)
	cmd.Stdin = strings.NewReader(string(data))
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("jq processing failed: %w", err)
	}
	
	return output, nil
}

// ProcessString processes a JSON string with a JQ program
func ProcessString(data string, program string) (string, error) {
	result, err := Process([]byte(data), program)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// ProcessObject processes a JSON object with a JQ program
func ProcessObject(data interface{}, program string) (interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	
	result, err := Process(jsonData, program)
	if err != nil {
		return nil, err
	}
	
	var output interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	return output, nil
}

// IsAvailable checks if JQ is available on the system
func IsAvailable() bool {
	_, err := exec.LookPath("jq")
	return err == nil
}