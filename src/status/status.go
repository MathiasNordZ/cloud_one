package status

import (
	"assignment_one/src/request"
	"assignment_one/src/structs"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var StartTime time.Time

func GetStatus(w http.ResponseWriter, r *http.Request) {
	restCountries := os.Getenv("COUNTRY_API")
	currencies := os.Getenv("CURRENCY_API")

	countryRes, err := request.Get(restCountries)
	if err != nil {
		http.Error(w, "Error contacting countries API: "+err.Error(), http.StatusBadGateway)
		return
	}

	currencyRes, err := request.Get(currencies)
	if err != nil {
		http.Error(w, "Error contacting currency API: "+err.Error(), http.StatusBadGateway)
		return
	}

	if StartTime.IsZero() {
		StartTime = time.Now()
	}

	statusResponse := structs.Status{
		RestCountriesApi: countryRes.Status,
		CurrenciesApi:    currencyRes.Status,
		Version:          "v1",
		Uptime:           time.Since(StartTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusResponse); err != nil {
		log.Println("Failed to encode:", err)
	}
}
