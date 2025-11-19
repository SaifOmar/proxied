package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/SaifOmar/proxied/proxy"
)

func Exc() {
	port, url, err := parseArgs()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(port, url)
	startServer(fmt.Sprintf(":%d", port), url)
}

func parseArgs() (int, string, error) {
	url := flag.String("url", "", "Target URL to proxy")
	port := flag.Int("port", 8080, "Port to run the proxy on")

	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url is required")
		flag.Usage()
		return 0, "", fmt.Errorf("Error: --url is required")
	}

	fmt.Printf("Starting proxy server...\nURL: %s\nPort: %d\n", *url, *port)
	return *port, *url, nil
}

func startServer(port string, u string) {
	target, _ := url.Parse(u)
	fmt.Println("Target is: ", target)
	proxy := proxy.NewProxy(target)
	http.Handle("/", proxy)

	err := http.ListenAndServe(port, nil)

	fmt.Println("Starting server on port " + port)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
