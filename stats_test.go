package main

import (
    "reflect"
    "testing"
)

func TestStats_ToRow(t *testing.T) {
    tests := []struct {
        name string
        s    *Stats
        want Row
    }{
        {
            name: "Test 1",
            s: &Stats{
                ChineseChars: 10,
                SpaceChars:   20,
                TotalChars:   30,
            },
            want: Row{
                "10",
                "20",
                "30",
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
        s    *Stats
        want Row
    }{
        {
            name: "Test 1",
            s: &Stats{
                ChineseChars: 10,
                SpaceChars:   20,
                TotalChars:   30,
            },
            want: Row{
                "ChineseChars",
                "SpaceChars",
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
        s    *Stats
        want []Row
    }{
        {
            name: "Test 1",
            s: &Stats{
                ChineseChars: 10,
                SpaceChars:   20,
                TotalChars:   30,
            },
            want: []Row{
                {
                    "ChineseChars",
                    "SpaceChars",
                    "TotalChars",
                },
                {
                    "10",
                    "20",
                    "30",
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
