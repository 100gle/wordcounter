package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
)

// current working directory
var wd string

func TestMain(m *testing.M) {
	// Perform setup actions before running the tests
	// Generate `test.txt` before testing
	err := ioutil.WriteFile("testdata/test.txt", []byte("你好 世界！Hello, world!"), 0644)
	if err != nil {
		log.Fatalf("Failed to generate test file, unexpected error: %v", err)
	}

	wd, err = os.Getwd()
	if err != nil {
		log.Fatalf("NewFileCounter() failed, unexpected error: %v", err)
	}

	// Run the tests
	exitCode := m.Run()

	// Perform teardown actions after running the tests
	err = os.Remove("testdata/test.txt")
	if err != nil {
		log.Fatalf("Failed to remove test file, unexpected error: %v", err)
	}

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestNewFileCounter(t *testing.T) {
	// Test creating a FileCounter instance with a valid filename and no ignore patterns
	fc := NewFileCounter("filename.txt")
	if fc.filename != filepath.Join(wd, "filename.txt") {
		t.Errorf("NewFileCounter() failed, expected filename: %s, got: %s", "filename.txt", fc.filename)
	}

	// Test creating a FileCounter instance with a valid filename and one or more ignore patterns
	fc = NewFileCounter("filename.txt", "*.txt", "*.csv")
	if fc.filename != filepath.Join(wd, "filename.txt") {
		t.Errorf("NewFileCounter() failed, expected filename: %s, got: %s", "filename.txt", fc.filename)
	}
	if len(fc.ignoreList) != 2 {
		t.Errorf("NewFileCounter() failed, expected ignoreList length: %d, got: %d", 2, len(fc.ignoreList))
	}

	// Test creating a FileCounter instance with an empty filename and no ignore patterns
	fc = NewFileCounter("")
	if fc.filename != "" {
		t.Errorf("NewFileCounter() failed, expected filename: %s, got: %s", "", fc.filename)
	}
}

func TestFileCounter_Count(t *testing.T) {
	filename := "testdata/test.txt"
	// Test counting the words in a valid file
	fc := NewFileCounter(filename)
	err := fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}
	expectedRow := Row{filepath.Join(wd, filename), "4", "2", "19"}
	row := fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}

	// Test counting the words in a non-existent file
	fc = NewFileCounter("testdata/nonexistent.txt")
	err = fc.Count()
	if err == nil {
		t.Error("FileCounter.Count() failed, expected error for non-existent file")
	}

	// Test counting the words in a file that should be ignored based on the ignore patterns
	fc = NewFileCounter(filename, "*.txt")
	err = fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}

	// Test counting the words in a long Chinese markdown content string
	longString := `这是一个长的中文字符串，用于测试。它应该包含足够的单词，以便我们可以测试 FileCounter.Count() 函数是否能够正确地计算这个字符串中的单词数。`
	filename = "testdata/long_chinese_string.txt"
	err = ioutil.WriteFile(filename, []byte(longString), 0644)
	if err != nil {
		t.Errorf("Failed to generate test file: %v", err)
	}
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			t.Errorf("Failed to delete test file: %v", err)
		}
	}()

	fc = NewFileCounter(filename)
	expectedRow = Row{filepath.Join(wd, filename), "54", "2", "79"}
	err = fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}
	row = fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}
}

func TestFileCounter_isIgnored(t *testing.T) {
	filename := "testdata/test.txt"
	fc := NewFileCounter(filename, "*.txt", "otherfile.txt")

	// Test checking if a file should be ignored based on an exact match ignore pattern
	result := fc.isIgnored("otherfile.txt")
	if !result {
		t.Error("FileCounter.isIgnored() failed, expected true for exact match ignore pattern")
	}

	// Test checking if a file should be ignored based on a wildcard ignore pattern
	result = fc.isIgnored("example.txt")
	if !result {
		t.Error("FileCounter.isIgnored() failed, expected true for wildcard ignore pattern")
	}

	// Test checking if a file should not be ignored
	result = fc.isIgnored("testfile.csv")
	if result {
		t.Error("FileCounter.isIgnored() failed, expected false for non-ignored file")
	}
}

