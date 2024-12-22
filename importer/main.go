package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Transaction struct represents the database model and parsed transactions
type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	Date        time.Time `gorm:"index"`
	Description string
	Type        string
	Amount      float64
	NetAmount   float64
	SettleDate  time.Time
}

func main() {
	// Initialize the database
	db, err := initDB("transactions.db")
	if err != nil {
		panic(err)
	}

	// Process transaction files
	for i := 1; i <= 19; i++ {
		filename := fmt.Sprintf("./importer/files/checking_%d.txt", i)
		filename, _ = filepath.Abs(filename)
		fmt.Printf("Processing file: %s\n", filename)
		processFile(filename, db)
	}

	fmt.Println("All files processed successfully!")
}

// initDB initializes the database and creates the table if it doesn't exist
func initDB(dbPath string) (*gorm.DB, error) {
	// Connect to SQLite
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the Transaction model
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}

// processFile parses a single transaction file and inserts data into the database
func processFile(filename string, db *gorm.DB) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filename, err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file %s: %v\n", filename, err)
		}
	}(file)

	var transactions []Transaction
	re := regexp.MustCompile(`^(\d{1,2}/\d{1,2}/\d{4})\s+(.*?)\s+(Transfer|Purchase|Direct Debit|ATM Withdrawal|Fee|Deposit|Round Up)\s+(-?\$\d+\.\d{2})\s+(-?\$\d+\.\d{2})\s+(\d{1,2}/\d{1,2}/\d{4})$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			amount, _ := strconv.ParseFloat(strings.ReplaceAll(match[4], "$", ""), 64)
			netAmount, _ := strconv.ParseFloat(strings.ReplaceAll(match[5], "$", ""), 64)

			date, err := time.Parse("1/02/2006", match[1])

			if err != nil {
				fmt.Printf("Error parsing date: %v\n", err)
				continue
			}

			settleDate, err := time.Parse("1/02/2006", match[6])

			if err != nil {
				fmt.Printf("Error parsing settlement date: %v\n", err)
				continue
			}

			transaction := Transaction{
				Date:        date,
				Description: match[2],
				Type:        match[3],
				Amount:      amount,
				NetAmount:   netAmount,
				SettleDate:  settleDate,
			}

			// Check if the transaction already exists
			var existing Transaction
			if err := db.Where("date = ? AND description = ? AND amount = ? AND net_amount = ? AND settle_date = ?",
				transaction.Date,
				transaction.Description,
				transaction.Amount,
				transaction.NetAmount,
				transaction.SettleDate,
			).
				First(&existing).
				Error; err == nil {

				fmt.Printf("Duplicate transaction found, skipping: %v\n", transaction)
				continue
			}

			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) > 0 {
		if err := db.Create(&transactions).Error; err != nil {
			fmt.Printf("Error inserting transactions: %v\n", err)
		} else {
			fmt.Printf("Inserted %d transactions from %s\n", len(transactions), filename)
		}
	} else {
		fmt.Printf("No transactions found in %s\n", filename)
	}
}
