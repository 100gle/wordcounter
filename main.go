package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type Counter struct {
	ChineseChars int
	SpaceChars   int
	TotalChars   int
	IgnoreList   []string
}

func NewCounter(ignores []string) *Counter {
	c := &Counter{}

	if len(ignores) > 0 {
		c.IgnoreList = append(c.IgnoreList, ignores...)
	}

	return c
}

func (c *Counter) Count(input interface{}) error {
	str := ""
	switch v := input.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	}
	if str == "" {
		return errors.New("no input provided")
	}
	for _, r := range str {
		c.TotalChars++
		if unicode.In(r, unicode.Scripts["Han"]) {
			c.ChineseChars++
		}
		if unicode.IsSpace(r) {
			c.SpaceChars++
		}
	}
	return nil
}

func (c *Counter) CountFile(filename string) error {
	// Check if the file should be ignored
	if c.isIgnored(filename) {
		return nil
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read each line of the file and count the words
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = c.Count(scanner.Bytes())
		if err != nil {
			return nil
		}
	}

	// Handle any errors that occurred while reading the file
	if err := scanner.Err(); err != nil {
		return err
	}

	// Return nil if everything was successful
	return nil
}

func (c *Counter) CountDir(dirname string) error {
	var wg sync.WaitGroup

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if c.isIgnored(path) {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.CountFile(path)
		}()

		return nil
	})

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (c *Counter) isIgnored(filename string) bool {
	for _, pattern := range c.IgnoreList {
		if strings.HasPrefix(pattern, "/") {
			if pattern[1:] == filename {
				return true
			}
		} else {
			match, err := filepath.Match(pattern, filename)
			if err != nil {
				return false
			}
			if match {
				return true
			}
		}
	}
	return false
}

func (c *Counter) Ignore(pattern string) {
	c.IgnoreList = append(c.IgnoreList, pattern)
}

func main() {
	counter := &Counter{}

	// Add some ignore patterns
	counter.Ignore(".gitignore")
	counter.Ignore("/example.txt")
	counter.Ignore("\\.txt$")

	// Count from a string
	err := counter.Count("你好，世界！Hello, world!  ")
	if err != nil {
		fmt.Println(err)
	}

	// Count from a file
	err = counter.CountFile("example.txt")
	if err != nil {
		fmt.Println(err)
	}

	// Count from a directory (with concurrent file counting)
	err = counter.CountDir("testdata")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Chinese characters: %d\n", counter.ChineseChars)
	fmt.Printf("Space characters: %d\n", counter.SpaceChars)
	fmt.Printf("Total characters: %d\n", counter.TotalChars)
}
