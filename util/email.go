package util

import (
	"log"

	"github.com/go-redis/redis/v8"
)

// SendEmail creates a unique token for a given userID, stores that
// token in redis and sends an email to the email address of that user to ask
// them to visit that link and do what is required
func SendEmail(rdb *redis.Client, userID, email, prefix, endpoint string) error {
	link, err := getLink(rdb, userID, prefix, endpoint)
	if err != nil {
		return err
	}

	log.Printf("\nNew confirm email link generated for %s: %s\n\n", email, link)

	// TODO: Send an email to the email address

	return nil
}
