package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("transactions.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
