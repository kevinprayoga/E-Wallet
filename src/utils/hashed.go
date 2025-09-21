package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashString(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing string: %v", err)
	}
	return string(hash)
}

func ValidateHashedString(hashed, input string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	if err != nil {
		return errors.New("password or pin does not match")
	}
	return nil
}
