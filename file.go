package main

import (
	"io"
	"os"
)

type FileCounter struct {
	tc         *TextCounter
	exporter   *Exporter
	filename   string
}

func NewFileCounter(filename string) *FileCounter {
	tc := NewTextCounter()
	exporter := NewExporter()
	absPath := ToAbsolutePath(filename)

	fc := &FileCounter{
		filename: absPath,
		tc:       tc,
		exporter: exporter,
	}

	return fc
}
func (fc *FileCounter) Count() error {
	file, err := os.Open(fc.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		err = fc.tc.Count(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (fc *FileCounter) GetRow() Row {
	row := append(Row{fc.filename}, fc.tc.s.ToRow()...)
	return row
}

func (fc *FileCounter) GetHeader() Row {
	headers := append(Row{"File"}, fc.tc.s.Header()...)
	return headers
}

func (fc *FileCounter) GetHeaderAndRow() []Row {
	headers := fc.GetHeader()
	row := fc.GetRow()
	return []Row{headers, row}
}

func (fc *FileCounter) ExportCSV(filename ...string) (string, error) {
	data := fc.GetHeaderAndRow()
	csvData, err := fc.exporter.ExportCSV(data, filename...)
	if err != nil {
		return "", err
	}
	return csvData, nil
}

func (fc *FileCounter) ExportExcel(filename ...string) error {
	data := fc.GetHeaderAndRow()
	return fc.exporter.ExportExcel(data, filename...)
}

func (fc *FileCounter) ExportTable() string {
	data := fc.GetHeaderAndRow()
	return fc.exporter.ExportTable(data)
}
