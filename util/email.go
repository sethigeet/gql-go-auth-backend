package util

import (
	"log"

	"github.com/go-redis/redis/v8"
)

// SendConfirmEmailEmail creates a unique token for a given userID, stores that
// token in redis and sends an email to the email address of that user to ask
// them to visit that link and confirm their email
func SendConfirmEmailEmail(rdb *redis.Client, userID, email string) error {
	link, err := getConfirmEmailLink(rdb, userID)
	if err != nil {
		return err
	}

	log.Printf("\nNew confirm email link generated for %s: %s\n\n", email, link)

	// TODO: Send an email to the email address

	return nil
}
