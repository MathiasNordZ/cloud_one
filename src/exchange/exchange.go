package exchange

import (
	"assignment_one/src/structs"
	"assignment_one/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// httpClient is the shared HTTP client used for upstream API calls.
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// GetExchange is an HTTP handler for the exchange endpoint.
// Expected route shape: /v1/exchange/{countryCode}
// Example: /v1/exchange/no
func GetExchange(w http.ResponseWriter, r *http.Request) {
	cc := countryCodeFromRequest(r)
	if cc == "" {
		http.Error(w, "missing country code", http.StatusBadRequest)
		return
	}

	countryURL, ok := utils.BuildAPIURL("COUNTRY_API", "v3.1/alpha/", cc)
	if !ok {
		http.Error(w, "failed to build country api url", http.StatusInternalServerError)
		return
	}

	countryJSON, err := getCountry(countryURL, httpClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch country: %v", err), http.StatusBadGateway)
		return
	}

	countryAPI, err := createCountry(countryJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode country response: %v", err), http.StatusBadGateway)
		return
	}

	borderCurrenciesJSON, err := getBorderCurrencies(countryAPI, httpClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch border currencies: %v", err), http.StatusBadGateway)
		return
	}

	borderCurrencies, err := decodeBorderCurrencies(borderCurrenciesJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode border currencies: %v", err), http.StatusBadGateway)
		return
	}

	code := getFirstCurrencyCode(countryAPI)
	currencyURL, ok := utils.BuildAPIURL("CURRENCY_API", code)
	fmt.Println(currencyURL)
	if !ok {
		http.Error(w, "failed to build currency api url", http.StatusInternalServerError)
		return
	}

	currencyJSON, err := getCurrencies(currencyURL, httpClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch currency rates: %v", err), http.StatusBadGateway)
		return
	}

	currencyAPI, err := decodeCurrencies(currencyJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode currency response: %v", err), http.StatusBadGateway)
		return
	}

	matchedRates := matchRates(currencyAPI, borderCurrencies)

	resp := structs.Exchange{
		Country:       countryAPI[0].Name.Common,
		Currency:      currencyAPI.BaseCode,
		ExchangeRates: matchedRates,
	}
	encodeStruct(resp, w)
}

// getCountry sends a GET request to the Country API and returns the raw JSON response body.
func getCountry(apiURL string, client *http.Client) ([]byte, error) {
	if apiURL == "" {
		return nil, errors.New("apiURL is empty")
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream error: %s (url=%s, body=%s)", resp.Status, apiURL, body)
	}

	return body, nil
}

// createCountry does unmarshal a Country API JSON response into a structs.CountryAPI.
func createCountry(input []byte) ([]structs.CountryAPI, error) {
	var out []structs.CountryAPI
	if err := json.Unmarshal(input, &out); err != nil {
		return out, err
	}
	return out, nil
}

// getBorderCurrencies fetches border-country currency information from the Country API.
//
// It builds a query based on countryApi.Borders and requests a reduced response using fields.
// Returned JSON is expected to match structs.CountryCurrencyResponse.
func getBorderCurrencies(countryAPI []structs.CountryAPI, client *http.Client) ([]byte, error) {
	if len(countryAPI) == 0 {
		return nil, errors.New("country api is empty")
	}
	if len(countryAPI[0].Borders) == 0 {
		return []byte("[]"), nil
	}

	query := strings.Join(countryAPI[0].Borders, ",")
	fmt.Println("Query: ", query)

	apiURL, ok := utils.BuildAPIURL(
		"COUNTRY_API",
		"v3.1/alpha",
	)
	apiURL += "?codes=" + query + "&fields=cca2,currencies"
	fmt.Println("API URL: ", apiURL)
	if !ok || apiURL == "" {
		return nil, errors.New("failed to build api url")
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream error: %s", resp.Status)
	}
	return body, nil
}

// decodeBorderCurrencies does unmarshal border-currency JSON into a structs.CountryCurrencyResponse.
func decodeBorderCurrencies(input []byte) (structs.CountryCurrencyResponse, error) {
	var out structs.CountryCurrencyResponse
	if err := json.Unmarshal(input, &out); err != nil {
		return out, err
	}
	return out, nil
}

// getCurrencies sends a GET request to the Currency API and returns the raw JSON response body.
//
// The Currency API response is expected to contain rates for many currencies.
func getCurrencies(apiURL string, client *http.Client) ([]byte, error) {
	if apiURL == "" {
		return nil, errors.New("apiURL is empty")
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// decodeCurrencies does unmarshal Currency API JSON into a structs.CurrencyAPI.
func decodeCurrencies(input []byte) (structs.CurrencyAPI, error) {
	var out structs.CurrencyAPI
	if err := json.Unmarshal(input, &out); err != nil {
		return out, err
	}
	return out, nil
}

// matchRates matches currency codes discovered in bordering countries with the exchange rates.
//
// It returns a map where keys are currency codes (e.g., "SEK") and values are the corresponding rate
// from rates.Rates (base is rates.BaseCode).
//
// Note: This assumes one currency per country and will pick the first encountered code.
func matchRates(rates structs.CurrencyAPI, countries structs.CountryCurrencyResponse) map[string]float64 {
	result := make(map[string]float64)

	for _, country := range countries {
		for currencyCode := range country.Currencies {
			if rate, ok := rates.Rates[currencyCode]; ok {
				result[currencyCode] = rate
			}
			break // assume one currency per country
		}
	}

	return result
}

// countryCodeFromRequest extracts and normalizes the country code from the request path.
//
// For a route like /v1/exchange/no, it returns "NO".
func countryCodeFromRequest(r *http.Request) string {
	country := strings.TrimPrefix(r.URL.Path, "/v1/exchange/")
	country = strings.Trim(country, "/")
	country = strings.ToUpper(country)
	return country
}

// encodeStruct writes the exchange struct as JSON to the response writer.
func encodeStruct(exchange structs.Exchange, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exchange); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func getFirstCurrencyCode(c []structs.CountryAPI) string {
	for code := range c[0].Currencies {
		return code
	}
	return ""
}
