package main

import (
	"testing"
	"time"
)

func TestServerConfiguration(t *testing.T) {
	t.Setenv("ADDR", ":9090")
	t.Setenv("CORS_ALLOWED_ORIGIN", "https://calculator.example")
	server, err := newServer()
	if err != nil {
		t.Fatal(err)
	}
	if server.Addr != ":9090" || server.ReadHeaderTimeout != 5*time.Second || server.ReadTimeout != 10*time.Second || server.WriteTimeout != 10*time.Second || server.IdleTimeout != 60*time.Second || server.MaxHeaderBytes != 8192 {
		t.Fatalf("unexpected server configuration: %+v", server)
	}
	for _, origin := range []string{"*", "null", "https://example.com/path", "https://example.com,https://other.com"} {
		t.Setenv("CORS_ALLOWED_ORIGIN", origin)
		if _, err := newServer(); err == nil {
			t.Errorf("accepted %q", origin)
		}
	}
}
