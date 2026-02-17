package main

import (
	"assignment_one/src/info"
	"assignment_one/src/status"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// Log error if .env file is not found, but don't exit if it's optional
		log.Println("Error loading .env file, assuming production environment variables are set")
	}

	http.Handle("/v1/status", http.HandlerFunc(status.GetStatus))
	http.Handle("/v1/info/", http.HandlerFunc(info.GetInfo))

	log.Println("Service is up and running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
