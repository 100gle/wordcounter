// Package wordcounter provides tools for counting Chinese characters in text files and directories.
//
// This package is designed primarily for Chinese text analysis, offering both single file
// and directory-based counting capabilities. It supports various export formats including
// ASCII tables, CSV, and Excel files.
//
// Key features:
//   - Count lines, Chinese characters, non-Chinese characters, and total characters
//   - Support for single files and recursive directory scanning
//   - Flexible ignore patterns similar to .gitignore
//   - Multiple export formats (table, CSV, Excel)
//   - HTTP server mode for API access
//   - Concurrent processing for improved performance
//   - Comprehensive error handling with structured error types
//
// Basic usage for single file:
//
//	counter := wordcounter.NewFileCounter("document.md")
//	if err := counter.Count(); err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(counter.ExportTable())
//
// Basic usage for directory:
//
//	counter := wordcounter.NewDirCounter("./docs", "*.tmp", "node_modules")
//	if err := counter.Count(); err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(counter.ExportTable())
package wordcounter

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TextCounter provides character counting functionality for text content.
// It implements the CharacterCounter interface and tracks statistics
// including lines, Chinese characters, non-Chinese characters, and total characters.
type TextCounter struct {
	S *Stats // Statistics collected during counting
}

// NewTextCounter creates a new TextCounter instance with initialized statistics.
// The returned counter is ready to use for counting operations.
func NewTextCounter() *TextCounter {
	return &TextCounter{S: &Stats{}}
}

// GetStats returns the counting statistics
func (c *TextCounter) GetStats() *Stats {
	return c.S
}

// Count analyzes the provided input and updates the character statistics.
// It accepts either string or []byte input and delegates to CountBytes for processing.
//
// Supported input types:
//   - string: converted to []byte for processing
//   - []byte: processed directly
//
// Returns an error if the input is empty or of an unsupported type.
func (c *TextCounter) Count(input interface{}) error {
	switch v := input.(type) {
	case string:
		if v == "" {
			return NewInvalidInputError("input string cannot be empty")
		}
		return c.CountBytes([]byte(v))
	case []byte:
		if len(v) == 0 {
			return NewInvalidInputError("input byte slice cannot be empty")
		}
		return c.CountBytes(v)
	default:
		return NewInvalidInputError(fmt.Sprintf("unsupported input type: %T, expected string or []byte", input))
	}
}

// CountBytes efficiently counts characters from a byte slice with minimal memory allocation.
// This method processes UTF-8 encoded text and updates the following statistics:
//   - Lines: counted by scanning for newline characters
//   - Chinese characters: identified using Unicode Han script ranges
//   - Non-Chinese characters: all other characters except newlines
//   - Total characters: sum of Chinese and non-Chinese characters (excluding newlines)
//
// The method uses utf8.DecodeRune for proper UTF-8 character boundary handling
// and avoids unnecessary string conversions for optimal performance.
//
// Returns an error if the input data is empty.
func (c *TextCounter) CountBytes(data []byte) error {
	if len(data) == 0 {
		return NewInvalidInputError("input data cannot be empty")
	}

	// Count lines by scanning for newline characters
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	// If there's content but no newlines, it's still one line
	if lines == 0 && len(data) > 0 {
		lines = 1
	}
	c.S.Lines += lines

	// Process runes directly from byte slice to avoid string conversion
	// Skip newline characters to match original behavior
	i := 0
	for i < len(data) {
		r, size := utf8.DecodeRune(data[i:])
		if r != '\n' { // Skip newline characters
			c.S.TotalChars++
			if unicode.In(r, unicode.Han) {
				c.S.ChineseChars++
			} else {
				c.S.NonChineseChars++
			}
		}
		i += size
	}

	return nil
}
