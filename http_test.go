package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestContextSet checks if the context is properly set
func TestContextSet(t *testing.T) {
	const key ContextKey = ContextKey("key")
	expectedValue := 123

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	reqWithCtx := ContextSet(req, key, expectedValue)

	value := reqWithCtx.Context().Value(key)
	if value != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}
}

// TestContextGet checks if the value can be retrieved from the context
func TestContextGet(t *testing.T) {
	const key ContextKey = ContextKey("key")
	expectedValue := 123

	// Create a request and set a value in the context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = ContextSet(req, key, expectedValue)

	// Retrieve the value using ContextGet
	value := ContextGet[int](req, key)
	if value != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}
}

// TestContextGet_Panic checks if the function panics when key doesn't exist in context
func TestContextGet_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic, but no panic occurred")
		}
	}()

	const key ContextKey = ContextKey("key")
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ContextGet[any](req, key)
}

// TestExtractBearerToken checks if the bearer token is correctly extracted
func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{"Bearer token123", "token123", false},
		{"Bearer", "", true},
		{"", "", true},
		{"Invalid token", "", true},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", tt.authHeader)

		token, err := ExtractBearerToken(req)
		if tt.expectError {
			if err == nil {
				t.Errorf("Expected error for header: %v, but got none", tt.authHeader)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for header: %v, but got: %v", tt.authHeader, err)
			}
			if token != tt.expectedToken {
				t.Errorf("Expected token %v, got %v", tt.expectedToken, token)
			}
		}
	}
}
