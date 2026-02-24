package utils

import (
	"net/http"
	"regexp"
)

// InputValidation ensures that input is 2 letters, a-zA-z.
func InputValidation(w http.ResponseWriter, country string) bool {
	if !regexp.MustCompile(`^[A-Za-z]{2}$`).MatchString(country) {
		http.Error(w, "Invalid country code. Use ISO3166 alpha-2 (two letters).", http.StatusBadRequest)
		return true
	}
	return false
}