func TestFileCounter_Ignore(t *testing.T) {
	fc := NewFileCounter("testdata/test.txt")

	// Test adding a new ignore pattern to the FileCounter instance
	fc.Ignore("*.txt")
	if len(fc.ignoreList) != 1 {
		t.Errorf("FileCounter.Ignore() failed, expected ignoreList length: %d, got: %d", 1, len(fc.ignoreList))
	}

	// Test adding multiple ignore patterns to the FileCounter instance
	fc.Ignore("*.csv", "*.xlsx")
	if len(fc.ignoreList) != 3 {
		t.Errorf("FileCounter.Ignore() failed, expected ignoreList length: %d, got: %d", 3, len(fc.ignoreList))
	}
}

func TestFileCounter_GetRow(t *testing.T) {
	filename := "testdata/test.txt"
	fc := NewFileCounter(filename)
	fc.Count()

	// Test getting the row data for a FileCounter instance with a valid filename and word counts
	expectedRow := Row{filepath.Join(wd, filename), "4", "2", "19"}
	row := fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}
}

func TestFileCounter_GetHeader(t *testing.T) {
	fc := NewFileCounter("testdata/test.txt")

	// Test getting the header row data for a FileCounter instance
	expectedHeader := Row{"File", "ChineseChars", "SpaceChars", "TotalChars"}
	header := fc.GetHeader()
	if !reflect.DeepEqual(header, expectedHeader) {
		t.Errorf("FileCounter.GetHeader() failed, expected header: %v, got: %v", expectedHeader, header)
	}
}

func TestFileCounter_GetHeaderAndRow(t *testing.T) {
	fc := NewFileCounter("testdata/test.txt")
	fc.Count()

	// Test getting both the header row and data row for a FileCounter instance
	expectedHeader := Row{"File", "ChineseChars", "SpaceChars", "TotalChars"}
	expectedRow := Row{filepath.Join(wd, "testdata/test.txt"), "4", "2", "19"}
	expectedData := []Row{expectedHeader, expectedRow}
	data := fc.GetHeaderAndRow()
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("FileCounter.GetHeaderAndRow() failed, expected data: %v, got: %v", expectedData, data)
	}
}

func TestFileCounter_ExportCSV(t *testing.T) {
	fc := NewFileCounter("testdata/test.txt")
	fc.Count()

	// Test exporting the word count data as a CSV string for a FileCounter instance
	expectedCSV := fmt.Sprintf("File,ChineseChars,SpaceChars,TotalChars\n%s,4,2,19", filepath.Join(wd, "testdata/test.txt"))
	csv := fc.ExportCSV()
	if csv != expectedCSV {
		t.Errorf("FileCounter.ExportCSV() failed, expected CSV: %s, got: %s", expectedCSV, csv)
	}
}

func TestFileCounter_ExportExcel(t *testing.T) {
	fc := NewFileCounter("testdata/test.txt")
	fc.Count()

	// Export the word count data to an Excel file for a FileCounter instance and check for errors
	if err := fc.ExportExcel("testdata/test.xlsx"); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// remove test.xlsx after testing
	if err := os.Remove("testdata/test.xlsx"); err != nil {
		t.Fatalf("Unexpected error while removing test.xlsx: %v", err)
	}
}

func TestFileCounter_ExportTable(t *testing.T) {
	filename := "testdata/test.txt"
	fc := NewFileCounter(filename)
	fc.Count()

	// Test exporting the word count data as a formatted table string for a FileCounter instance

	expectedTable := table.NewWriter()
	expectedTable.AppendHeader(Row{"File", "ChineseChars", "SpaceChars", "TotalChars"})
	expectedTable.AppendRow(Row{filepath.Join(wd, filename), "4", "2", "19"})

	table := fc.ExportTable()
	if table != expectedTable.Render() {
		t.Errorf("FileCounter.ExportTable() failed, expected table: %s, got: %s", expectedTable, table)
	}
}
