package wordcounter

import (
	"bufio"
	"os"
	"strings"
)

func DiscoverIgnoreFile() []string {
	ignores := []string{}
	file, err := os.Open(IgnoreFileName)
	if err != nil {
		return ignores
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") || scanner.Text() == "" {
			continue
		}
		ignores = append(ignores, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return ignores
	}

	return ignores
}
