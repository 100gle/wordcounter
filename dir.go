package main

import (
	"os"
	"path/filepath"
	"sync"
)

type DirCounter struct {
	ignoreList []string
	fcs        []*FileCounter
	exporter   *Exporter
}

func NewDirCounter(ignores ...string) *DirCounter {
	exporter := NewExporter()

	return &DirCounter{
		ignoreList: ignores,
		fcs:        []*FileCounter{},
		exporter:   exporter,
	}
}

func (dc *DirCounter) Count(dirname string) error {
	var wg sync.WaitGroup

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fc := NewFileCounter(path, dc.ignoreList...)

		if info.IsDir() {
			if fc.isIgnored(path) {
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

		dc.fcs = append(dc.fcs, fc)
		return nil

	})

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (dc *DirCounter) Ignore(pattern string) {
	dc.ignoreList = append(dc.ignoreList, pattern)
}

func (dc *DirCounter) GetRows() []Row {
	data := make([]Row, 0, len(dc.fcs))

	for _, fc := range dc.fcs {
		row := fc.GetRow()
		data = append(data, row)
	}

	return data
}

func (dc *DirCounter) GetHeaderAndRows() []Row {
	data := make([]Row, 0, len(dc.fcs))
	header := dc.fcs[0].GetHeader()
	data = append(data, header)
	data = append(data, dc.GetRows()...)

	return data
}

func (dc *DirCounter) ExportCSV() string {
	data := dc.GetHeaderAndRows()
	return dc.exporter.ExportCSV(data)
}

func (dc *DirCounter) ExportExcel(filename ...string) error {
	data := dc.GetHeaderAndRows()
	return dc.exporter.ExportExcel(data, filename...)
}

func (dc *DirCounter) ExportTable() string {
	data := dc.GetHeaderAndRows()
	return dc.exporter.ExportTable(data)
}
