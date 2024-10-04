package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLogger_PrintInfo tests the PrintInfo method
func TestLogger_PrintInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, LevelInfo)

	properties := map[string]string{"key": "value"}
	logger.PrintInfo("This is an info message", properties)

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
	if entry.Message != "This is an info message" {
		t.Errorf("Expected message 'This is an info message', got %s", entry.Message)
	}
	if entry.Properties["key"] != "value" {
		t.Errorf("Expected property 'key' to have value 'value', got %s", entry.Properties["key"])
	}
}

// TestLogger_PrintError tests the PrintError method
func TestLogger_PrintError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, LevelError)

	properties := map[string]string{"errorKey": "errorValue"}
	testErr := errors.New("This is an error message")
	logger.PrintError(testErr, properties)

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "ERROR" {
		t.Errorf("Expected level ERROR, got %s", entry.Level)
	}
	if entry.Message != testErr.Error() {
		t.Errorf("Expected message '%s', got %s", testErr.Error(), entry.Message)
	}
	if entry.Properties["errorKey"] != "errorValue" {
		t.Errorf("Expected property 'errorKey' to have value 'errorValue', got %s", entry.Properties["errorKey"])
	}
	if entry.Trace == "" {
		t.Error("Expected trace to be non-empty for ERROR level")
	}
}

// TestLoggerMiddleware tests the LoggerMiddleware method
func TestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, LevelInfo)

	handler := logger.LoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.RemoteAddr = "127.0.0.1"

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
	if entry.Message != "request" {
		t.Errorf("Expected message 'request', got %s", entry.Message)
	}
	if entry.Properties["remote_addr"] != "127.0.0.1" {
		t.Errorf("Expected remote_addr to be '127.0.0.1', got %s", entry.Properties["remote_addr"])
	}
	if entry.Properties["method"] != http.MethodGet {
		t.Errorf("Expected method to be '%s', got %s", http.MethodGet, entry.Properties["method"])
	}
}
