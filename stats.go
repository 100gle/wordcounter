package main

import (
	"fmt"
	"reflect"
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
		s.Header(),
		s.ToRow(),
	}
}
