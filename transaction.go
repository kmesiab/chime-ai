package main

import "time"

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
