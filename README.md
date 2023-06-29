![example workflow](https://github.com/100gle/wordcounter/actions/workflows/test-and-coverage.yml/badge.svg)
[![codecov](https://codecov.io/gh/100gle/wordcounter/branch/main/graph/badge.svg?token=WO50205PUY)](https://codecov.io/gh/100gle/wordcounter)

# wordcounter

wordcounter is a tool to mainly count *Chinese* words in a file like Markdown, Plain Text, etc. I create it for my writing word count stats purpose.

```plain
$ ./wcg.exe ./testdata
+---------------------------------------+-------+--------------+-----------------+------------+
| FILE                                  | LINES | CHINESECHARS | NONCHINESECHARS | TOTALCHARS |
+---------------------------------------+-------+--------------+-----------------+------------+
| D:\Repos\wordcounter\testdata\foo.md  |     1 |           12 |               1 |         13 |
| D:\Repos\wordcounter\testdata\test.md |     1 |            4 |               1 |          5 |
+---------------------------------------+-------+--------------+-----------------+------------+
```

## how to use?

### use as command line

you can download the binary from release.

```shell
$ wcg --help
wordcounter is a simple tool that counts the chinese characters in a file

Usage:
  wcg [flags]

Flags:
      --exclude stringArray   you can specify multiple patterns by call multiple times
  -e, --export string         export type: table, csv, or excel. table is default (default "table")
      --exportPath string     export path only for excel (default "counter.xlsx")
  -h, --help                  help for wcg
  -m, --mode string           count from file or directory: dir or file (default "dir")
```

or clone the repository and just build from source

note:

- If you try to build from source, please ensure your OS has installed Go 1.19 or later.
- If you are in China, highly recommend use [goproxy] to config your Go proxy firstly before installation and building.

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

### use as library

```shell
go get -u github.com/100gle/wordcounter
```

there is two optional counter for you:

`FileCounter` is a counter for single file.

```go
package main

import (
    wcg "github.com/100gle/wordcounter"
)

func main() {
    ignores := []string{"*.png", "*.jpg", "**/*.js"}
    counter := wcg.NewFileCounter("testdata/foo.md", ignores...)
    counter.Count()

    // will output ascii table in console
    tbl := counter.ExportTable()
    fmt.Println(tbl)    

    // there are other optional export methods for you
    // counter.ExportCSV()
    // counter.ExportExcel("counter.xlsx")
}
```

`DirCounter` is a counter based on `FileCounter` for directory. It will recursively count the item if it is a directory.

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
    // counter.ExportExcel("counter.xlsx")
}
```

## Licence

MIT
