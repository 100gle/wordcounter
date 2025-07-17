package wordcounter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileCounter provides character counting functionality for individual files.
// It implements the Counter interface and combines file I/O operations
// with text analysis capabilities.
type FileCounter struct {
	tc              *Counter // Internal text counter for character analysis (private)
	FileName        string   // Absolute path to the file being analyzed
	originalPath    string   // Original path as provided by user
	pathDisplayMode string   // Path display mode: absolute or relative
}

// NewFileCounter creates a new FileCounter instance for the specified file.
// The file path is automatically converted to an absolute path for consistency.
// The file is not read until Count() is called, allowing for lazy evaluation
// and better error handling.
//
// Parameters:
//   - filename: path to the file to be analyzed (relative or absolute)
//
// Returns a configured FileCounter ready for counting operations.
func NewFileCounter(filename string) *FileCounter {
	return NewFileCounterWithPathMode(filename, PathDisplayAbsolute)
}

// NewFileCounterWithPathMode creates a new FileCounter with specified path display mode.
func NewFileCounterWithPathMode(filename string, pathDisplayMode string) *FileCounter {
	tc := NewCounter()
	absPath := ToAbsolutePath(filename)

	fc := &FileCounter{
		FileName:        absPath,
		originalPath:    filename,
		pathDisplayMode: pathDisplayMode,
		tc:              tc,
	}

	return fc
}

// Count reads the file and performs character analysis.
// This method opens the file, reads its entire content into memory,
// and delegates the character counting to the internal Counter.
//
// The method uses io.ReadAll for optimal performance, reading the entire
// file at once to avoid issues with UTF-8 character boundaries that could
// occur with buffered reading.
//
// Returns structured errors for different failure scenarios:
//   - FileNotFoundError: if the file doesn't exist
//   - FileReadError: if there are I/O errors during reading or counting
func (fc *FileCounter) Count() error {
	file, err := os.Open(fc.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileNotFoundError(fc.FileName, err)
		}
		return NewFileReadError(fc.FileName, err)
	}
	defer file.Close()

	// Read entire file at once for better performance and simpler logic
	// This avoids issues with splitting lines/characters across buffer boundaries
	data, err := io.ReadAll(file)
	if err != nil {
		return NewFileReadError(fc.FileName, err)
	}

	// Check for empty file and issue warning
	if len(data) == 0 {
		displayPath := fc.getDisplayPath()
		fmt.Fprintf(os.Stderr, "Warning: Empty file detected: %s\n", displayPath)
	}

	if err := fc.tc.CountBytes(data); err != nil {
		return NewFileReadError(fc.FileName, err)
	}

	return nil
}

// GetStats returns the counting statistics from the internal Counter.
// This method provides access to the detailed character counting results
// after Count() has been called.
func (fc *FileCounter) GetStats() *Stats {
	return fc.tc.GetStats()
}

func (fc *FileCounter) GetRow() Row {
	displayPath := fc.getDisplayPath()
	row := append(Row{displayPath}, fc.tc.S.ToRow()...)
	return row
}

// getDisplayPath returns the path to display based on the path display mode
func (fc *FileCounter) getDisplayPath() string {
	if fc.pathDisplayMode == PathDisplayRelative {
		if fc.originalPath != "" {
			return fc.originalPath
		}
		// Fallback to basename if original path is not available
		return filepath.Base(fc.FileName)
	}
	return fc.FileName
}

func (fc *FileCounter) GetHeader() Row {
	headers := append(Row{"File"}, fc.tc.S.Header()...)
	return headers
}

func (fc *FileCounter) ExportCSV(filename ...string) (string, error) {
	return ExportCounterCSV(fc, filename...)
}

func (fc *FileCounter) ExportExcel(filename ...string) error {
	return ExportCounterExcel(fc, filename...)
}

func (fc *FileCounter) ExportTable() string {
	return ExportCounterTable(fc)
}

// GetRows returns the data rows (implements Counter interface)
func (fc *FileCounter) GetRows() []Row {
	return []Row{fc.GetRow()}
}
