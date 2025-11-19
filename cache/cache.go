package cache

import (
	"net/http"
	"net/url"
)

type Cache struct {
	Origin  *url.URL
	Path    string // can be nil
	Body    []byte
	Status  int
	Method  string
	Headers http.Header
}

func NewCache(origin *url.URL, path string, body []byte, status int, method string, headers http.Header) Cache {
	return Cache{
		Origin:  origin,
		Path:    path,
		Body:    body,
		Status:  status,
		Method:  method,
		Headers: headers,
	}
}
