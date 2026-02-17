package info

import (
	"assignment_one/src/structs"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	country := strings.TrimPrefix(r.URL.Path, "/v1/info/")
	country = strings.Trim(country, "/")
	countryUrl := os.Getenv("COUNTRY_API")

	if inputValidation(w, country) {
		return
	}

	res, err := http.NewRequest(http.MethodGet, countryUrl+"v3.1/alpha/"+country, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error when creating request.", err.Error())
		return
	}
	c := http.Client{
		Timeout: time.Second * 5,
	}
	defer c.CloseIdleConnections()

	resp, err := c.Do(res)
	if err != nil {
		log.Println("Error during request execution.", err)
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error with upstream API. "+resp.Status, http.StatusBadGateway)
		return
	}

	decoder := json.NewDecoder(resp.Body)

	var countryRes []structs.Country
	err = decoder.Decode(&countryRes)
	if err != nil {
		log.Println("Error during json decoding.", err)
	}

	if len(countryRes) == 0 {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	single := countryRes[0]

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(single); err != nil {
		log.Println("Failed to encode:", err)
	}
}

/*
Method that validates the inputted country code according to ISO3166.
This is to prevent illegal inputs into the url.
*/
func inputValidation(w http.ResponseWriter, country string) bool {
	if country == "" {
		http.Error(w, "No country specified.", http.StatusBadRequest)
		return true
	} else if len(country) > 3 {
		http.Error(w, "Invalid country code.", http.StatusBadRequest)
		return true
	} else if !regexp.MustCompile("^[a-zA-Z]{2,3}$").MatchString(country) {
		http.Error(w, "Wrong input format.", http.StatusBadRequest)
		return true
	}
	return false
}
