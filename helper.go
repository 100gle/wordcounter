package main

import (
	"fmt"
	"path/filepath"
)

// ToAbsolutePath detects if a path is absolute or not. If not, it converts path to absolute.
func ToAbsolutePath(path string) string {
	if path == "" {
		return path
	}

	if !filepath.IsAbs(path) {
		absPath, _ := filepath.Abs(path)
		path = absPath
	}
	return path
}

func ConvertToSliceOfString(input [][]interface{}) [][]string {
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

func GetTotal(fcs []*FileCounter) Row {
	AllLines := 0
	AllChineseChars := 0
	AllNonChineseChars := 0
	AllTotalChars := 0

	for _, fc := range fcs {
		AllLines += fc.tc.s.Lines
		AllChineseChars += fc.tc.s.ChineseChars
		AllNonChineseChars += fc.tc.s.NonChineseChars
		AllTotalChars += fc.tc.s.TotalChars
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
