package session

import (
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

func encodeCookie(s *securecookie.SecureCookie, value string) (*http.Cookie, error) {
	var secure bool
	if os.Getenv("GO_ENV") == "production" {
		secure = true
	} else {
		secure = false
	}

	encodedValue, err := s.Encode(CookieName, value)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     CookieName,
		Value:    encodedValue,
		MaxAge:   int(ExpirationDuration),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	return &cookie, nil
}

func decodeCookie(s *securecookie.SecureCookie, cookie *http.Cookie) (string, error) {
	var value string

	err := s.Decode(CookieName, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}
