package wordcounter_test

import (
	"os"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/jedib0t/go-pretty/table"
)

func TestExportToCSV(t *testing.T) {
	data := []wcg.Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}
	expected := "Name,Age,Gender\nAlice,25,Female\nBob,30,Male"

	result, err := wcg.ExportToCSV(data)
	if err != nil {
		t.Errorf("ExportToCSV failed with error: %v", err)
	}

	if result != expected {
		t.Errorf("ExportToCSV failed. Expected:\n%v, got:\n%v", expected, result)
	}
}

func TestExportToCSVWithFilename(t *testing.T) {
	data := []wcg.Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}

	filename := "test.csv"
	expected := "Name,Age,Gender\nAlice,25,Female\nBob,30,Male"
	csvData, err := wcg.ExportToCSV(data, filename)
	if err != nil {
		t.Errorf("ExportToCSV failed with error: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("ExportToCSV did not create file: %v", err)
	}

	if csvData != expected {
		t.Errorf("ExportToCSV failed. Expected:\n%v, got:\n%v", expected, csvData)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("ExportToCSV could not delete file: %v", err)
	}
}

func TestExportToExcel(t *testing.T) {
	data := []wcg.Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}
	filename := "test.xlsx"

	err := wcg.ExportToExcel(data, filename)

	if err != nil {
		t.Errorf("ExportToExcel failed with error: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("ExportToExcel did not create file: %v", err)
	}

	// Clean up file
	err = os.Remove(filename)
	if err != nil {
		t.Errorf("ExportToExcel could not delete file: %v", err)
	}
}

func TestExportToTable(t *testing.T) {
	data := []wcg.Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}

	expectedTbl := table.NewWriter()
	expectedTbl.AppendHeader(data[0])
	for _, row := range data[1:] {
		expectedTbl.AppendRow(row)
	}

	expected := expectedTbl.Render()
	result := wcg.ExportToTable(data)

	if result != expected {
		t.Errorf("ExportToTable failed. Expected: \n%v\nGot: \n%v", expected, result)
	}
}
