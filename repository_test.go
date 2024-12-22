package main

import (
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
