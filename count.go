package main

import (
	"bufio"
	"errors"
	"strings"
	"unicode"
)

type TextCounter struct {
	s Stats
}

func NewTextCounter() *TextCounter {
	return &TextCounter{}
}

func (c *TextCounter) Count(input interface{}) error {
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
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		c.s.Lines++
		line := scanner.Text()
		for _, r := range line {
			c.s.TotalChars++
			if unicode.In(r, unicode.Han) {
				c.s.ChineseChars++
			} else {
				c.s.NonChineseChars++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
