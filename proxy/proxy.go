package proxy

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"sync"

	"github.com/SaifOmar/proxied/cache"
)

type Proxy struct {
	Origin *url.URL
	Cached map[string]cache.Cache
	mu     sync.RWMutex
}

func NewProxy(origin *url.URL) *Proxy {

	return &Proxy{
		Origin: origin,
		Cached: make(map[string]cache.Cache),
	}

}

// TODO: improve
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var body []byte
	var status int
	var headers http.Header
	rTarget := r.URL.Path
	rMethod := r.Method
	// rQuery := r.URL.Query()
	ss := rMethod + ":" + rTarget
	// path one // cached value and we can get
	value, ok := GetCache(ss, p)
	if ok {
		fmt.Println("Cached value exists for " + ss)
		body = value.Body
		status = value.Status
		headers = value.Headers
	} else {

		fmt.Println("New target is: ", rTarget)
		fmt.Println("Making GET request to " + p.Origin.Host)

		resp, err := http.Get("https://" + p.Origin.Hostname() + rTarget)
		if err != nil {
			log.Fatalf("Error making GET request: %v", err)
		}

		body, err = io.ReadAll(resp.Body)
		defer r.Body.Close()
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		c := cache.NewCache(p.Origin, rTarget, body, resp.StatusCode, rMethod, resp.Header)
		SetCached(ss, p, c)

		fmt.Println("Couldn't find cached value for " + ss)
		fmt.Println("Created new cache for " + ss)

		if err != nil {
			log.Fatalf("Error making GET request: %v", err)
		}

		fmt.Printf("Status: %s\n", resp.Status)
		status = resp.StatusCode
		headers = resp.Header

	}
	// write to the client
	WriteResponseWithHeaders(w, body, status, headers)
}

func WriteResponseWithHeaders(w http.ResponseWriter, body []byte, status int, headers http.Header) {
	w.WriteHeader(status)
	maps.Copy(w.Header(), headers)

	w.Write(body)
}
func SetCached(s string, p *Proxy, c cache.Cache) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Cached[s] = c
}

func GetCache(s string, p *Proxy) (cache.Cache, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.Cached[s]
	return v, ok
}
