package wordcounter

import (
	"reflect"
)

type Stats struct {
	Lines           int `json:"lines,omitempty"`
	ChineseChars    int `json:"chinese_chars,omitempty"`
	NonChineseChars int `json:"non_chinese_chars,omitempty"`
	TotalChars      int `json:"total_chars,omitempty"`
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
