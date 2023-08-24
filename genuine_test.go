package sqlmock

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/panhongrainbow/go-sqlxmock/genuine/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test_Convert_ConvertStringFormats tests ConvertStringFormats function with different cases.
func Test_Check_ConvertStringFormats(t *testing.T) {
	// Define a set of test cases.
	tests := []struct {
		input    string
		caseFlag uint8
		expected string
	}{
		{"convertThisString", Case_Upper, "CONVERTTHISSTRING"},
		{"convertThisString", Case_Lower, "convertthisstring"},
		{"convertThisString", Case_Snake, "convert_this_string"},
	}

	// Iterate through the test cases and perform the conversion, then compare with expected output.
	for _, test := range tests {
		converted := ConvertStringFormats(test.input, test.caseFlag)
		assert.Equal(t, test.expected, converted, "Conversion result is incorrect.")
	}
}

// Test_Check_MakeCreateTableSQLStr is a unit test to verify the functionality of MakeCreateTableSQLStr.
func Test_Check_MakeCreateTableSQLStr(t *testing.T) {
	testCases := []struct {
		name        string
		model       interface{}
		text        uint8
		err         error
		expectedSQL string
	}{
		{
			name: "hotel",
			model: []struct {
				ID            int
				Name          string
				City          string
				Rating        float64
				PricePerNight float64
				CreatedAt     time.Time
			}{},
			text:        Case_No_Change,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS hotel (ID INT, Name VARCHAR(255), City VARCHAR(255), Rating DECIMAL(10, 2), PricePerNight DECIMAL(10, 2), CreatedAt TIMESTAMP);",
		},
		{
			name: "zoo",
			model: []struct {
				ID        int
				Name      string
				Size      int
				CreatedAt time.Time
			}{},
			text:        Case_Snake,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS zoo (id INT, name VARCHAR(255), size INT, created_at TIMESTAMP);",
		},
		{
			name: "gallery",
			model: []struct {
				ID        int
				Name      string
				Theme     string
				CreatedAt time.Time
			}{},
			text:        Case_Upper,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS gallery (ID INT, NAME VARCHAR(255), THEME VARCHAR(255), CREATEDAT TIMESTAMP);",
		},
		{
			name: "student_score",
			model: []struct {
				ID        int
				Name      string
				Subject   string
				Score     int
				CreatedAt time.Time
			}{},
			text:        Case_Lower,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS student_score (id INT, name VARCHAR(255), subject VARCHAR(255), score INT, createdat TIMESTAMP);",
		},
		{
			name: "restaurant",
			model: []struct {
				ID        int
				Name      string
				Cuisine   string
				Rating    float64
				CreatedAt time.Time `create:"CREATEAT"`
			}{},
			text:        Case_Lower,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS restaurant (id INT, name VARCHAR(255), cuisine VARCHAR(255), rating DECIMAL(10, 2), CREATEAT TIMESTAMP);",
		},
		{
			name: "kde_theme",
			model: []struct {
				ID        int
				Name      string
				Theme     string
				Rating    float64
				Quantity  int
				CreatedAt time.Time
			}{},
			text:        Case_No_Change,
			err:         nil,
			expectedSQL: "CREATE TABLE IF NOT EXISTS kde_theme (ID INT, Name VARCHAR(255), Theme VARCHAR(255), Rating DECIMAL(10, 2), Quantity INT, CreatedAt TIMESTAMP);",
		},
	}

	// Iterate through each test case.
	for _, testCase := range testCases {
		// Loop through each test case to generate SQL strings and compare them.
		sqlStr, err := MakeCreateTableSQLStr(testCase.name, testCase.model, testCase.text)
		require.NoError(t, err)
		require.Equal(t, testCase.expectedSQL, sqlStr)
	}
}

