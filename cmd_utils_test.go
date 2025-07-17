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
