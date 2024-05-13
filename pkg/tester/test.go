package tester

// Test is an interface that represents a test.
type Test interface {
	// Name returns the name of the test.
	Name() string

	// Config returns the configuration of the test.
	Config() map[string]any

	// RunTest runs the test and returns the result.
	RunTest() (*Result, error)
}

// Result represents the result of a test.
type Result struct {
	// Labels contains additional labels associated with the result.
	Labels map[string]string

	// Data contains the data associated with the result.
	Data map[string]any
}
