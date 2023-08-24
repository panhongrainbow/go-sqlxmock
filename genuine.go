package sqlmock

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// >>>>> >>>>> >>>>> Creating Data Table Functions

// MakeCreateTableSQLStr generates a SQL string for creating a table with the given table name and model.
// Database column names prioritize tags first, followed by other settings.
// (栏位优先考虑标签)
func MakeCreateTableSQLStr(tableName string, model interface{}, convertCase uint8) (sqlStr string, err error) {
	// tableFields := make([]string, 0) // 其他常用的写法先保留
	t := reflect.TypeOf(model)

	// Enhancing Performance with String Building.
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(tableName)
	sb.WriteString(" (")

	// Iterate through the fields of the model.
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)

		columnName := ConvertStringFormats(field.Name, convertCase)

		tag := field.Tag.Get("create") // Using the 'db' Tag as the Preferred Setting Value (field 为获取结构原素资讯)

		// Check if the 'db' tag is empty, indicating that this field should be used for the table.
		if tag != "" {
			columnName = tag
		}

		// tableFields = append(tableFields, fmt.Sprintf("%s %s", field.Name, fieldType)) // 其他常用的写法先保留
		sb.WriteString(columnName)
		sb.WriteByte(' ')
		var fieldType string
		fieldType, err = mapCreateFieldType(field.Type) // Fetching Information about Structure Elements
		if err != nil {
			return
		}

		sb.WriteString(fieldType)

		// Add a comma if it's not the last field.
		if i < t.Elem().NumField()-1 {
			sb.WriteString(", ")
		}
	}

	// sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(tableFields, ", ")) // 其他常用的写法先保留

	// Close the CREATE TABLE statement.
	sb.WriteString(");")
	sqlStr = sb.String()

	return
}

// mapCreateFieldType takes a reflect.Type as input and returns the corresponding SQL data type.
func mapCreateFieldType(t reflect.Type) (typed string, err error) {
	switch t.Kind() {
	case reflect.Int, reflect.Int64:
		typed = "INT" // Map int type to SQL INT data type.
	case reflect.String:
		typed = "VARCHAR(255)" // Map string type to SQL VARCHAR(255) data type.
	case reflect.Float32, reflect.Float64:
		typed = "DECIMAL(10, 2)" // Map float type to SQL DECIMAL(10, 2) data type.
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			typed = "TIMESTAMP" // Map time.Time struct to SQL TIMESTAMP data type.
		}
	default:
		err = fmt.Errorf("database column data type is undefined")
		// Return an empty string if the type does not match any mapping.
	}

	// Return here.
	return
}

// >>>>> >>>>> >>>>> Inserting Test Data Functions

// MakeInsertTableSQLStr generates SQL INSERT statements for a given table name and model.
func MakeInsertTableSQLStr(tableName string, model interface{}, convertCase uint8) (sqlStr string) {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	// Loop through each data entry in the model.
	for j := 0; j < v.Len(); j++ {

		var column []string
		var values []string

		// Loop through each field in the model.
		for i := 0; i < t.Elem().NumField(); i++ {
			tag := t.Elem().Field(i).Tag.Get("insert") // Get the "insert" tag value.
			if tag != "skip" {
				value := v.Index(j).Field(i).Interface() // Get the field's value.

				if value != nil {
					// Convert the value to the correct format and append to the values slice.
					values = append(values,
						mapCreateFieldValue(t.Elem().Field(i).Type, value),
					)

					// Convert the field name to the desired format and append to the column names slice.
					column = append(column, fmt.Sprintf("%v",
						ConvertStringFormats(t.Elem().Field(i).Name, convertCase)),
					)
				}
			}
		}

		// Generate the INSERT SQL statement.
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(column, ", "), strings.Join(values, ", "))
		} else {
			sqlStr += ", (" + strings.Join(values, ", ") + ")"
		}
	}

	sqlStr += ";"

	return sqlStr
}

// mapCreateFieldValue converts different data types to their SQL-friendly string representations.
func mapCreateFieldValue(t reflect.Type, value interface{}) (corrected string) {
	switch t.Kind() {
	case reflect.Int:
		corrected = strconv.Itoa(value.(int)) // Convert integer to string.
	case reflect.Int64:
		corrected = strconv.FormatInt(value.(int64), 10) // Convert int64 to string.
	case reflect.String:
		corrected = "'" + value.(string) + "'" // Surround string with single quotes.
	case reflect.Float32:
		corrected = strconv.FormatFloat(float64(value.(float32)), 'f', 3, 64) // Convert float32 to string with 3 decimal places.
	case reflect.Float64:
		corrected = strconv.FormatFloat(value.(float64), 'f', 3, 64) // Convert float64 to string with 3 decimal places.
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			corrected = "'" + value.(time.Time).UTC().Format("2006-01-02 15:04:05") + "'" // Convert time.Time to UTC formatted string.
		}
	default:
		corrected = fmt.Sprintf("%v", value) // Convert other types to string using default formatting.
	}

	// Return here.
	return
}

// >>>>> >>>>> >>>>> Reading Entire Data Functions

// FetchResultsFromRows fetches results from a SQL rows object and converts them into a 2D string slice.
func FetchResultsFromRows(rows *sql.Rows) (results [][]string, err error) {
	// Retrieve the column names from the SQL rows object.
	columns, err := rows.Columns()
	if err != nil {
		return
	}

	// Prepare a slice to store the values of the current row.
	values := make([]interface{}, len(columns))
	rowValues := make([]string, len(columns))
	for i := range values {
		values[i] = &rowValues[i]
	}

	// Loop through each row in the rows object.
	for rows.Next() {
		// Scan the current row's values into the values slice.
		if err = rows.Scan(values...); err != nil {
			return
		}

		// Append the rowValues slice to the results slice, forming a 2D string slice.
		results = append(results, append([]string{}, rowValues...)) // <<<<< <<<<< <<<<< There's no type conversion here. (这里没型态转换)
		// There probably won't be any use of reflection here, because append is a built-in function.
		// Considering this, this one is better.
		// (应不会有 reflect，较好)
	}

	return
}

