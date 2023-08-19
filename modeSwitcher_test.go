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
		expectedSQL string
	}{
		{
			name: "hotel",
			model: struct {
				ID            int
				Name          string
				City          string
				Rating        float64
				PricePerNight float64
				CreatedAt     time.Time
			}{},
			text:        Case_No_Change,
			expectedSQL: "CREATE TABLE IF NOT EXISTS hotel (ID INT, Name VARCHAR(255), City VARCHAR(255), Rating DECIMAL(10, 2), PricePerNight DECIMAL(10, 2), CreatedAt TIMESTAMP);",
		},
		{
			name: "zoo",
			model: struct {
				ID        int
				Name      string
				Size      int
				CreatedAt time.Time
			}{},
			text:        Case_Snake,
			expectedSQL: "CREATE TABLE IF NOT EXISTS zoo (id INT, name VARCHAR(255), size INT, created_at TIMESTAMP);",
		},
		{
			name: "gallery",
			model: struct {
				ID        int
				Name      string
				Theme     string
				CreatedAt time.Time
			}{},
			text:        Case_Upper,
			expectedSQL: "CREATE TABLE IF NOT EXISTS gallery (ID INT, NAME VARCHAR(255), THEME VARCHAR(255), CREATEDAT TIMESTAMP);",
		},
		{
			name: "student_score",
			model: struct {
				ID        int
				Name      string
				Subject   string
				Score     int
				CreatedAt time.Time
			}{},
			text:        Case_Lower,
			expectedSQL: "CREATE TABLE IF NOT EXISTS student_score (id INT, name VARCHAR(255), subject VARCHAR(255), score INT, createdat TIMESTAMP);",
		},
		{
			name: "restaurant",
			model: struct {
				ID        int
				Name      string
				Cuisine   string
				Rating    float64
				CreatedAt time.Time `db:"CREATEAT"`
			}{},
			text:        Case_Lower,
			expectedSQL: "CREATE TABLE IF NOT EXISTS restaurant (id INT, name VARCHAR(255), cuisine VARCHAR(255), rating DECIMAL(10, 2), CREATEAT TIMESTAMP);",
		},
	}

	// Iterate through each test case.
	for _, testCase := range testCases {
		// Loop through each test case to generate SQL strings and compare them.
		sqlStr := MakeCreateTableSQLStr(testCase.name, testCase.model, testCase.text)
		require.Equal(t, testCase.expectedSQL, sqlStr)
	}
}
