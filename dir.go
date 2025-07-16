package wordcounter

import (
	"os"
	"path/filepath"
	"runtime"
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
	absPath := ToAbsolutePath(dc.dirname)

	// First pass: collect all files to process
	var filePaths []string
	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if dc.IsIgnored(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if !dc.IsIgnored(path) {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Second pass: process files concurrently with worker pool
	return dc.processFilesConcurrently(filePaths)
}

// processFilesConcurrently processes files using a worker pool pattern while preserving order
func (dc *DirCounter) processFilesConcurrently(filePaths []string) error {
	// Determine optimal number of workers
	numWorkers := runtime.NumCPU()
	if numWorkers < MinWorkers {
		numWorkers = MinWorkers
	}
	if numWorkers > MaxWorkers {
		numWorkers = MaxWorkers
	}
	if len(filePaths) < numWorkers {
		numWorkers = len(filePaths)
	}

	// Create a job structure that includes index to preserve order
	type job struct {
		index    int
		filePath string
	}

	type result struct {
		index int
		fc    *FileCounter
		err   error
	}

	// Create channels for work distribution and result collection
	jobs := make(chan job, len(filePaths))
	results := make(chan result, len(filePaths))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				fc := NewFileCounter(j.filePath)
				err := fc.Count()
				results <- result{index: j.index, fc: fc, err: err}
			}
		}()
	}

	// Send jobs to workers
	go func() {
		defer close(jobs)
		for i, filePath := range filePaths {
			jobs <- job{index: i, filePath: filePath}
		}
	}()

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and preserve order
	resultMap := make(map[int]*FileCounter)
	for i := 0; i < len(filePaths); i++ {
		res := <-results
		if res.err != nil {
			return res.err
		}
		resultMap[res.index] = res.fc
	}

	// Build final slice in correct order
	dc.Fcs = make([]*FileCounter, len(filePaths))
	for i := 0; i < len(filePaths); i++ {
		dc.Fcs[i] = resultMap[i]
	}

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
				// Log the error but don't fail the entire operation
				// Invalid patterns are treated as non-matching
				continue
			}
			if match {
				return true
			}
		}
	}
	return false
}

// IsIgnoredWithError checks if a file should be ignored and returns any pattern matching errors
func (dc *DirCounter) IsIgnoredWithError(filename string) (bool, error) {
	for _, pattern := range dc.IgnoreList {
		if strings.HasPrefix(pattern, "/") {
			if pattern[1:] == filename {
				return true, nil
			}
		} else {
			match, err := filepath.Match(pattern, filepath.Base(filename))
			if err != nil {
				return false, NewPatternMatchError(pattern, err)
			}
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func (dc *DirCounter) Ignore(pattern string) {
	dc.IgnoreList = append(dc.IgnoreList, pattern)
}

// AddIgnorePattern adds a new ignore pattern (implements IgnoreChecker interface)
func (dc *DirCounter) AddIgnorePattern(pattern string) {
	dc.Ignore(pattern)
}

// GetHeader returns the header row (implements Counter interface)
func (dc *DirCounter) GetHeader() Row {
	if len(dc.Fcs) == 0 {
		return Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"}
	}
	return dc.Fcs[0].GetHeader()
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