// Test_Check_MakeInsertTableSQLStr showcases detailed cases of SQL statement generation for movie data.
func Test_Check_MakeInsertTableSQLStr(t *testing.T) {
	type Movie struct {
		ID          int       `insert:"skip"` // ID field will be skipped during insertion.
		Title       string    // Title of the movie.
		Genre       string    // Genre of the movie.
		Rating      float64   // Rating of the movie.
		Duration    int       // Duration of the movie in minutes.
		ReleaseDate time.Time // Release date of the movie.
	}

	// Create a sample slice of Movie instances.
	model := []Movie{
		{1, "Inception", "Science Fiction", 8.8, 148, time.Now()},
		{2, "The Shawshank Redemption", "Drama", 9.3, 142, time.Now()},
		{3, "Pulp Fiction", "Crime", 8.9, 154, time.Now()},
		{4, "The Dark Knight", "Action", 9.0, 152, time.Now()},
		{5, "Forrest Gump", "Drama", 8.8, 142, time.Now()},
	}

	// Construct the expected SQL statement.
	expectedSQL := "INSERT INTO movies (Title, Genre, Rating, Duration, ReleaseDate) VALUES " +
		"('Inception', 'Science Fiction', 8.800, 148, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('The Shawshank Redemption', 'Drama', 9.300, 142, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('Pulp Fiction', 'Crime', 8.900, 154, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('The Dark Knight', 'Action', 9.000, 152, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('Forrest Gump', 'Drama', 8.800, 142, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "');"

	// Generate the SQL statement using the MakeInsertTableSQLStr function..
	sqlStr := MakeInsertTableSQLStr("movies", model, Case_No_Change)
	// Assert that the generated SQL matches the expected SQL.
	assert.Equal(t, expectedSQL, sqlStr)
}

// Test_Check_MakeSelectTableSQLStr validates SQL statement generation for an empty struct slice representing table schema.
func Test_Check_MakeSelectTableSQLStr(t *testing.T) {
	// Define an empty struct slice to represent the table schema.
	var m []struct {
		ID   int
		Name string
		Age  int
	}

	// Specify the table name and the expected SQL statement.
	tableName := "sample"
	expectedSQL := "SELECT ID, Name, Age FROM sample;"

	// Generate the SQL statement using MakeSelectTableSQLStr function.
	sqlStr := MakeSelectTableSQLStr(tableName, m, Case_No_Change)

	// Compare the generated SQL with the expected SQL.
	assert.Equal(t, expectedSQL, sqlStr, fmt.Sprintf("Expected SQL: %s, but got: %s", expectedSQL, sqlStr))
}

// Test_Check_FetchResultsFromRows validates FetchResultsFromRows function with mock database results and assertions for correctness.
func Test_Check_FetchResultsFromRows(t *testing.T) {
	// Create a mock database connection and obtain a mock database handle
	db, sqlMock, err := New()
	assert.NoError(t, err, "Failed to create mock database connection")
	defer db.Close()

	// Define the columns and rows you want to use for testing
	columns := []string{"col1", "col2"}
	sqlMock.ExpectQuery("SELECT").WillReturnRows(
		sqlMock.NewRows(columns).AddRow("value1-1", "value1-2").AddRow("value2-1", "value2-2"),
	)

	// Call the function to be tested
	rows, err := db.Query("SELECT ...")
	assert.NoError(t, err, "Failed to execute query")
	defer rows.Close()

	results, err := FetchResultsFromRows(rows)

	// Assertions
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, [][]string{{"value1-1", "value1-2"}, {"value2-1", "value2-2"}}, results, "Mismatch in results")

	// Make sure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "Not all expectations were met")
}

