package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level represents the severity level for a log entry.
type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

// Return a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// LogEntry represents a log entry
type LogEntry struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

// Logger holds various information about the logger
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// NewLogger returns a new Logger instance
func NewLogger(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// PrintInfo prints a INFO level log
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	if _, err := l.print(LevelInfo, message, properties); err != nil {
		l.error(err)
	}
}

// PrintInfo prints a ERROR level log
func (l *Logger) PrintError(err error, properties map[string]string) {
	if _, err := l.print(LevelError, err.Error(), properties); err != nil {
		l.error(err)
	}
}

// PrintInfo prints a FATAL level log
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	if _, err := l.print(LevelFatal, err.Error(), properties); err != nil {
		l.error(err)
	}
	os.Exit(1)
}

// LoggerMiddleware is a HTTP middleware to log incoming requests
func (l *Logger) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		properties := map[string]string{
			"protocol":    r.Proto,
			"uri":         r.RequestURI,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
		}

		l.PrintInfo("request", properties)
		next.ServeHTTP(w, r)
	})
}

// print creates and writes a log entry to the Logger.out
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	entry := LogEntry{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		entry.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(entry)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// error handles print errors
func (l *Logger) error(err error) {
	fmt.Fprintf(os.Stderr, "Failed to log message: %v\n", err)
}
