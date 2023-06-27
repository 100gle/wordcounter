package main

import (
	"fmt"
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
				fcs:        fc,
				exporter:   exporter,
			}
			if err := dc.Count(tt.args.dirname); (err != nil) != tt.wantErr {
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
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: []Row{
				{"File", "ChineseChars", "SpaceChars", "TotalChars"},
				{filepath.Join(testDir, "foo.md"), "12", "1", "13"},
				{filepath.Join(testDir, "test.md"), "4", "1", "6"},
				{filepath.Join(testDir, "test.txt"), "4", "2", "19"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count(testDir)
			if got := tt.dc.GetHeaderAndRows(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DirCounter.GetHeaderAndRows() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportCSV(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedCSV := fmt.Sprintf("File,ChineseChars,SpaceChars,TotalChars\n%s,12,1,13\n%s,4,1,6\n%s,4,2,19",
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
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: expectedCSV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dc.Count(testDir)
			if got := tt.dc.ExportCSV(); got != tt.want {
				t.Errorf("DirCounter.ExportCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirCounter_ExportTable(t *testing.T) {
	testDir := filepath.Join(wd, "testdata")
	expectedTbl := table.NewWriter()
	expectedTbl.AppendHeader(Row{"File", "ChineseChars", "SpaceChars", "TotalChars"})
	rows := []table.Row{
		{filepath.Join(testDir, "foo.md"), "12", "1", "13"},
		{filepath.Join(testDir, "test.md"), "4", "1", "6"},
		{filepath.Join(testDir, "test.txt"), "4", "2", "19"},
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
				ignoreList: []string{},
				fcs:        []*FileCounter{},
				exporter:   NewExporter(),
			},
			want: expectedTbl.Render(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            tt.dc.Count(testDir)
			if got := tt.dc.ExportTable(); got != tt.want {
				t.Errorf("DirCounter.ExportTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
