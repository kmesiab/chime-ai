package main

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTransactionsByDescription_NoMatchingDescription(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Define the test start and end dates
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Call the method with a description that doesn't exist
	transactions, err := repo.GetTransactionsByDescription(startDate, endDate, "Nonexistent Description")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the returned slice is empty
	if len(transactions) != 0 {
		t.Errorf("expected no transactions, got %d", len(transactions))
	}
}

func TestGetTransactionsByDescription_DBConnectionLost(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Simulate a database connection loss by closing the database
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get generic database object: %v", err)
	}
	sqlDB.Close()

	// Define the test start and end dates
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Call the method expecting an error due to closed connection
	_, err = repo.GetTransactionsByDescription(startDate, endDate, "Any Description")
	if err == nil {
		t.Errorf("expected an error due to closed database connection, got nil")
	}
}

func TestGetTransactionsByDescription_EmptyDescription(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with some transactions
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Grocery Shopping", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), Description: "Salary", Type: "Deposit", Amount: 1500.00, NetAmount: 1500.00, SettleDate: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Define the test start and end dates
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Call the method with an empty description
	result, err := repo.GetTransactionsByDescription(startDate, endDate, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that all transactions are returned
	if len(result) != len(transactions) {
		t.Errorf("expected %d transactions, got %d", len(transactions), len(result))
	}
}

func TestGetTransactionsByDescription_SpecialCharacters(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with transactions containing special characters
	transactions := []Transaction{
		{Date: time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC), Description: "Grocery @ Store #123", Type: "Purchase", Amount: 75.00, NetAmount: 75.00, SettleDate: time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC), Description: "Dinner at Joe's!", Type: "Purchase", Amount: 120.00, NetAmount: 120.00, SettleDate: time.Date(2023, 3, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Define the test start and end dates
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Call the method with a description containing special characters
	result, err := repo.GetTransactionsByDescription(startDate, endDate, "@ Store #123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct transaction is returned
	if len(result) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(result))
	} else if result[0].Description != "Grocery @ Store #123" {
		t.Errorf("expected transaction description to be 'Grocery @ Store #123', got '%s'", result[0].Description)
	}
}

func TestGetTransactionsByDescription_SQLWildcardCharacters(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with transactions containing SQL wildcard characters
	transactions := []Transaction{
		{Date: time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC), Description: "50% Discount Sale", Type: "Purchase", Amount: 100.00, NetAmount: 100.00, SettleDate: time.Date(2023, 4, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC), Description: "Special _Offer_ Today", Type: "Purchase", Amount: 200.00, NetAmount: 200.00, SettleDate: time.Date(2023, 4, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Define the test start and end dates
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Call the method with a description containing SQL wildcard characters
	result, err := repo.GetTransactionsByDescription(startDate, endDate, "% Discount %")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct transaction is returned
	if len(result) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(result))
	} else if result[0].Description != "50% Discount Sale" {
		t.Errorf("expected transaction description to be '50 Discount Sale', got '%s'", result[0].Description)
	}
}

func TestGetTransactionsByDescription_SameStartAndEndDate(t *testing.T) {
	// Initialize the in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with a transaction on a specific date
	transactionDate := time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC)
	transactions := []Transaction{
		{Date: transactionDate, Description: "Coffee Shop", Type: "Purchase", Amount: 5.00, NetAmount: 5.00, SettleDate: transactionDate},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	// Define the test start and end date as the same
	startDate := transactionDate
	endDate := transactionDate

	// Call the method with the same start and end date
	result, err := repo.GetTransactionsByDescription(startDate, endDate, "Coffee Shop")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert that the correct transaction is returned
	if len(result) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(result))
	} else if result[0].Description != "Coffee Shop" {
		t.Errorf("expected transaction description to be 'Coffee Shop', got '%s'", result[0].Description)
	}
}

