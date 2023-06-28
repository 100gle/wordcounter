/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	mode           string
	exportType     string
	exportPath     string
	excludePattern []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wcg",
	Short: "wordcounter is a simple tool that counts the chinese characters in a file",
	Run:   runWordCounter,
}

func runWordCounter(cmd *cobra.Command, args []string) {
	switch mode {
	case "dir":
		runDirCounter(args[0])
	case "file":
		runFileCounter(args[0])
	default:
		log.Fatal("Invalid mode. Choose either dir or file")
	}
}

func runDirCounter(filePath string) {
	ignores := DiscoverIgnoreFile()
	ignores = append(ignores, excludePattern...)

	counter := NewDirCounter(ignores...)
	if err := counter.Count(filePath); err != nil {
		log.Fatal(err)
	}

	switch exportType {
	case "csv":
		fmt.Println(counter.ExportCSV())
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(counter.ExportTable())
	}
}

func runFileCounter(filePath string) {
	ignores := DiscoverIgnoreFile()
	ignores = append(ignores, excludePattern...)

	counter := NewFileCounter(filePath, ignores...)
	if err := counter.Count(); err != nil {
		log.Fatal(err)
	}

	switch exportType {
	case "csv":
		fmt.Println(counter.ExportCSV())
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(counter.ExportTable())
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "dir", "count from file or directory: dir or file")
	rootCmd.Flags().StringVarP(&exportType, "export", "e", "table", "export type: table, csv, or excel. table is default")
	rootCmd.Flags().StringVarP(&exportPath, "exportPath", "", "counter.xlsx", "export path only for excel")
	rootCmd.Flags().StringArrayVarP(&excludePattern, "exclude", "", []string{}, "you can specify multiple patterns by call multiple times")
}
