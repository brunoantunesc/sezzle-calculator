package api

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxClients = 10000

type clientWindow struct {
	expires time.Time
	count   int
}

type limiter struct {
	nextCleanup time.Time
	mu          sync.Mutex
	entries     map[string]clientWindow
	now         func() time.Time
}

func newLimiter(now func() time.Time) *limiter {
	return &limiter{entries: make(map[string]clientWindow), now: now}
}

// allow returns the remaining wait when the request is denied.
func (l *limiter) allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if !now.Before(l.nextCleanup) {
		for key, entry := range l.entries {
			if !now.Before(entry.expires) {
				delete(l.entries, key)
			}
		}
		l.nextCleanup = now.Add(time.Minute)
	}
	entry, exists := l.entries[ip]
	if !exists || !now.Before(entry.expires) {
		if !exists && len(l.entries) >= maxClients {
			return false, l.nextCleanup.Sub(now)
		}
		entry = clientWindow{expires: now.Add(time.Minute)}
	}
	if entry.count >= 60 {
		return false, entry.expires.Sub(now)
	}
	entry.count++
	l.entries[ip] = entry
	return true, 0
}

// sharedSecurityHeaders can also be used for future static frontend routes.
func sharedSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func apiSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	w.Header().Set("Cache-Control", "no-store")
}

// Secure applies headers, an explicit CORS allowlist origin, and per-peer limits.
// Forwarded IP headers are deliberately ignored; proxy trust requires deployment
// configuration. Each handler instance owns its limiter and has no worker goroutines.
func Secure(next http.Handler, origin string, now func() time.Time) http.Handler {
	limit := newLimiter(now)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sharedSecurityHeaders(w)
		if r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/api/") {
			apiSecurityHeaders(w)
		}
		w.Header().Add("Vary", "Origin")
		allowed := origin != "" && r.Header.Get("Origin") == origin
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if parsed := net.ParseIP(ip); parsed != nil {
			ip = parsed.String()
		}
		if ok, wait := limit.allow(ip); !ok {
			// Round up so clients do not retry before the window expires.
			seconds := int64(wait / time.Second)
			if wait%time.Second != 0 {
				seconds++
			}
			w.Header().Set("Retry-After", strconv.FormatInt(seconds, 10))
			writeError(w, 429, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded")
			return
		}
		if r.Method == http.MethodOptions && allowed {
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			method := ""
			switch r.URL.Path {
			case "/health":
				method = "GET"
			case "/api/v1/calculations":
				method = "POST"
			}
			headersOK := true
			for _, header := range strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",") {
				if header = strings.TrimSpace(header); header != "" && !strings.EqualFold(header, "Content-Type") {
					headersOK = false
				}
			}
			if method != "" && r.Header.Get("Access-Control-Request-Method") == method && headersOK {
				w.Header().Set("Access-Control-Allow-Methods", method)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
