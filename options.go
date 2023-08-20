package sqlmock

// The following design utilizes [Function Options Pattern].
// Reference: https://www.sohamkamani.com/golang/options-pattern/

// SetMockOptsFunc is a type of function that sets options for MockOptions.
type SetMockOptsFunc func(*MockerOptions)

// WithBasicOptions is a function that creates a SetOptsFunc to set BasicOptions.
func WithBasicOptions(basicOpts BasicOptions) SetMockOptsFunc {
	return func(lockerOpts *MockerOptions) {
		lockerOpts.Basic = basicOpts
	}
}

// WithMockOptions is a function that creates a SetOptsFunc to set BasicOptions.
func WithMockOptions(mockOpts MockOptions) SetMockOptsFunc {
	return func(lockerOpts *MockerOptions) {
		lockerOpts.Mock = mockOpts
	}
}

// WithDBOptions is a function that creates a SetOptsFunc to set BasicOptions.
func WithDBOptions(dbOpts MockOptions) SetMockOptsFunc {
	return func(lockerOpts *MockerOptions) {
		lockerOpts.Mock = dbOpts
	}
}

// MockerOptions is the collection of configuration.
type MockerOptions struct {
	Basic BasicOptions // Struct field to hold Basic configuration options.
	Mock  MockOptions  // Struct field to hold Mocking configuration options.
	DB    DBOptions    // Struct field to hold Database (DB) configuration options.
}

// NewMockerOptions is a function that creates a new instance of MockerOptions with the provided options.
func NewMockerOptions(funcs ...SetMockOptsFunc) MockerOptions {
	optCollection := MockerOptions{}

	// Apply each SetOptsFunc to the optCollection to set the corresponding options
	for _, eachFunc := range funcs {
		eachFunc(&optCollection)
	}

	return optCollection
}

// BasicOptions holds basic configuration options.
type BasicOptions struct {
	UseDB bool // UseDB indicates whether to use a database.
}

// MockOptions holds options related to mocking.
type MockOptions struct {
	// Placeholder for any specific mocking configuration options.
}

// DBOptions holds database configuration options.
type DBOptions struct {
	TruncateTable bool // TruncateTable indicates whether to truncate database tables.
}
