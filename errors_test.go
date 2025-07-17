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
