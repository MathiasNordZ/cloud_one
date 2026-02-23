package structs

type Exchange struct {
	Country       string             `json:"country"`
	Currency      string             `json:"base-currency"`
	ExchangeRates map[string]float64 `json:"exchange-rates"`
}

type CountryAPI struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Borders    []string `json:"borders"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

type CurrencyAPI struct {
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

type CountryCurrencyResponse []struct {
	CCA2       string `json:"cca2"`
	Currencies map[string]struct {
		Name string `json:"name"`
	} `json:"currencies"`
}
