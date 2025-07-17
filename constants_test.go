package wordcounter_test

import (
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestConstants(t *testing.T) {
	// Test export type constants
	if wcg.ExportTypeTable != "table" {
		t.Errorf("ExportTypeTable = %v, want 'table'", wcg.ExportTypeTable)
	}
	if wcg.ExportTypeCSV != "csv" {
		t.Errorf("ExportTypeCSV = %v, want 'csv'", wcg.ExportTypeCSV)
	}
	if wcg.ExportTypeExcel != "excel" {
		t.Errorf("ExportTypeExcel = %v, want 'excel'", wcg.ExportTypeExcel)
	}

	// Test mode constants
	if wcg.ModeDir != "dir" {
		t.Errorf("ModeDir = %v, want 'dir'", wcg.ModeDir)
	}
	if wcg.ModeFile != "file" {
		t.Errorf("ModeFile = %v, want 'file'", wcg.ModeFile)
	}

	// Test default values
	if wcg.DefaultExportPath != "counter.xlsx" {
		t.Errorf("DefaultExportPath = %v, want 'counter.xlsx'", wcg.DefaultExportPath)
	}
	if wcg.DefaultHost != "127.0.0.1" {
		t.Errorf("DefaultHost = %v, want '127.0.0.1'", wcg.DefaultHost)
	}
	if wcg.DefaultPort != 8080 {
		t.Errorf("DefaultPort = %v, want 8080", wcg.DefaultPort)
	}
	if wcg.DefaultMode != wcg.ModeDir {
		t.Errorf("DefaultMode = %v, want %v", wcg.DefaultMode, wcg.ModeDir)
	}
	if wcg.DefaultExportType != wcg.ExportTypeTable {
		t.Errorf("DefaultExportType = %v, want %v", wcg.DefaultExportType, wcg.ExportTypeTable)
	}

	// Test server configuration
	if wcg.ServerAppName != "WordCounter" {
		t.Errorf("ServerAppName = %v, want 'WordCounter'", wcg.ServerAppName)
	}
	if wcg.APIVersion != "v1" {
		t.Errorf("APIVersion = %v, want 'v1'", wcg.APIVersion)
	}
	if wcg.APIBasePath != "/v1/wordcounter" {
		t.Errorf("APIBasePath = %v, want '/v1/wordcounter'", wcg.APIBasePath)
	}
	if wcg.PingEndpoint != "/v1/wordcounter/ping" {
		t.Errorf("PingEndpoint = %v, want '/v1/wordcounter/ping'", wcg.PingEndpoint)
	}
	if wcg.CountEndpoint != "/v1/wordcounter/count" {
		t.Errorf("CountEndpoint = %v, want '/v1/wordcounter/count'", wcg.CountEndpoint)
	}

	// Test file patterns
	if wcg.IgnoreFileName != ".wcignore" {
		t.Errorf("IgnoreFileName = %v, want '.wcignore'", wcg.IgnoreFileName)
	}

	// Test worker pool configuration
	if wcg.MinWorkers != 1 {
		t.Errorf("MinWorkers = %v, want 1", wcg.MinWorkers)
	}
	if wcg.MaxWorkers != 32 {
		t.Errorf("MaxWorkers = %v, want 32", wcg.MaxWorkers)
	}
}
