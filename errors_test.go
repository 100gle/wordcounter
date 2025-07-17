package wordcounter_test

import (
	"errors"
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestWordCounterError(t *testing.T) {
	tests := []struct {
		name     string
		err      *wcg.WordCounterError
		wantMsg  string
		wantType wcg.ErrorType
	}{
		{
			name:     "File not found error",
			err:      wcg.NewFileNotFoundError("/path/to/file", errors.New("no such file")),
			wantMsg:  "file or directory not found: /path/to/file: no such file",
			wantType: wcg.ErrorTypeFileNotFound,
		},
		{
			name:     "Invalid input error",
			err:      wcg.NewInvalidInputError("invalid input"),
			wantMsg:  "invalid input",
			wantType: wcg.ErrorTypeInvalidInput,
		},
		{
			name:     "File read error",
			err:      wcg.NewFileReadError("/path/to/file", errors.New("permission denied")),
			wantMsg:  "failed to read file: /path/to/file: permission denied",
			wantType: wcg.ErrorTypeFileRead,
		},
		{
			name:     "File write error",
			err:      wcg.NewFileWriteError("/path/to/output", errors.New("disk full")),
			wantMsg:  "failed to write file: /path/to/output: disk full",
			wantType: wcg.ErrorTypeFileWrite,
		},
		{
			name:     "Invalid path error",
			err:      wcg.NewInvalidPathError("/invalid/path", errors.New("invalid characters")),
			wantMsg:  "invalid path: /invalid/path: invalid characters",
			wantType: wcg.ErrorTypeInvalidPath,
		},
		{
			name:     "Pattern match error",
			err:      wcg.NewPatternMatchError("*.{", errors.New("invalid regex")),
			wantMsg:  "invalid pattern: *.{: invalid regex",
			wantType: wcg.ErrorTypePatternMatch,
		},
		{
			name:     "Export error",
			err:      wcg.NewExportError("CSV export", errors.New("encoding error")),
			wantMsg:  "export failed: CSV export: encoding error",
			wantType: wcg.ErrorTypeExport,
		},
		{
			name:     "Server error",
			err:      wcg.NewServerError("failed to start server", errors.New("port in use")),
			wantMsg:  "failed to start server: port in use",
			wantType: wcg.ErrorTypeServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("WordCounterError.Error() = %v, want %v", tt.err.Error(), tt.wantMsg)
			}
			if tt.err.Type != tt.wantType {
				t.Errorf("WordCounterError.Type = %v, want %v", tt.err.Type, tt.wantType)
			}
		})
	}
}

func TestWordCounterError_WithContext(t *testing.T) {
	err := wcg.NewInvalidInputError("test error")
	err = err.WithContext("key1", "value1")
	err = err.WithContext("key2", 42)

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected context key1 = value1, got %v", err.Context["key1"])
	}
	if err.Context["key2"] != 42 {
		t.Errorf("Expected context key2 = 42, got %v", err.Context["key2"])
	}
}

func TestWordCounterError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := wcg.NewFileReadError("/path", cause)

	if unwrapped := err.Unwrap(); unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

