package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashItem(item string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(item), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash item %v", err)
	}
	return string(hashedPassword), nil
}

func VerifyItem(hashedItem string, candidateItem string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedItem), []byte(candidateItem))
}
