package utils

import (
	"fmt"
	"net/http"
)

// HandleErr handles errors im an uniform way.
func HandleErr(w http.ResponseWriter, err error, msg string, code int) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), code)
		return true
	}
	return false
}
