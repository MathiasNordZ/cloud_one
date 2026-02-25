package structs

// Exchange represents the struct being output by the exchange endpoint.
type Exchange struct {
	Country       string             `jsonHandling:"country"`
	Currency      string             `jsonHandling:"base-currency"`
	ExchangeRates map[string]float64 `jsonHandling:"exchange-rates"`
}

// CountryAPI is being used as a temporary transferring struct between the input and output.
type CountryAPI struct {
	Name struct {
		Common string `jsonHandling:"common"`
	} `jsonHandling:"name"`
	Borders    []string `jsonHandling:"borders"`
	Currencies map[string]struct {
		Name   string `jsonHandling:"name"`
		Symbol string `jsonHandling:"symbol"`
	} `jsonHandling:"currencies"`
}

// CurrencyAPI is being used as a temporary transferring struct between the input and output.
type CurrencyAPI struct {
	BaseCode string             `jsonHandling:"base_code"`
	Rates    map[string]float64 `jsonHandling:"rates"`
}

// CountryCurrencyResponse is being used as a temporary transferring struct between the input and output.
type CountryCurrencyResponse []struct {
	CCA2       string `jsonHandling:"cca2"`
	Currencies map[string]struct {
		Name string `jsonHandling:"name"`
	} `jsonHandling:"currencies"`
}
