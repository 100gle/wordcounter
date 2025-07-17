package wordcounter_test

import (
	"strings"
	"testing"

	"github.com/100gle/wordcounter"
)

func TestCounter_Count(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Valid string",
			input:   "Hello 世界",
			wantErr: false,
		},
		{
			name:    "Valid byte slice",
			input:   []byte("Hello 世界"),
			wantErr: false,
		},
		{
			name:    "Invalid input type",
			input:   42,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := wordcounter.NewCounter()

			err := tc.Count(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Counter.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tc.S.TotalChars != 8 {
					t.Errorf("Counter.Count() total chars = %d, want 8", tc.S.TotalChars)
				}

				if tc.S.ChineseChars != 2 {
					t.Errorf("Counter.Count() chinese chars = %d, want 2", tc.S.ChineseChars)
				}

				if tc.S.Lines != 1 {
					t.Errorf("Counter.Count() space chars = %d, want 1", tc.S.Lines)
				}
			}
		})
	}
}

func TestCounter_CountBytes_Comprehensive(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedLines      int
		expectedChinese    int
		expectedNonChinese int
		expectedTotal      int
	}{
		{
			name:               "Simple Chinese text",
			input:              "你好世界",
			expectedLines:      1,
			expectedChinese:    4,
			expectedNonChinese: 0,
			expectedTotal:      4,
		},
		{
			name:               "Mixed Chinese and English",
			input:              "Hello 你好 World 世界",
			expectedLines:      1,
			expectedChinese:    4,
			expectedNonChinese: 13,
			expectedTotal:      17,
		},
		{
			name:               "Multiple lines with Chinese",
			input:              "第一行\n第二行\n第三行",
			expectedLines:      3,
			expectedChinese:    9,
			expectedNonChinese: 0,
			expectedTotal:      9,
		},
		{
			name:               "Text with Chinese punctuation",
			input:              "你好，世界！",
			expectedLines:      1,
			expectedChinese:    6, // Including Chinese punctuation
			expectedNonChinese: 0,
			expectedTotal:      6,
		},
		{
			name:               "Empty lines",
			input:              "第一行\n\n第三行",
			expectedLines:      3,
			expectedChinese:    6,
			expectedNonChinese: 0,
			expectedTotal:      6,
		},
		{
			name:               "Only English",
			input:              "Hello World",
			expectedLines:      1,
			expectedChinese:    0,
			expectedNonChinese: 11,
			expectedTotal:      11,
		},
		{
			name:               "Numbers and symbols",
			input:              "123 + 456 = 579",
			expectedLines:      1,
			expectedChinese:    0,
			expectedNonChinese: 15,
			expectedTotal:      15,
		},
		{
			name:               "Unicode symbols",
			input:              "😀😃😄 emoji test",
			expectedLines:      1,
			expectedChinese:    0,
			expectedNonChinese: 14, // Emojis count as non-Chinese (3 emojis + 11 other chars)
			expectedTotal:      14,
		},
		{
			name:               "Traditional Chinese",
			input:              "繁體中文測試",
			expectedLines:      1,
			expectedChinese:    6,
			expectedNonChinese: 0,
			expectedTotal:      6,
		},
		{
			name:               "Single character",
			input:              "中",
			expectedLines:      1,
			expectedChinese:    1,
			expectedNonChinese: 0,
			expectedTotal:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := wordcounter.NewCounter()
			err := tc.CountBytes([]byte(tt.input))

			if err != nil {
				t.Errorf("Counter.CountBytes() error = %v", err)
				return
			}

			stats := tc.GetStats()
			if stats.Lines != tt.expectedLines {
				t.Errorf("Lines = %d, want %d", stats.Lines, tt.expectedLines)
			}
			if stats.ChineseChars != tt.expectedChinese {
				t.Errorf("ChineseChars = %d, want %d", stats.ChineseChars, tt.expectedChinese)
			}
			if stats.NonChineseChars != tt.expectedNonChinese {
				t.Errorf("NonChineseChars = %d, want %d", stats.NonChineseChars, tt.expectedNonChinese)
			}
			if stats.TotalChars != tt.expectedTotal {
				t.Errorf("TotalChars = %d, want %d", stats.TotalChars, tt.expectedTotal)
			}
		})
	}
}

func TestCounter_CountBytes_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Empty byte slice",
			input:   []byte{},
			wantErr: false,
		},
		{
			name:    "Nil byte slice",
			input:   nil,
			wantErr: false,
		},
		{
			name:    "Single newline",
			input:   []byte("\n"),
			wantErr: false,
		},
		{
			name:    "Multiple newlines",
			input:   []byte("\n\n\n"),
			wantErr: false,
		},
		{
			name:    "Invalid UTF-8 sequence",
			input:   []byte{0xFF, 0xFE, 0xFD},
			wantErr: false, // Should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := wordcounter.NewCounter()
			err := tc.CountBytes(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Counter.CountBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCounter_CountBytes_EmptyData(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "Empty byte slice",
			input: []byte{},
		},
		{
			name:  "Nil byte slice",
			input: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := wordcounter.NewCounter()
			err := tc.CountBytes(tt.input)

			if err != nil {
				t.Errorf("Counter.CountBytes() error = %v, expected no error for empty data", err)
				return
			}

			stats := tc.GetStats()
			if stats.Lines != 0 || stats.ChineseChars != 0 || stats.NonChineseChars != 0 || stats.TotalChars != 0 {
				t.Errorf("Counter.CountBytes() for empty data, expected all zero stats, got: Lines=%d, ChineseChars=%d, NonChineseChars=%d, TotalChars=%d",
					stats.Lines, stats.ChineseChars, stats.NonChineseChars, stats.TotalChars)
			}
		})
	}
}

