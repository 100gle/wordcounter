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
