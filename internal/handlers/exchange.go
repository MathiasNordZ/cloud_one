package handlers

import (
	"assignment_one/internal/exchange"
	"net/http"
)

func InitExchange() {
	http.Handle("/countryinfo/v1/exchange/", http.HandlerFunc(exchange.GetExchange))
}
