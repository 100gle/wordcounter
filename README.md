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

Features:

- **Generate statistics for content**: About the number of lines, Chinese characters, non-Chinese characters, and total characters in the document content (with the option to include a total count).
- **Support for single file**: if you prefer not to install plugins or rely on specific writing applications, this feature should suit your needs.
- **Support exporting tabular statistics**: You can export the results to CSV or Excel. this feature is useful when you have more file or directory to count.
- **Support multiple platforms**: You can use this tool on Linux, macOS, and Windows as well.
- **Support ignore mode**: you can create a file named `.wcignore` which is similar to `.gitignore` can record patterns of files or folders that should not be scanned for counting, or specific patterns can be specified directly in the command line.
- **Provide count API for automation**: In server mode, you can combine your automation tools like Automator, [Keyboard Maestro](https://www.keyboardmaestro.com/main/) to count content of a file.

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
package wordcounter

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
package wordcounter

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
