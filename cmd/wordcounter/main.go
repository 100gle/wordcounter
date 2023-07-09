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

func runDirCounter(dirPath string) {
	ignores := wcg.DiscoverIgnoreFile()
	ignores = append(ignores, excludePattern...)

	counter := wcg.NewDirCounter(dirPath, ignores...)
	if withTotal {
		counter.EnableTotal()
	}
	if err := counter.Count(); err != nil {
		log.Fatal(err)
	}

	switch exportType {
	case "csv":
		csvData, err := counter.ExportCSV(exportPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(csvData)
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(counter.ExportTable())
	}
}

func runFileCounter(filePath string) {
	counter := wcg.NewFileCounter(filePath)
	if err := counter.Count(); err != nil {
		log.Fatal(err)
	}

	switch exportType {
	case "csv":
		csvData, err := counter.ExportCSV(exportPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(csvData)
	case "excel":
		if err := counter.ExportExcel(exportPath); err != nil {
			log.Fatal(err)
		}
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
	err := serverCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "dir", "count from file or directory: dir or file")
	rootCmd.Flags().StringVarP(&exportType, "export", "e", "table", "export type: table, csv, or excel. table is default")
	rootCmd.Flags().StringVarP(&exportPath, "exportPath", "", "counter.xlsx", "export path only for csv and excel")
	rootCmd.Flags().StringArrayVarP(&excludePattern, "exclude", "", []string{}, "you can specify multiple patterns by call multiple times")
	rootCmd.Flags().BoolVarP(&withTotal, "total", "", false, "enable total count only work for mode=dir")

	serverCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "host")
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "port")
	rootCmd.AddCommand(serverCmd)
}
