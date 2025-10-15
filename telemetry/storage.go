// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileStorage implements the Storage interface using a local file.
type FileStorage struct {
	filename string
	mu       sync.RWMutex
}

// NewFileStorage creates a new file storage.
func NewFileStorage(filename string) Storage {
	return &FileStorage{
		filename: filename,
	}
}

// Store stores an event in the file.
func (s *FileStorage) Store(event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Read existing events
	events, err := s.readEvents()
	if err != nil {
		return fmt.Errorf("failed to read existing events: %w", err)
	}
	
	// Add new event
	events = append(events, event)
	
	// Write back to file
	return s.writeEvents(events)
}

// Retrieve retrieves all stored events.
func (s *FileStorage) Retrieve() ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.readEvents()
}

// Clear removes all stored events.
func (s *FileStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Write empty array to file
	return s.writeEvents([]*Event{})
}

// Count returns the number of stored events.
func (s *FileStorage) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	events, err := s.readEvents()
	if err != nil {
		return 0, err
	}
	
	return len(events), nil
}

// readEvents reads events from the file.
func (s *FileStorage) readEvents() ([]*Event, error) {
	// Check if file exists
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return []*Event{}, nil
	}
	
	// Read file
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Unmarshal JSON
	var events []*Event
	if len(data) > 0 {
		if err := json.Unmarshal(data, &events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}
	}
	
	return events, nil
}

// writeEvents writes events to the file.
func (s *FileStorage) writeEvents(events []*Event) error {
	// Ensure directory exists
	dir := filepath.Dir(s.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(s.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// MemoryStorage implements the Storage interface using in-memory storage.
type MemoryStorage struct {
	events []*Event
	mu     sync.RWMutex
}

// NewMemoryStorage creates a new memory storage.
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		events: make([]*Event, 0),
	}
}

// Store stores an event in memory.
func (s *MemoryStorage) Store(event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.events = append(s.events, event)
	return nil
}

// Retrieve retrieves all stored events.
func (s *MemoryStorage) Retrieve() ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Create a copy to avoid reference issues
	events := make([]*Event, len(s.events))
	copy(events, s.events)
	
	return events, nil
}

// Clear removes all stored events.
func (s *MemoryStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.events = s.events[:0]
	return nil
}

// Count returns the number of stored events.
func (s *MemoryStorage) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return len(s.events), nil
}

// MockStorage is a mock storage for testing.
type MockStorage struct {
	events []*Event
	errors []error
	mu     sync.RWMutex
}

// NewMockStorage creates a new mock storage.
func NewMockStorage() *MockStorage {
	return &MockStorage{
		events: make([]*Event, 0),
		errors: make([]error, 0),
	}
}

// Store stores an event and returns any configured error.
func (s *MockStorage) Store(event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.events = append(s.events, event)
	
	// Return any configured error
	if len(s.errors) > 0 {
		err := s.errors[0]
		s.errors = s.errors[1:]
		return err
	}
	
	return nil
}

// Retrieve retrieves all stored events and returns any configured error.
func (s *MockStorage) Retrieve() ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return any configured error
	if len(s.errors) > 0 {
		err := s.errors[0]
		s.errors = s.errors[1:]
		return nil, err
	}
	
	// Create a copy to avoid reference issues
	events := make([]*Event, len(s.events))
	copy(events, s.events)
	
	return events, nil
}

// Clear clears all stored events and returns any configured error.
func (s *MockStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Return any configured error
	if len(s.errors) > 0 {
		err := s.errors[0]
		s.errors = s.errors[1:]
		return err
	}
	
	s.events = s.events[:0]
	return nil
}

// Count returns the number of stored events and any configured error.
func (s *MockStorage) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return any configured error
	if len(s.errors) > 0 {
		err := s.errors[0]
		s.errors = s.errors[1:]
		return 0, err
	}
	
	return len(s.events), nil
}

// SetError configures the next operation to return an error.
func (s *MockStorage) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.errors = append(s.errors, err)
}

// GetEvents returns all stored events.
func (s *MockStorage) GetEvents() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Create a copy to avoid reference issues
	events := make([]*Event, len(s.events))
	copy(events, s.events)
	
	return events
}

// ClearAll clears all stored events and errors.
func (s *MockStorage) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.events = s.events[:0]
	s.errors = s.errors[:0]
}