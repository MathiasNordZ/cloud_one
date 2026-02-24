package utils

import (
	"net/http"
	"time"
)

// HttpClient is the shared HTTP client used for upstream API calls.
var HttpClient = &http.Client{
	Timeout: 5 * time.Second,
}
