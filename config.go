package sqlmock

import (
	"database/sql/driver"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> Specify the location of the mock folder first.

var mockConfigLocation string

// SetMockLocationByTriggerMain sets mockConfigLocation and gets the directory path of the caller's location.
func SetMockLocationByTriggerMain() {
	if mockConfigLocation == "" {
		_, file, _, _ := runtime.Caller(1)
		mockConfigLocation = filepath.Dir(file)
	}
	return
}

// SetMockLocationByManual sets mockConfigLocation by manual.
func SetMockLocationByManual(path string) {
	mockConfigLocation = path
	return
}

// GetMockLocation returns mockConfigLocation.
func GetMockLocation() (path string) {
	return mockConfigLocation
}

// >>>>> >>>>> >>>>> >>>>> Set the storage format for the mock.

// ConfigSet can be used to set the basic response
type ConfigSet struct {
	QueryString string         `json:"qureyString"`
	QueryArgs   []driver.Value `json:"queryArgs"`
	ReturnRows  []ConfigRows   `json:"returnRows"`
}

// ConfigRows can be used to set the corresponding database schema.
type ConfigRows struct {
	Columns []string         `json:"columns"`
	Rows    [][]driver.Value `json:"rows"`
}

// >>>>> >>>>> >>>>> >>>>> Start using the function to load mock configurations.

// LoadMockConfig is used to load the configuration values for Mock.
// It contains JSON data and requires UseNumber to preserve numbers as strings.
// If there are specific performance requirements, you may need to write your own parser.
// For now, we'll use the simplest method using json.NewDecoder to handle it.
func LoadMockConfig(sqlMock Sqlmock, jsonFile string) (error error) {
	// Join the directory and the JSON file name to get the full file path
	mockFile := filepath.Join(mockConfigLocation, jsonFile)

	// Read the JSON data from the configuration file.
	data, err := os.ReadFile(mockFile)
	if err != nil {
		return
	}

	// Use json.Decoder to parse and preserve the original numeric types. Unfortunately,
	// json.Unmarshal converts all integers and decimals to float64 in the interface, which cannot be used.
	// (json.Unmarshal 执行之后，数据 []driver.Value 都会变成 float64，会很困扰)
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.UseNumber()

	// Using json.NewDecoder instead of json.Unmarshal.
	var mockData []ConfigSet
	err = decoder.Decode(&mockData)
	if err != nil {
		return
	}

	// Iterate through the mock data and print the values.
	for _, mock := range mockData {
		for _, returnRows := range mock.ReturnRows {
			response := NewRows(returnRows.Columns)
			for _, row := range returnRows.Rows {
				response = response.AddRow(convertNumbers(row)...)
			}
			sqlMock.ExpectQuery(
				// Using QuoteMeta simplifies the configuration file and makes the setup more convenient.
				regexp.QuoteMeta(mock.QueryString),
			).WithArgs(convertNumbers(mock.QueryArgs)...).WillReturnRows(response)
		}
	}

	// If no errors occur, it returns.
	return
}

// >>>>> >>>>> >>>>> >>>>> Using a helper function to assist the LoadMockConfig function.

// convertNumbers is used to convert data in json.Number format to integers or decimals,
// as json.NewDecoder returns data in json.Number format.
func convertNumbers(data []driver.Value) []driver.Value {
	convertedData := make([]driver.Value, len(data))
	for i, v := range data {
		switch value := v.(type) {
		case json.Number:
			if isInteger(value) {
				intVal, _ := value.Int64()
				convertedData[i] = intVal
			} else {
				floatVal, _ := value.Float64()
				convertedData[i] = floatVal
			}
		default:
			convertedData[i] = v
		}
	}
	return convertedData
}

// isInteger is helper function to check if a JSON Number is an integer.
func isInteger(num json.Number) bool {
	// Use reflect to check if the number can be converted to int64 without errors.
	_, err := num.Int64()
	return err == nil
}
