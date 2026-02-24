package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// GetJSON sends a GET request to the Country API and returns the raw JSON response body.
func GetJSON(apiURL string, client *http.Client, w http.ResponseWriter) ([]byte, error) {
	if apiURL == "" {
		return nil, errors.New("apiURL is empty")
	}

	resp, err := client.Get(apiURL)
	if HandleErr(w, err, "failed to fetch response", http.StatusInternalServerError) {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if HandleErr(w, err, "failed to read body", http.StatusInternalServerError) {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream error: %s (url=%s, body=%s)", resp.Status, apiURL, body)
	}
	return body, nil
}

// DecodeJSON decodes JSON bytes into the given target.
func DecodeJSON[T any](b []byte) (T, error) {
	var out T
	if len(b) == 0 {
		return out, errors.New("empty json")
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}

// EncodeJSON writes v to JSON
func EncodeJSON(w http.ResponseWriter, status int, v any) {
	if status == 0 {
		status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("failed to encode:", err)
	}
}