// Benchmark_Performance_FetchResultsFromRows Benchmarks two row-fetching methods with mock data and
// comparing FetchResultsFromRows function using a comparator.
func Benchmark_Performance_FetchResultsFromRows(b *testing.B) {
	// Create a mock database connection and obtain a mock database handle
	db, sqlMock, err := New()
	require.NoError(b, err)
	defer db.Close()

	// Define the columns and rows you want to use for testing
	columns := []string{"col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"}
	sqlMock.ExpectQuery("SELECT").WillReturnRows(
		sqlMock.NewRows(columns).
			AddRow("value1-1", "value1-2", "value1-3", "value1-4", "value1-5", "value1-6", "value1-7", "value1-8", "value1-9", "value1-10").
			AddRow("value2-1", "value2-2", "value2-3", "value2-4", "value2-5", "value2-6", "value2-7", "value2-8", "value2-9", "value2-10").
			AddRow("value3-1", "value3-2", "value3-3", "value3-4", "value3-5", "value3-6", "value3-7", "value3-8", "value3-9", "value3-10").
			AddRow("value4-1", "value4-2", "value4-3", "value4-4", "value4-5", "value4-6", "value4-7", "value4-8", "value4-9", "value4-10").
			AddRow("value5-1", "value5-2", "value5-3", "value5-4", "value5-5", "value5-6", "value5-7", "value5-8", "value5-9", "value5-10").
			AddRow("value6-1", "value6-2", "value6-3", "value6-4", "value6-5", "value6-6", "value6-7", "value6-8", "value6-9", "value6-10").
			AddRow("value7-1", "value7-2", "value7-3", "value7-4", "value7-5", "value7-6", "value7-7", "value7-8", "value7-9", "value7-10").
			AddRow("value8-1", "value8-2", "value8-3", "value8-4", "value8-5", "value8-6", "value8-7", "value8-8", "value8-9", "value8-10").
			AddRow("value9-1", "value9-2", "value9-3", "value9-4", "value9-5", "value9-6", "value9-7", "value9-8", "value9-9", "value9-10").
			AddRow("value10-1", "value10-2", "value10-3", "value10-4", "value10-5", "value10-6", "value10-7", "value10-8", "value10-9", "value10-10"),
	)

	// Call the function to be tested
	rows, err := db.Query("SELECT ...")
	require.NoError(b, err)
	defer rows.Close()

	b.ResetTimer()
	b.ReportAllocs()

	// In fact, the performance is nearly the same. (效能一样)

	b.Run("FetchResultsFromRows Comparator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = FetchResultsFromRowsComparator(rows)
		}
	})

	b.ResetTimer()

	b.Run("FetchResultsFromRows In Use", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = FetchResultsFromRows(rows)
		}
	})
}

// Test_Check_CompareResults validates comparison of results with detailed scenarios, ensuring correctness of differences.
func Test_Check_CompareResults(t *testing.T) {
	// Define a test table with different scenarios for result comparisons.
	tests := []struct {
		name               string
		results1           [][]string
		results2           [][]string
		expectedSame       bool
		expectedCondition  uint8
		expectedDifference []DiffPlace
	}{
		{
			name: "EqualResults",
			results1: [][]string{
				{"a", "b"},
				{"c", "d"},
			},
			results2: [][]string{
				{"a", "b"},
				{"c", "d"},
			},
			expectedSame:       true,
			expectedCondition:  Condition_The_Same,
			expectedDifference: nil,
		},
		{
			name: "DifferentResults",
			results1: [][]string{
				{"a", "b"},
				{"c", "d"},
			},
			results2: [][]string{
				{"a", "b"},
				{"e", "d"},
			},
			expectedSame:      false,
			expectedCondition: Condition_Diff_In_Value,
			expectedDifference: []DiffPlace{
				{1, 0, "c", "e"},
			},
		},
		{
			name: "ThreeDifferentElements",
			results1: [][]string{
				{"a", "b", "c", "d", "e"},
				{"f", "g", "h", "i", "j"},
				{"k", "l", "m", "n", "o"},
				{"p", "q", "r", "s", "t"},
				{"u", "v", "w", "x", "y"},
			},
			results2: [][]string{
				{"a", "b", "x", "d", "e"},
				{"f", "g", "h", "i", "j"},
				{"k", "z", "m", "n", "o"},
				{"p", "q", "r", "s", "t"},
				{"u", "v", "w", "x", "z"},
			},
			expectedSame:      false,
			expectedCondition: Condition_Diff_In_Value,
			expectedDifference: []DiffPlace{
				{0, 2, "c", "x"},
				{2, 1, "l", "z"},
				{4, 4, "y", "z"},
			},
		},
	}

	// Loop through the test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the function to be tested.
			same, condition, differences := CompareResults(test.results1, test.results2)
			// Assertions.
			assert.Equal(t, test.expectedSame, same, fmt.Sprintf("Expected same=%v, but got same=%v", test.expectedSame, same))
			assert.Equal(t, test.expectedCondition, condition)
			assert.Equal(t, test.expectedDifference, differences)
		})
	}
}

