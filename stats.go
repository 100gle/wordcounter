package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Stats struct {
	ChineseChars int
	SpaceChars   int
	TotalChars   int
}

func (s *Stats) ToRow() Row {
	return Row{
		fmt.Sprintf("%d", s.ChineseChars),
		fmt.Sprintf("%d", s.SpaceChars),
		fmt.Sprintf("%d", s.TotalChars),
	}
}


func (s *Stats) Header() Row {
	var headers Row
	t := reflect.TypeOf(*s)
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}
	return headers
}

func (s *Stats) HeaderAndRows() []Row {
	return []Row{
		{"Chinese Chars", "Space Chars", "Total Chars"},
		s.ToRow(),
	}
}

func (s *Stats) ToTable() string {
	var sb strings.Builder

	sb.WriteString("+---------------+---------+\n")
	sb.WriteString("| Character Type|  Count  |\n")
	sb.WriteString("+---------------+---------+\n")
	sb.WriteString(fmt.Sprintf("| Chinese Chars | %7d |\n", s.ChineseChars))
	sb.WriteString(fmt.Sprintf("|   Space Chars | %7d |\n", s.SpaceChars))
	sb.WriteString(fmt.Sprintf("|   Total Chars | %7d |\n", s.TotalChars))
	sb.WriteString("+---------------+---------+\n")

	return sb.String()
}

func (s *Stats) ToCsv(delimiter ...rune) (string, error) {
	var sb strings.Builder

	writer := csv.NewWriter(&sb)
	if len(delimiter) > 0 {
		writer.Comma = delimiter[0]
	}

	rows := [][]string{
		{"Chinese Chars", "Space Chars", "Total Chars"},
		{fmt.Sprintf("%d", s.ChineseChars), fmt.Sprintf("%d", s.SpaceChars), fmt.Sprintf("%d", s.TotalChars)},
	}

	for _, row := range rows {
		err := writer.Write(row)
		if err != nil {
			return "", err
		}
	}

	writer.Flush()

	return sb.String(), nil
}

// ToExcel generates an Excel file from the Stats struct using Excelize library.
func (s *Stats) ToExcel(filename ...string) error {
	f := excelize.NewFile()
	defer f.Close()

	defaultFilename := "counter.xlsx"
	if len(filename) > 0 {
		defaultFilename = filename[0]
	}

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	csvData, err := s.ToCsv()
	if err != nil {
		return err
	}

	rows := [][]string{}
	reader := csv.NewReader(strings.NewReader(csvData))
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		rows = append(rows, row)
	}

	for index, row := range rows {
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
