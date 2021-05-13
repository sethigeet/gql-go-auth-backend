package util

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// ExpirationDuration is the duration after which the email link expires
const ExpirationDuration = 24 * time.Hour

// ConfirmEmailPrefix is the string that is prefixed before the actual confirm email token while storing in redis
// so that it will be easily identifiable and it will not collide with any other key
const ConfirmEmailPrefix = "confirm-email:"

var ctx = context.Background()

// getConfirmEmailLink create a link that can be put in an email and sent to the
// user through which the user can verify their email
// It creates an entry for the user id in redis and returns the encrypted key in a url back
func getConfirmEmailLink(rdb *redis.Client, userID string) (string, error) {
	token := uuid.New().String()

	var err error
	err = rdb.Set(ctx, ConfirmEmailPrefix+token, userID, ExpirationDuration).Err()
	if err != nil {
		return "", err
	}

	encryptedToken, err := Encrypt(token)
	if err != nil {
		return "", err
	}

	link := os.Getenv("FRONTEND_HOST") + "/confirm-email/" + encryptedToken

	return link, nil
}

// GetUserIDFromEmailToken first decrypts the provided token with the secret key and
// then looks up the user id of the key in redis and returns it
func GetUserIDFromEmailToken(rdb *redis.Client, token string) (string, error) {
	decryptedToken, err := Decrypt(token)
	if err != nil {
		return "", err
	}

	userID, err := rdb.Get(ctx, ConfirmEmailPrefix+decryptedToken).Result()
	switch {
	case err == redis.Nil:
		// Value does not exist
		return "", nil
	case err != nil:
		// Unable to get the value from redis
		return "", err
	}

	return userID, nil
}
