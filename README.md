![example workflow](https://github.com/100gle/wordcounter/actions/workflows/test-and-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/100gle/wordcounter/branch/main/graph/badge.svg?token=WO50205PUY)](https://codecov.io/gh/100gle/wordcounter)

# Wordcounter

wordcounter is a tool mainly for *Chinese* characters count in a file like Markdown, Plain Text, etc. I create it for my writing word count stats purpose.

```plain
$ ./wcg.exe ./testdata --total
+---------------------------------------+-------+--------------+-----------------+------------+
| FILE                                  | LINES | CHINESECHARS | NONCHINESECHARS | TOTALCHARS |
+---------------------------------------+-------+--------------+-----------------+------------+
| D:\Repos\wordcounter\testdata\foo.md  |     1 |           12 |               1 |         13 |
| D:\Repos\wordcounter\testdata\test.md |     1 |            4 |               1 |          5 |
| Total                                 |     2 |           16 |               2 |         18 |
+---------------------------------------+-------+--------------+-----------------+------------+
```

## How to use?

### Use as command line

you can download the binary from release.

```shell
$ wcg --help
wordcounter is a simple tool that counts the chinese characters in a file

Usage:
  wcg [flags]

Flags:
      --exclude stringArray   you can specify multiple patterns by call multiple times
  -e, --export string         export type: table, csv, or excel. table is default (default "table")
      --exportPath string     export path only for csv and excel (default "counter.xlsx")
  -h, --help                  help for wcg
  -m, --mode string           count from file or directory: dir or file (default "dir")
      --total                 enable total count only work for mode=dir
```

or clone the repository and just build from source

note:

- If you try to build from source, please ensure your OS has installed Go 1.19 or later.
- If you are in China, highly recommend use [goproxy](https://goproxy.cn/) to config your Go proxy firstly before installation and building.

```shell
git clone https://github.com/100gle/wordcounter
cd wordcounter

# config goproxy as you need.
# go env -w GO111MODULE=on
# go env -w GOPROXY=https://goproxy.cn,direct

go mod download

# linux/macOS
go build -o wcg .

# windows
go build -o wcg.exe .
```

### Use as library

```shell
go get -u github.com/100gle/wordcounter
```

there are two optional counters for you:

`FileCounter` is a counter for single file.

```go
package main

import (
    wcg "github.com/100gle/wordcounter"
)

func main() {
    counter := wcg.NewFileCounter("testdata/foo.md")
    counter.Count()

    // will output ascii table in console
    tbl := counter.ExportTable()
    fmt.Println(tbl)    

    // there are other optional export methods for you
    // counter.ExportCSV()
    // counter.ExportCSV("counter.csv") // Export to specific file

    // counter.ExportExcel("counter.xlsx")
}
```

`DirCounter` is a counter based on `FileCounter` for directory. It will recursively count the item if it is a directory. It supports to set some patterns like `.gitignore` to exclude some directories or files.

```go
package main

import (
    wcg "github.com/100gle/wordcounter"
)

func main() {
    ignores := []string{"*.png", "*.jpg", "**/*.js"}
    counter := wcg.NewDirCounter("testdata", ignores...)
    counter.Count()

    // will output ascii table in console
    tbl := counter.ExportTable()
    fmt.Println(tbl)

    // there are other optional export methods for you
    // counter.ExportCSV()
    // counter.ExportCSV("counter.csv") // Export to specific file

    // counter.ExportExcel("counter.xlsx")
}
```

## Licence

MIT