func TestGetTransactionsByDescription_CaseInsensitive(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	transactions := []Transaction{
		{Date: time.Date(2023, 7, 10, 0, 0, 0, 0, time.UTC), Description: "Online Shopping", Type: "Purchase", Amount: 100.00, NetAmount: 100.00, SettleDate: time.Date(2023, 7, 11, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC), Description: "ONLINE SHOPPING", Type: "Purchase", Amount: 150.00, NetAmount: 150.00, SettleDate: time.Date(2023, 7, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	result, err := repo.GetTransactionsByDescription(startDate, endDate, "online shopping")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(result))
	}
}

func TestGetTransactionsByDescription_StartDateAfterEndDate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	startDate, _ := time.Parse("2006-01-02", "2023-12-31")
	endDate, _ := time.Parse("2006-01-02", "2023-01-01")

	transactions, err := repo.GetTransactionsByDescription(startDate, endDate, "Any Description")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(transactions) != 0 {
		t.Errorf("expected no transactions, got %d", len(transactions))
	}
}

func TestGetDistinctTransactionDescriptions_DBError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Simulate a database error by closing the database connection
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get generic database object: %v", err)
	}
	sqlDB.Close()

	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-12-31")

	_, err = repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err == nil {
		t.Errorf("expected an error due to closed database connection, got nil")
	}
}

func TestGetDistinctTransactionDescriptions_DuplicateDescriptions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with transactions having duplicate descriptions
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Grocery Shopping", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Description: "Grocery Shopping", Type: "Purchase", Amount: 75.00, NetAmount: 75.00, SettleDate: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Gas Station", Type: "Purchase", Amount: 40.00, NetAmount: 40.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	descriptions, err := repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedDescriptions := []string{"Grocery Shopping", "Gas Station"}
	if len(descriptions) != len(expectedDescriptions) {
		t.Errorf("expected %d distinct descriptions, got %d", len(expectedDescriptions), len(descriptions))
	}

	for _, expected := range expectedDescriptions {
		found := false
		for _, actual := range descriptions {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected description '%s' not found in result", expected)
		}
	}
}

func TestGetDistinctTransactionDescriptions_NoTransactions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	descriptions, err := repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(descriptions) != 0 {
		t.Errorf("expected empty slice, got slice with %d elements", len(descriptions))
	}
}

func TestGetDistinctTransactionDescriptions_ExactStartAndEndDates(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC)

	// Seed the database with transactions exactly on start and end dates
	transactions := []Transaction{
		{Date: startDate, Description: "Start Date Transaction", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: startDate},
		{Date: endDate, Description: "End Date Transaction", Type: "Purchase", Amount: 100.00, NetAmount: 100.00, SettleDate: endDate},
		{Date: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC), Description: "Middle Transaction", Type: "Purchase", Amount: 75.00, NetAmount: 75.00, SettleDate: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	descriptions, err := repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedDescriptions := []string{"Start Date Transaction", "Middle Transaction", "End Date Transaction"}
	if len(descriptions) != len(expectedDescriptions) {
		t.Errorf("expected %d distinct descriptions, got %d", len(expectedDescriptions), len(descriptions))
	}

	for _, expected := range expectedDescriptions {
		if !contains(descriptions, expected) {
			t.Errorf("expected description '%s' not found in result", expected)
		}
	}
}

func TestGetDistinctTransactionDescriptions_AlphabeticalOrder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with transactions having descriptions in non-alphabetical order
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Zebra Store", Type: "Purchase", Amount: 50.00, NetAmount: 50.00, SettleDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Description: "Apple Market", Type: "Purchase", Amount: 75.00, NetAmount: 75.00, SettleDate: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Mango Shop", Type: "Purchase", Amount: 40.00, NetAmount: 40.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	descriptions, err := repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedOrder := []string{"Apple Market", "Mango Shop", "Zebra Store"}
	if !reflect.DeepEqual(descriptions, expectedOrder) {
		t.Errorf("descriptions not in alphabetical order. Expected %v, got %v", expectedOrder, descriptions)
	}
}

