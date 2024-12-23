package database

import (
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) ExecuteRawQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Raw(query, args...).Scan(&result).Error
	return result, err
}
