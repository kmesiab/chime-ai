package main

import (
	"time"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetTransactionsByDate(startDate, endDate time.Time) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetTransactionsByType(startDate, endDate time.Time, transactionType string) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Where("date BETWEEN? AND? AND type = ?", startDate, endDate, transactionType).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetTransactionsByDescription(startDate, endDate time.Time, description string) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Where("date BETWEEN? AND? AND description LIKE ?", startDate, endDate, "%"+description+"%").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetDistinctTransactionDescriptions(startDate, endDate time.Time) ([]string, error) {
	var descriptions []string
	if err := r.db.Model(&Transaction{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Distinct("description").
		Order("description ASC").
		Pluck("description", &descriptions).Error; err != nil {
		return nil, err
	}
	return descriptions, nil
}

func (r *TransactionRepository) GetDistinctTransactionDescriptionsAndTotal(startDate, endDate time.Time) (*[]DescriptionTotal, error) {
	var results []DescriptionTotal

	if err := r.db.Model(&Transaction{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select("description, SUM(net_amount) as total_spent").
		Group("description").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *TransactionRepository) GetDistinctTransactionDescriptionsAndCount(startDate, endDate time.Time) (*[]DescriptionCount, error) {
	var results []DescriptionCount

	if err := r.db.Model(&Transaction{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select("description, COUNT(*) as total_transactions").
		Group("description").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	return &results, nil
}
