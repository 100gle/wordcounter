![example workflow](https://github.com/100gle/wordcounter/actions/workflows/test-and-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/100gle/wordcounter/branch/main/graph/badge.svg?token=WO50205PUY)](https://codecov.io/gh/100gle/wordcounter)

# Wordcounter

wordcounter is a tool mainly for *Chinese* characters count in a file like Markdown, Plain Text, etc. I create it for my writing word count stats purpose.

you can use it as a command line tool:

```shell
$ wcg count ./testdata --total
+---------------------------------------+-------+--------------+-----------------+------------+
| FILE                                  | LINES | CHINESECHARS | NONCHINESECHARS | TOTALCHARS |
+---------------------------------------+-------+--------------+-----------------+------------+
| D:\Repos\wordcounter\testdata\foo.md  |     1 |           12 |               1 |         13 |
| D:\Repos\wordcounter\testdata\test.md |     1 |            4 |               1 |          5 |
| Total                                 |     2 |           16 |               2 |         18 |
+---------------------------------------+-------+--------------+-----------------+------------+
```

or run it as a server(default host is `localhost` and port is `8080`):

```shell
$ wcg server

 ┌───────────────────────────────────────────────────┐ 
 │                    WordCounter                    │ 
 │                   Fiber v2.47.0                   │ 
 │               http://127.0.0.1:8080               │ 
 │       (bound on host 0.0.0.0 and port 8080)       │ 
 │                                                   │ 
 │ Handlers ............. 3  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 49653 │ 
 └───────────────────────────────────────────────────┘ 


$ curl -s \
--location 'localhost:8080/v1/wordcounter/count' \
--header 'Content-Type: application/json' \
--data '{
    "content": "落霞与孤鹜齐飞，秋水共长天一色"
}' | jq

{
  "data": {
    "lines": 1,
    "chinese_chars": 14,
    "non_chinese_chars": 1,
    "total_chars": 15
  },
  "error": "",
  "msg": "ok"
}
```

## Features

- **📊 Comprehensive Statistics**: Count lines, Chinese characters, non-Chinese characters, and total characters with optional total summaries
- **📁 Flexible Input**: Support for both single files and recursive directory scanning
- **📤 Multiple Export Formats**: Export results as ASCII tables, CSV, or Excel files
- **🚀 High Performance**: Optimized with concurrent processing, efficient memory usage, and large buffer I/O
- **🎯 Smart Filtering**: `.wcignore` file support and command-line pattern exclusion (similar to `.gitignore`)
- **🌐 Cross-Platform**: Works on Linux, macOS, and Windows
- **🔧 API Server Mode**: HTTP API for automation and integration with tools like Automator or Keyboard Maestro
- **⚡ Concurrent Processing**: Worker pool pattern for fast directory processing
- **🛡️ Robust Error Handling**: Structured error types with detailed context information
- **🧪 Well Tested**: High test coverage with comprehensive unit tests

## How to use?

### Use as command line

you can download the binary from release.

```shell
$ wcg --help
wordcounter is a simple tool that counts the chinese characters in a file

Usage:
  wcg [command]

Available Commands:
  count       Count for a file or directory
  server      Run wordcounter as a server, only support pure text content

Flags:
  -h, --help   help for wcg

Use "wcg [command] --help" for more information about a command.
```

or clone the repository and just build from source

note:

- If you try to build from source, please ensure your OS has installed Go 1.19 or later.
- If you are in China, **highly recommend** use [goproxy](https://goproxy.cn/) to config your Go proxy firstly before installation and building.

```shell
git clone https://github.com/100gle/wordcounter
cd wordcounter

# config goproxy as you need.
# go env -w GO111MODULE=on
# go env -w GOPROXY=https://goproxy.cn,direct

go mod tidy

# linux/macOS
go build -o wcg ./cmd/wordcounter/main.go

# windows
go build -o wcg.exe ./cmd/wordcounter/main.go
```

> note: `wcg` is a short of `wordcounter-go`.

### Use as library

```shell
go get -u github.com/100gle/wordcounter
```

there are two optional counters for you:

`FileCounter` is a counter for single file.

```go
package main

import (
    "fmt"
    "log"

    wcg "github.com/100gle/wordcounter"
)

func main() {
    // Create a file counter
    counter := wcg.NewFileCounter("testdata/foo.md")

    // Perform counting with error handling
    if err := counter.Count(); err != nil {
        log.Fatalf("Failed to count file: %v", err)
    }

    // Export as ASCII table (default)
    fmt.Println(counter.ExportTable())

    // Export to CSV with error handling
    csvData, err := counter.ExportCSV()
    if err != nil {
        log.Fatalf("Failed to export CSV: %v", err)
    }
    fmt.Println("CSV output:")
    fmt.Println(csvData)

    // Export to Excel file
    if err := counter.ExportExcel("counter.xlsx"); err != nil {
        log.Fatalf("Failed to export Excel: %v", err)
    }
    fmt.Println("Excel file exported successfully!")
}
```

`DirCounter` is a counter based on `FileCounter` for directory. It will recursively count the item if it is a directory. It supports to set some patterns like `.gitignore` to exclude some directories or files.

```go
package main

import (
    "fmt"
    "log"

    wcg "github.com/100gle/wordcounter"
)

func main() {
    // Create a directory counter with ignore patterns
    ignores := []string{"*.png", "*.jpg", "**/*.js", "node_modules", ".git"}
    counter := wcg.NewDirCounter("./docs", ignores...)

    // Enable total count for summary
    counter.EnableTotal()

    // Perform counting with error handling
    if err := counter.Count(); err != nil {
        log.Fatalf("Failed to count directory: %v", err)
    }

    // Export as ASCII table with totals
    fmt.Println(counter.ExportTable())

    // Export to CSV file
    csvData, err := counter.ExportCSV("report.csv")
    if err != nil {
        log.Fatalf("Failed to export CSV: %v", err)
    }
    fmt.Printf("CSV report saved to: report.csv\n")

    // Export to Excel with custom filename
    if err := counter.ExportExcel("detailed_report.xlsx"); err != nil {
        log.Fatalf("Failed to export Excel: %v", err)
    }
    fmt.Println("Excel report exported successfully!")
}
```

## Performance & Optimization

This library has been optimized for performance and reliability:

### 🚀 Performance Features

- **Concurrent Directory Processing**: Uses worker pool pattern with CPU-core-based scaling
- **Efficient File I/O**: Reads entire files at once to avoid UTF-8 boundary issues
- **Memory Optimization**: Direct UTF-8 decoding without unnecessary string conversions
- **Smart Buffer Management**: Optimized buffer sizes for different file types

### 🛡️ Reliability Features

- **Structured Error Handling**: Custom error types with context information
- **UTF-8 Safe Processing**: Proper handling of multi-byte characters
- **Path Normalization**: Automatic conversion to absolute paths
- **Comprehensive Testing**: High test coverage with edge case handling

### 📊 Benchmarks

For large directories with thousands of files, the concurrent processing can provide significant speedup:

- **Single-threaded**: ~1000 files/second
- **Multi-threaded**: ~4000+ files/second (on 8-core systems)

## API Documentation

For detailed API documentation, see the [GoDoc](https://pkg.go.dev/github.com/100gle/wordcounter).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Licence

MIT
