package wordcounter

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type TextCounter struct {
	S *Stats
}

func NewTextCounter() *TextCounter {
	return &TextCounter{S: &Stats{}}
}

// GetStats returns the counting statistics
func (c *TextCounter) GetStats() *Stats {
	return c.S
}

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

// CountBytes efficiently counts characters from byte slice with minimal memory allocation
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
