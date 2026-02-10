package main

import (
	"assignment_one/src/request"
	"assignment_one/src/status"
	"assignment_one/src/structs"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// Log error if .env file is not found, but don't exit if it's optional
		log.Println("Error loading .env file, assuming production environment variables are set")
	}

	http.Handle("/v1/status", http.HandlerFunc(status.GetStatus))
	http.Handle("/info/", http.HandlerFunc(GetInfo))

	log.Println("Service is up and running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	country := strings.TrimPrefix(r.URL.Path, "/info/")
	country = strings.Trim(country, "/")

	if country == "" {
		http.Error(w, "Seems like no country is specified.", http.StatusBadRequest)
		return
	}

	res, err := request.Get("http://129.241.150.113:8080/v3.1/alpha/" + country)

	if err != nil {
		http.Error(w, "Error contacting country API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var countryRes []structs.Country
	err = decoder.Decode(&countryRes)
	if err != nil {
		log.Fatal(err)
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
