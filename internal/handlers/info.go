package handlers

import (
	"assignment_one/internal/errorHandling"
	"assignment_one/internal/jsonHandling"
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"net/http"
	"net/url"
	"os"
)

// InitInfo serves as the handler for the info endpoint.
func InitInfo() {
	http.Handle("/countryinfo/v1/info/", http.HandlerFunc(getInfo))
}

// getInfo assembles all methods into one method that is served as the handler.
func getInfo(w http.ResponseWriter, r *http.Request) {
	utils.CheckGET(w, r)

	country := utils.CountryCode(r)
	err := utils.InputValidation(country)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}

	countryURL := os.Getenv("COUNTRY_API")
	if countryURL == "" {
		http.Error(w, "COUNTRY_API not set", http.StatusInternalServerError)
		return
	}

	full, err := url.JoinPath(countryURL, "v3.1/alpha", country)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}

	resp, err := jsonHandling.GetJSON(full, utils.HttpClient)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}

	countryRes, err := jsonHandling.DecodeJSON[[]structs.Country](resp)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}

	if len(countryRes) == 0 {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}
	err = jsonHandling.EncodeJSON(w, 0, countryRes)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
	}
}
