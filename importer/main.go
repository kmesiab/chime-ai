package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	// Accept directory path as a command-line argument
	dir := flag.String("dir", "./importer/files", "Directory containing PDFs and text files")
	flag.Parse()

	db, err := initDB("transactions.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Convert PDFs to text files if `pdftotext` is installed
	if isCommandAvailable("pdftotext") {
		log.Println("pdftotext found, converting PDFs to text...")
		convertPDFsToText(*dir)
	} else {
		log.Println("pdftotext not found, skipping PDF conversion. Please install it to process PDFs.")
	}

	files, err := filepath.Glob(filepath.Join(*dir, "*.txt"))
	if err != nil {
		log.Fatalf("Failed to list text files: %v", err)
	}

	if len(files) == 0 {
		log.Println("No text files found for processing.")
		return
	}

	log.Printf("Found %d text files for processing.", len(files))

	processFilesConcurrently(files, db)
	cleanupTxtFiles(files)
	log.Println("All files processed and cleaned up successfully!")
}

// initDB initializes the database and creates the table if it doesn't exist
func initDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&Transaction{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}

// isCommandAvailable checks if a command is available in the system
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// convertPDFsToText converts all PDFs in the given directory to text files
func convertPDFsToText(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.pdf"))
	if err != nil {
		log.Printf("Failed to list PDF files: %v", err)
		return
	}

	for _, pdfFile := range files {
		baseName := strings.TrimSuffix(filepath.Base(pdfFile), ".pdf")
		txtFile := filepath.Join(dir, baseName+".txt")

		cmd := exec.Command("pdftotext", "-layout", pdfFile, txtFile)
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to convert %s to text: %v", pdfFile, err)
		} else {
			log.Printf("Converted %s to %s", pdfFile, txtFile)
		}
	}
}

// processFilesConcurrently processes multiple files in parallel
func processFilesConcurrently(filenames []string, db *gorm.DB) {
	var wg sync.WaitGroup

	for _, filename := range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			processFile(f, db)
		}(filename)
	}

	wg.Wait()
}

// processFile parses a single transaction file and inserts data into the database
func processFile(filename string, db *gorm.DB) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %v", filename, err)
		return
	}
	defer file.Close()

	log.Printf("Processing file: %s", filename)
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
				log.Printf("Error parsing date in file %s: %v", filename, err)
				continue
			}

			settleDate, err := time.Parse("1/02/2006", match[6])
			if err != nil {
				log.Printf("Error parsing settlement date in file %s: %v", filename, err)
				continue
			}

			transaction := Transaction{
				Date:        date,
				Description: strings.TrimSpace(match[2]),
				Type:        match[3],
				Amount:      amount,
				NetAmount:   netAmount,
				SettleDate:  settleDate,
			}

			var existing Transaction
			if err := db.Where("date = ? AND description = ? AND amount = ? AND net_amount = ? AND settle_date = ?",
				transaction.Date,
				transaction.Description,
				transaction.Amount,
				transaction.NetAmount,
				transaction.SettleDate).
				First(&existing).
				Error; err == nil {
				log.Printf("Duplicate transaction found, skipping: %+v", transaction)
				continue
			}

			transactions = append(transactions, transaction)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return
	}

	if len(transactions) > 0 {
		if err := db.Create(&transactions).Error; err != nil {
			log.Printf("Error inserting transactions from file %s: %v", filename, err)
		} else {
			log.Printf("Inserted %d transactions from %s", len(transactions), filename)
		}
	} else {
		log.Printf("No transactions found in %s", filename)
	}
}

// cleanupTxtFiles removes all .txt files in the provided list
func cleanupTxtFiles(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to delete %s: %v", file, err)
		} else {
			log.Printf("Deleted %s", file)
		}
	}
}
