# wordcounter

wordcounter is a tool to mainly count *Chinese* words in a file like Markdown, Plain Text, etc. I create it for my writing word count stats purpose.

```plain
$ ./wcg.exe ./testdata
+---------------------------------------+--------------+------------+------------+
| FILE                                  | CHINESECHARS | SPACECHARS | TOTALCHARS |
+---------------------------------------+--------------+------------+------------+
| D:\Repos\wordcounter\testdata\foo.md  | 12           | 1          | 13         |
| D:\Repos\wordcounter\testdata\test.md | 4            | 1          | 6          |
+---------------------------------------+--------------+------------+------------+
```

## how to use?

### use command line

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

```shell
git clone https://github.com/100gle/wordcounter
cd wordcounter

# linux/macOS
go build -o wcg .

# windows
go build -o wcg.exe .
```

### use library

```shell
go get -u github.com/100gle/wordcounter
```

there is three optional counter for you:

`TextCounter` is a fundamental counter for `FileCounter` and `DirCounter`.

```go
package main

import (
    wcg "github.com/100gle/wordcounter"
)

func main() {
    counter := wcg.NewTextCounter()

    // Count string
    counter.Count("你好，世界")
    // Count []byte
    counter.Count([]byte{"你好，世界"})
}
```

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
```

`DirCounter` is a counter based on `FileCounter` for directory. It will recursively count the item if it is a directory.

```go
package main

import (
    wcg "github.com/100gle/wordcounter"
)

func main() {
    ignores := []string{"*.png", "*.jpg", "**/*.js"}
    counter := wcg.NewDirCounter(ignores...)
    counter.Count("testdata")

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
