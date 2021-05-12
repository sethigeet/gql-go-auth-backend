package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwd string) (string, error) {
	bytePasswd := []byte(passwd)
	hash, err := bcrypt.GenerateFromPassword(bytePasswd, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)

	return err == nil
}
