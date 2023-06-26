package main

import "path/filepath"

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