// FetchResultsFromRowsComparator fetches results from a SQL rows object and converts them into a 2D string slice.
func FetchResultsFromRowsComparator(rows *sql.Rows) (results [][]string, err error) {
	// Retrieve the column names from the SQL rows object.
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return
	}

	// Loop through each row in the rows object.
	for rows.Next() {
		// Create a slice to store the values of the current row.
		values := make([]interface{}, len(columns))
		for i := range values {
			// Each value is represented as a pointer to a sql.RawBytes type.
			values[i] = &sql.RawBytes{}
		}

		// Scan the current row's values into the values slice.
		if err = rows.Scan(values...); err != nil {
			return
		}

		// Convert the scanned values to strings and store them in a new rowValues slice.
		rowValues := make([]string, len(values))
		for i, val := range values {
			rowValues[i] = string(*val.(*sql.RawBytes)) // <<<<< <<<<< <<<<< Here's type conversion. (这里有型态转换)
		}

		// Append the rowValues slice to the results slice, forming a 2D string slice.
		results = append(results, rowValues)
	}

	return
}

const (
	Condition_The_Same       uint8 = iota // Condition_The_Same indicates that the state is the same; all initialized variables have not been modified.
	Condition_Diff_In_Length              // Condition_Diff_In_Length indicates that two slices have different lengths, and comparison cannot proceed due to length mismatch, resulting in direct interruption.
	Condition_Diff_In_Value               // Condition_Diff_In_Value indicates that two slices have different values.
)

// DiffPlace represents the differences between two slices.
type DiffPlace struct {
	RowIndex    int    // Index of the differing row
	ColumnIndex int    // Index of the differing column
	BeforeValue string // Value before modification
	AfterValue  string // Value after modification
}

// CompareResults compares two sets of string slices and identifies any differences.
func CompareResults(results1, results2 [][]string) (same bool, condition uint8, differences []DiffPlace) {
	// Compare the lengths of the slices
	if len(results1) != len(results2) {
		differences = append(differences, DiffPlace{
			RowIndex:    -1,
			ColumnIndex: -1,
		})
		same = false
		condition = Condition_Diff_In_Length
		return
	}

	// Loop through the rows of both sets of slices.
	for i := 0; i < len(results1); i++ {
		// Compare the lengths of the rows
		if len(results1[i]) != len(results2[i]) {
			// If the lengths of the rows are different, record the difference and set the conditions accordingly.
			differences = append(differences, DiffPlace{
				RowIndex:    i,
				ColumnIndex: -1,
			})
			same = false
			condition = Condition_Diff_In_Length
			return
		}

		// Loop through the columns of the current row.
		for j := 0; j < len(results1[i]); j++ {
			// Compare the values in the corresponding positions of both sets of slices.
			if results1[i][j] != results2[i][j] {
				// If the values are different, record the difference and set the conditions accordingly.
				difference := DiffPlace{
					RowIndex:    i,
					ColumnIndex: j,
					BeforeValue: results1[i][j],
					AfterValue:  results2[i][j],
				}
				condition = Condition_Diff_In_Value
				differences = append(differences, difference)
			}
			// The loop must be completed to collect the entire set of differences.
			// (回圈一定要跑完，因为要记录所有不同的位置)
		}
	}

	// Determine whether the two sets of slices are the same based on the number of differences found.
	same = len(differences) == 0
	return
}

// >>>>> >>>>> >>>>> Shared Functions Functions

// MakeSelectTableSQLStr generates SQL INSERT statements for a given table name and model.
func MakeSelectTableSQLStr(tableName string, model interface{}, convertCase uint8) (sqlStr string) {
	// Get the reflection type of the model.
	t := reflect.TypeOf(model)
	var columns []string

	// Loop through the fields of the model.
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		// Convert field names to the desired case format and append to the columns slice.
		columns = append(columns,
			ConvertStringFormats(field.Name, convertCase),
		)
	}

	// Join the column names with commas to form the columnsStr.
	columnsStr := strings.Join(columns, ", ")

	// Generate the final SELECT SQL statement.
	sqlStr = fmt.Sprintf("SELECT %s FROM %s;", columnsStr, tableName)
	return
}

// Defines text case conversion constants.
const (
	Case_Upper     uint8 = iota + 1 // Represents converting text to uppercase.
	Case_Lower                      // Represents converting text to lowercase.
	Case_Snake                      // Represents converting text to snake_case.
	Case_No_Change                  // Represents no change in text case. (Keep Camel)
)

// ConvertStringFormats converts the input string to different case formats.
// It takes the input string and a conversion case flag.
// Supported conversion cases: Case_Upper, Case_Lower, Case_Snake.
func ConvertStringFormats(input string, convertCase uint8) (formatted string) {
	switch convertCase {
	case Case_Upper:
		formatted = strings.ToUpper(input) // Convert the input string to uppercase.
	case Case_Lower:
		formatted = strings.ToLower(input) // Convert the input string to lowercase.
	case Case_Snake:
		re := regexp.MustCompile(`([a-z0-9])([A-Z])`) // Convert camelCase to snake_case.
		formatted = strings.ToLower(
			re.ReplaceAllString(input, "${1}_${2}"),
		)
	case Case_No_Change:
		formatted = input
	default:
		// Handle other cases
	}
	return formatted
}
