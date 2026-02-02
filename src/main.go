package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status struct {
	RestCountriesApi string `json:"RestCountriesApi"`
	CurrenciesApi    string `json:"CurrenciesApi"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	const restCountries string = "http://129.241.150.113:8080/"
	const currencies string = "http://129.241.150.113:9090/currency/"

	countryRes := requestCountry(restCountries)
	currencyRes := requestCurrency(currencies)

	status := Status{RestCountriesApi: countryRes.Status, CurrenciesApi: currencyRes.Status, Version: "v1", Uptime: "0000"}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func requestCountry(restCountries string) *http.Response {
	res, err := http.Get(restCountries)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func requestCurrency(currencyUrl string) *http.Response {
	res, err := http.Get(currencyUrl)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

/*
* https://medium.com/moesif/building-a-restful-api-with-go-dbd6e7aecf87
 */
func main() {
	http.Handle("/v1/status", http.HandlerFunc(getStatus))
	log.Println("Service is up and running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
