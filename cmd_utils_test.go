package wordcounter_test

import (
	"os"
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestValidatePath(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_validate_path")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid file path",
			path:    tmpFile.Name(),
			wantErr: false,
		},
		{
			name:    "Empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "Non-existent path",
			path:    "/path/that/does/not/exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wcg.ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateExportType(t *testing.T) {
	tests := []struct {
		name       string
		exportType string
		wantErr    bool
	}{
		{
			name:       "Valid table type",
			exportType: "table",
			wantErr:    false,
		},
		{
			name:       "Valid csv type",
			exportType: "csv",
			wantErr:    false,
		},
		{
			name:       "Valid excel type",
			exportType: "excel",
			wantErr:    false,
		},
		{
			name:       "Invalid type",
			exportType: "invalid",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wcg.ValidateExportType(tt.exportType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExportType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{
			name:    "Valid dir mode",
			mode:    "dir",
			wantErr: false,
		},
		{
			name:    "Valid file mode",
			mode:    "file",
			wantErr: false,
		},
		{
			name:    "Invalid mode",
			mode:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wcg.ValidateMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCounterExporter_Export(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_counter_export")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some test content
	content := "测试内容 test content"
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	counter := wcg.NewFileCounter(tmpFile.Name())
	if err := counter.Count(); err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	tests := []struct {
		name       string
		exportType string
		wantErr    bool
	}{
		{
			name:       "Export table",
			exportType: "table",
			wantErr:    false,
		},
		{
			name:       "Export CSV",
			exportType: "csv",
			wantErr:    false,
		},
		{
			name:       "Export Excel",
			exportType: "excel",
			wantErr:    false,
		},
		{
			name:       "Invalid export type",
			exportType: "invalid",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outputPath string
			switch tt.exportType {
			case "excel":
				outputPath = "test_output.xlsx"
			case "csv":
				outputPath = "test_output.csv"
			default:
				outputPath = "test_output"
			}

			config := wcg.ExportConfig{
				Type: tt.exportType,
				Path: outputPath,
			}
			exporter := wcg.NewCounterExporter(counter, config)
			err := exporter.Export()
			if (err != nil) != tt.wantErr {
				t.Errorf("CounterExporter.Export() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up created files
			if !tt.wantErr && (tt.exportType == "excel" || tt.exportType == "csv") {
				if _, err := os.Stat(outputPath); err == nil {
					os.Remove(outputPath)
				}
			}
		})
	}
}

// TestCounterExporter_ExportErrorHandling tests error handling in export functions
func TestCounterExporter_ExportErrorHandling(t *testing.T) {
	// Create a test file counter
	tmpFile, err := os.CreateTemp("", "test_export_error")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "测试内容 test content"
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	counter := wcg.NewFileCounter(tmpFile.Name())
	if err := counter.Count(); err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	tests := []struct {
		name       string
		exportType string
		path       string
		wantErr    bool
		setupFunc  func() func() // setup function that returns cleanup function
	}{
		{
			name:       "CSV export to invalid directory",
			exportType: "csv",
			path:       "/invalid/directory/test.csv",
			wantErr:    true,
			setupFunc:  func() func() { return func() {} },
		},
		{
			name:       "Excel export to invalid directory",
			exportType: "excel",
			path:       "/invalid/directory/test.xlsx",
			wantErr:    true,
			setupFunc:  func() func() { return func() {} },
		},
		{
			name:       "CSV export to read-only directory",
			exportType: "csv",
			path:       "/tmp/readonly/test.csv",
			wantErr:    true,
			setupFunc: func() func() {
				// Create read-only directory
				if err := os.MkdirAll("/tmp/readonly", 0444); err != nil {
					return func() {}
				}
				return func() { os.RemoveAll("/tmp/readonly") }
			},
		},
		{
			name:       "Excel export to read-only directory",
			exportType: "excel",
			path:       "/tmp/readonly_excel/test.xlsx",
			wantErr:    true,
			setupFunc: func() func() {
				// Create read-only directory
				if err := os.MkdirAll("/tmp/readonly_excel", 0444); err != nil {
					return func() {}
				}
				return func() { os.RemoveAll("/tmp/readonly_excel") }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupFunc()
			defer cleanup()

			config := wcg.ExportConfig{
				Type: tt.exportType,
				Path: tt.path,
			}
			exporter := wcg.NewCounterExporter(counter, config)
			err := exporter.Export()
			if (err != nil) != tt.wantErr {
				t.Errorf("CounterExporter.Export() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up any created files
			if !tt.wantErr {
				os.Remove(tt.path)
			}
		})
	}
}

// TestCounterExporter_ExportWithEmptyData tests export with empty data
func TestCounterExporter_ExportWithEmptyData(t *testing.T) {
	// Create an empty file
	tmpFile, err := os.CreateTemp("", "test_empty_export")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create counter with empty file
	counter := wcg.NewFileCounter(tmpFile.Name())
	if err := counter.Count(); err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	tests := []struct {
		name       string
		exportType string
		path       string
		wantErr    bool
	}{
		{
			name:       "CSV export with empty data",
			exportType: "csv",
			path:       "empty_test.csv",
			wantErr:    false,
		},
		{
			name:       "Excel export with empty data",
			exportType: "excel",
			path:       "empty_test.xlsx",
			wantErr:    false,
		},
		{
			name:       "Table export with empty data",
			exportType: "table",
			path:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := wcg.ExportConfig{
				Type: tt.exportType,
				Path: tt.path,
			}
			exporter := wcg.NewCounterExporter(counter, config)
			err := exporter.Export()
			if (err != nil) != tt.wantErr {
				t.Errorf("CounterExporter.Export() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up created files
			if !tt.wantErr && tt.path != "" {
				os.Remove(tt.path)
			}
		})
	}
}

// TestCounterExporter_ExportPathEdgeCases tests edge cases with export paths
func TestCounterExporter_ExportPathEdgeCases(t *testing.T) {
	// Create a test file counter
	tmpFile, err := os.CreateTemp("", "test_path_edge")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "测试 test"
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	counter := wcg.NewFileCounter(tmpFile.Name())
	if err := counter.Count(); err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	tests := []struct {
		name       string
		exportType string
		path       string
		wantErr    bool
	}{
		{
			name:       "CSV export with relative path",
			exportType: "csv",
			path:       "./relative_test.csv",
			wantErr:    false,
		},
		{
			name:       "Excel export with relative path",
			exportType: "excel",
			path:       "./relative_test.xlsx",
			wantErr:    false,
		},
		{
			name:       "CSV export with long filename",
			exportType: "csv",
			path:       "very_long_filename_that_should_still_work_fine_test.csv",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := wcg.ExportConfig{
				Type: tt.exportType,
				Path: tt.path,
			}
			exporter := wcg.NewCounterExporter(counter, config)
			err := exporter.Export()
			if (err != nil) != tt.wantErr {
				t.Errorf("CounterExporter.Export() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up created files
			if !tt.wantErr {
				os.Remove(tt.path)
			}
		})
	}
}
