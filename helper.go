package wordcounter

import (
	"fmt"
	"path/filepath"
)

// ToAbsolutePath detects if a path is absolute or not. If not, it converts path to absolute.
// Returns the original path if conversion fails.
func ToAbsolutePath(path string) string {
	if path == "" {
		return path
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			// Return original path if conversion fails
			// This maintains backward compatibility while being more robust
			return path
		}
		path = absPath
	}
	return path
}

// toAbsolutePathWithError detects if a path is absolute or not. If not, it converts path to absolute.
// Returns an error if the conversion fails.
func toAbsolutePathWithError(path string) (string, error) {
	if path == "" {
		return "", NewInvalidInputError("path cannot be empty")
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", NewInvalidPathError(path, err)
		}
		return absPath, nil
	}
	return path, nil
}

func convertToSliceOfString(input [][]interface{}) [][]string {
	result := make([][]string, len(input))

	for i, row := range input {
		result[i] = make([]string, len(row))
		for j, value := range row {
			if value == nil {
				result[i][j] = ""
			} else {
				result[i][j] = fmt.Sprintf("%v", value)
			}
		}
	}

	return result
}

func getTotal(fcs []*FileCounter) Row {
	AllLines := 0
	AllChineseChars := 0
	AllNonChineseChars := 0
	AllTotalChars := 0

	for _, fc := range fcs {
		AllLines += fc.tc.S.Lines
		AllChineseChars += fc.tc.S.ChineseChars
		AllNonChineseChars += fc.tc.S.NonChineseChars
		AllTotalChars += fc.tc.S.TotalChars
	}

	row := Row{
		"Total",
		AllLines,
		AllChineseChars,
		AllNonChineseChars,
		AllTotalChars,
	}

	return row
}
