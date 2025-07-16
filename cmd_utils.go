package wordcounter

import (
	"fmt"
	"os"
)

// ExportConfig holds configuration for export operations
type ExportConfig struct {
	Type string
	Path string
}

// CounterExporter provides common export functionality for counters
type CounterExporter struct {
	counter Counter
	config  ExportConfig
}

// NewCounterExporter creates a new CounterExporter
func NewCounterExporter(counter Counter, config ExportConfig) *CounterExporter {
	return &CounterExporter{
		counter: counter,
		config:  config,
	}
}

// Export performs the export operation based on configuration
func (ce *CounterExporter) Export() error {
	switch ce.config.Type {
	case ExportTypeCSV:
		return ce.exportCSV()
	case ExportTypeExcel:
		return ce.exportExcel()
	case ExportTypeTable:
		return ce.exportTable()
	default:
		return NewInvalidInputError(fmt.Sprintf("unsupported export type: %s", ce.config.Type))
	}
}

func (ce *CounterExporter) exportCSV() error {
	var csvData string
	var err error

	switch c := ce.counter.(type) {
	case *FileCounter:
		csvData, err = c.ExportCSV(ce.config.Path)
	case *DirCounter:
		csvData, err = c.ExportCSV(ce.config.Path)
	default:
		return NewInvalidInputError("unsupported counter type for CSV export")
	}

	if err != nil {
		return NewExportError("CSV export", err)
	}

	fmt.Println(csvData)
	return nil
}

func (ce *CounterExporter) exportExcel() error {
	var err error

	switch c := ce.counter.(type) {
	case *FileCounter:
		err = c.ExportExcel(ce.config.Path)
	case *DirCounter:
		err = c.ExportExcel(ce.config.Path)
	default:
		return NewInvalidInputError("unsupported counter type for Excel export")
	}

	if err != nil {
		return NewExportError("Excel export", err)
	}

	fmt.Printf("Excel file exported to: %s\n", ce.config.Path)
	return nil
}

func (ce *CounterExporter) exportTable() error {
	var table string

	switch c := ce.counter.(type) {
	case *FileCounter:
		table = c.ExportTable()
	case *DirCounter:
		table = c.ExportTable()
	default:
		return NewInvalidInputError("unsupported counter type for table export")
	}

	fmt.Println(table)
	return nil
}

// ValidatePath validates if a path exists
func ValidatePath(path string) error {
	if path == "" {
		return NewInvalidInputError("path cannot be empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewFileNotFoundError(path, err)
	}

	return nil
}

// ValidateExportType validates if an export type is supported
func ValidateExportType(exportType string) error {
	switch exportType {
	case ExportTypeTable, ExportTypeCSV, ExportTypeExcel:
		return nil
	default:
		return NewInvalidInputError(fmt.Sprintf("unsupported export type: %s, supported types: %s, %s, %s",
			exportType, ExportTypeTable, ExportTypeCSV, ExportTypeExcel))
	}
}

// ValidateMode validates if a mode is supported
func ValidateMode(mode string) error {
	switch mode {
	case ModeDir, ModeFile:
		return nil
	default:
		return NewInvalidInputError(fmt.Sprintf("unsupported mode: %s, supported modes: %s, %s",
			mode, ModeDir, ModeFile))
	}
}