// >>>>> >>>>> >>>>> Integrated Testing

// Test_Check_Integrated_Testing tests the integrated process with specific settings.
func Test_Check_Integrated_Testing(t *testing.T) {
	// Set primary path (主要路径) ❗️.
	SetMockLocationByManual("./mock")

	t.Run("Query Once", func(t *testing.T) {
		// Define basic options for the test.
		basicOpts := BasicOptions{
			UseDB: false,
		}

		// Define options for mocking.
		mockOpts := MockOptions{
			ConfigSubFolder: "/basic", // Set sub path (次要路径) ❗️.
			ConfigFile:      []string{"select_once.json"},
		}

		// Define database options.
		dbOpts := DBOptions{}

		// Create a new mocker instance using specified options.
		mocker, err := NewMocker(
			NewMockerOptions(
				WithBasicOptions(basicOpts),
				WithMockOptions(mockOpts),
				WithDBOptions(dbOpts),
			))
		require.NoError(t, err)

		// Make sure the mocker is closed at the end of the test.
		defer mocker.Close()

		// SQL query to be executed.
		selectSql := "SELECT id, name, city, rating, price_per_night FROM hotels WHERE city = '?' AND rating >= ?;"

		// Execute the query using the mocker.
		rows, err := mocker.Query(selectSql, "New York", 4)
		require.NoError(t, err)

		// Record the current data in the database.
		var data [][]string
		data, err = FetchResultsFromRows(rows)
		require.NoError(t, err)

		// Compare the fetched data with the expected data.
		require.Equal(t, [][]string{
			{"1", "Grand Hotel", "New York", "4.5", "150"},
			{"2", "Luxury Inn", "New York", "4.2", "120"},
		}, data)
	})
	t.Run("Query Twice", func(t *testing.T) {
		// Define basic options for the test.
		basicOpts := BasicOptions{
			UseDB: false,
		}

		// Define options for mocking.
		mockOpts := MockOptions{
			ConfigSubFolder: "/basic", // Set sub path (次要路径) ❗️
			ConfigFile:      []string{"select_twice.json"},
		}

		// Define database options.
		dbOpts := DBOptions{}

		// Create a new mocker instance using specified options.
		mocker, err := NewMocker(
			NewMockerOptions(
				WithBasicOptions(basicOpts),
				WithMockOptions(mockOpts),
				WithDBOptions(dbOpts),
			))
		require.NoError(t, err)

		// Make sure the mocker is closed at the end of the test.
		defer mocker.Close()

		// SQL query to be executed.
		selectSql := "SELECT id, name, city, rating, price_per_night FROM hotels WHERE city = '?' AND rating >= ?;"

		// Execute the query using the mocker for the first city.
		rows, err := mocker.Query(selectSql, "New York", 4)
		require.NoError(t, err)

		// Record the current data in the database.
		var data [][]string
		data, err = FetchResultsFromRows(rows)
		require.NoError(t, err)

		// Compare the fetched data with the expected data.
		require.Equal(t, [][]string{
			{"1", "Grand Hotel", "New York", "4.5", "150"},
			{"2", "Luxury Inn", "New York", "4.2", "120"},
		}, data)

		// Execute the query using the mocker for the second city.
		rows, err = mocker.Query(selectSql, "Beijing", 4)
		require.NoError(t, err)

		// Record the current data in the database.
		data, err = FetchResultsFromRows(rows)
		require.NoError(t, err)

		// Compare the fetched data with the expected data.
		require.Equal(t, [][]string{
			{"3", "Grand Hotel", "Beijing", "4.5", "150"},
			{"4", "Luxury Inn", "Beijing", "4.2", "120"},
		}, data)
	})
}

