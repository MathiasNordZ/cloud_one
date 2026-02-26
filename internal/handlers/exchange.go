package handlers

import (
	"assignment_one/internal/errorHandling"
	"assignment_one/internal/jsonHandling"
	"assignment_one/internal/structs"
	"assignment_one/internal/utils"
	"net/http"
	"net/url"
	"strings"
)

// InitExchange creates an http handler for the getExchange method.
// The handler serves on the endpoint /countryinfo/v1/exchange/{country}
func InitExchange() {
	http.Handle("/countryinfo/v1/exchange/", http.HandlerFunc(getExchange))
}

// getExchange is the main method of the endpoint.
// It assembles all methods into one method for the http handler.
func getExchange(w http.ResponseWriter, r *http.Request) {
	cc, err := validate(r)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
		return
	}

	resp, err := buildExchange(cc)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
		return
	}
	err = jsonHandling.EncodeJSON(w, http.StatusOK, resp)
	if err != nil {
		errorHandling.WriteHTTPError(w, err)
		return
	}
}

// validate handles the validation logic.
// It validates that correct method is used, as well as the inputted country code.
func validate(r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		return "", errorHandling.NewHTTPError("Method not allowed.", http.StatusMethodNotAllowed)
	}
	cc := extractCC(r)
	if cc == "" {
		return "", errorHandling.NewHTTPError("Missing country code.", http.StatusBadRequest)
	}

	if err := utils.InputValidation(cc); err != nil {
		return "", errorHandling.NewHTTPError("Invalid input.", http.StatusBadRequest)
	}
	return cc, nil
}

// buildExchange handles the process of building the exchange struct before returning it.
func buildExchange(cc string) (structs.Exchange, error) {
	countryAPI, err := fetchCountry(cc)
	if err != nil {
		return structs.Exchange{}, err
	}

	borderCurrencies, err := fetchBorders(countryAPI)
	if err != nil {
		return structs.Exchange{}, err
	}

	code := getFirstCurrencyCode(countryAPI)

	currencyAPI, err := fetchCurrency(countryAPI)
	if err != nil {
		return structs.Exchange{}, err
	}

	rates := matchRates(currencyAPI, borderCurrencies)
	if len(rates) == 0 {
		return structs.Exchange{}, errorHandling.NewHTTPError("No matching border currencies.", http.StatusInternalServerError)
	}
	return structs.Exchange{
		Country:       countryAPI[0].Name.Common,
		Currency:      code,
		ExchangeRates: rates,
	}, nil
}

// fetchCountry handles the fetching and decoding logic of a single country.
func fetchCountry(cc string) ([]structs.CountryAPI, error) {
	countryUrl, err := utils.BuildAPIURL("COUNTRY_API", "v3.1/alpha/", cc)
	if err != nil {
		return nil, err
	}
	json, err := jsonHandling.GetJSON(countryUrl, utils.HttpClient)
	if err != nil {
		return nil, errorHandling.NewHTTPError("country api unavailable", http.StatusBadGateway)
	}
	return jsonHandling.DecodeJSON[[]structs.CountryAPI](json)
}

// fetchBorders fetches information of the bordering countries of the provided country.
// The information about the bordering countries is decoded into a struct.
func fetchBorders(countryAPI []structs.CountryAPI) (structs.CountryCurrencyResponse, error) {
	if len(countryAPI) == 0 {
		return nil, errorHandling.NewHTTPError("CountryAPI is empty.", http.StatusInternalServerError)
	}

	if len(countryAPI[0].Borders) == 0 {
		return structs.CountryCurrencyResponse{}, nil
	}

	baseURL, err := utils.BuildAPIURL("COUNTRY_API", "v3.1/alpha")
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("codes", strings.Join(countryAPI[0].Borders, ","))
	q.Set("fields", "cca2,currencies")
	u.RawQuery = q.Encode()

	json, err := jsonHandling.GetJSON(u.String(), utils.HttpClient)
	if err != nil {
		return nil, errorHandling.NewHTTPError("Failed to get JSON from API", http.StatusBadGateway)
	}
	return jsonHandling.DecodeJSON[structs.CountryCurrencyResponse](json)
}

// fetchCountry fetches information about a country's currency and decodes it into a CurrencyAPI struct.
func fetchCurrency(countryAPI []structs.CountryAPI) (structs.CurrencyAPI, error) {
	code := getFirstCurrencyCode(countryAPI)
	if code == "" {
		return structs.CurrencyAPI{}, errorHandling.NewHTTPError("missing currency code", http.StatusBadGateway)
	}

	currencyUrl, err := utils.BuildAPIURL("CURRENCY_API", code)
	if err != nil {
		return structs.CurrencyAPI{}, err
	}

	json, err := jsonHandling.GetJSON(currencyUrl, utils.HttpClient)
	if err != nil {
		return structs.CurrencyAPI{}, errorHandling.NewHTTPError("currency api unavailable", http.StatusBadGateway)
	}
	return jsonHandling.DecodeJSON[structs.CurrencyAPI](json)
}

// matchRates matches the rates in the CurrencyAPI struct against the wanted currencies. It then returns a map of the currency and rate.
func matchRates(rates structs.CurrencyAPI, countries structs.CountryCurrencyResponse) map[string]float64 {
	result := make(map[string]float64)

	for _, country := range countries {
		for code := range country.Currencies {
			if rate, ok := rates.Rates[code]; ok {
				result[code] = rate
			}
			break
		}
	}
	return result
}

// extractCC extracts the country code.
func extractCC(r *http.Request) string {
	country := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/exchange/")
	country = strings.Trim(country, "/")
	return strings.ToUpper(country)
}

// getFirstCurrencyCode extracts the first currency of the country.
func getFirstCurrencyCode(c []structs.CountryAPI) string {
	if len(c) == 0 {
		return ""
	}
	for code := range c[0].Currencies {
		return code
	}
	return ""
}
