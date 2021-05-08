// Package util provides some useful utility functions
package util

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv Loads the environment variables defined in the .env file according to
// the environment running
func LoadEnv(verify bool) error {

	env := os.Getenv("GO_ENV")

	var envFileExt string
	switch env {
	case "production":
		envFileExt = "prod"
	case "testing":
		envFileExt = "test"
	default:
		envFileExt = "local"
	}

	if verify {
		err := verifyEnv(".env." + envFileExt)
		if err != nil {
			return err
		}
	}

	err := godotenv.Load(".env." + envFileExt)
	if err != nil {
		return err
	}

	return nil
}

func verifyEnv(filename string) error {
	var exampleEnv, actualEnv map[string]string
	var err error

	actualEnv, err = godotenv.Read(filename)
	if err != nil {
		return fmt.Errorf("an error occured while reading the '%s' file:\n%s", filename, err)
	}

	exampleEnv, err = godotenv.Read(".env.example")
	if err != nil {
		return fmt.Errorf("an error occured while reading the '.env.example' file:\n%s", err)
	}

	for envVar := range exampleEnv {
		if actualEnv[envVar] == "" {
			return fmt.Errorf("the env var %s is not present but is defined in the .env.example file", envVar)
		}
	}

	return nil
}
