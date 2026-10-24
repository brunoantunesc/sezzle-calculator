package web

import (
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestStaticFiles(t *testing.T) {
	files := fstest.MapFS{
		"index.html":     {Data: []byte("<html>calculator</html>")},
		"assets/app.js":  {Data: []byte("console.log('app')")},
		"assets/app.css": {Data: []byte("body{}")},
	}
	handler, err := NewHandler(files)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, path string
		status       int
		body         string
	}{
		{"GET", "/", 200, "calculator"},
		{"GET", "/calculator/history", 200, "calculator"},
		{"HEAD", "/", 200, ""},
		{"GET", "/assets/app.js", 200, "console.log"},
		{"GET", "/assets/app.css", 200, "body{}"},
		{"GET", "/assets/missing.js", 404, ""},
		{"GET", "/assets", 404, ""},
		{"GET", "/assets/", 404, ""},
		{"GET", "/.env", 404, ""},
		{"GET", "/../secret", 404, ""},
		{"POST", "/", 405, ""},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.body != "" && !strings.Contains(w.Body.String(), tc.body) {
				t.Fatal(w.Body.String())
			}
			if tc.method == "HEAD" && w.Body.Len() != 0 {
				t.Fatal("HEAD has a body")
			}
			if tc.status == 200 {
				csp := w.Header().Get("Content-Security-Policy")
				if !strings.Contains(csp, "script-src 'self'") || !strings.Contains(csp, "style-src 'self' 'unsafe-inline'") {
					t.Fatalf("incompatible CSP: %s", csp)
				}
				if strings.Contains(tc.path, ".js") && !strings.Contains(w.Header().Get("Content-Type"), "javascript") {
					t.Fatal("invalid script MIME")
				}
			}
		})
	}
}

func TestMissingIndex(t *testing.T) {
	if _, err := NewHandler(fstest.MapFS{}); err == nil {
		t.Fatal("missing build accepted")
	}
}
