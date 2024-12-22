package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {

	var (
		db  *gorm.DB
		err error
	)

	// Get the database connection
	if db, err = getDBConnection(); err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return
	}

	// Get a repository instance
	repository := NewTransactionRepository(db)

	// Bracket november transactions
	startDate, _ := time.Parse("2006-01-02", "2024-03-24")
	endDate, _ := time.Parse("2006-01-02", "2024-03-29")

	// Get transactions by date
	transactions, err := repository.GetTransactionsByDate(startDate, endDate)

	if err != nil {
		log.Printf("Error getting transactions by date: %v\n", err)
		return
	}

	totalSpent := float64(0)

	fmt.Printf("Transactions between %s and %s:\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	fmt.Println("-----------------------------------------")
	fmt.Printf("Description\tNet Amount\n")
	// Print the transactions
	for _, transaction := range transactions {
		totalSpent += transaction.NetAmount
		fmt.Printf("%s\t$%.2f\n", transaction.Description, transaction.NetAmount)
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Total spent: $%.2f\n", totalSpent*-1)
}
