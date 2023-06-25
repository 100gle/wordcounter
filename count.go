package main

import (
	"errors"
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
