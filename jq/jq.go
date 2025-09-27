package jq

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

// Process processes JSON data with a JQ program
func Process(data []byte, program string) ([]byte, error) {
	if program == "" {
		return data, nil
	}

	// Parse the jq program
	query, err := gojq.Parse(program)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jq program: %w", err)
	}

	// Parse input JSON
	var input interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %w", err)
	}

	// Execute the query
	iter := query.Run(input)
	var results []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq execution error: %w", err)
		}
		results = append(results, v)
	}

	// Marshal results back to JSON
	if len(results) == 1 {
		// Single result - return as-is
		return json.Marshal(results[0])
	}
	// Multiple results - return as array
	return json.Marshal(results)
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
