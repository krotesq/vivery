package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generates and signes a new JWT with given secret and exp (minutes)
func GenerateJWT(sub, iss string, secret []byte, exp int) (string, error) {

	now := time.Now()

	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": iss,
			"sub": sub,
			"iat": now.Unix(),
			"exp": now.Add(time.Minute * time.Duration(exp)).Unix(),
			"jti": jti,
		},
	)
	st, err := t.SignedString(secret)
	return st, err
}

func ValidateJWT(token string, secret []byte) (string, error) {

	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		},
	)

	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("missing subject")
	}

	return sub, nil
}

// creates a cryptographically random unique token
func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate jti: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
