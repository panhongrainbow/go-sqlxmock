package sqlmock

import (
	"database/sql"
	"fmt"
	"strings"
)

// Mocker represents a struct for managing database or mocking.
type Mocker struct {
	// for mock and genuine
	db *sql.DB // Pointer to the active database connection.
	// for mock
	sqlMock Sqlmock
	// for genuine
	dsn string // Data Source Name (DSN) for the database connection.
	// for config
	config MockerOptions // Configuration options for the mocker.
}

// NewMocker creates a new Mocker instance with the provided options.
func NewMocker(mOpts MockerOptions) (mocker *Mocker, err error) {
	mocker = &Mocker{
		config: mOpts,
	}

	if mOpts.Basic.UseDB {
		// for genuine situation

		// Construct the Data Source Name (DSN) for the database connection.
		ds := mocker.config.DB.DS
		mocker.dsn = ds.User + ":" + ds.Password +
			"@" + ds.Protocal +
			"(" + ds.IP + ":" + ds.Port + ")/" +
			ds.DbName

		// Open a database connection using the specified driver and DSN.
		mocker.db, err = sql.Open(mocker.config.DB.DS.Driver, mocker.dsn)
		if err != nil {
			return
		}
	} else {
		// Create a new SQL mock for testing.
		mocker.db, mocker.sqlMock, err = New()

		// Prepare SQL mock data.
		for i := 0; i < len(mOpts.Mock.ConfigFile); i++ {
			err = LoadMockConfig(mocker.sqlMock, mOpts.Mock.ConfigSubFolder, mOpts.Mock.ConfigFile[i]) // Load the mock configuration for testing.
			if err != nil {
				return
			}
		}

	}

	// for mock situation
	return
}

// DSN returns the Data Source Name (DSN) used for the database connection.
func (m *Mocker) DSN() string {
	return m.dsn
}

// Query executes a SQL query and returns the resulting rows.
//
//go:inline
func (m *Mocker) Query(query string, args ...any) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

// Exec executes SQL statements.
//
//go:inline
func (m *Mocker) Exec(query string, args ...any) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

// Close closes the database connection.
func (m *Mocker) Close() {
	if m.db != nil {
		_ = m.db.Close()
	}
	m = nil
}

// 定义一个 Erace Action 避免混用
type EraseTableAction string

const (
	EraseDropTableAction     EraseTableAction = "DROP TABLE IF EXISTS "
	EraseTruncateTableAction EraseTableAction = "TRUNCATE TABLE "
)

// EraseTable drops specified tables in the given database.
// (Data in the database will be deleted or cleared here, please be cautious ☢️.)
func (m *Mocker) EraseTable(action EraseTableAction, db string, tables ...string) (err error) {
	// When using mock, refrain from performing table deletion actions.
	if m.config.Basic.UseDB {
		// Join the table names and formulate the SQL query.
		all := strings.Join(tables, ", ")
		query := fmt.Sprintf(string(action)+"%s.%s;", db, all)

		// Execute the DROP TABLE query.
		_, err = m.db.Exec(query)
		if err != nil {
			return
		}
	}

	return
}
