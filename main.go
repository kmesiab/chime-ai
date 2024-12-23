package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {

	var (
		sqlDB *sql.DB
		db    *gorm.DB
		err   error
	)

	// Get the database connection
	if db, err = getDBConnection(); err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return
	}

	// Defer the closing of the database connection
	if sqlDB, err = db.DB(); err == nil {
		defer sqlDB.Close()
	}

	// Get a repository instance
	repository := NewTransactionRepository(db)

	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2024-12-31")

	printTransactionsByDateRange(repository, startDate, endDate)
	printDistinctTransactionDescriptions(repository, startDate, endDate)
}

func printDistinctTransactionDescriptions(repository *TransactionRepository, startDate, endDate time.Time) {
	descriptions, err := repository.GetDistinctTransactionDescriptionsAndTotal(startDate, endDate)
	if err != nil {
		log.Printf("Error getting distinct transaction descriptions: %v\n", err)
		return
	}

	fmt.Printf("Distinct transaction descriptions between %s and %s:\n",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	for _, d := range *descriptions {
		fmt.Printf("%s\t%2f\n", d.Description, d.TotalSpent)
	}
}

func printTransactionsByDateRange(repository *TransactionRepository, startDate, endDate time.Time) {
	transactions, err := repository.GetTransactionsByDate(startDate, endDate)

	if err != nil {
		log.Printf("Error getting transactions by date: %v\n", err)
		return
	}

	totalSpent := float64(0)

	fmt.Printf("Transactions between %s and %s:\n",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	fmt.Println("-----------------------------------------")
	fmt.Printf("Description\tNet Amount\n")

	for _, transaction := range transactions {
		totalSpent += transaction.NetAmount
		fmt.Printf("%s\t$%.2f\n", transaction.Description, transaction.NetAmount)
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Total spent: $%.2f\n", totalSpent*-1)
}