// Test_Genuine_Integrated_Testing tests the process of automatically generating test databases.
// (Data in the database will be deleted or cleared here, please be cautious ☢️.)
func Test_Genuine_Integrated_Testing(t *testing.T) {
	// Start with general settings.
	basicOpts := BasicOptions{
		UseDB: true,
	}

	// Define mock options.
	mockOpts := MockOptions{}

	// Set the location for genuine files manually.
	SetGenuineLocationByManual("./genuine") // 主路径 ❗️

	// Load the dbOpts variable from the configuration file.
	var dbOpts DBOptions
	var err error
	dbOpts, err = LoadGenuineConfig("/mock_db", "mock_db.json") // 次要路径 ❗️

	// Declare a variable to hold DBOptions.
	mocker, err := NewMocker(
		NewMockerOptions(
			WithBasicOptions(basicOpts),
			WithMockOptions(mockOpts),
			WithDBOptions(dbOpts),
		))

	defer mocker.Close()

	require.NoError(t, err)
	require.Equal(t, "root:12345@tcp(127.0.0.1:3306)/mock", mocker.DSN())

	// Drop specified data in the mock database.
	// (Data in the database will be deleted or cleared here, please be cautious ☢️.)
	err = mocker.EraseTable(EraseDropTableAction, "mock", []string{"hotel"}...)
	require.NoError(t, err)

	// Generate CREATE TABLE SQL.
	sqlCreateStr, err := MakeCreateTableSQLStr("hotel", testdata.HotelExample, Case_Snake)
	require.NoError(t, err)
	expectedCreateSQL := "CREATE TABLE IF NOT EXISTS hotel (id INT, name VARCHAR(255), city VARCHAR(255), rating DECIMAL(10, 2), price_per_night DECIMAL(10, 2), description VARCHAR(255), facilities VARCHAR(255), contact_email VARCHAR(255), phone VARCHAR(255), website VARCHAR(255), created_at TIMESTAMP);"
	assert.Equal(t, expectedCreateSQL, sqlCreateStr)
	// Execute CREATE TABLE SQL.
	_, err = mocker.Exec(sqlCreateStr)
	require.NoError(t, err)

	// Generate INSERT SQL.
	var res sql.Result
	var affected int64
	sqlInsertStr := MakeInsertTableSQLStr("hotel", testdata.HotelExample, Case_Snake)
	// Execute INSERT SQL.
	res, err = mocker.Exec(sqlInsertStr)
	affected, err = res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(15), affected)

	// Generate SELECT SQL.
	sqlSelectStr := MakeSelectTableSQLStr("hotel", testdata.HotelExample, Case_Snake)
	expectedSelectSQL := "SELECT id, name, city, rating, price_per_night, description, facilities, contact_email, phone, website, created_at FROM hotel;"
	assert.Equal(t, expectedSelectSQL, sqlSelectStr)
	// Execute SELECT SQL.
	var rows *sql.Rows
	rows, err = mocker.Query(sqlSelectStr)
	require.NoError(t, rows.Err())

	// Record the current data in the database.
	var data, copyArray [][]string
	data, err = FetchResultsFromRows(rows)
	assert.Equal(t, "4.20", data[2][3])

	// Assume the database is modified, then record the current data again.
	copyArray = make([][]string, len(data))
	for i := range data {
		copyArray[i] = make([]string, len(data[i]))
		copy(copyArray[i], data[i])
	}
	// Only one piece of data is different.
	copyArray[2][3] = "X"

	// Find where the database has been modified.
	same, condition, differences := CompareResults(data, copyArray)
	// Assertions
	assert.Equal(t, false, same)
	assert.Equal(t, Condition_Diff_In_Value, condition)
	assert.Equal(t, []DiffPlace{{RowIndex: 2, ColumnIndex: 3, BeforeValue: "4.20", AfterValue: "X"}}, differences)

	// Truncate specified data in the mock database.
	// (Data in the database will be deleted or cleared here, please be cautious ☢️.)
	err = mocker.EraseTable(EraseTruncateTableAction, "mock", []string{"hotel"}...)
	require.NoError(t, err)

	// Query database records again after clearing.
	rows, err = mocker.Query(sqlSelectStr)
	require.NoError(t, rows.Err())

	// The entire query result should be empty.
	data, err = FetchResultsFromRows(rows)
	require.NoError(t, rows.Err())
	assert.Equal(t, [][]string(nil), data)
}
