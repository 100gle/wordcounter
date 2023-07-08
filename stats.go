package wordcounter

import (
	"reflect"
)

type Stats struct {
	Lines           int
	ChineseChars    int
	NonChineseChars int
	TotalChars      int
}

func (s *Stats) ToRow() Row {
	return Row{
		s.Lines,
		s.ChineseChars,
		s.NonChineseChars,
		s.TotalChars,
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
