package wordcounter_test

import (
	"os"
	"strings"
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestExportToCSV(t *testing.T) {
	// Test CSV export through FileCounter interface
	// Create a test file
	testContent := "Hello 世界"
	testFile := "testdata/export_test.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Create FileCounter and test CSV export
	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	result, err := fc.ExportCSV()
	if err != nil {
		t.Errorf("ExportCSV failed with error: %v", err)
	}

	// Check if result contains expected headers
	if !strings.Contains(result, "File,Lines,ChineseChars,NonChineseChars,TotalChars") {
		t.Errorf("ExportCSV result doesn't contain expected headers: %v", result)
	}
}

func TestExportToCSVWithFilename(t *testing.T) {
	// Test CSV export with filename through FileCounter interface
	testContent := "Hello 世界"
	testFile := "testdata/export_test2.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	filename := "test.csv"
	csvData, err := fc.ExportCSV(filename)
	if err != nil {
		t.Errorf("ExportCSV failed with error: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("ExportCSV did not create file: %v", err)
	}

	// Check if result contains expected headers
	if !strings.Contains(csvData, "File,Lines,ChineseChars,NonChineseChars,TotalChars") {
		t.Errorf("ExportCSV result doesn't contain expected headers: %v", csvData)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("ExportCSV could not delete file: %v", err)
	}
}

func TestExportToExcel(t *testing.T) {
	// Test Excel export through FileCounter interface
	testContent := "Hello 世界"
	testFile := "testdata/export_test3.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	filename := "test.xlsx"
	err = fc.ExportExcel(filename)

	if err != nil {
		t.Errorf("ExportExcel failed with error: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("ExportExcel did not create file: %v", err)
	}

	// Clean up file
	err = os.Remove(filename)
	if err != nil {
		t.Errorf("ExportExcel could not delete file: %v", err)
	}
}

func TestExportToTable(t *testing.T) {
	// Test Table export through FileCounter interface
	testContent := "Hello 世界"
	testFile := "testdata/export_test4.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	result := fc.ExportTable()

	// Check if result contains expected headers and data
	if !strings.Contains(result, "FILE") || !strings.Contains(result, "LINES") {
		t.Errorf("ExportTable result doesn't contain expected content: %v", result)
	}
}

func TestExportCounterFunctions(t *testing.T) {
	// Test the standalone export functions
	testContent := "Hello 世界\nSecond line"
	testFile := "testdata/export_counter_test.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	// Test ExportCounterCSV
	csvResult, err := wcg.ExportCounterCSV(fc)
	if err != nil {
		t.Errorf("ExportCounterCSV failed: %v", err)
	}
	if !strings.Contains(csvResult, "File,Lines,ChineseChars,NonChineseChars,TotalChars") {
		t.Errorf("ExportCounterCSV result doesn't contain expected headers")
	}

	// Test ExportCounterCSV with filename
	csvFile := "test_counter.csv"
	_, err = wcg.ExportCounterCSV(fc, csvFile)
	if err != nil {
		t.Errorf("ExportCounterCSV with filename failed: %v", err)
	}
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Errorf("ExportCounterCSV did not create file")
	}
	defer os.Remove(csvFile)

	// Test ExportCounterExcel
	excelFile := "test_counter.xlsx"
	err = wcg.ExportCounterExcel(fc, excelFile)
	if err != nil {
		t.Errorf("ExportCounterExcel failed: %v", err)
	}
	if _, err := os.Stat(excelFile); os.IsNotExist(err) {
		t.Errorf("ExportCounterExcel did not create file")
	}
	defer os.Remove(excelFile)

	// Test ExportCounterTable
	tableResult := wcg.ExportCounterTable(fc)
	if !strings.Contains(tableResult, "FILE") || !strings.Contains(tableResult, "LINES") {
		t.Errorf("ExportCounterTable result doesn't contain expected content")
	}
}

func TestExportErrorHandling(t *testing.T) {
	// Test export to invalid path for Excel
	testContent := "Hello world"
	testFile := "testdata/export_error_test.txt"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	fc := wcg.NewFileCounter(testFile)
	err = fc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	// Test Excel export to invalid directory
	invalidPath := "/invalid/directory/test.xlsx"
	err = wcg.ExportCounterExcel(fc, invalidPath)
	if err == nil {
		t.Errorf("Expected error when exporting to invalid path, but got none")
	}

	// Test CSV export to invalid directory
	invalidCSVPath := "/invalid/directory/test.csv"
	_, err = wcg.ExportCounterCSV(fc, invalidCSVPath)
	if err == nil {
		t.Errorf("Expected error when exporting CSV to invalid path, but got none")
	}
}
