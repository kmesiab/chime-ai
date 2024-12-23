package database

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestExecuteRawQuery_EmptyQueryString(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Execute the method with an empty query string
	result, err := repo.ExecuteRawQuery("")
	if err == nil {
		t.Errorf("expected an error due to empty query string, got nil")
	}

	// Assert that the result is empty
	if len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}

func TestExecuteRawQuery_NoArguments(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	// Seed the database with some transactions
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Grocery Shopping", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), Description: "Salary", Type: "Deposit", Amount: 1500.00, NetAmount: 1500.00, SettleDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Execute a raw query with no arguments
	query := "SELECT * FROM transactions"
	results, err := repo.ExecuteRawQuery(query)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct number of results are returned
	if len(results) != len(transactions) {
		t.Errorf("expected %d results, got %d", len(transactions), len(results))
	}
}

func TestExecuteRawQuery_LargeNumberOfArguments(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	// Seed the database with transactions
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Transaction 1", Type: "Purchase", Amount: 10.00, NetAmount: 10.00, SettleDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Description: "Transaction 2", Type: "Purchase", Amount: 20.00, NetAmount: 20.00, SettleDate: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
		// Add more transactions as needed
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Construct a query with a large number of arguments
	query := "SELECT * FROM transactions WHERE description IN (?, ?, ?)"
	args := []interface{}{"Transaction 1", "Transaction 2", "Transaction 3"}

	// Execute the raw query
	result, err := repo.ExecuteRawQuery(query, args...)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct number of results are returned
	expectedCount := 2
	if len(result) != expectedCount {
		t.Errorf("expected %d results, got %d", expectedCount, len(result))
	}
}

func TestExecuteRawQuery_NoResults(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Define a query that will return no results
	query := "SELECT * FROM transactions WHERE description = ?"
	args := []interface{}{"Nonexistent Description"}

	result, err := repo.ExecuteRawQuery(query, args...)
	if err == nil {
		t.Errorf("expected error but got none")
	}

	// Assert that the result is an empty slice
	if len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}

func TestExecuteRawQuery_SpecialSQLKeywords(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with some transactions
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "SELECT * FROM", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), Description: "DROP TABLE", Type: "Deposit", Amount: 1500.00, NetAmount: 1500.00, SettleDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Execute a raw query with special SQL keywords
	query := "SELECT description FROM transactions WHERE description IN ('SELECT * FROM', 'DROP TABLE')"
	result, err := repo.ExecuteRawQuery(query)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct number of results is returned
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}

	// Assert that the correct descriptions are returned
	expectedDescriptions := map[string]bool{
		"SELECT * FROM": true,
		"DROP TABLE":    true,
	}

	for _, row := range result {
		desc, ok := row["description"].(string)
		if !ok || !expectedDescriptions[desc] {
			t.Errorf("unexpected description found: %v", desc)
		}
	}
}

func TestExecuteRawQuery_InvalidSyntax(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Define a query with invalid SQL syntax
	query := "SELEC * FROM transactions"

	// Execute the raw query
	_, err = repo.ExecuteRawQuery(query)
	if err == nil {
		t.Errorf("expected an error due to invalid SQL syntax, got nil")
	}
}

func TestExecuteRawQuery_LargeResultSet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with a large number of transactions
	var transactions []Transaction
	for i := 0; i < 1000; i++ {
		transactions = append(transactions, Transaction{
			Date:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i),
			Description: "Transaction " + strconv.Itoa(i),
			Type:        "Purchase",
			Amount:      float64(i),
			NetAmount:   float64(i),
			SettleDate:  time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i),
		})
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Execute a raw query to retrieve all transactions
	query := "SELECT * FROM transactions"
	result, err := repo.ExecuteRawQuery(query)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that all transactions are returned
	if len(result) != 1000 {
		t.Errorf("expected 1000 transactions, got %d", len(result))
	}
}

func TestExecuteRawQuery_HandleNullValues(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Create a table for testing
	type TestTable struct {
		ID    int
		Name  string
		Value *string
	}
	if err := db.AutoMigrate(&TestTable{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	// Seed the table with data including NULL values
	value := "Some Value"
	records := []TestTable{
		{ID: 1, Name: "Record1", Value: &value},
		{ID: 2, Name: "Record2", Value: nil},
	}
	if err := db.Create(&records).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Execute a raw query that will return NULL values
	query := "SELECT id, name, value FROM test_tables"
	result, err := repo.ExecuteRawQuery(query)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the result is as expected
	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	expectedResults := []map[string]interface{}{
		{"id": int64(1), "name": "Record1", "value": "Some Value"},
		{"id": int64(2), "name": "Record2", "value": nil},
	}

	for i, expected := range expectedResults {
		for key, expectedValue := range expected {
			if !reflect.DeepEqual(result[i][key], expectedValue) {
				t.Errorf("expected %v for key '%s', got %v", expectedValue, key, result[i][key])
			}
		}
	}
}

func TestExecuteRawQuery_MixedDataTypes(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Create a temporary table for testing
	if err := db.Exec(`
        CREATE TABLE test_data (
            id INTEGER PRIMARY KEY,
            name TEXT,
            amount REAL,
            created_at TEXT
        );
    `).Error; err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// Insert mixed data types into the table
	if err := db.Exec(`
        INSERT INTO test_data (name, amount, created_at) VALUES
        ('Sample Item', 123.45, '2023-10-01T12:00:00Z'),
        ('Another Item', 678.90, '2023-10-02T15:30:00Z');
    `).Error; err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Execute a raw query to retrieve mixed data types
	query := "SELECT id, name, amount, created_at FROM test_data"
	result, err := repo.ExecuteRawQuery(query)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct number of rows are returned
	if len(result) != 2 {
		t.Errorf("expected 2 rows, got %d", len(result))
	}

	// Assert that the data types are correctly returned
	expected := []map[string]interface{}{
		{"id": int64(1), "name": "Sample Item", "amount": 123.45, "created_at": "2023-10-01T12:00:00Z"},
		{"id": int64(2), "name": "Another Item", "amount": 678.90, "created_at": "2023-10-02T15:30:00Z"},
	}

	for i, row := range result {
		for key, expectedValue := range expected[i] {
			if !reflect.DeepEqual(row[key], expectedValue) {
				t.Errorf("expected %v for key '%s', got %v", expectedValue, key, row[key])
			}
		}
	}
}
