package structs

// Status is the struct being used for the output in the status endpoint.
type Status struct {
	RestCountriesApi string `jsonHandling:"restCountriesApi"`
	CurrenciesApi    string `jsonHandling:"currenciesApi"`
	Version          string `jsonHandling:"version"`
	Uptime           string `jsonHandling:"uptime"`
}
