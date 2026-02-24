package status

import (
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	StartTime time.Time
	startOnce sync.Once
)

// GetStatus is an http handler for the status endpoint.
// Usage: /v1/status/
// Returns status of the API.
func GetStatus(w http.ResponseWriter, r *http.Request) {
	utils.CheckGET(w, r)

	restCountries := os.Getenv("COUNTRY_API")
	currencies := os.Getenv("CURRENCY_API_BASE")

	countryRes, err := utils.HttpClient.Get(restCountries)
	if utils.HandleErr(w, err, "failed to fetch country api", http.StatusInternalServerError) {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(countryRes.Body)

	currencyRes, err := utils.HttpClient.Get(currencies)
	if utils.HandleErr(w, err, "failed to fetch currency api", http.StatusInternalServerError) {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(currencyRes.Body)

	startOnce.Do(func() {
		StartTime = time.Now()
	})

	utils.EncodeJSON(w, http.StatusOK, structs.Status{
		RestCountriesApi: countryRes.Status,
		CurrenciesApi:    currencyRes.Status,
		Version:          "v1",
		Uptime:           time.Since(StartTime).String(),
	})
}
