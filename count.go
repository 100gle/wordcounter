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
	s          Stats
	IgnoreList []string
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
		c.s.TotalChars++
		if unicode.In(r, unicode.Scripts["Han"]) {
			c.s.ChineseChars++
		}
		if unicode.IsSpace(r) {
			c.s.SpaceChars++
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

// TODO: Implement Table method to output a table with ChineseChars, SpaceChars, TotalChars stats
// Table method outputs a table with ChineseChars, SpaceChars, TotalChars stats.
func (c *Counter) Table() {
	// data, _ := c.s.ToCsv()

	// fmt.Println(data)
	err := c.s.ToExcel()
	if err != nil {
		fmt.Println(err)
	}
}
