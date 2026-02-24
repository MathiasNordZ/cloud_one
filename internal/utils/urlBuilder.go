package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// BuildAPIURL builds url with env for api calls.
func BuildAPIURL(envKey string, parts ...string) (string, error) {
	base := os.Getenv(envKey)
	if base == "" {
		return "", fmt.Errorf("environment variable %s is not set", envKey)
	}
	full, err := url.JoinPath(base, parts...)
	if err != nil {
		return "", fmt.Errorf("failed to join path: %w", err)
	}
	return full, nil
}

// CountryCode Method that extracts country code from url.
func CountryCode(r *http.Request) string {
	country := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/info/")
	country = strings.Trim(country, "/")
	country = strings.ToUpper(country)
	return country
}
