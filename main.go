package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	mode := flag.String(
		"mode",
		"dir",
		"count from file or directory: dir or file",
	)

	filePath := flag.String(
		"file",
		"",
		"file to count if mode is file",
	)

	exportType := flag.String(
		"export",
		"table",
		"export type: table, csv, or excel. table is default",
	)

	exportPath := flag.String(
		"exportPath",
		"counter.xlsx",
		"export path only for excel",
	)

	flag.Parse()

	switch *mode {
	case "dir":
		counter := NewDirCounter()
		if err := counter.Count(*filePath); err != nil {
			log.Fatal(err)
		}

		switch *exportType {
		case "csv":
			fmt.Println(counter.ExportCSV())
			return
		case "excel":
			if err := counter.ExportExcel(*exportPath); err != nil {
				log.Fatal(err)
			}
			return
		default:
			fmt.Println(counter.ExportTable())
			return
		}

	case "file":
		counter := NewFileCounter(*filePath)
		if err := counter.Count(); err != nil {
			log.Fatal(err)
		}

		switch *exportType {
		case "csv":
			fmt.Println(counter.ExportCSV())
			return
		case "excel":
			if err := counter.ExportExcel(*exportPath); err != nil {
				log.Fatal(err)
			}
			return
		default:
			fmt.Println(counter.ExportTable())
			return
		}
	default:
		log.Fatal("Invalid mode. Choose either dir or file")
	}
}