func TestGetDistinctTransactionDescriptions_UnicodeCharacters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	// Seed the database with transactions containing Unicode characters
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Café ☕", Type: "Purchase", Amount: 5.00, NetAmount: 5.00, SettleDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Description: "Sushi 🍣", Type: "Purchase", Amount: 20.00, NetAmount: 20.00, SettleDate: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC)},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Bücher 📚", Type: "Purchase", Amount: 30.00, NetAmount: 30.00, SettleDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	descriptions, err := repo.GetDistinctTransactionDescriptions(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Expected descriptions should be in alphabetical order
	expectedDescriptions := []string{"Bücher 📚", "Café ☕", "Sushi 🍣"}
	if !reflect.DeepEqual(descriptions, expectedDescriptions) {
		t.Errorf("expected descriptions %v, got %v", expectedDescriptions, descriptions)
	}
}

func TestGetDistinctTransactionDescriptionsAndTotal_ValidData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Grocery", NetAmount: 100.00},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Description: "Grocery", NetAmount: 200.00},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Gas", NetAmount: 50.00},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	results, err := repo.GetDistinctTransactionDescriptionsAndTotal(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedResults := []DescriptionTotal{
		{Description: "Grocery", TotalSpent: 300.00},
		{Description: "Gas", TotalSpent: 50.00},
	}

	if len(*results) != len(expectedResults) {
		t.Errorf("expected %d results, got %d", len(expectedResults), len(*results))
		return
	}

	for _, expected := range expectedResults {
		found := false
		for _, actual := range *results {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected result %v not found in actual results", expected)
		}
	}
}

func TestGetDistinctTransactionDescriptionsAndTotal_EmptyDatabase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	results, err := repo.GetDistinctTransactionDescriptionsAndTotal(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Check for nil to avoid dereferencing a nil pointer
	if results == nil || len(*results) != 0 {
		t.Errorf("expected empty result, got %d", len(*results))
	}
}

func TestGetDistinctTransactionDescriptionsAndTotal_SameStartAndEndDate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	transactionDate := time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC)
	transactions := []Transaction{
		{Date: transactionDate, Description: "Coffee Shop", NetAmount: 5.00},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := transactionDate
	endDate := transactionDate

	results, err := repo.GetDistinctTransactionDescriptionsAndTotal(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedResults := []DescriptionTotal{
		{Description: "Coffee Shop", TotalSpent: 5.00},
	}
	if !reflect.DeepEqual(*results, expectedResults) {
		t.Errorf("expected %v, got %v", expectedResults, *results)
	}
}

func TestGetDistinctTransactionDescriptionsAndTotal_DBConnectionLost(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get generic database object: %v", err)
	}
	sqlDB.Close()

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err = repo.GetDistinctTransactionDescriptionsAndTotal(startDate, endDate)
	if err == nil {
		t.Errorf("expected an error due to closed database connection, got nil")
	}
}

func TestGetDistinctTransactionDescriptionsAndCount_ValidData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Description: "Grocery", NetAmount: 100.00},
		{Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Description: "Grocery", NetAmount: 200.00},
		{Date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Description: "Gas", NetAmount: 50.00},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	results, err := repo.GetDistinctTransactionDescriptionsAndCount(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedResults := []DescriptionCount{
		{Description: "Grocery", TotalTransactions: 2},
		{Description: "Gas", TotalTransactions: 1},
	}

	if len(*results) != len(expectedResults) {
		t.Errorf("expected %d results, got %d", len(expectedResults), len(*results))
		return
	}

	for _, expected := range expectedResults {
		found := false
		for _, actual := range *results {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected result %v not found in actual results", expected)
		}
	}
}

func TestGetDistinctTransactionDescriptionsAndCount_EmptyDatabase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	results, err := repo.GetDistinctTransactionDescriptionsAndCount(startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if results == nil || len(*results) != 0 {
		t.Errorf("expected empty result, got %d", len(*results))
	}
}

func TestGetDistinctTransactionDescriptionsAndCount_DBConnectionLost(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewTransactionRepository(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get generic database object: %v", err)
	}
	sqlDB.Close()

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	_, err = repo.GetDistinctTransactionDescriptionsAndCount(startDate, endDate)
	if err == nil {
		t.Errorf("expected an error due to closed database connection, got nil")
	}
}

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
	if err != nil {
		t.Errorf("unexpected error: %v", err)
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
