package wordcounter

import (
	"io"
	"os"
)

type FileCounter struct {
	tc       *TextCounter
	exporter *Exporter
	FileName string
}

func NewFileCounter(filename string) *FileCounter {
	tc := NewTextCounter()
	exporter := NewExporter()
	absPath := ToAbsolutePath(filename)

	fc := &FileCounter{
		FileName: absPath,
		tc:       tc,
		exporter: exporter,
	}

	return fc
}
func (fc *FileCounter) Count() error {
	file, err := os.Open(fc.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileNotFoundError(fc.FileName, err)
		}
		return NewFileReadError(fc.FileName, err)
	}
	defer file.Close()

	// Read entire file at once for better performance and simpler logic
	// This avoids issues with splitting lines/characters across buffer boundaries
	data, err := io.ReadAll(file)
	if err != nil {
		return NewFileReadError(fc.FileName, err)
	}

	if err := fc.tc.CountBytes(data); err != nil {
		return NewFileReadError(fc.FileName, err)
	}

	return nil
}

func (fc *FileCounter) GetRow() Row {
	row := append(Row{fc.FileName}, fc.tc.S.ToRow()...)
	return row
}

func (fc *FileCounter) GetHeader() Row {
	headers := append(Row{"File"}, fc.tc.S.Header()...)
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
