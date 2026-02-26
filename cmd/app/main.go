package main

import (
	"assignment_one/internal/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// main the main entrypoint of the program.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading enviromental variables.")
	}

	handlers.InitStatus()
	handlers.InitInfo()
	handlers.InitExchange()

	log.Println("Service is up and running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
