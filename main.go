package main

import (
	"fmt"
)

func main() {
	counter := &Counter{}

	// Add some ignore patterns
	counter.Ignore(".gitignore")
	counter.Ignore("/example.txt")
	counter.Ignore("\\.txt$")

	// Count from a file
	// err = counter.CountFile("example.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Count from a directory (with concurrent file counting)
	err := counter.CountDir("testdata")
	if err != nil {
		fmt.Println(err)
	}

	counter.Table()
}
