// Package database provides the logic for interacting with the database
package database

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Connect connects to the database and returns the db object
func Connect(migrate bool) (*gorm.DB, *redis.Client, error) {
	var err error
	db, err := connectPostgres(migrate)
	if err != nil {
		return nil, nil, err
	}

	rdb, err := connectRedis()
	if err != nil {
		return nil, nil, err
	}

	return db, rdb, nil
}
