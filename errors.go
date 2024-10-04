package common

import (
	"errors"
	"fmt"
)

/*
################################
||            HTTP            ||
################################
*/
var (
	ErrInternalServer    = errors.New("internal server error")            // 500 Internal Server Error
	ErrNotFound          = errors.New("resource not found")               // 404 Not Found
	ErrNoAuthHeader      = errors.New("no authorization header provided") // 403 Forbidden
	ErrForbidden         = errors.New("forbidden")                        // 403 Forbidden
	ErrUnauthorized      = errors.New("unauthorized")                     // 401 Unauthorized
	ErrInvalidAuthHeader = errors.New("invalid authorization header")     // 401 Unauthorized
	ErrBadRequest        = errors.New("bad request")                      // 400 Bad Request
)

// APIError represents an API error
type APIError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

// NewAPIError create a new APIError
func NewAPIError(code int, reason string) *APIError {
	return &APIError{
		Code:   code,
		Reason: reason,
	}
}

/*
################################
||            JSON            ||
################################
*/
var (
	ErrJSONParsing     = errors.New("body contains badly-formed JSON")                 // Malformed JSON body
	ErrJSONEmpty       = errors.New("body must not be empty")                          // Empty JSON body
	ErrJSONSingleValue = errors.New("body must only contain a single JSON value")      // JSON body has more than one value
	ErrJSONTooLarge    = fmt.Errorf("body must not be larger than %d bytes", MaxBytes) // JSON body exceeds max size
)

// JSONError is a custom type that represents various JSON error categories
type JSONError int8

const (
	InvalidType JSONError = iota //  JSON field has an incorrect data type
	UnknownKey                   // JSON contains a field not expected by the application
)

// JSONCustomError represents a custom error for JSON validation.
// Use NewJSONCustomError to create a new instance.
type JSONCustomError struct {
	Type JSONError // Type of the error
	Key  string    // The field key related to the error
}

// NewJSONCustomError creates a new instance of JSONCustomError
func NewJSONCustomError(t JSONError, k string) *JSONCustomError {
	return &JSONCustomError{
		Type: t,
		Key:  k,
	}
}

// Error returns a descriptive error message for a JSONCustomError based on its type and key
func (e *JSONCustomError) Error() string {
	if e.Key == "" {
		panic("JSONCustomError initialized without a key")
	}

	switch e.Type {
	case InvalidType:
		return fmt.Sprintf("body contains incorrect JSON type for field %q", e.Key)
	case UnknownKey:
		return fmt.Sprintf("body contains unknown key %s", e.Key)
	default:
		return "body contains incorrect JSON type"
	}
}
