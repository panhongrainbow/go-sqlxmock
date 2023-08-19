package sqlmock

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test_Convert_ConvertStringFormats tests ConvertStringFormats function with different cases.
func Test_Convert_ConvertStringFormats(t *testing.T) {
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
	model := []struct {
		ID        int
		Name      string
		Subject   string
		Score     float64
		CreatedAt time.Time // `insert:"skip"`
	}{
		{
			ID:        1,
			Name:      "A",
			Subject:   "B",
			Score:     2.1,
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			Name:      "C",
			Subject:   "D",
			Score:     4.2,
			CreatedAt: time.Now(),
		},
	}

	MakeInsertTableSQLStr("student_score", model, Case_No_Change)
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
