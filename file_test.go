package wordcounter_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/jedib0t/go-pretty/v6/table"
)

// current working directory
var wd string

func TestMain(m *testing.M) {
	// Perform setup actions before running the tests
	// Generate `test.txt` before testing
	err := os.WriteFile("testdata/test.txt", []byte("你好 世界！Hello, world!"), 0644)
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
	fc := wcg.NewFileCounter("filename.txt")
	if fc.FileName != filepath.Join(wd, "filename.txt") {
		t.Errorf("wcg.NewFileCounter() failed, expected filename: %s, got: %s", "filename.txt", fc.FileName)
	}

	// Test creating a FileCounter instance with a valid filename and one or more ignore patterns
	fc = wcg.NewFileCounter("filename.txt")
	if fc.FileName != filepath.Join(wd, "filename.txt") {
		t.Errorf("NewFileCounter() failed, expected filename: %s, got: %s", "filename.txt", fc.FileName)
	}
	// Test creating a FileCounter instance with an empty filename and no ignore patterns
	fc = wcg.NewFileCounter("")
	if fc.FileName != "" {
		t.Errorf("NewFileCounter() failed, expected filename: %s, got: %s", "", fc.FileName)
	}
}

func TestFileCounter_Count(t *testing.T) {
	filename := "testdata/test.txt"
	// Test counting the words in a valid file
	fc := wcg.NewFileCounter(filename)
	err := fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}
	expectedRow := wcg.Row{filepath.Join(wd, filename), 1, 4, 15, 19}
	row := fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}

	// Test counting the words in a non-existent file
	fc = wcg.NewFileCounter("testdata/nonexistent.txt")
	err = fc.Count()
	if err == nil {
		t.Error("FileCounter.Count() failed, expected error for non-existent file")
	}

	// Test counting the words in a file that should be ignored based on the ignore patterns
	fc = wcg.NewFileCounter(filename)
	err = fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}

	// Test counting the words in a long Chinese markdown content string
	longString := `这是一个长的中文字符串，用于测试。它应该包含足够的单词，以便我们可以测试 FileCounter.Count() 函数是否能够正确地计算这个字符串中的单词数。`
	filename = "testdata/long_chinese_string.txt"
	err = os.WriteFile(filename, []byte(longString), 0644)
	if err != nil {
		t.Errorf("Failed to generate test file: %v", err)
	}
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			t.Errorf("Failed to delete test file: %v", err)
		}
	}()

	fc = wcg.NewFileCounter(filename)
	expectedRow = wcg.Row{filepath.Join(wd, filename), 1, 54, 25, 79}
	err = fc.Count()
	if err != nil {
		t.Errorf("FileCounter.Count() failed, unexpected error: %v", err)
	}
	row = fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}
}

func TestFileCounter_GetRow(t *testing.T) {
	filename := "testdata/test.txt"
	fc := wcg.NewFileCounter(filename)
	fc.Count()

	// Test getting the row data for a FileCounter instance with a valid filename and word counts
	expectedRow := wcg.Row{filepath.Join(wd, filename), 1, 4, 15, 19}
	row := fc.GetRow()
	if !reflect.DeepEqual(row, expectedRow) {
		t.Errorf("FileCounter.GetRow() failed, expected row: %v, got: %v", expectedRow, row)
	}
}

func TestFileCounter_GetHeader(t *testing.T) {
	fc := wcg.NewFileCounter("testdata/test.txt")

	// Test getting the header row data for a FileCounter instance
	expectedHeader := wcg.Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"}
	header := fc.GetHeader()
	if !reflect.DeepEqual(header, expectedHeader) {
		t.Errorf("FileCounter.GetHeader() failed, expected header: %v, got: %v", expectedHeader, header)
	}
}

func TestFileCounter_GetHeaderAndRows(t *testing.T) {
	fc := wcg.NewFileCounter("testdata/test.txt")
	fc.Count()

	// Test getting both the header row and data row for a FileCounter instance
	expectedHeader := wcg.Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"}
	expectedRow := wcg.Row{filepath.Join(wd, "testdata/test.txt"), 1, 4, 15, 19}
	expectedData := []wcg.Row{expectedHeader, expectedRow}

	// Use the public interface methods instead of the private helper
	header := fc.GetHeader()
	rows := fc.GetRows()
	data := make([]wcg.Row, 0, len(rows)+1)
	data = append(data, header)
	data = append(data, rows...)

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("FileCounter GetHeaderAndRows() failed, expected data: %v, got: %v", expectedData, data)
	}
}

func TestFileCounter_ExportCSV(t *testing.T) {
	fc := wcg.NewFileCounter("testdata/test.txt")
	fc.Count()

	// Test exporting the word count data as a CSV string for a FileCounter instance
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,1,4,15,19", filepath.Join(wd, "testdata/test.txt"))
	csv, err := fc.ExportCSV()
	if err != nil {
		t.Fatalf("Unexpected error when export to csv: %v", err)
	}
	if csv != expectedCSV {
		t.Errorf("FileCounter.ExportCSV() failed, expected CSV: %s, got: %s", expectedCSV, csv)
	}
}

func TestFileCounter_ExportCSVWithFileName(t *testing.T) {
	fc := wcg.NewFileCounter("testdata/test.txt")
	fc.Count()

	// Test exporting the word count data as a CSV string for a FileCounter instance
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,1,4,15,19", filepath.Join(wd, "testdata/test.txt"))
	csv, err := fc.ExportCSV("test.csv")
	if err != nil {
		t.Fatalf("Unexpected error when export to csv: %v", err)
	}

	if _, err := os.Stat("test.csv"); os.IsNotExist(err) {
		t.Fatalf("Expected file test.csv does not exist")
	}

	if csv != expectedCSV {
		t.Errorf("FileCounter.ExportCSV() failed, expected CSV: %s, got: %s", expectedCSV, csv)
	}

	err = os.Remove("test.csv")
	if err != nil {
		t.Fatalf("Unexpected error while removing test.csv: %v", err)
	}
}

func TestFileCounter_ExportExcel(t *testing.T) {
	fc := wcg.NewFileCounter("testdata/test.txt")
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
	fc := wcg.NewFileCounter(filename)
	fc.Count()

	// Test exporting the word count data as a formatted table string for a FileCounter instance

	expectedTable := table.NewWriter()
	expectedTable.AppendHeader(wcg.Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"})
	expectedTable.AppendRow(wcg.Row{filepath.Join(wd, filename), 1, 4, 15, 19})

	table := fc.ExportTable()
	if table != expectedTable.Render() {
		t.Errorf("FileCounter.ExportTable() failed, expected table: %s, got: %s", expectedTable, table)
	}
}

func TestFileCounter_GetStats(t *testing.T) {
	filename := "testdata/test.txt"
	fc := wcg.NewFileCounter(filename)
	fc.Count()

	// Test the new GetStats method
	stats := fc.GetStats()
	if stats == nil {
		t.Error("FileCounter.GetStats() returned nil")
		return
	}

	if stats.Lines != 1 || stats.ChineseChars != 4 || stats.NonChineseChars != 15 || stats.TotalChars != 19 {
		t.Errorf("FileCounter.GetStats() returned incorrect stats: %+v", stats)
	}
}
