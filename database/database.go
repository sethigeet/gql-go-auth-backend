// Package database provides the logic for interacting with the database
package database

import (
	"fmt"
	"os"

	"github.com/sethigeet/gql-go-auth-backend/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database and returns the db object
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DBNAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := automigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func automigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})

	return err
}
