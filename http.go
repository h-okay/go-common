package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ContextKey defines a context key
type ContextKey string

// ContextSet adds a value to the request context, associating it with a specific key
func ContextSet(r *http.Request, key ContextKey, value any) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

// ContextGet retrieves a value from the request context using the provided key
func ContextGet[V any](r *http.Request, key ContextKey) V {
	value, exists := r.Context().Value(key).(V)
	if !exists {
		panic(fmt.Sprintf("key: %s doesn't exist on context", key))
	}
	return value
}

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
