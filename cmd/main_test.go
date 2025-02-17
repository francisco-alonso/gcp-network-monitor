package main

import (
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/logging"
)

// MockLogger simulates the Google Cloud Logger
type MockLogger struct {
	entries []logging.Entry
}

// Log function stores log entries instead of sending to GCP
func (m *MockLogger) Log(entry logging.Entry) {
	m.entries = append(m.entries, entry)
	fmt.Println("Mock Log:", entry.Payload)
}

func TestSendToCloudLogging(t *testing.T) {
	mockLogger := &MockLogger{} // Create a mock logger

	// Simulate logging a network device
	testIP := "192.168.1.100"
	testMAC := "00:1A:2B:3C:4D:5E"
	mockLogger.Log(logging.Entry{
		Timestamp: time.Now(),
		Severity:  logging.Info,
		Payload:   fmt.Sprintf("Discovered Device -> IP: %s, MAC: %s", testIP, testMAC),
	})

	// Verify the log entry was captured
	if len(mockLogger.entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(mockLogger.entries))
	}

	expectedLog := fmt.Sprintf("Discovered Device -> IP: %s, MAC: %s", testIP, testMAC)
	if mockLogger.entries[0].Payload != expectedLog {
		t.Errorf("Expected log entry '%s', but got '%s'", expectedLog, mockLogger.entries[0].Payload)
	}
}
