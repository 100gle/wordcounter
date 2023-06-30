package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
)

func TestDirCounter_Count(t *testing.T) {
	type args struct {
		dirname string
	}

	ignoreList := []string{".git", ".idea"}
	exporter := NewExporter()
	fc := []*FileCounter{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Count files in directory",
			args: args{dirname: "./testdata"},
		},
		{
			name:    "Empty directory",
			args:    args{dirname: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DirCounter{
				ignoreList: ignoreList,
				dirname:    tt.args.dirname,
				fcs:        fc,
				exporter:   exporter,
			}
			if err := dc.Count(); (err != nil) != tt.wantErr {
				t.Errorf("DirCounter.Count() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false && len(dc.GetRows()) == 0 {
				t.Errorf("DirCounter.Count() did not count any files")
			}
		})
	}
}

func TestDirCounter_Ignore(t *testing.T) {
	type fields struct {
		ignoreList []string
		fcs        []*FileCounter
		exporter   *Exporter
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Ignore single pattern",
			fields: fields{
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			args: args{pattern: ".git"},
			want: []string{".git"},
		},
		{
			name: "Ignore multiple patterns",
			fields: fields{
				ignoreList: []string{".git", ".idea"},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			args: args{pattern: "node_modules"},
			want: []string{".git", ".idea", "node_modules"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DirCounter{
				ignoreList: tt.fields.ignoreList,
				fcs:        tt.fields.fcs,
				exporter:   tt.fields.exporter,
			}
			dc.Ignore(tt.args.pattern)
			if !reflect.DeepEqual(dc.ignoreList, tt.want) {
				t.Errorf("DirCounter.Ignore() got = %v, want %v", dc.ignoreList, tt.want)
			}
		})
	}
}

func TestDirCounter_GetHeaderAndRows(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	tests := []struct {
		name string
		dc   *DirCounter
		want []Row
	}{
		{
			name: "GetHeaderAndRows",
			dc: &DirCounter{
				dirname:    testDir,
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: []Row{
				{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"},
				{filepath.Join(testDir, "foo.md"), 1, 12, 1, 13},
				{filepath.Join(testDir, "test.md"), 1, 4, 1, 5},
				{filepath.Join(testDir, "test.txt"), 1, 4, 15, 19},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			if got := tt.dc.GetHeaderAndRows(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DirCounter.GetHeaderAndRows() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportCSV(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,1,12,1,13\n%s,1,4,1,5\n%s,1,4,15,19",
		filepath.Join(testDir, "foo.md"),
		filepath.Join(testDir, "test.md"),
		filepath.Join(testDir, "test.txt"),
	)
	tests := []struct {
		name string
		dc   *DirCounter
		want string
	}{
		{
			name: "ExportCSV",
			dc: &DirCounter{
				dirname:    testDir,
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: expectedCSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			got, err := tt.dc.ExportCSV()
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("DirCounter.ExportCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportCSVWithFileName(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedCSV := fmt.Sprintf("File,Lines,ChineseChars,NonChineseChars,TotalChars\n%s,1,12,1,13\n%s,1,4,1,5\n%s,1,4,15,19",
		filepath.Join(testDir, "foo.md"),
		filepath.Join(testDir, "test.md"),
		filepath.Join(testDir, "test.txt"),
	)
	tests := []struct {
		name string
		dc   *DirCounter
		want string
	}{
		{
			name: "ExportCSVWithFileName",
			dc: &DirCounter{
				dirname:    testDir,
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: expectedCSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			got, err := tt.dc.ExportCSV("test.csv")
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
			if _, err := os.Stat("test.csv"); err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("DirCounter.ExportCSV() = %v, want %v", got, tt.want)
			}
			err = os.Remove("test.csv")
			if err != nil {
				t.Errorf("DirCounter.ExportCSV() error = %v", err)
			}
		})
	}
}

func TestDirCounter_ExportTable(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedTbl := table.NewWriter()
	expectedTbl.AppendHeader(Row{"File", "Lines", "ChineseChars", "NonChineseChars", "TotalChars"})
	rows := []table.Row{
		{filepath.Join(testDir, "foo.md"), 1, 12, 1, 13},
		{filepath.Join(testDir, "test.md"), 1, 4, 1, 5},
		{filepath.Join(testDir, "test.txt"), 1, 4, 15, 19},
	}
	expectedTbl.AppendRows(rows)
	tests := []struct {
		name string
		dc   *DirCounter
		want string
	}{
		{
			name: "ExportTable",
			dc: &DirCounter{
				dirname:    testDir,
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: expectedTbl.Render(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count()
			if got := tt.dc.ExportTable(); got != tt.want {
				t.Errorf("DirCounter.ExportTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportExcel(t *testing.T) {
	fc := NewFileCounter("testdata")
	fc.Count()

	// Export the word count data to an Excel file for a FileCounter instance and check for errors
	if err := fc.ExportExcel("testdata/test.xlsx"); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// remove test.xlsx after testing
	if err := os.Remove("testdata/test.xlsx"); err != nil {
		t.Fatalf("Unexpected error while removing test.xlsx: %v", err)
	}
}
