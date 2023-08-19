package sqlmock

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

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

// MakeCreateTableSQLStr generates a SQL string for creating a table with the given table name and model.
// Database column names prioritize tags first, followed by other settings.
// (栏位优先考虑标签)
func MakeCreateTableSQLStr(tableName string, model interface{}, convertCase uint8) (sqlStr string) {
	// tableFields := make([]string, 0) // 其他常用的写法先保留
	t := reflect.TypeOf(model)

	// Enhancing Performance with String Building.
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(tableName)
	sb.WriteString(" (")

	// Iterate through the fields of the model.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		columnName := ConvertStringFormats(field.Name, convertCase)

		tag := field.Tag.Get("db") // Using the 'db' Tag as the Preferred Setting Value (field 为获取结构原素资讯)

		// Check if the 'db' tag is empty, indicating that this field should be used for the table.
		if tag != "" {
			columnName = tag
		}

		// tableFields = append(tableFields, fmt.Sprintf("%s %s", field.Name, fieldType)) // 其他常用的写法先保留
		sb.WriteString(columnName)
		sb.WriteByte(' ')
		fieldType := mapFieldType(field.Type) // Fetching Information about Structure Elements
		sb.WriteString(fieldType)

		// Add a comma if it's not the last field.
		if i < t.NumField()-1 {
			sb.WriteString(", ")
		}

	}

	// sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(tableFields, ", ")) // 其他常用的写法先保留

	// Close the CREATE TABLE statement.
	sb.WriteString(");")
	sqlStr = sb.String()

	return
}

// mapFieldType takes a reflect.Type as input and returns the corresponding SQL data type.
func mapFieldType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int:
		return "INT" // Map int type to SQL INT data type.
	case reflect.String:
		return "VARCHAR(255)" // Map string type to SQL VARCHAR(255) data type.
	case reflect.Float64:
		return "DECIMAL(10, 2)" // Map float64 type to SQL DECIMAL(10, 2) data type.
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP" // Map time.Time struct to SQL TIMESTAMP data type.
		}
	}

	// Return an empty string if the type does not match any mapping.
	return ""
}
