package structs

type Status struct {
	RestCountriesApi string `json:"restCountriesApi"`
	CurrenciesApi    string `json:"currenciesApi"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}
