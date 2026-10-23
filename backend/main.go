package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/brunoantunesc/sezzle-calculator/backend/api"
)

func newServer() (*http.Server, error) {
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin != "" {
		parsed, err := url.Parse(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGIN must be one http or https origin without a path")
		}
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return &http.Server{
		Addr:              addr,
		Handler:           api.Secure(api.NewHandler(), origin, time.Now),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    8192,
	}, nil
}

func main() {
	server, err := newServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
