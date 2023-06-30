package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestToAbsolutePath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	if runtime.GOOS == "windows" {
		testCases := []struct {
			desc           string
			input          string
			expectedOutput string
		}{
			{
				desc:           "Testing absolute path on Windows",
				input:          "C:\\Windows",
				expectedOutput: "C:\\Windows",
			},
			{
				desc:           "Testing relative path on Windows",
				input:          "Documents\\file.txt",
				expectedOutput: filepath.Join(wd, "Documents", "file.txt"),
			},
			{
				desc:           "Testing relative path based on current directory on Windows",
				input:          "file.txt",
				expectedOutput: filepath.Join(wd, "file.txt"),
			},
			{
				desc:           "Testing empty string on Windows",
				input:          "",
				expectedOutput: "",
			},
		}

		for _, tC := range testCases {
			absPath := ToAbsolutePath(tC.input)
			if absPath != tC.expectedOutput {
				t.Errorf("Test case: %s - ToAbsolutePath(\"%s\") = %s; want %s", tC.desc, tC.input, absPath, tC.expectedOutput)
			}
		}
	} else {
		testCases := []struct {
			desc           string
			input          string
			expectedOutput string
		}{
			{
				desc:           "Testing absolute path on Linux or macOS",
				input:          "/usr/local",
				expectedOutput: "/usr/local",
			},
			{
				desc:           "Testing relative path on Linux or macOS",
				input:          "README.md",
				expectedOutput: filepath.Join(wd, "README.md"),
			},
			{
				desc:           "Testing empty string on Linux or macOS",
				input:          "",
				expectedOutput: "",
			},
		}

		for _, tC := range testCases {
			absPath := ToAbsolutePath(tC.input)
			if absPath != tC.expectedOutput {
				t.Errorf("Test case: %s - ToAbsolutePath(\"%s\") = %s; want %s", tC.desc, tC.input, absPath, tC.expectedOutput)
			}
		}
	}
}

func TestConvertToSliceOfString(t *testing.T) {
	type args struct {
		input [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Testing convert to slice of string",
			args: args{
				input: [][]interface{}{
					{"1", "2", "3"},
					{"4", "5", "6"},
					{"7", "8", "9"},
				},
			},
			want: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			},
		},
		{
			name: "Testing convert to slice of string with empty rows",
			args: args{
				input: [][]interface{}{
					{},
					{},
					{},
				},
			},
			want: [][]string{
				{},
				{},
				{},
			},
		},
		{
			name: "Testing convert to slice of string with empty columns",
			args: args{
				input: [][]interface{}{
					{1, nil, nil},
					{nil, 2, nil},
					{nil, nil, nil},
				},
			},
			want: [][]string{
				{"1", "", ""},
				{"", "2", ""},
				{"", "", ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToSliceOfString(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToSliceOfString() = %v, want %v", got, tt.want)
			}
		})
	}
}