func TestErrorConstructorsWithContext(t *testing.T) {
	tests := []struct {
		name        string
		err         *wcg.WordCounterError
		wantContext map[string]any
	}{
		{
			name: "File write error with path context",
			err:  wcg.NewFileWriteError("/output/file.txt", errors.New("disk full")),
			wantContext: map[string]any{
				"path": "/output/file.txt",
			},
		},
		{
			name: "Invalid path error with path context",
			err:  wcg.NewInvalidPathError("/invalid/path", errors.New("invalid chars")),
			wantContext: map[string]any{
				"path": "/invalid/path",
			},
		},
		{
			name: "Pattern match error with pattern context",
			err:  wcg.NewPatternMatchError("*.{", errors.New("invalid regex")),
			wantContext: map[string]any{
				"pattern": "*.{",
			},
		},
		{
			name: "Export error with operation context",
			err:  wcg.NewExportError("Excel export", errors.New("format error")),
			wantContext: map[string]any{
				"operation": "Excel export",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, expectedValue := range tt.wantContext {
				if actualValue, exists := tt.err.Context[key]; !exists {
					t.Errorf("Expected context key %s to exist", key)
				} else if actualValue != expectedValue {
					t.Errorf("Expected context[%s] = %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestWordCounterError_WithContextEdgeCases tests edge cases for WithContext
func TestWordCounterError_WithContextEdgeCases(t *testing.T) {
	// Test WithContext on error with nil context initially
	err := &wcg.WordCounterError{
		Type:    wcg.ErrorTypeInvalidInput,
		Message: "test error",
		Cause:   nil,
		Context: nil, // Start with nil context
	}

	// Add context to nil context map
	err = err.WithContext("first_key", "first_value")
	if err.Context == nil {
		t.Errorf("Expected context to be initialized, got nil")
	}
	if err.Context["first_key"] != "first_value" {
		t.Errorf("Expected first_key = first_value, got %v", err.Context["first_key"])
	}

	// Add more context
	err = err.WithContext("second_key", 123)
	err = err.WithContext("third_key", true)
	err = err.WithContext("fourth_key", nil) // Test nil value

	// Verify all values
	expected := map[string]any{
		"first_key":  "first_value",
		"second_key": 123,
		"third_key":  true,
		"fourth_key": nil,
	}

	for key, expectedValue := range expected {
		if actualValue, exists := err.Context[key]; !exists {
			t.Errorf("Expected context key %s to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected context[%s] = %v, got %v", key, expectedValue, actualValue)
		}
	}

	// Test overwriting existing key
	err = err.WithContext("first_key", "overwritten_value")
	if err.Context["first_key"] != "overwritten_value" {
		t.Errorf("Expected first_key to be overwritten, got %v", err.Context["first_key"])
	}
}

// TestWordCounterError_WithContextChaining tests method chaining
func TestWordCounterError_WithContextChaining(t *testing.T) {
	err := wcg.NewInvalidInputError("test error").
		WithContext("step", "validation").
		WithContext("input", "user_data").
		WithContext("timestamp", "2023-01-01")

	// Verify chaining worked
	if err.Context["step"] != "validation" {
		t.Errorf("Expected step = validation, got %v", err.Context["step"])
	}
	if err.Context["input"] != "user_data" {
		t.Errorf("Expected input = user_data, got %v", err.Context["input"])
	}
	if err.Context["timestamp"] != "2023-01-01" {
		t.Errorf("Expected timestamp = 2023-01-01, got %v", err.Context["timestamp"])
	}
}

// TestWordCounterError_WithContextComplexTypes tests complex data types in context
func TestWordCounterError_WithContextComplexTypes(t *testing.T) {
	err := wcg.NewExportError("test operation", errors.New("test cause"))

	// Test with slice
	slice := []string{"item1", "item2", "item3"}
	err = err.WithContext("items", slice)

	// Test with map
	mapData := map[string]int{"count": 42, "total": 100}
	err = err.WithContext("stats", mapData)

	// Test with struct
	type TestStruct struct {
		Name  string
		Value int
	}
	structData := TestStruct{Name: "test", Value: 123}
	err = err.WithContext("config", structData)

	// Verify complex types
	if actualSlice, ok := err.Context["items"].([]string); !ok {
		t.Errorf("Expected items to be []string")
	} else if len(actualSlice) != 3 || actualSlice[0] != "item1" {
		t.Errorf("Expected slice with 3 items starting with item1, got %v", actualSlice)
	}

	if actualMap, ok := err.Context["stats"].(map[string]int); !ok {
		t.Errorf("Expected stats to be map[string]int")
	} else if actualMap["count"] != 42 {
		t.Errorf("Expected count = 42, got %v", actualMap["count"])
	}

	if actualStruct, ok := err.Context["config"].(TestStruct); !ok {
		t.Errorf("Expected config to be TestStruct")
	} else if actualStruct.Name != "test" || actualStruct.Value != 123 {
		t.Errorf("Expected struct with Name=test, Value=123, got %+v", actualStruct)
	}
}
