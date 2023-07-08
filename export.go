package wordcounter

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xuri/excelize/v2"
)

type Row = []interface{}

type TabularExporter interface {
	ExportCSV(data []Row, filename ...string) (string, error)
	ExportExcel(data []Row, filename ...string) error
	ExportTable(data []Row) string
}

type Exporter struct {
	w table.Writer
}

func NewExporter() *Exporter {
	w := table.NewWriter()
	return &Exporter{w: w}
}

func (e *Exporter) ExportCSV(data []Row, filename ...string) (string, error) {

	e.w.AppendHeader(data[0])
	for _, row := range data[1:] {
		e.w.AppendRow(row)
	}

	csvData := e.w.RenderCSV()
	if len(filename) > 0 && filename[0] != "" {
		absPath := ToAbsolutePath(filename[0])
		file, err := os.Create(absPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		records := ConvertToSliceOfString(data)
		err = writer.WriteAll(records)
		if err != nil {
			return "", err
		}
	}
	return csvData, nil
}

func (e *Exporter) ExportExcel(data []Row, filename ...string) error {
	f := excelize.NewFile()
	defer f.Close()

	defaultFilename := "counter.xlsx"
	if len(filename) > 0 {
		defaultFilename = ToAbsolutePath(filename[0])
	}

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}

	for index, row := range data {
		err = f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", index+1), &row)
		if err != nil {
			return err
		}
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(defaultFilename); err != nil {
		fmt.Println(err)
	}
	return nil
}

func (e *Exporter) ExportTable(data []Row) string {
	e.w.AppendHeader(data[0])
	for _, row := range data[1:] {
		e.w.AppendRow(row)
	}

	return e.w.Render()

}
