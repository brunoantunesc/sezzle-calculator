package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestServerConfiguration(t *testing.T) {
	t.Setenv("STATIC_DIR", "")
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

func TestDefaultServerConfiguration(t *testing.T) {
	t.Setenv("STATIC_DIR", "")
	t.Setenv("ADDR", "")
	t.Setenv("CORS_ALLOWED_ORIGIN", "")
	server, err := newServer()
	if err != nil {
		t.Fatal(err)
	}
	if server.Addr != ":8080" || server.Handler == nil {
		t.Fatalf("invalid defaults: %+v", server)
	}
}

func TestCombinedServer(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "index.html"), []byte("<html>calculator</html>"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("STATIC_DIR", directory)
	t.Setenv("CORS_ALLOWED_ORIGIN", "")
	server, err := newServer()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, path, body string
		status             int
		want               string
	}{
		{"GET", "/", "", 200, "calculator"},
		{"GET", "/history", "", 200, "calculator"},
		{"GET", "/health", "", 200, `{"status":"ok"}`},
		{"POST", "/api/v1/calculations", `{"operation":"divide","operands":[10,2]}`, 200, `{"result":5}`},
		{"GET", "/api/missing", "", 404, "404"},
		{"GET", "/api", "", 404, "404"},
	} {
		request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.want) {
			t.Fatalf("%s: %d %s", tc.path, response.Code, response.Body.String())
		}
		if response.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatal("security headers lost")
		}
		if tc.path == "/health" || strings.HasPrefix(tc.path, "/api/") {
			if response.Header().Get("Content-Security-Policy") != "default-src 'none'; frame-ancestors 'none'" {
				t.Fatal("API CSP changed")
			}
		}
	}
	t.Setenv("STATIC_DIR", filepath.Join(directory, "missing"))
	if _, err := newServer(); err == nil {
		t.Fatal("invalid static directory accepted")
	}
}
