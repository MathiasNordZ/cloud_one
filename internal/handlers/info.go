package handlers

import (
	"assignment_one/internal/info"
	"net/http"
)

func InitInfo() {
	http.Handle("/countryinfo/v1/info/", http.HandlerFunc(info.GetInfo))
}
