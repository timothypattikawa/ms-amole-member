package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GenerateHashPassword(password string) (string, error) {
	hashResult, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashResult), err
}

func ValidatePassword(hash, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	return true
}
