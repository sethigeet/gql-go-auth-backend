// Package util provides some useful utility functions
package util

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv Loads the environment variables defined in the .env file according to
// the environment running
func LoadEnv() error {

	env := os.Getenv("GO_ENV")

	var err error
	switch env {
	case "production":
		err = godotenv.Load(".env.prod")
	case "testing":
		err = godotenv.Load(".env.test")
	case "development":
		err = godotenv.Load(".env.local")
	default:
		err = godotenv.Load(".env.local")
	}

	if err != nil {
		return err
	}

	return nil
}
