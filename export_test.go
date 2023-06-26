package main

import (
	"os"
	"testing"

	"github.com/jedib0t/go-pretty/table"
)

func TestExporter_ExportCSV(t *testing.T) {
	e := NewExporter()
	data := []Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}
	expected := "Name,Age,Gender\nAlice,25,Female\nBob,30,Male"

	result := e.ExportCSV(data)

	if result != expected {
		t.Errorf("ExportCSV failed. Expected:\n%v, got:\n%v", expected, result)
	}
}

func TestExporter_ExportExcel(t *testing.T) {
	e := NewExporter()
	data := []Row{
		{"Name", "Age", "Gender"},
		{"Alice", 25, "Female"},
		{"Bob", 30, "Male"},
	}
	filename := "test.xlsx"

	err := e.ExportExcel(data, filename)

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

func TestExporter_ExportTable(t *testing.T) {
	e := NewExporter()
	data := []Row{
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
	result := e.ExportTable(data)

	if result != expected {
		t.Errorf("ExportTable failed. Expected: \n%v\nGot: \n%v", expected, result)
	}
}
