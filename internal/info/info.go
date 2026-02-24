package info

import (
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	utils.CheckGET(w, r)

	country := utils.CountryCode(r)
	if utils.InputValidation(w, country) {
		return
	}

	countryURL := os.Getenv("COUNTRY_API")
	if countryURL == "" {
		http.Error(w, "COUNTRY_API not set", http.StatusInternalServerError)
		return
	}

	full, err := url.JoinPath(countryURL, "v3.1/alpha", country)
	if utils.HandleErr(w, err, "Failed to join url", http.StatusInternalServerError) {
		return
	}

	req, err := http.NewRequest(http.MethodGet, full, nil)
	if utils.HandleErr(w, err, "Failed to create request", http.StatusInternalServerError) {
		return
	}

	resp, err := utils.HttpClient.Do(req)
	if utils.HandleErr(w, err, "Failed to do request", http.StatusInternalServerError) {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error with upstream API. "+resp.Status, http.StatusBadGateway)
		return
	}

	countryRes, err := CountryDecoder(resp.Body)
	if utils.HandleErr(w, err, "Failed to decode response body", http.StatusInternalServerError) {
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

func JsonEncoder(w http.ResponseWriter, single structs.Country) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(single); err != nil {
		log.Println("Failed to encode:", err)
	}
}
