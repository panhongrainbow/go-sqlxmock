package sqlmock

import "database/sql/driver"

// ValueConverterOption allows to create a sqlmock connection
// with a custom ValueConverter to support drivers with special data types.
func ValueConverterOption(converter driver.ValueConverter) func(*sqlmock) error {
	return func(s *sqlmock) error {
		s.converter = converter
		return nil
	}
}

// QueryMatcherOption allows to customize SQL query matcher
// and match SQL query strings in more sophisticated ways.
// The default QueryMatcher is QueryMatcherRegexp.
func QueryMatcherOption(queryMatcher QueryMatcher) func(*sqlmock) error {
	return func(s *sqlmock) error {
		s.queryMatcher = queryMatcher
		return nil
	}
}

// MonitorPingsOption determines whether calls to Ping on the driver should be
// observed and mocked.
//
// If true is passed, we will check these calls were expected. Expectations can
// be registered using the ExpectPing() method on the mock.
//
// If false is passed or this option is omitted, calls to Ping will not be
// considered when determining expectations and calls to ExpectPing will have
// no effect.
func MonitorPingsOption(monitorPings bool) func(*sqlmock) error {
	return func(s *sqlmock) error {
		s.monitorPings = monitorPings
		return nil
	}
}

// >>>>> >>>>> >>>>> for mocker

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
func WithDBOptions(dbOpts DBOptions) SetMockOptsFunc {
	return func(lockerOpts *MockerOptions) {
		lockerOpts.DB = dbOpts
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
	DS DataSource
	OP Operate
}

type DataSource struct {
	Driver   string
	User     string
	Password string
	Protocal string
	IP       string
	Port     string
	DbName   string
}

type Operate struct {
	DropTable     bool // DropTable indicates whether to drop the whole table.
	TruncateTable bool // TruncateTable indicates whether to truncate database tables.
}
