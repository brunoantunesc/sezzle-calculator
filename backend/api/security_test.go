package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	now := time.Unix(0, 0)
	limiter := newLimiter(func() time.Time { return now })
	for i := 0; i < 60; i++ {
		if ok, _ := limiter.allow("a"); !ok {
			t.Fatal("early limit")
		}
	}
	if ok, _ := limiter.allow("a"); ok {
		t.Fatal("accepted request 61")
	}
	if ok, _ := limiter.allow("b"); !ok {
		t.Fatal("shared IP limit")
	}
	now = now.Add(time.Minute)
	if ok, _ := limiter.allow("c"); !ok {
		t.Fatal("new window rejected")
	}
	if len(limiter.entries) != 1 {
		t.Fatalf("expired entries retained: %d", len(limiter.entries))
	}
	if ok, _ := limiter.allow("a"); !ok {
		t.Fatal("expired limit retained")
	}
}

func TestRateLimitConcurrent(t *testing.T) {
	limiter := newLimiter(func() time.Time { return time.Unix(0, 0) })
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ok, _ := limiter.allow("a"); ok {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	if accepted.Load() != 60 {
		t.Fatalf("accepted %d", accepted.Load())
	}
}

func TestRateLimitCapacity(t *testing.T) {
	limiter := newLimiter(func() time.Time { return time.Unix(0, 0) })
	for i := 0; i < maxClients; i++ {
		limiter.entries[string(rune(i))] = clientWindow{expires: time.Unix(0, 0).Add(time.Minute)}
	}
	if ok, _ := limiter.allow("new client"); ok {
		t.Fatal("unbounded entries")
	}
	if len(limiter.entries) != maxClients {
		t.Fatal("capacity exceeded")
	}
}

func TestSecurity(t *testing.T) {
	handler := Secure(NewHandler(), "https://calculator.example", func() time.Time { return time.Unix(0, 0) })
	for _, tc := range []struct {
		origin, method string
		status         int
		cors           string
	}{
		{"", "GET", 200, ""},
		{"https://calculator.example", "GET", 200, "https://calculator.example"},
		{"https://evil.example", "GET", 200, ""},
		{"https://calculator.example", "OPTIONS", 204, "https://calculator.example"},
		{"https://evil.example", "OPTIONS", 405, ""},
	} {
		t.Run(tc.origin+tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/health", nil)
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("Access-Control-Request-Method", "GET")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != tc.status || w.Header().Get("Access-Control-Allow-Origin") != tc.cors {
				t.Fatalf("status %d headers %v", w.Code, w.Header())
			}
			for _, key := range []string{"Content-Security-Policy", "X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy"} {
				if w.Header().Get(key) == "" {
					t.Errorf("missing %s", key)
				}
			}
		})
	}
}

func TestRateLimitHTTP(t *testing.T) {
	h := Secure(NewHandler(), "", func() time.Time { return time.Unix(0, 0) })
	for i := 0; i < 61; i++ {
		r := httptest.NewRequest(http.MethodGet, "/health", nil)
		r.RemoteAddr = "192.0.2.1:1234"
		r.Header.Set("X-Forwarded-For", "203.0.113.1")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if i == 60 && (w.Code != 429 || w.Header().Get("Retry-After") != "60" || !strings.Contains(w.Body.String(), `"code":"RATE_LIMIT_EXCEEDED"`)) {
			t.Fatalf("rate response: %d %s", w.Code, w.Body.String())
		}
	}
	r := httptest.NewRequest("GET", "/health", nil)
	r.RemoteAddr = "192.0.2.2:1234"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal("different IP limited")
	}
}

func TestPreflightRestrictions(t *testing.T) {
	for _, tc := range []struct {
		origin, method, headers string
		want                    int
	}{
		{"https://calculator.example", "POST", "Content-Type", 204},
		{"https://calculator.example", "DELETE", "Content-Type", 405},
		{"https://calculator.example", "POST", "Authorization", 405},
		{"https://evil.example", "POST", "Content-Type", 405},
	} {
		h := Secure(NewHandler(), "https://calculator.example", func() time.Time { return time.Unix(0, 0) })
		r := httptest.NewRequest("OPTIONS", "/api/v1/calculations", nil)
		r.Header.Set("Origin", tc.origin)
		r.Header.Set("Access-Control-Request-Method", tc.method)
		r.Header.Set("Access-Control-Request-Headers", tc.headers)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Errorf("preflight %+v: status %d", tc, w.Code)
		}
		if tc.want == 204 && (w.Header().Get("Access-Control-Allow-Methods") != "POST" || w.Header().Get("Access-Control-Allow-Headers") != "Content-Type") {
			t.Fatal("incorrect preflight headers")
		}
	}
	h := Secure(NewHandler(), "", func() time.Time { return time.Unix(0, 0) })
	r := httptest.NewRequest("GET", "/health", nil)
	r.Header.Set("Origin", "https://calculator.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("CORS enabled by default")
	}
}

func TestScheduledCleanup(t *testing.T) {
	now := time.Unix(0, 0)
	l := newLimiter(func() time.Time { return now })
	l.allow("first")
	scheduled := l.nextCleanup
	now = now.Add(30 * time.Second)
	l.allow("later")
	if l.nextCleanup != scheduled {
		t.Fatal("cleanup rescheduled before due")
	}
	now = now.Add(30 * time.Second)
	l.allow("active")
	if _, exists := l.entries["first"]; exists {
		t.Fatal("expired entry retained during cleanup")
	}
	if l.nextCleanup != now.Add(time.Minute) {
		t.Fatal("cleanup not rescheduled")
	}
	now = now.Add(30 * time.Second)
	for i := 0; i < 60; i++ {
		if ok, _ := l.allow("later"); !ok {
			t.Fatal("expired window not reset before cleanup")
		}
	}
	if ok, wait := l.allow("later"); ok || wait != time.Minute {
		t.Fatalf("allow=%v wait=%v", ok, wait)
	}
	if l.nextCleanup != time.Unix(120, 0) {
		t.Fatal("cleanup ran early")
	}
}

func TestRetryAfterRemainingWindow(t *testing.T) {
	now := time.Unix(0, 0)
	h := Secure(NewHandler(), "", func() time.Time { return now })
	request := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
		return w
	}
	for i := 0; i < 60; i++ {
		request()
	}
	now = now.Add(40200 * time.Millisecond)
	w := request()
	if w.Code != 429 || w.Header().Get("Retry-After") != "20" {
		t.Fatalf("status=%d retry=%s", w.Code, w.Header().Get("Retry-After"))
	}
	now = now.Add(19800 * time.Millisecond)
	if w = request(); w.Code != 200 {
		t.Fatal("window did not reset")
	}
}

func TestFrontendHeaders(t *testing.T) {
	h := Secure(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }), "", func() time.Time { return time.Unix(0, 0) })
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Header().Get("Content-Security-Policy") != "" || w.Header().Get("Cache-Control") != "" {
		t.Fatal("API-only headers applied to frontend")
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing shared headers")
	}
}
