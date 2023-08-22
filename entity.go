package sqlmock

import (
	"database/sql"
	"fmt"
	"strings"
)

// Mocker represents a struct for managing database or mocking.
type Mocker struct {
	db     *sql.DB       // Pointer to the active database connection.
	dsn    string        // Data Source Name (DSN) for the database connection.
	config MockerOptions // Configuration options for the mocker.
}

// NewMocker creates a new Mocker instance with the provided options.
func NewMocker(mOpts MockerOptions) (mocker *Mocker, err error) {
	mocker = &Mocker{
		config: mOpts,
	}

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

	return
}

// DSN returns the Data Source Name (DSN) used for the database connection.
func (m *Mocker) DSN() string {
	return m.dsn
}

// Query executes a SQL query and returns the resulting rows.
//
//go:inline
func (m *Mocker) Query(sqlStr string) (*sql.Rows, error) {
	return m.db.Query(sqlStr)
}

// Exec executes SQL statements.
//
//go:inline
func (m *Mocker) Exec(sqlStr string) (sql.Result, error) {
	return m.db.Exec(sqlStr)
}

// Close closes the database connection.
func (m *Mocker) Close() {
	if m.db != nil {
		_ = m.db.Close()
	}
}

// DropTable drops specified tables in the given database.
func (m *Mocker) DropTable(db string, tables ...string) (err error) {
	// Join the table names and formulate the SQL query.
	all := strings.Join(tables, ", ")
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s;", db, all)

	// Execute the DROP TABLE query.
	_, err = m.db.Exec(query)
	if err != nil {
		return
	}

	return
}
