package wordcounter

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xuri/excelize/v2"
)

type Row = []any

// ExportToCSV exports data to CSV format
func ExportToCSV(data []Row, filename ...string) (string, error) {
	if len(data) == 0 {
		return "", NewInvalidInputError("no data to export")
	}

	w := table.NewWriter()
	w.AppendHeader(data[0])
	for _, row := range data[1:] {
		w.AppendRow(row)
	}

	csvData := w.RenderCSV()
	if len(filename) > 0 && filename[0] != "" {
		absPath, err := ToAbsolutePathWithError(filename[0])
		if err != nil {
			return "", NewExportError("CSV export", err)
		}

		file, err := os.Create(absPath)
		if err != nil {
			return "", NewFileWriteError(absPath, err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		records := ConvertToSliceOfString(data)
		if err := writer.WriteAll(records); err != nil {
			return "", NewFileWriteError(absPath, err)
		}
	}
	return csvData, nil
}

// ExportToExcel exports data to Excel format
func ExportToExcel(data []Row, filename ...string) error {
	if len(data) == 0 {
		return NewInvalidInputError("no data to export")
	}

	f := excelize.NewFile()
	defer f.Close()

	defaultFilename := "counter.xlsx"
	if len(filename) > 0 {
		absPath, err := ToAbsolutePathWithError(filename[0])
		if err != nil {
			return NewExportError("Excel export", err)
		}
		defaultFilename = absPath
	}

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return NewExportError("Excel export - create sheet", err)
	}

	for rowIndex, row := range data {
		if err := f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", rowIndex+1), &row); err != nil {
			return NewExportError(fmt.Sprintf("Excel export - set row %d", rowIndex+1), err)
		}
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(defaultFilename); err != nil {
		return NewFileWriteError(defaultFilename, err)
	}
	return nil
}

// ExportToTable exports data to table format
func ExportToTable(data []Row) string {
	if len(data) == 0 {
		return ""
	}

	w := table.NewWriter()
	w.AppendHeader(data[0])
	for _, row := range data[1:] {
		w.AppendRow(row)
	}

	return w.Render()
}

// GetHeaderAndRows is a helper function that combines header and rows from a Counter
func GetHeaderAndRows(counter Counter) []Row {
	header := counter.GetHeader()
	rows := counter.GetRows()

	result := make([]Row, 0, len(rows)+1)
	result = append(result, header)
	result = append(result, rows...)

	return result
}

// ExportCounterCSV exports a Counter to CSV format
func ExportCounterCSV(counter Counter, filename ...string) (string, error) {
	data := GetHeaderAndRows(counter)
	return ExportToCSV(data, filename...)
}

// ExportCounterExcel exports a Counter to Excel format
func ExportCounterExcel(counter Counter, filename ...string) error {
	data := GetHeaderAndRows(counter)
	return ExportToExcel(data, filename...)
}

// ExportCounterTable exports a Counter to table format
func ExportCounterTable(counter Counter) string {
	data := GetHeaderAndRows(counter)
	return ExportToTable(data)
}
