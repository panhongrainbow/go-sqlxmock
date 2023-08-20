package sqlmock

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
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
		{"convertThisString", Case_Snake, "convert_This_String"},
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

func Test_Check_MakeInsertTableSQLStr(t *testing.T) {
	type KDEtheme struct {
		ID        int `insert:"skip"`
		Name      string
		Theme     string
		Rating    float64
		Quantity  int
		CreatedAt time.Time
	}

	model := []KDEtheme{
		{1, "A", "B", 4.5, 10, time.Now()},
		{2, "C", "D", 3.2, 5, time.Now()},
		{3, "E", "F", 2.8, 8, time.Now()},
		{4, "G", "H", 4.0, 12, time.Now()},
		{5, "I", "J", 2.5, 6, time.Now()},
		{6, "K", "L", 3.7, 7, time.Now()},
		{7, "M", "N", 1.9, 15, time.Now()},
		{8, "O", "P", 4.8, 9, time.Now()},
		{9, "Q", "R", 3.6, 11, time.Now()},
		{10, "S", "T", 2.1, 3, time.Now()},
	}

	expectedSQL := "INSERT INTO kde_theme (Name, Theme, Rating, Quantity, CreatedAt) VALUES " +
		"('A', 'B', 4.500, 10, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('C', 'D', 3.200, 5, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('E', 'F', 2.800, 8, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('G', 'H', 4.000, 12, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('I', 'J', 2.500, 6, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('K', 'L', 3.700, 7, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('M', 'N', 1.900, 15, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('O', 'P', 4.800, 9, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('Q', 'R', 3.600, 11, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "'), " +
		"('S', 'T', 2.100, 3, '" + time.Now().UTC().Format("2006-01-02 15:04:05") + "');"

	sqlStr := MakeInsertTableSQLStr("kde_theme", model, Case_No_Change)
	assert.Equal(t, expectedSQL, sqlStr)
}

func Test_Check_MakeSelectTableSQLStr(t *testing.T) {
	var m []struct {
		ID   int
		Name string
		Age  int
	}

	tableName := "sample"
	expectedSQL := "SELECT ID, Name, Age FROM sample;"

	sqlStr := MakeSelectTableSQLStr(tableName, m, Case_No_Change)

	if sqlStr != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, sqlStr)
	}
}

// >>>>> >>>>> >>>>> Integrated Testing

