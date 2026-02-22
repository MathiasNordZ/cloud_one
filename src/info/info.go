package info

import (
	"assignment_one/src/structs"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	country := CountryCode(r)
	if InputValidation(w, country) {
		return
	}

	countryURL := os.Getenv("COUNTRY_API")
	if countryURL == "" {
		http.Error(w, "COUNTRY_API not set", http.StatusInternalServerError)
		return
	}

	full, err := url.JoinPath(countryURL, "v3.1/alpha", country)
	if err != nil {
		http.Error(w, "Invalid COUNTRY_API base URL", http.StatusInternalServerError)
		log.Println("JoinPath error:", err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, full, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error when creating request:", err)
		return
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Upstream request failed", http.StatusBadGateway)
		log.Println("Error during request execution:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error with upstream API. "+resp.Status, http.StatusBadGateway)
		return
	}

	countryRes, err := CountryDecoder(resp.Body)
	if err != nil {
		http.Error(w, "Failed to decode upstream response", http.StatusBadGateway)
		log.Println("Decode error:", err)
		return
	}

	if len(countryRes) == 0 {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}
	JsonEncoder(w, countryRes[0])
}

func CountryDecoder(r io.Reader) ([]structs.Country, error) {
	var out []structs.Country
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func InputValidation(w http.ResponseWriter, country string) bool {
	if !regexp.MustCompile(`^[A-Za-z]{2}$`).MatchString(country) {
		http.Error(w, "Invalid country code. Use ISO3166 alpha-2 (two letters).", http.StatusBadRequest)
		return true
	}
	return false
}

func JsonEncoder(w http.ResponseWriter, single structs.Country) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(single); err != nil {
		log.Println("Failed to encode:", err)
	}
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
