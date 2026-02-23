package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func BuildAPIURL(envKey string, parts ...string) (string, bool) {
	base := os.Getenv(envKey)
	if base == "" {
		fmt.Println("Environment variable " + envKey + " is not set")
		return "", false
	}

	full, err := url.JoinPath(base, parts...)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	return full, true
}

/*
Method that extracts country code from url.
*/
func CountryCode(r *http.Request) string {
	country := strings.TrimPrefix(r.URL.Path, "/v1/info/")
	country = strings.Trim(country, "/")
	country = strings.ToUpper(country)
	return country
}
