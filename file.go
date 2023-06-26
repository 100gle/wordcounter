package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileCounter struct {
	tc         *TextCounter
	exporter   *Exporter
	filename   string
	ignoreList []string
}

func NewFileCounter(filename string, ignores ...string) *FileCounter {

	tc := NewTextCounter()
	exporter := NewExporter()
	absPath := ToAbsolutePath(filename)

	fc := &FileCounter{
		filename: absPath,
		tc:       tc,
		exporter: exporter,
	}

	if len(ignores) > 0 {
		fc.ignoreList = append(fc.ignoreList, ignores...)
	}

	return fc
}
func (fc *FileCounter) Count() error {
	if fc.isIgnored(fc.filename) {
		return nil
	}

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

func (fc *FileCounter) isIgnored(filename string) bool {
	for _, pattern := range fc.ignoreList {
		if strings.HasPrefix(pattern, "/") {
			if pattern[1:] == filename {
				return true
			}
		} else {
			match, err := filepath.Match(pattern, filename)
			if err != nil {
				return false
			}
			if match {
				return true
			}
		}
	}
	return false
}

func (fc *FileCounter) Ignore(pattern ...string) {
	fc.ignoreList = append(fc.ignoreList, pattern...)
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

func (fc *FileCounter) ExportCSV() string {
	data := fc.GetHeaderAndRow()
	return fc.exporter.ExportCSV(data)
}

func (fc *FileCounter) ExportExcel(filename ...string) error {
	data := fc.GetHeaderAndRow()
	return fc.exporter.ExportExcel(data, filename...)
}

func (fc *FileCounter) ExportTable() string {
	data := fc.GetHeaderAndRow()
	return fc.exporter.ExportTable(data)
}
