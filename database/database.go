// Package database provides the logic for interacting with the database
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sethigeet/gql-go-auth-backend/graph/model"
)

// Connect connects to the database and returns the db object
func Connect(migrate bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DBNAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database
	var logLevel int
	env := os.Getenv("GO_ENV")
	if env == "testing" || env == "production" {
		logLevel = int(logger.Error)
	} else {
		logLevel = int(logger.Info)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.LogLevel(logLevel),
			Colorful:      true,
		}),
	})
	if err != nil {
		return nil, err
	}

	if migrate {
		if err := automigrate(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func automigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})

	return err
}