// Test_Genuine_Integrated_Testing tests the process of automatically generating test databases.
func Test_Genuine_Integrated_Testing(t *testing.T) {
	// >>>>> DB Client

	dsn := "root:12345@tcp(127.0.0.1:3306)/mock"

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	err = db.Ping()
	require.NoError(t, err)

	_, _ = db.Exec("DROP TABLE " + "hotel")
	// require.NoError(t, err)

	// >>>>> Test Data

	type Hotel struct {
		ID            int
		Name          string
		City          string
		Rating        float64
		PricePerNight float64
		Description   string
		Facilities    string
		ContactEmail  string
		Phone         string
		Website       string
		CreatedAt     time.Time
	}

	hotel := []Hotel{
		{1, "Luxury Resort", "Miami", 4.500, 250.000, "A luxurious beachside resort", "Pool, Spa, Private Beach", "info@luxuryresort.com", "+1-123-456-7890", "https://www.luxuryresort.com", time.Date(2025, 3, 14, 9, 23, 45, 0, time.UTC)},
		{2, "Cozy Inn", "New York", 3.800, 120.000, "A charming inn in the heart of the city", "Free Wi-Fi, Breakfast, Lounge", "info@cozyinn.com", "+1-987-654-3210", "https://www.cozyinn.com", time.Date(1984, 7, 9, 9, 12, 32, 0, time.UTC)},
		{3, "Seaside Lodge", "Los Angeles", 4.200, 180.000, "A cozy lodge with ocean views", "Ocean View, Fireplace, Restaurant", "info@seasidelodge.com", "+1-555-123-4567", "https://www.seasidelodge.com", time.Date(2012, 2, 27, 9, 56, 8, 0, time.UTC)},
		{4, "Mountain Retreat", "Denver", 4.000, 150.000, "A retreat in the mountains", "Hiking Trails, Spa, Scenic Views", "info@mountainretreat.com", "+1-888-567-8901", "https://www.mountainretreat.com", time.Date(2001, 11, 2, 9, 34, 17, 0, time.UTC)},
		{5, "Urban Hotel", "Chicago", 3.700, 200.000, "A modern hotel in the city center", "Fitness Center, Rooftop Bar", "info@urbanhotel.com", "+1-333-456-7890", "https://www.urbanhotel.com", time.Date(1993, 10, 15, 10, 42, 19, 0, time.UTC)},
		{6, "Sunny Getaway", "Miami", 4.600, 280.000, "A sunny paradise by the beach", "Beachfront, Poolside Bar", "info@sunnygetaway.com", "+1-555-789-1234", "https://www.sunnygetaway.com", time.Date(2017, 5, 2, 10, 18, 41, 0, time.UTC)},
		{7, "Downtown Suites", "New York", 4.300, 180.000, "Luxury suites in the heart of the city", "Spa, Concierge, Sky Lounge", "info@downtownsuites.com", "+1-999-888-7777", "https://www.downtownsuites.com", time.Date(2029, 12, 31, 10, 3, 52, 0, time.UTC)},
		{8, "Beachfront Resort", "Los Angeles", 4.800, 320.000, "A beachfront oasis with stunning views", "Private Beach, Oceanfront Dining", "info@beachfrontresort.com", "+1-444-555-6666", "https://www.beachfrontresort.com", time.Date(2015, 9, 4, 10, 24, 37, 0, time.UTC)},
		{9, "Mountain Chalet", "Denver", 4.400, 210.000, "Charming chalets nestled in the mountains", "Ski Access, Fireplace, Spa", "info@mountainchalet.com", "+1-222-333-4444", "https://www.mountainchalet.com", time.Date(1990, 1, 23, 11, 55, 13, 0, time.UTC)},
		{10, "City Center Hotel", "Chicago", 3.900, 190.000, "Conveniently located hotel in the city center", "Business Center, On-site Restaurant", "info@citycenterhotel.com", "+1-777-888-9999", "https://www.citycenterhotel.com", time.Date(1991, 6, 30, 11, 11, 22, 0, time.UTC)},
		{11, "Lakeview Lodge", "Seattle", 4.700, 260.000, "Lakeside retreat with stunning lake views", "Boating, Fishing, Lakeside Dining", "info@lakeviewlodge.com", "+1-111-222-3333", "https://www.lakeviewlodge.com", time.Date(1980, 4, 8, 11, 9, 45, 0, time.UTC)},
		{12, "Riverside Inn", "San Francisco", 4.100, 170.000, "Charming inn along the riverside", "Scenic Views, Garden Patio", "info@riversideinn.com", "+1-444-333-2222", "https://www.riversideinn.com", time.Date(2022, 1, 1, 11, 59, 57, 0, time.UTC)},
		{13, "Historic Mansion", "Boston", 4.500, 300.000, "Elegant mansion with a rich history", "Antique Furnishings, Ballroom", "info@historicmansion.com", "+1-555-666-7777", "https://www.historicmansion.com", time.Date(2000, 8, 13, 12, 15, 0, 0, time.UTC)},
		{14, "Desert Oasis", "Phoenix", 3.600, 140.000, "A tranquil oasis in the desert", "Spa, Desert Gardens, Pool", "info@desertoasis.com", "+1-777-888-9999", "https://www.desertoasis.com", time.Date(2021, 11, 11, 12, 38, 2, 0, time.UTC)},
		{15, "Skyline Tower", "Las Vegas", 4.400, 240.000, "Modern tower with breathtaking city views", "Rooftop Pool, Casino, Nightclub", "info@skylinetower.com", "+1-111-222-3333", "https://www.skylinetower.com", time.Date(1996, 12, 31, 12, 17, 36, 0, time.UTC)},
	}

	// >>>>> Create Table

	// expectedCreateSQL := "CREATE TABLE IF NOT EXISTS hotel (id INT, name VARCHAR(255), city VARCHAR(255), rating DECIMAL(10, 2), price_per_night DECIMAL(10, 2), created_at TIMESTAMP, description VARCHAR(255), facilities VARCHAR(255), contact_email VARCHAR(255), phone VARCHAR(255), website VARCHAR(255));"

	sqlCreateStr, err := MakeCreateTableSQLStr("hotel", hotel, Case_Snake)
	require.NoError(t, err)
	// require.Equal(t, expectedCreateSQL, sqlCreateStr)

	_, err = db.Exec(sqlCreateStr)
	require.NoError(t, err)

	// >>>>> Insert Data

	sqlInsertStr := MakeInsertTableSQLStr("hotel", hotel, Case_Snake)
	var res sql.Result
	res, err = db.Exec(sqlInsertStr)
	require.NoError(t, err)
	var affected int64
	affected, _ = res.RowsAffected()
	require.Equal(t, int64(15), affected)

	// >>>>> Select Data

	sqlSelectStr := MakeSelectTableSQLStr("hotel", hotel, Case_Snake)
	var rows *sql.Rows
	rows, err = db.Query(sqlSelectStr)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice of interface{} for storing column values
	values := make([]interface{}, len(columns))
	for i := range values {
		var val sql.RawBytes
		values[i] = &val
	}

	var results [][]string

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}

		rowValues := make([]string, len(values))
		for i, val := range values {
			rowValues[i] = string(*val.(*sql.RawBytes))
		}
		results = append(results, rowValues)
	}

	require.NoError(t, rows.Err())
}
