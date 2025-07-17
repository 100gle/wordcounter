package wordcounter

// Countable defines the interface for counting operations
type Countable interface {
	// Count performs the counting operation
	Count() error
	// GetHeader returns the header row for export
	GetHeader() Row
	// GetRows returns the data rows for export
	GetRows() []Row
}

// CharacterCounter defines the interface for character counting
type CharacterCounter interface {
	// Count counts characters in the given input
	Count(input any) error
	// CountBytes counts characters from byte slice
	CountBytes(data []byte) error
	// GetStats returns the counting statistics
	GetStats() *Stats
}

// IgnoreChecker defines the interface for checking if files should be ignored
type IgnoreChecker interface {
	// IsIgnored checks if a file should be ignored
	IsIgnored(filename string) bool
	// IsIgnoredWithError checks if a file should be ignored and returns any errors
	IsIgnoredWithError(filename string) (bool, error)
	// AddIgnorePattern adds a new ignore pattern
	AddIgnorePattern(pattern string)
}

// Server defines the interface for server operations
type Server interface {
	// Run starts the server on the specified port
	Run(port int) error
	// Count handles the count request
	Count(ctx any) error
}
