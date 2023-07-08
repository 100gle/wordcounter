package wordcounter_test

import (
	"testing"

	"github.com/100gle/wordcounter"
)

func TestTextCounter_Count(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
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
			tc := wordcounter.NewTextCounter()

			err := tc.Count(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("TextCounter.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tc.S.TotalChars != 8 {
					t.Errorf("TextCounter.Count() total chars = %d, want 8", tc.S.TotalChars)
				}

				if tc.S.ChineseChars != 2 {
					t.Errorf("TextCounter.Count() chinese chars = %d, want 2", tc.S.ChineseChars)
				}

				if tc.S.Lines != 1 {
					t.Errorf("TextCounter.Count() space chars = %d, want 1", tc.S.Lines)
				}
			}
		})
	}
}
