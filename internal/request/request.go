package request

import "net/http"

func Get(url string) (*http.Response, error) {
	return http.Get(url)
}
