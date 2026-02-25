package json

import (
	"assignment_one/internal/errorHandling"
	"encoding/json"
	"io"
	"net/http"
)

// GetJSON sends a GET request to the Country API and returns the raw JSON response body.
func GetJSON(apiURL string, client *http.Client) ([]byte, error) {
	if apiURL == "" {
		return nil, errorHandling.NewHTTPError("apiURL is required", http.StatusBadRequest)
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, errorHandling.NewHTTPError("Error fetching JSON", http.StatusBadGateway)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorHandling.NewHTTPError("Error reading JSON", http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errorHandling.NewHTTPError("Error fetching JSON", resp.StatusCode)
	}
	return body, nil
}

// DecodeJSON decodes JSON bytes into the given target.
func DecodeJSON[T any](b []byte) (T, error) {
	var out T
	if len(b) == 0 {
		return out, errorHandling.NewHTTPError("Empty JSON", http.StatusInternalServerError)
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, errorHandling.NewHTTPError("Error decoding JSON", http.StatusBadRequest)
	}
	return out, nil
}

// EncodeJSON writes v to JSON
func EncodeJSON(w http.ResponseWriter, status int, v any) error {
	if status == 0 {
		status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errorHandling.NewHTTPError("Error encoding JSON", http.StatusInternalServerError)
	}
	return nil
}
