package auth

import (
	"fmt"
	"slices"
)

func ValidateUsername(username string) error {
	if len(username) < 4 {
		return fmt.Errorf("username needs to be at least 4 characters long")
	}

	if len(username) > 20 {
		return fmt.Errorf("username can only be max. 20 characters long")
	}

	notAllowed := []string{
		"admin",
		"administrator",
		"support",
		"help",
		"vivery",
	}
	if slices.Contains(notAllowed, username) {
		return fmt.Errorf("username not allowed")
	}

	return nil
}
