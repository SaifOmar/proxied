package proxy

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"

	"github.com/SaifOmar/proxied/cache"
)

type Proxy struct {
	Origin *url.URL
	Cached map[string]cache.Cache
}

func NewProxy(origin *url.URL) *Proxy {

	return &Proxy{
		Origin: origin,
		Cached: make(map[string]cache.Cache),
	}

}

// TODO: improve
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rTarget := r.URL.Path
	rMethod := r.Method
	// rQuery := r.URL.Query()
	ss := rMethod + ":" + rTarget
	// path one // cached value and we can get
	value, exists := p.Cached[ss]
	if exists {
		fmt.Println("Cached value exists for " + ss)
		WriteResponseWithHeaders(w, value.Body, value.Status, value.Headers)
		return
	}

	fmt.Println("New target is: ", rTarget)
	fmt.Println("Making GET request to " + p.Origin.Host)

	resp, err := http.Get("https://" + p.Origin.Hostname() + rTarget)

	body, err := io.ReadAll(resp.Body)
	defer r.Body.Close() // Ensure the response body is closed
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	p.Cached[ss] = cache.NewCache(p.Origin, rTarget, body, resp.StatusCode, rMethod, resp.Header)

	fmt.Println("Couldn't find cached value for " + ss)
	fmt.Println("Created new cache for " + ss)

	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)

	// write to the client
	defer WriteResponseWithHeaders(w, body, resp.StatusCode, resp.Header)
}

func WriteResponseWithHeaders(w http.ResponseWriter, body []byte, status int, headers http.Header) {
	w.WriteHeader(status)
	maps.Copy(w.Header(), headers)

	w.Write(body)
}
