package common

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		dst            any
		expectedErr    error
		expectedErrMsg string
		expectedResult any
	}{
		{
			name:           "Valid JSON",
			body:           `{"key":"value"}`,
			dst:            &map[string]string{},
			expectedErr:    nil,
			expectedResult: &map[string]string{"key": "value"},
		},
		{
			name:        "Empty Body",
			body:        ``,
			dst:         &map[string]string{},
			expectedErr: ErrJSONEmpty,
		},
		{
			name:        "Syntax Error",
			body:        `{"key":`,
			dst:         &map[string]string{},
			expectedErr: ErrJSONParsing,
		},
		{
			name:           "Unknown Field",
			body:           `{"unknownField":"value"}`,
			dst:            &struct{ Key string }{},
			expectedErrMsg: `body contains unknown key "unknownField"`,
		},
		{
			name:        "Too Large Body",
			body:        fmt.Sprintf(`{"tooBig": "%s"}`, strings.Repeat("a", MaxBytes+1)),
			dst:         &map[string]string{},
			expectedErr: ErrJSONTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			err := ReadJSON(rec, req, tt.dst)

			if tt.expectedErrMsg != "" {
				if err == nil || err.Error() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, err)
				}
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				if !reflect.DeepEqual(tt.dst, tt.expectedResult) {
					t.Errorf("expected result %v, got %v", tt.expectedResult, tt.dst)
				}
			}
		})
	}
}

// TestWriteJSON tests the Write1JSON function
func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		data           Envelope
		status         int
		headers        http.Header
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Response",
			data:           Envelope{"key": "value"},
			status:         http.StatusOK,
			headers:        http.Header{"Custom-Header": []string{"custom-value"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"key":"value"}`,
		},
		{
			name:           "Internal Server Error",
			data:           Envelope{"error": "something went wrong"},
			status:         http.StatusInternalServerError,
			headers:        nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			err := WriteJSON(rec, tt.status, tt.data, tt.headers)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if strings.TrimSpace(rec.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}

			for key, value := range tt.headers {
				if rec.Header().Get(key) != value[0] {
					t.Errorf("expected header %q to have value %q, got %q", key, value[0], rec.Header().Get(key))
				}
			}
		})
	}
}
