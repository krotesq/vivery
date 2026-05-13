package auth

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// checks that the password meets complexity requirements
func ValidatePassword(password, username string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if strings.EqualFold(password, username) {
		return errors.New("password must not be the same as the username")
	}

	if strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		return errors.New("password must not contain the username")
	}

	var hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsDigit(ch):
			hasDigit = true
		case !unicode.IsLetter(ch) && !unicode.IsDigit(ch):
			hasSpecial = true
		}
	}

	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// get hash from password
func HashPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

// verify password hash
func ComparePassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}