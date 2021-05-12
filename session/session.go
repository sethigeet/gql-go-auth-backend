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

// CookieName is the name of the cookie that is stored in the browser
const CookieName = "qid"

// ExpirationDuration is the duration after which the cookie expires in the browser
const ExpirationDuration = 365 * 24 * time.Hour

// SessionIDPrefix is the string that is prefixed before the actual sessionID while storing in redis
// so that it will be easily identifiable and it will not collide with any other key
const SessionIDPrefix = "sess:"

var ctx = context.Background()

// Create creates a session of the user and stores the session ID in a cookie in the user's
// browser and in redis
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

// Retrieve takes a boolean argument and return the sessionID or userID accoring to it
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

	userID, err := manager.RDB.Get(ctx, SessionIDPrefix+sessionID).Result()
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

// Delete deletes the session using the session ID from redis and also deletes
// the cookie from the user's browser
func (manager SessionManager) Delete(sessionID string) error {
	cookie := http.Cookie{
		Name:   CookieName,
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(manager.Writer, &cookie)

	err := manager.RDB.Del(ctx, SessionIDPrefix+sessionID).Err()
	if err != nil {
		return err
	}

	return nil
}
