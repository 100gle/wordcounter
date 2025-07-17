package main

import (
	"fmt"
	"log"
	"os"

	wcg "github.com/100gle/wordcounter"
	"github.com/spf13/cobra"
)

var (
	mode           string
	exportType     string
	exportPath     string
	excludePattern []string
	withTotal      bool
	relativePath   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wcg",
	Short: "wordcounter is a simple tool that counts the chinese characters in a file",
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count for a file or directory",
	Run:   runWordCounter,
}

func runWordCounter(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Error: path argument is required")
	}

	path := args[0]
	if path == "" {
		log.Fatal("Error: path cannot be empty")
	}

	switch mode {
	case "dir":
		runDirCounter(path)
	case "file":
		runFileCounter(path)
	default:
		log.Fatal("Error: Invalid mode. Choose either 'dir' or 'file'")
	}
}

func runDirCounter(dirPath string) {
	// Validate directory path
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Fatalf("Error: Directory does not exist: %s", dirPath)
	}

	ignores := wcg.DiscoverIgnoreFile()
	ignores = append(ignores, excludePattern...)

	pathDisplayMode := wcg.PathDisplayAbsolute
	if relativePath {
		pathDisplayMode = wcg.PathDisplayRelative
	}

	counter := wcg.NewDirCounterWithPathMode(dirPath, pathDisplayMode, ignores...)
	if withTotal {
		counter.EnableTotal()
	}
	if err := counter.Count(); err != nil {
		log.Fatalf("Error counting files in directory: %v", err)
	}

	switch exportType {
	case "csv":
		csvData, err := counter.ExportCSV(exportPath)
		if err != nil {
			log.Fatalf("Error exporting to CSV: %v", err)
		}
		fmt.Println(csvData)
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatalf("Error exporting to Excel: %v", err)
		}
		fmt.Printf("Excel file exported to: %s\n", exportPath)
	default:
		fmt.Println(counter.ExportTable())
	}
}

func runFileCounter(filePath string) {
	// Validate file path
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("Error: File does not exist: %s", filePath)
	}

	pathDisplayMode := wcg.PathDisplayAbsolute
	if relativePath {
		pathDisplayMode = wcg.PathDisplayRelative
	}

	counter := wcg.NewFileCounterWithPathMode(filePath, pathDisplayMode)
	if err := counter.Count(); err != nil {
		log.Fatalf("Error counting characters in file: %v", err)
	}

	switch exportType {
	case "csv":
		csvData, err := counter.ExportCSV(exportPath)
		if err != nil {
			log.Fatalf("Error exporting to CSV: %v", err)
		}
		fmt.Println(csvData)
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatalf("Error exporting to Excel: %v", err)
		}
		fmt.Printf("Excel file exported to: %s\n", exportPath)
	default:
		fmt.Println(counter.ExportTable())
	}
}

var (
	host string
	port int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run wordcounter as a server, only support pure text content",
	Run:   runWordCounterServer,
}

func runWordCounterServer(cmd *cobra.Command, args []string) {
	srv := wcg.NewWordCounterServer()
	if err := srv.Run(port); err != nil {
		log.Fatal(err)
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
	countCmd.Flags().StringVarP(&mode, "mode", "m", "dir", "count from file or directory: dir or file")
	countCmd.Flags().StringVarP(&exportType, "export", "e", "table", "export type: table, csv, or excel. table is default")
	countCmd.Flags().StringVarP(&exportPath, "exportPath", "", "counter.xlsx", "export path only for csv and excel")
	countCmd.Flags().StringArrayVarP(&excludePattern, "exclude", "", []string{}, "you can specify multiple patterns by call multiple times")
	countCmd.Flags().BoolVarP(&withTotal, "total", "", false, "enable total count only work for mode=dir")
	countCmd.Flags().BoolVarP(&relativePath, "relative", "r", false, "show relative paths instead of absolute paths")

	serverCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "host")
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "port")

	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(serverCmd)
}
