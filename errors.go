package wordcounter

import (
	"fmt"
)

// WordCounterError represents different types of errors that can occur in wordcounter
type WordCounterError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]any
}

// ErrorType represents the category of error
type ErrorType int

const (
	// ErrorTypeFileNotFound indicates a file or directory was not found
	ErrorTypeFileNotFound ErrorType = iota
	// ErrorTypeFileRead indicates an error reading a file
	ErrorTypeFileRead
	// ErrorTypeFileWrite indicates an error writing a file
	ErrorTypeFileWrite
	// ErrorTypeInvalidInput indicates invalid input was provided
	ErrorTypeInvalidInput
	// ErrorTypeInvalidPath indicates an invalid file path
	ErrorTypeInvalidPath
	// ErrorTypePatternMatch indicates an error in pattern matching
	ErrorTypePatternMatch
	// ErrorTypeExport indicates an error during export operations
	ErrorTypeExport
	// ErrorTypeServer indicates a server-related error
	ErrorTypeServer
)

// Error implements the error interface
func (e *WordCounterError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *WordCounterError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error
func (e *WordCounterError) WithContext(key string, value any) *WordCounterError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// NewError creates a new WordCounterError
func NewError(errorType ErrorType, message string, cause error) *WordCounterError {
	return &WordCounterError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]any),
	}
}

// NewFileNotFoundError creates a file not found error
func NewFileNotFoundError(path string, cause error) *WordCounterError {
	return NewError(ErrorTypeFileNotFound, fmt.Sprintf("file or directory not found: %s", path), cause).
		WithContext("path", path)
}

// NewFileReadError creates a file read error
func NewFileReadError(path string, cause error) *WordCounterError {
	return NewError(ErrorTypeFileRead, fmt.Sprintf("failed to read file: %s", path), cause).
		WithContext("path", path)
}

// NewFileWriteError creates a file write error
func NewFileWriteError(path string, cause error) *WordCounterError {
	return NewError(ErrorTypeFileWrite, fmt.Sprintf("failed to write file: %s", path), cause).
		WithContext("path", path)
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string) *WordCounterError {
	return NewError(ErrorTypeInvalidInput, message, nil)
}

// NewInvalidPathError creates an invalid path error
func NewInvalidPathError(path string, cause error) *WordCounterError {
	return NewError(ErrorTypeInvalidPath, fmt.Sprintf("invalid path: %s", path), cause).
		WithContext("path", path)
}

// NewPatternMatchError creates a pattern matching error
func NewPatternMatchError(pattern string, cause error) *WordCounterError {
	return NewError(ErrorTypePatternMatch, fmt.Sprintf("invalid pattern: %s", pattern), cause).
		WithContext("pattern", pattern)
}

// NewExportError creates an export error
func NewExportError(operation string, cause error) *WordCounterError {
	return NewError(ErrorTypeExport, fmt.Sprintf("export failed: %s", operation), cause).
		WithContext("operation", operation)
}

// NewServerError creates a server error
func NewServerError(message string, cause error) *WordCounterError {
	return NewError(ErrorTypeServer, message, cause)
}
