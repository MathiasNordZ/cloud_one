package main

import (
	"assignment_one/internal/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assuming production environment variables are set")
	}

	handlers.InitStatus()
	handlers.InitInfo()
	handlers.InitExchange()

	log.Println("Service is up and running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
