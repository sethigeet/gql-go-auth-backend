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

// ForgotPasswordPrefix is the string that is prefixed before the actual forgot password token while storing in redis
// so that it will be easily identifiable and it will not collide with any other key
const ForgotPasswordPrefix = "forgot-password:"

var ctx = context.Background()

// getLink create a link that can be put in an email and sent to the
// user through which the user can verify their email
// It creates an entry for the user id in redis and returns the encrypted key in a url back
func getLink(rdb *redis.Client, userID string, prefix string, endpoint string) (string, error) {
	token := uuid.New().String()

	var err error
	err = rdb.Set(ctx, prefix+token, userID, ExpirationDuration).Err()
	if err != nil {
		return "", err
	}

	encryptedToken, err := Encrypt(token)
	if err != nil {
		return "", err
	}

	link := os.Getenv("FRONTEND_HOST") + endpoint + "/" + encryptedToken

	return link, nil
}

// GetUserIDFromToken first decrypts the provided token with the secret key and
// then looks up the user id of the key in redis and returns it
func GetUserIDFromToken(rdb *redis.Client, token string, prefix string) (string, func() error, error) {
	decryptedToken, err := Decrypt(token)
	if err != nil {
		return "", nil, err
	}

	userID, err := rdb.Get(ctx, prefix+decryptedToken).Result()
	switch {
	case err == redis.Nil:
		// Value does not exist
		return "", nil, nil
	case err != nil:
		// Unable to get the value from redis
		return "", nil, err
	}

	return userID, func() error {
		err := rdb.Del(ctx, prefix+decryptedToken).Err()
		return err
	}, nil
}
