package wordcounter

import (
	"bufio"
	"errors"
	"strings"
	"unicode"
)

type TextCounter struct {
	S *Stats
}

func NewTextCounter() *TextCounter {
	return &TextCounter{S: &Stats{}}
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
		c.S.Lines++
		line := scanner.Text()
		for _, r := range line {
			c.S.TotalChars++
			if unicode.In(r, unicode.Han) {
				c.S.ChineseChars++
			} else {
				c.S.NonChineseChars++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
