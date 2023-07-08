package wordcounter

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DirCounter struct {
	dirname    string
	IgnoreList []string
	Fcs        []*FileCounter
	Exporter   *Exporter
	WithTotal  bool
}

func NewDirCounter(dirname string, ignores ...string) *DirCounter {
	exporter := NewExporter()

	return &DirCounter{
		IgnoreList: ignores,
		dirname:    dirname,
		Fcs:        []*FileCounter{},
		Exporter:   exporter,
		WithTotal:  false,
	}
}

func (dc *DirCounter) EnableTotal() {
	dc.WithTotal = true
}

func (dc *DirCounter) Count() error {
	var wg sync.WaitGroup

	absPath := ToAbsolutePath(dc.dirname)

	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fc := NewFileCounter(path)

		if info.IsDir() {
			if dc.IsIgnored(path) {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fc.Count()
		}()

		dc.Fcs = append(dc.Fcs, fc)
		return nil

	})

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (dc *DirCounter) IsIgnored(filename string) bool {
	for _, pattern := range dc.IgnoreList {
		if strings.HasPrefix(pattern, "/") {
			if pattern[1:] == filename {
				return true
			}
		} else {
			match, err := filepath.Match(pattern, filepath.Base(filename))
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

func (dc *DirCounter) Ignore(pattern string) {
	dc.IgnoreList = append(dc.IgnoreList, pattern)
}

func (dc *DirCounter) GetRows() []Row {
	data := make([]Row, 0, len(dc.Fcs))

	for _, fc := range dc.Fcs {
		row := fc.GetRow()
		data = append(data, row)
	}

	if dc.WithTotal {
		data = append(data, GetTotal(dc.Fcs))
	}

	return data
}

func (dc *DirCounter) GetHeaderAndRows() []Row {
	data := make([]Row, 0, len(dc.Fcs))
	header := dc.Fcs[0].GetHeader()
	data = append(data, header)
	data = append(data, dc.GetRows()...)

	return data
}

func (dc *DirCounter) ExportCSV(filename ...string) (string, error) {
	data := dc.GetHeaderAndRows()
	csvData, err := dc.Exporter.ExportCSV(data, filename...)
	if err != nil {
		return "", err
	}
	return csvData, nil
}

func (dc *DirCounter) ExportExcel(filename ...string) error {
	data := dc.GetHeaderAndRows()
	return dc.Exporter.ExportExcel(data, filename...)
}

func (dc *DirCounter) ExportTable() string {
	data := dc.GetHeaderAndRows()
	return dc.Exporter.ExportTable(data)
}