func TestCounter_LargeText(t *testing.T) {
	// Test with large text to ensure performance
	largeText := strings.Repeat("这是一个测试文本。", 10000) // 10,000 repetitions
	tc := wordcounter.NewCounter()

	err := tc.Count(largeText)
	if err != nil {
		t.Errorf("Counter.Count() error = %v", err)
		return
	}

	stats := tc.GetStats()
	// "这是一个测试文本。" has 9 characters total:
	// 8 Chinese characters: 这是一个测试文本
	// 1 Chinese punctuation: 。(U+3002) - now counted as Chinese
	expectedChinese := 90000 // 9 Chinese chars (including punctuation) * 10,000 repetitions
	expectedNonChinese := 0  // No non-Chinese characters
	expectedTotal := 90000

	if stats.ChineseChars != expectedChinese {
		t.Errorf("ChineseChars = %d, want %d", stats.ChineseChars, expectedChinese)
	}
	if stats.NonChineseChars != expectedNonChinese {
		t.Errorf("NonChineseChars = %d, want %d", stats.NonChineseChars, expectedNonChinese)
	}
	if stats.TotalChars != expectedTotal {
		t.Errorf("TotalChars = %d, want %d", stats.TotalChars, expectedTotal)
	}
}

// Benchmark tests for Counter performance

// BenchmarkCounter_CountBytes_SmallText benchmarks counting small text
func BenchmarkCounter_CountBytes_SmallText(b *testing.B) {
	text := "Hello 世界! This is a test 这是一个测试。"
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_MediumText benchmarks counting medium-sized text
func BenchmarkCounter_CountBytes_MediumText(b *testing.B) {
	// Create a medium-sized text (about 1KB)
	text := strings.Repeat("这是一个中文测试文本，包含了各种字符。Hello World! 123456789.\n", 20)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_LargeText benchmarks counting large text
func BenchmarkCounter_CountBytes_LargeText(b *testing.B) {
	// Create a large text (about 100KB)
	text := strings.Repeat("这是一个中文测试文本，包含了各种字符。Hello World! 123456789.\n", 2000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_VeryLargeText benchmarks counting very large text
func BenchmarkCounter_CountBytes_VeryLargeText(b *testing.B) {
	// Create a very large text (about 1MB)
	text := strings.Repeat("这是一个中文测试文本，包含了各种字符。Hello World! 123456789.\n", 20000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_ChineseOnly benchmarks counting Chinese-only text
func BenchmarkCounter_CountBytes_ChineseOnly(b *testing.B) {
	text := strings.Repeat("这是一个完全由中文组成的测试文本，用于测试中文字符识别的性能。\n", 1000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_EnglishOnly benchmarks counting English-only text
func BenchmarkCounter_CountBytes_EnglishOnly(b *testing.B) {
	text := strings.Repeat("This is a test text composed entirely of English characters for performance testing.\n", 1000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_CountBytes_MixedContent benchmarks counting mixed content
func BenchmarkCounter_CountBytes_MixedContent(b *testing.B) {
	text := strings.Repeat("Mixed content 混合内容: English 英文, Numbers 123456, Symbols !@#$%^&*(), Emojis 😀😃😄\n", 1000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}

// BenchmarkCounter_Count_String benchmarks the Count method with string input
func BenchmarkCounter_Count_String(b *testing.B) {
	text := strings.Repeat("这是一个测试文本 This is a test text.\n", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.Count(text)
	}
}

// BenchmarkCounter_Count_ByteSlice benchmarks the Count method with byte slice input
func BenchmarkCounter_Count_ByteSlice(b *testing.B) {
	text := strings.Repeat("这是一个测试文本 This is a test text.\n", 1000)
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.Count(data)
	}
}

// BenchmarkCounter_MultipleOperations benchmarks multiple counting operations on the same counter
func BenchmarkCounter_MultipleOperations(b *testing.B) {
	texts := []string{
		"第一段文本 First text segment",
		"第二段文本 Second text segment",
		"第三段文本 Third text segment",
		"第四段文本 Fourth text segment",
		"第五段文本 Fifth text segment",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		for _, text := range texts {
			tc.Count(text)
		}
	}
}

// BenchmarkIsChinese benchmarks the isChinese function indirectly through character counting
func BenchmarkIsChinese(b *testing.B) {
	// Create text with various character types to test isChinese performance
	text := "中文English123!@#😀"
	data := []byte(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := wordcounter.NewCounter()
		tc.CountBytes(data)
	}
}
