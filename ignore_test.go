package wordcounter_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	wcg "github.com/100gle/wordcounter"
)

func TestDiscoverIgnoreFile(t *testing.T) {
	const ignoreFilename = ".wcignore"

	tests := []struct {
		name    string
		ignores []string
		want    []string
	}{
		{
			name:    "Successfully discover ignore file",
			ignores: []string{"**/*.js", "*.png", "*.jpg", "assets/", "image/", ".git/"},
			want:    []string{"**/*.js", "*.png", "*.jpg", "assets/", "image/", ".git/"},
		},
		{
			name: "Filter blank line and comment which starts with #",
			ignores: []string{
				"**/*.js",
				"*.png",
				"*.jpg",
				"# ignore assets relative path",
				"assets/",
				"image/",
				"# ignore git file",
				".git/",
			},
			want: []string{
				"**/*.js",
				"*.png",
				"*.jpg",
				"assets/",
				"image/",
				".git/",
			},
		},
		{
			name:    "Empty ignore file",
			ignores: []string{""},
			want:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(ignoreFilename, []byte(strings.Join(tt.ignores, "\n")), 0644)
			if err != nil {
				t.Fatalf("can't create testing ignore file: %v\n", err)
			}

			if got := wcg.DiscoverIgnoreFile(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscoverIgnoreFile() = %v, want %v", got, tt.want)
			}

			os.Remove(".wcignore")
		})
	}
	t.Run("Not exists ignore file", func(t *testing.T) {
		want := []string{}
		if got := wcg.DiscoverIgnoreFile(); !reflect.DeepEqual(got, want) {
			t.Errorf("DiscoverIgnoreFile() = %v, want %v", got, want)
		}
	})

	// Test with mixed content including empty lines and comments
	t.Run("Mixed content with empty lines", func(t *testing.T) {
		ignoreContent := []string{
			"*.log",
			"",
			"# This is a comment",
			"*.tmp",
			"",
			"# Another comment",
			"build/",
			"",
		}

		err := os.WriteFile(ignoreFilename, []byte(strings.Join(ignoreContent, "\n")), 0644)
		if err != nil {
			t.Fatalf("can't create testing ignore file: %v\n", err)
		}
		defer os.Remove(ignoreFilename)

		want := []string{"*.log", "*.tmp", "build/"}
		if got := wcg.DiscoverIgnoreFile(); !reflect.DeepEqual(got, want) {
			t.Errorf("DiscoverIgnoreFile() = %v, want %v", got, want)
		}
	})
}
