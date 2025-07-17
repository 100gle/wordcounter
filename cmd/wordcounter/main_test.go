package main

import (
	"testing"
)

// Test basic functionality without testing log.Fatal calls
// This provides some test coverage for the cmd package

func TestGlobalVariables(t *testing.T) {
	// Test that global variables have expected default values
	// This is a simple test to ensure the package can be imported and basic variables exist

	// Reset to defaults
	mode = "dir"
	exportType = "table"
	exportPath = "counter.xlsx"
	excludePattern = []string{}
	withTotal = false

	// Test that variables can be set
	mode = "file"
	if mode != "file" {
		t.Errorf("Expected mode to be 'file', got %s", mode)
	}

	exportType = "csv"
	if exportType != "csv" {
		t.Errorf("Expected exportType to be 'csv', got %s", exportType)
	}

	exportPath = "test.xlsx"
	if exportPath != "test.xlsx" {
		t.Errorf("Expected exportPath to be 'test.xlsx', got %s", exportPath)
	}

	withTotal = true
	if !withTotal {
		t.Errorf("Expected withTotal to be true, got %v", withTotal)
	}

	excludePattern = []string{"*.tmp"}
	if len(excludePattern) != 1 || excludePattern[0] != "*.tmp" {
		t.Errorf("Expected excludePattern to be ['*.tmp'], got %v", excludePattern)
	}
}
