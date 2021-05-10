// Package database provides the logic for interacting with the database
package database

import "gorm.io/gorm"

// Connect connects to the database and returns the db object
func Connect(migrate bool) (*gorm.DB, error) {
	db, err := connectPostgres(migrate)
	if err != nil {
		return nil, err
	}

	return db, nil
}
