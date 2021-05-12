// Package session provides a session manager to easily manage saving, retreiving
// and deleting sessions
package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type SessionManager struct {
	RDB     *redis.Client
	Writer  http.ResponseWriter
	Request *http.Request
}

const CookieName = "qid"
const ExpirationDuration = 365 * 24 * time.Hour
const SessionIDPrefix = "sess:"

var ctx = context.Background()

func (manager SessionManager) Create(userID string) error {
	// Create a new session ID
	sessionID := uuid.New().String()

	// Set the value of the seesion ID to the user ID in redis
	var err error
	err = manager.RDB.Set(ctx, SessionIDPrefix+sessionID, userID, ExpirationDuration).Err()
	if err != nil {
		return err
	}

	// Encode the cookie for safety
	cookie, err := encodeCookie(sessionID)
	if err != nil {
		// Delete the session ID from redis as the cookie could not be encoded and hence is not used
		manager.RDB.Del(ctx, SessionIDPrefix+sessionID)
		return err
	}

	// Set the cookie
	http.SetCookie(manager.Writer, cookie)

	return nil
}

func (manager SessionManager) Retrieve(onlySessionID bool) (string, error) {
	var err error
	cookie, err := manager.Request.Cookie(CookieName)
	if err != nil {
		return "", err
	}

	sessionID, err := decodeCookie(cookie)
	if err != nil {
		return "", err
	}

	if onlySessionID {
		return sessionID, nil
	}

	userID, err := manager.RDB.Get(ctx, sessionID).Result()
	switch {
	case err == redis.Nil:
		// key does not exist
		return "", fmt.Errorf("this session does not exist")
	case err != nil:
		// failed to get the value
		return "", err
	}

	return userID, nil
}

func (manager SessionManager) Delete(sessionID string) error {
	cookie := http.Cookie{
		Name:   CookieName,
		MaxAge: -1, // -1 deletes the cookie
	}

	http.SetCookie(manager.Writer, &cookie)

	err := manager.RDB.Del(ctx, sessionID).Err()
	if err != nil {
		return err
	}

	return nil
}
