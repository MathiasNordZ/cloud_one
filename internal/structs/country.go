package structs

// Country represents the struct being output from the info endpoint.
type Country struct {
	Name struct {
		Common string `jsonHandling:"common"`
	} `jsonHandling:"name"`

	Continents []string          `jsonHandling:"continents"`
	Population int               `jsonHandling:"population"`
	Area       float64           `jsonHandling:"area"`
	Languages  map[string]string `jsonHandling:"languages"`
	Borders    []string          `jsonHandling:"borders"`
	Flag       string            `jsonHandling:"flag"`
	Capital    []string          `jsonHandling:"capital"`
}
