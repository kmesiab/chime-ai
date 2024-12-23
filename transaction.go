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

type DescriptionTotal struct {
	Description string
	TotalSpent  float64 `json:"total_spent"`
}

type DescriptionCount struct {
	Description       string
	TotalTransactions int `json:"total_transactions"`
}
