package sqlmock

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Test_Check_Config_Path tests mock configuration path setup.
func Test_Check_Config_Path(t *testing.T) {
	SetMockLocationByManual("./config")
	require.Equal(t, "./config", GetMockLocation())
}

// Test_Check_Config_Select tests SQL query with mock data for hotel selection.
func Test_Check_Config_Select(t *testing.T) {
	// Create a new SQL mock for testing.
	sqlDB, sqlMock, err := New()

	// Prepare SQL mock data
	SetMockLocationByManual("./config")          // Set the mock configuration location manually.
	err = LoadMockConfig(sqlMock, "select.json") // Load the mock configuration for testing.
	require.NoError(t, err)

	// Create a new SQLx DB instance.
	sqlxDB := sqlx.NewDb(sqlDB, "mysql")

	// Execute the SQL query to retrieve hotel data.
	rows, err := sqlxDB.Query("SELECT id, name, city, rating, price_per_night FROM hotels WHERE city = '?' AND rating >= ?;", "New York", 4)
	require.NoError(t, err)
	defer func() {
		_ = rows.Close()
	}()

	// Define test cases with expected hotel data.
	tests := []struct {
		ID            int
		Name          string
		City          string
		Rating        float64
		PricePerNight float64
	}{
		{1, "Grand Hotel", "New York", 4.5, 150.0},
		{2, "Luxury Inn", "New York", 4.2, 120.0},
	}

	// Loop through each row in the query result.
	var index int
	for rows.Next() {
		var id int
		var name, city string
		var rating float64
		var pricePerNight float64

		// Scan the row's values into variables.
		err = rows.Scan(&id, &name, &city, &rating, &pricePerNight)
		require.NoError(t, err)

		// Perform assertions to compare actual values with expected values.
		assert.Equal(t, tests[index].ID, id)
		assert.Equal(t, tests[index].Name, name)
		assert.Equal(t, tests[index].City, city)
		assert.Equal(t, tests[index].Rating, rating)
		assert.Equal(t, tests[index].PricePerNight, pricePerNight)

		index++
	}
}
