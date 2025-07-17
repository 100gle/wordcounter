package wordcounter

// Export types
const (
	ExportTypeTable = "table"
	ExportTypeCSV   = "csv"
	ExportTypeExcel = "excel"
)

// Mode types
const (
	ModeDir  = "dir"
	ModeFile = "file"
)

// Default values
const (
	DefaultExportPath = "counter.xlsx"
	DefaultHost       = "127.0.0.1"
	DefaultPort       = 8080
	DefaultMode       = ModeDir
	DefaultExportType = ExportTypeTable
)

// Server configuration
const (
	ServerAppName    = "WordCounter"
	APIVersion       = "v1"
	APIBasePath      = "/" + APIVersion + "/wordcounter"
	PingEndpoint     = APIBasePath + "/ping"
	CountEndpoint    = APIBasePath + "/count"
)

// File patterns
const (
	IgnoreFileName = ".wcignore"
)

// Worker pool configuration
const (
	// MinWorkers is the minimum number of workers in the pool
	MinWorkers = 1
	// MaxWorkers is the maximum number of workers in the pool
	MaxWorkers = 32
)
