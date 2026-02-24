package exchange

import (
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetExchange is an HTTP handler for the exchange endpoint.
// Expected route shape: /v1/exchange/{countryCode}
// Example: /v1/exchange/no
func GetExchange(w http.ResponseWriter, r *http.Request) {
	utils.CheckGET(w, r) // Check that method used is GET.

	cc := countryCodeFromRequest(r)
	if cc == "" {
		http.Error(w, "missing country code", http.StatusBadRequest)
		return
	}
	if utils.InputValidation(w, cc) {
		return
	}

	countryURL, err := utils.BuildAPIURL("COUNTRY_API", "v3.1/alpha/", cc)
	if utils.HandleErr(w, err, "failed to build country api url", http.StatusInternalServerError) {
		return
	}

	countryJSON, err := utils.GetJSON(countryURL, utils.HttpClient, w)
	if utils.HandleErr(w, err, "failed to fetch country api", http.StatusInternalServerError) {
		return
	}

	countryAPI, err := utils.DecodeJSON[[]structs.CountryAPI](countryJSON)
	if utils.HandleErr(w, err, "failed to create country", http.StatusInternalServerError) {
		return
	}

	borderCurrenciesJSON, err := getBorderCurrencies(countryAPI, w)
	if utils.HandleErr(w, err, "failed to fetch border currencies", http.StatusInternalServerError) {
		return
	}

	borderCurrencies, err := utils.DecodeJSON[structs.CountryCurrencyResponse](borderCurrenciesJSON)
	if utils.HandleErr(w, err, "failed to decode border currencies", http.StatusInternalServerError) {
		return
	}

	code := getFirstCurrencyCode(countryAPI)
	if code == "" {
		http.Error(w, fmt.Sprintf("no currency code for %s", countryAPI), http.StatusBadRequest)
	}

	currencyURL, err := utils.BuildAPIURL("CURRENCY_API", code)
	if utils.HandleErr(w, err, "failed to build currency api url", http.StatusInternalServerError) {
		return
	}

	currencyJSON, err := utils.GetJSON(currencyURL, utils.HttpClient, w)
	if utils.HandleErr(w, err, "failed to fetch currency api", http.StatusInternalServerError) {
		return
	}

	currencyAPI, err := utils.DecodeJSON[structs.CurrencyAPI](currencyJSON)
	if utils.HandleErr(w, err, "failed to decode currencies", http.StatusInternalServerError) {
		return
	}

	matchedRates := matchRates(currencyAPI, borderCurrencies)
	if len(matchedRates) == 0 || len(countryAPI) == 0 {
		http.Error(w, fmt.Sprintf("length of matched rates or country api is 0"), http.StatusBadGateway)
		return
	}

	resp := structs.Exchange{
		Country:       countryAPI[0].Name.Common,
		Currency:      currencyAPI.BaseCode,
		ExchangeRates: matchedRates,
	}
	encodeStruct(resp, w)
}

// getBorderCurrencies fetches border-country currency information from the Country API.
//
// It builds a query based on countryApi.Borders and requests a reduced response using fields.
// Returned JSON is expected to match structs.CountryCurrencyResponse.
func getBorderCurrencies(countryAPI []structs.CountryAPI, w http.ResponseWriter) ([]byte, error) {
	if len(countryAPI) == 0 {
		return nil, errors.New("country api is empty")
	}
	if len(countryAPI[0].Borders) == 0 {
		return []byte("[]"), nil
	}

	baseURL, err := utils.BuildAPIURL("COUNTRY_API", "v3.1/alpha")
	if utils.HandleErr(w, err, "failed to build api url", http.StatusInternalServerError) {
		return nil, err
	}

	u, err := url.Parse(baseURL)
	if utils.HandleErr(w, err, "failed to build api url", http.StatusInternalServerError) {
		return nil, err
	}

	q := u.Query()
	q.Set("codes", strings.Join(countryAPI[0].Borders, ","))
	q.Set("fields", "cca2,currencies")
	u.RawQuery = q.Encode()

	return utils.GetJSON(u.String(), utils.HttpClient, w)
}

// matchRates matches currency codes discovered in bordering countries with the exchange rates.
func matchRates(rates structs.CurrencyAPI, countries structs.CountryCurrencyResponse) map[string]float64 {
	result := make(map[string]float64)

	for _, country := range countries {
		for currencyCode := range country.Currencies {
			if rate, ok := rates.Rates[currencyCode]; ok {
				result[currencyCode] = rate
			}
			break
		}
	}
	return result
}

// countryCodeFromRequest extracts and normalizes the country code from the request path.
//
// For a route like /v1/exchange/no, it returns "NO".
func countryCodeFromRequest(r *http.Request) string {
	country := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/exchange/")
	country = strings.Trim(country, "/")
	country = strings.ToUpper(country)
	return country
}

// encodeStruct writes the exchange struct as JSON to the response writer.
func encodeStruct(exchange structs.Exchange, w http.ResponseWriter) {
	utils.EncodeJSON(w, http.StatusOK, exchange)
}

// getFirstCurrencyCode extracts the code of the first currency.
func getFirstCurrencyCode(c []structs.CountryAPI) string {
	for code := range c[0].Currencies {
		return code
	}
	return ""
}
