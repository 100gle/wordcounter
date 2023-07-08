package wordcounter_test

import (
	"reflect"
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestStats_ToRow(t *testing.T) {
	tests := []struct {
		name string
		s    *wcg.Stats
		want wcg.Row
	}{
		{
			name: "Test 1",
			s: &wcg.Stats{
				Lines:           20,
				NonChineseChars: 10,
				ChineseChars:    10,
				TotalChars:      30,
			},
			want: wcg.Row{
				20,
				10,
				10,
				30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToRow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stats.ToRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_Header(t *testing.T) {
	tests := []struct {
		name string
		s    *wcg.Stats
		want wcg.Row
	}{
		{
			name: "Test 1",
			s: &wcg.Stats{
				Lines:           20,
				ChineseChars:    10,
				NonChineseChars: 10,
				TotalChars:      30,
			},
			want: wcg.Row{
				"Lines",
				"ChineseChars",
				"NonChineseChars",
				"TotalChars",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stats.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_HeaderAndRows(t *testing.T) {
	tests := []struct {
		name string
		s    *wcg.Stats
		want []wcg.Row
	}{
		{
			name: "Test 1",
			s: &wcg.Stats{
				Lines:           20,
				ChineseChars:    10,
				NonChineseChars: 10,
				TotalChars:      30,
			},
			want: []wcg.Row{
				{
					"Lines",
					"ChineseChars",
					"NonChineseChars",
					"TotalChars",
				},
				{
					20,
					10,
					10,
					30,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.HeaderAndRows(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stats.HeaderAndRows() = %v, want %v", got, tt.want)
			}
		})
	}
}
