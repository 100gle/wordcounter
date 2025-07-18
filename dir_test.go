package wordcounter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/jedib0t/go-pretty/v6/table"
)

func TestDirCounter_Count(t *testing.T) {
	type args struct {
		dirname string
	}

	ignoreList := []string{".git", ".idea"}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Count files in directory",
			args: args{dirname: "./testdata"},
		},
		{
			name:    "Empty directory",
			args:    args{dirname: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := wcg.NewDirCounter(tt.args.dirname, ignoreList...)

			if err := dc.Count(); (err != nil) != tt.wantErr {
				t.Errorf("DirCounter.Count() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false && len(dc.GetRows()) == 0 {
				t.Errorf("DirCounter.Count() did not count any files")
			}
		})
	}
}
func TestFileCounter_IsIgnored(t *testing.T) {
	dirname := "testdata"
	dc := wcg.NewDirCounter(dirname, "*.txt", "otherfile.txt", ".*")

	// Test checking if a file should be ignored based on an exact match ignore pattern
	result := dc.IsIgnored("otherfile.txt")
	if !result {
		t.Error("FileCounter.isIgnored() failed, expected true for exact match ignore pattern")
	}

	// Test checking if a file should be ignored based on a wildcard ignore pattern
	result = dc.IsIgnored("example.txt")
	if !result {
		t.Error("FileCounter.isIgnored() failed, expected true for wildcard ignore pattern")
	}

	// Test checking if a file should not be ignored
	result = dc.IsIgnored("testfile.csv")
	if result {
		t.Error("FileCounter.isIgnored() failed, expected false for non-ignored file")
	}
	// Test checking if a file should not be ignored
	result = dc.IsIgnored(".git")
	if !result {
		t.Error("FileCounter.isIgnored() failed, expected false for ignored file")
	}
	// Test glob-like ignores with test table
	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{
			name:     "match one pattern",
			patterns: []string{"*.go"},
			path:     "main.go",
			want:     true,
		},
		{
			name:     "match multiple patterns",
			patterns: []string{"*.md", "*.txt"},
			path:     "README.md",
			want:     true,
		},
		{
			name:     "not match pattern",
			patterns: []string{"*.md"},
			path:     "main.go",
			want:     false,
		},
		{
			name:     "match single suffix pattern",
			patterns: []string{"*.js", "**/*.js"},
			path:     "foo.js",
			want:     true,
		},
		{
			name:     "match multiple suffix pattern",
			patterns: []string{"*.js.map", "**/*.js.map"},
			path:     "foo.js.map",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc = wcg.NewDirCounter(dirname, tt.patterns...)
			got := dc.IsIgnored(tt.path)
			if got != tt.want {
				t.Errorf("FileCounter.isIgnored(%v) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestDirCounter_Ignore(t *testing.T) {
	type fields struct {
		ignoreList []string
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Ignore single pattern",
			fields: fields{
				ignoreList: []string{},
			},
			args: args{pattern: ".git"},
			want: []string{".git"},
		},
		{
			name: "Ignore multiple patterns",
			fields: fields{
				ignoreList: []string{".git", ".idea"},
			},
			args: args{pattern: "node_modules"},
			want: []string{".git", ".idea", "node_modules"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := wcg.NewDirCounter(".", tt.fields.ignoreList...)
			dc.Ignore(tt.args.pattern)
			if !reflect.DeepEqual(dc.GetIgnoreList(), tt.want) {
				t.Errorf("DirCounter.Ignore() got = %v, want %v", dc.GetIgnoreList(), tt.want)
			}
		})
	}
}

func TestDirCounter_GetHeaderAndRows(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	tests := []struct {
		name string
		dc   *wcg.DirCounter
		want []wcg.Row
	}{
		{
			name: "GetHeaderAndRows",
			dc:   wcg.NewDirCounter(testDir),
			want: []wcg.Row{
				{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"},
				{filepath.Join(testDir, "empty.md"), 0, 0, 0, 0},
				{filepath.Join(testDir, "foo.md"), 1, 12, 1, 13},
				{filepath.Join(testDir, "test.md"), 2, 5, 0, 5},
				{filepath.Join(testDir, "test.txt"), 1, 5, 14, 19},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			// Use the public method instead of the private helper
			got := tt.dc.GetHeaderAndRows()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DirCounter GetHeaderAndRows() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportCSV(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,0,0,0,0\n%s,1,12,1,13\n%s,2,5,0,5\n%s,1,5,14,19",
		filepath.Join(testDir, "empty.md"),
		filepath.Join(testDir, "foo.md"),
		filepath.Join(testDir, "test.md"),
		filepath.Join(testDir, "test.txt"),
	)
	tests := []struct {
		name string
		dc   *wcg.DirCounter
		want string
	}{
		{
			name: "ExportCSV",
			dc:   wcg.NewDirCounter(testDir),
			want: expectedCSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			got, err := tt.dc.ExportCSV()
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("DirCounter.ExportCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportCSVWithFileName(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,0,0,0,0\n%s,1,12,1,13\n%s,2,5,0,5\n%s,1,5,14,19",
		filepath.Join(testDir, "empty.md"),
		filepath.Join(testDir, "foo.md"),
		filepath.Join(testDir, "test.md"),
		filepath.Join(testDir, "test.txt"),
	)
	tests := []struct {
		name string
		dc   *wcg.DirCounter
		want string
	}{
		{
			name: "ExportCSVWithFileName",
			dc:   wcg.NewDirCounter(testDir),
			want: expectedCSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			got, err := tt.dc.ExportCSV("test.csv")
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
			if _, err := os.Stat("test.csv"); err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("DirCounter.ExportCSV() = %v, want %v", got, tt.want)
			}
			err = os.Remove("test.csv")
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
		})
	}
}

func TestDirCounter_ExportTable(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedTbl := table.NewWriter()
	expectedTbl.AppendHeader(wcg.Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"})
	rows := []table.Row{
		{filepath.Join(testDir, "empty.md"), 0, 0, 0, 0},
		{filepath.Join(testDir, "foo.md"), 1, 12, 1, 13},
		{filepath.Join(testDir, "test.md"), 2, 5, 0, 5},
		{filepath.Join(testDir, "test.txt"), 1, 5, 14, 19},
	}
	expectedTbl.AppendRows(rows)
	tests := []struct {
		name string
		dc   *wcg.DirCounter
		want string
	}{
		{
			name: "ExportTable",
			dc:   wcg.NewDirCounter(testDir),
			want: expectedTbl.Render(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			if got := tt.dc.ExportTable(); got != tt.want {
				t.Errorf("DirCounter.ExportTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportExcel(t *testing.T) {
	fc := wcg.NewFileCounter("testdata")
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

func TestDirCounter_EnableTotal(t *testing.T) {
	tests := []struct {
		name string
		dc   *wcg.DirCounter
		want bool
	}{
		{
			name: "EnableTotal",
			dc:   wcg.NewDirCounter("testdata"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Count files first to have data
			tt.dc.Count()

			// Get rows before enabling total
			rowsBefore := tt.dc.GetRows()

			// Enable total
			tt.dc.EnableTotal()

			// Get rows after enabling total
			rowsAfter := tt.dc.GetRows()

			// Check if total row was added
			if tt.want && len(rowsAfter) != len(rowsBefore)+1 {
				t.Errorf("DirCounter.EnableTotal() should add total row, before: %d, after: %d", len(rowsBefore), len(rowsAfter))
			}
		})
	}
}

func TestDirCounter_GetFileCounters(t *testing.T) {
	dc := wcg.NewDirCounter("testdata")
	dc.Count()

	// Test the new GetFileCounters method
	fcs := dc.GetFileCounters()
	if fcs == nil {
		t.Error("DirCounter.GetFileCounters() returned nil")
	}

	if len(fcs) == 0 {
		t.Error("DirCounter.GetFileCounters() returned empty slice")
	}
}

func TestDirCounter_GetIgnoreList(t *testing.T) {
	ignorePatterns := []string{"*.tmp", "node_modules"}
	dc := wcg.NewDirCounter("testdata", ignorePatterns...)

	// Test the new GetIgnoreList method
	ignoreList := dc.GetIgnoreList()
	if !reflect.DeepEqual(ignoreList, ignorePatterns) {
		t.Errorf("DirCounter.GetIgnoreList() = %v, want %v", ignoreList, ignorePatterns)
	}
}

func TestDirCounter_GetHeader(t *testing.T) {
	// Test GetHeader with empty DirCounter
	dc := wcg.NewDirCounter("nonexistent")
	header := dc.GetHeader()
	expectedHeader := wcg.Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"}
	if !reflect.DeepEqual(header, expectedHeader) {
		t.Errorf("GetHeader() for empty DirCounter = %v, want %v", header, expectedHeader)
	}

	// Test GetHeader with files
	testDir := filepath.Join(wd, "testdata")
	dc = wcg.NewDirCounter(testDir)
	err := dc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	header = dc.GetHeader()
	if !reflect.DeepEqual(header, expectedHeader) {
		t.Errorf("GetHeader() for populated DirCounter = %v, want %v", header, expectedHeader)
	}
}

func TestDirCounter_ExportExcelMethod(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	dc := wcg.NewDirCounter(testDir)
	err := dc.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	// Test Excel export
	excelFile := "test_dir.xlsx"
	err = dc.ExportExcel(excelFile)
	if err != nil {
		t.Errorf("ExportExcel failed: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(excelFile); os.IsNotExist(err) {
		t.Errorf("ExportExcel did not create file")
	}
	defer os.Remove(excelFile)
}

func TestDirCounter_IsIgnoredWithError(t *testing.T) {
	// Test with valid patterns first
	dc1 := wcg.NewDirCounter("testdata")
	dc1.AddIgnorePattern("*.txt")
	dc1.AddIgnorePattern("/specific.md")

	tests1 := []struct {
		name        string
		filename    string
		wantIgnored bool
		wantError   bool
	}{
		{
			name:        "Match txt pattern",
			filename:    "test.txt",
			wantIgnored: true,
			wantError:   false,
		},
		{
			name:        "Match specific file",
			filename:    "specific.md",
			wantIgnored: true,
			wantError:   false,
		},
		{
			name:        "No match",
			filename:    "test.md",
			wantIgnored: false,
			wantError:   false,
		},
	}

	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			ignored, err := dc1.IsIgnoredWithError(tt.filename)
			if (err != nil) != tt.wantError {
				t.Errorf("IsIgnoredWithError() error = %v, wantError %v", err, tt.wantError)
			}
			if ignored != tt.wantIgnored {
				t.Errorf("IsIgnoredWithError() ignored = %v, want %v", ignored, tt.wantIgnored)
			}
		})
	}

	// Test with invalid pattern separately
	dc2 := wcg.NewDirCounter("testdata")
	dc2.AddIgnorePattern("[") // Invalid pattern that will cause filepath.Match to error

	t.Run("Invalid pattern", func(t *testing.T) {
		ignored, err := dc2.IsIgnoredWithError("test.invalid")
		if err == nil {
			t.Errorf("Expected error when using invalid pattern, but got none")
		}
		if ignored {
			t.Errorf("IsIgnoredWithError() ignored = %v, want false", ignored)
		}
	})
}

// TestDirCounter_IsIgnoredEdgeCases tests edge cases for IsIgnored function
func TestDirCounter_IsIgnoredEdgeCases(t *testing.T) {
	// Test with invalid glob patterns that should be silently ignored
	dc := wcg.NewDirCounter("testdata")
	dc.AddIgnorePattern("[")     // Invalid pattern
	dc.AddIgnorePattern("*.txt") // Valid pattern
	dc.AddIgnorePattern("\\")    // Another invalid pattern on some systems

	// Test that invalid patterns don't cause crashes and valid patterns still work
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "Valid pattern match",
			filename: "test.txt",
			want:     true,
		},
		{
			name:     "Invalid pattern should not match",
			filename: "test.invalid",
			want:     false,
		},
		{
			name:     "Another file with valid pattern",
			filename: "another.txt",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dc.IsIgnored(tt.filename)
			if got != tt.want {
				t.Errorf("IsIgnored(%v) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestDirCounter_IsIgnoredAbsolutePaths tests absolute path patterns
func TestDirCounter_IsIgnoredAbsolutePaths(t *testing.T) {
	dc := wcg.NewDirCounter("testdata")
	dc.AddIgnorePattern("/exact-match.txt")
	dc.AddIgnorePattern("/another-exact.md")

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "Exact match with leading slash",
			filename: "exact-match.txt",
			want:     true,
		},
		{
			name:     "Another exact match",
			filename: "another-exact.md",
			want:     true,
		},
		{
			name:     "No match",
			filename: "no-match.txt",
			want:     false,
		},
		{
			name:     "Partial match should not work",
			filename: "prefix-exact-match.txt",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dc.IsIgnored(tt.filename)
			if got != tt.want {
				t.Errorf("IsIgnored(%v) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestDirCounter_ProcessFilesConcurrently tests concurrent file processing
func TestDirCounter_ProcessFilesConcurrently(t *testing.T) {
	// Create multiple test files to test concurrent processing
	testFiles := []struct {
		name    string
		content string
	}{
		{"testdata/concurrent1.txt", "Hello 世界 1"},
		{"testdata/concurrent2.txt", "测试 content 2"},
		{"testdata/concurrent3.txt", "More 内容 here 3"},
		{"testdata/concurrent4.txt", "Final 测试 file 4"},
	}

	// Create test files
	for _, tf := range testFiles {
		err := os.WriteFile(tf.name, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.name, err)
		}
		defer os.Remove(tf.name)
	}

	// Test concurrent processing
	dc := wcg.NewDirCounter("testdata")
	err := dc.Count()
	if err != nil {
		t.Fatalf("Failed to count directory: %v", err)
	}

	// Verify that all files were processed
	fileCounters := dc.GetFileCounters()
	if len(fileCounters) == 0 {
		t.Errorf("Expected at least some file counters, got 0")
	}

	// Verify that concurrent processing produced correct results
	totalFiles := 0
	for _, fc := range fileCounters {
		if fc.Lines > 0 || fc.TotalChars > 0 {
			totalFiles++
		}
	}

	if totalFiles == 0 {
		t.Errorf("Expected at least some files to be counted")
	}
}

// TestDirCounter_ProcessFilesConcurrentlyRelativePathError tests relative path error handling
func TestDirCounter_ProcessFilesConcurrentlyRelativePathError(t *testing.T) {
	// Create test files in a complex directory structure to test relative path calculation
	testFiles := []struct {
		name    string
		content string
	}{
		{"testdata/subdir/deep/file1.txt", "Hello 世界 1"},
		{"testdata/subdir/file2.txt", "测试 content 2"},
		{"testdata/file3.txt", "More 内容 here 3"},
	}

	// Create directories and files
	for _, tf := range testFiles {
		dir := filepath.Dir(tf.name)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(tf.name, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.name, err)
		}
		defer os.Remove(tf.name)
	}
	defer os.RemoveAll("testdata/subdir")

	// Test with relative path mode to trigger relative path calculation
	dc := wcg.NewDirCounterWithPathMode("testdata", wcg.PathDisplayRelative)
	err := dc.Count()
	if err != nil {
		t.Fatalf("Failed to count directory: %v", err)
	}

	// Verify that files were processed correctly
	fileCounters := dc.GetFileCounters()
	if len(fileCounters) == 0 {
		t.Errorf("Expected at least some file counters, got 0")
	}

	// Check that relative paths are used in the display
	rows := dc.GetRows()
	foundRelativePath := false
	for _, row := range rows {
		if len(row) > 0 {
			filename := row[0].(string)
			// Should not contain absolute path prefix
			if !strings.HasPrefix(filename, "/") && strings.Contains(filename, ".txt") {
				foundRelativePath = true
				break
			}
		}
	}

	if !foundRelativePath {
		t.Errorf("Expected to find relative paths in output")
	}
}
