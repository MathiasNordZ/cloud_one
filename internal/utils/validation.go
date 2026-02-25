package utils

import (
	"errors"
	"regexp"
)

// InputValidation ensures that input is 2 letters, a-zA-z.
func InputValidation(country string) error {
	if !regexp.MustCompile(`^[A-Za-z]{2}$`).MatchString(country) {
		return errors.New("invalid country code")
	}
	return nil
}
