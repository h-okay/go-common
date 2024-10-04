package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const MaxBytes = 1_048_576 // 1MB

// Envelope is a type for structured JSON responses
type Envelope map[string]any

// ReadJSON decodes the JSON request body into the target destination
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(MaxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Syntax errors (bad JSON format)
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w (at character %d)", ErrJSONParsing, syntaxError.Offset)

		// Unexpected end of JSON input (truncated JSON)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return ErrJSONParsing

		// Wrong JSON type errors (type mismatch)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return NewJSONCustomError(InvalidType, unmarshalTypeError.Field)
			}
			return fmt.Errorf("%w (at character %d)", ErrJSONParsing, unmarshalTypeError.Offset)

		// Empty JSON body
		case errors.Is(err, io.EOF):
			return ErrJSONEmpty

		// Unknown fields (extra fields in JSON that aren't expected)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return NewJSONCustomError(UnknownKey, field)

		// Too large body
		case err.Error() == "http: request body too large":
			return ErrJSONTooLarge

		// Invalid JSON unmarshal errors (non-pointer or nil value passed as destination)
		case errors.As(err, &invalidUnmarshalError):
			panic(err) // should not happen

		default:
			return err
		}
	}

	// Single JSON object
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return ErrJSONSingleValue
	}

	return nil
}

// WriteJSON sends a JSON response to the client with the provided status code, data, and headers
func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}
