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
	"unicode/utf8"
)

// Counter provides character counting functionality for text content.
// It implements the CharacterCounter interface and tracks statistics
// including lines, Chinese characters, non-Chinese characters, and total characters.
type Counter struct {
	S *Stats // Statistics collected during counting
}

// NewCounter creates a new Counter instance with initialized statistics.
// The returned counter is ready to use for counting operations.
func NewCounter() *Counter {
	return &Counter{S: &Stats{}}
}

// GetStats returns the counting statistics
func (c *Counter) GetStats() *Stats {
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
func (c *Counter) Count(input any) error {
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

// isChinese checks if a rune is a Chinese character using direct Unicode range checks.
// This is more efficient than using unicode.In(r, unicode.Han) as it avoids
// the overhead of range table lookups.
//
// Covers the main CJK Unicode blocks:
//   - 0x4E00-0x9FFF: CJK Unified Ideographs (most common Chinese characters)
//   - 0x3400-0x4DBF: CJK Extension A
//   - 0x20000-0x2A6DF: CJK Extension B
//   - 0x2A700-0x2B73F: CJK Extension C
//   - 0x2B740-0x2B81F: CJK Extension D
//   - 0x2B820-0x2CEAF: CJK Extension E
//   - 0x2CEB0-0x2EBEF: CJK Extension F
//   - 0x3000-0x303F: CJK Symbols and Punctuation
//   - 0xFF00-0xFFEF: Halfwidth and Fullwidth Forms (Chinese punctuation)
func isChinese(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Extension A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK Extension E
		(r >= 0x2CEB0 && r <= 0x2EBEF) || // CJK Extension F
		(r >= 0x3000 && r <= 0x303F) || // CJK Symbols and Punctuation
		(r >= 0xFF00 && r <= 0xFFEF) // Halfwidth and Fullwidth Forms
}

// CountBytes efficiently counts characters from a byte slice with minimal memory allocation.
// This optimized version processes UTF-8 encoded text in a single pass and updates the following statistics:
//   - Lines: counted by scanning for newline characters (newlines + 1 for content)
//   - Chinese characters: identified using optimized Unicode range checks
//   - Non-Chinese characters: all other characters except newlines
//   - Total characters: sum of Chinese and non-Chinese characters (excluding newlines)
//
// Performance optimizations:
//   - Single-pass processing (combines line counting and character analysis)
//   - Direct Unicode range checks instead of unicode.In() for better performance
//   - Minimal function call overhead
//   - Local variables to reduce struct field access overhead
//
// Empty data is handled gracefully and returns zero counts for all statistics.
func (c *Counter) CountBytes(data []byte) error {

	// Use local variables to minimize struct field access overhead
	lines := 0
	chineseChars := 0
	nonChineseChars := 0

	// Single-pass processing: count lines and characters simultaneously
	for i := 0; i < len(data); {
		r, size := utf8.DecodeRune(data[i:])
		if r == '\n' {
			lines++
		} else {
			// Count non-newline characters
			if isChinese(r) {
				chineseChars++
			} else {
				nonChineseChars++
			}
		}
		i += size
	}

	// Line counting logic: number of newlines + 1 (if there's any content)
	// This correctly handles cases like "line1\nline2\nline3" (2 newlines = 3 lines)
	if len(data) > 0 {
		lines++ // Add 1 for the content itself
	}

	// Update statistics in batch to minimize memory writes
	c.S.Lines += lines
	c.S.ChineseChars += chineseChars
	c.S.NonChineseChars += nonChineseChars
	c.S.TotalChars += chineseChars + nonChineseChars

	return nil
}
