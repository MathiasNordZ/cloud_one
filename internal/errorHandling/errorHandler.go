package errorHandling

import (
	"errors"
	"net/http"
)

type httpError struct {
	message string
	status  int
}

func (e httpError) Error() string {
	return e.message
}

func NewHTTPError(msg string, status int) error {
	return httpError{message: msg, status: status}
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	var e httpError
	if errors.As(err, &e) {
		http.Error(w, e.message, e.status)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
