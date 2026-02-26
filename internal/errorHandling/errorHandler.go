package errorHandling

import (
	"errors"
	"net/http"
)

// httpError struct representing an error.
type httpError struct {
	message string
	status  int
}

// Error method that returns an error message.
func (e httpError) Error() string {
	return e.message
}

// NewHTTPError method that creates a new error.
func NewHTTPError(msg string, status int) error {
	return httpError{message: msg, status: status}
}

// WriteHTTPError writes the error inside the handlers.
func WriteHTTPError(w http.ResponseWriter, err error) {
	var e httpError
	if errors.As(err, &e) {
		http.Error(w, e.message, e.status)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
