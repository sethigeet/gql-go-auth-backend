package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashed the given password using bcrypt and returns the string form of the hashed password
func HashPassword(passwd string) (string, error) {
	bytePasswd := []byte(passwd)
	hash, err := bcrypt.GenerateFromPassword(bytePasswd, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

// ComparePasswords compares a hashed password with plain string password and
// returns a boolean representing whether the passwords match or not
func ComparePasswords(hashedPasswd string, plainPasswd []byte) bool {
	byteHash := []byte(hashedPasswd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPasswd)

	return err == nil
}
