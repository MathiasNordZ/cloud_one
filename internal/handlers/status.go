package handlers

import (
	"assignment_one/internal/errorHandling"
	"assignment_one/internal/jsonHandling"
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func InitStatus() {
	http.Handle("/countryinfo/v1/status", http.HandlerFunc(getStatus))
}

var (
	StartTime time.Time
	startOnce sync.Once
)

// getStatus is an http handler for the status endpoint.
func getStatus(w http.ResponseWriter, r *http.Request) {
	utils.CheckGET(w, r)

	restCountries := os.Getenv("COUNTRY_API")
	currencies := os.Getenv("CURRENCY_API_BASE")

	countryRes, err := utils.HttpClient.Get(restCountries)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(countryRes.Body)

	currencyRes, err := utils.HttpClient.Get(currencies)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}
	defer currencyRes.Body.Close()

	startOnce.Do(func() {
		StartTime = time.Now()
	})

	err = jsonHandling.EncodeJSON(w, http.StatusOK, structs.Status{
		RestCountriesApi: countryRes.Status,
		CurrenciesApi:    currencyRes.Status,
		Version:          "v1",
		Uptime:           time.Since(StartTime).String(),
	})
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}
}
