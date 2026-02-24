package handlers

import (
	"assignment_one/internal/status"
	"net/http"
)

func InitStatus() {
	http.Handle("/countryinfo/v1/status", http.HandlerFunc(status.GetStatus))
}
