# Sezzle Full-Stack Calculator

TEMPORARY README FILE
## Backend

From `backend`, run `go run .` to start the API or `go test ./...` to run tests.
No third-party dependencies are required. The module specifies Go 1.27.3.

- `POST /api/v1/calculations`: JSON with `operation` and numeric `operands`.
- `GET /health`: returns `{"status":"ok"}` for health checks.
- `ADDR` defaults to `:8080`.
- `CORS_ALLOWED_ORIGIN` optionally allows one origin, such as
  `http://localhost:5173`. CORS is disabled when unset. Wildcards, credentials,
  paths and multiple origins are rejected. Preflight permits the endpoint's
  method and the `Content-Type` request header; credentials are not enabled.

The server uses a 5-second header timeout, 10-second read/write timeouts,
60-second idle timeout and an 8 KiB header setting (Go's HTTP parser adds its
own buffer allowance). Calculation request bodies are limited to 4,096 bytes.
Shared headers restrict framing, MIME sniffing and referrer disclosure. API routes (`/api/` and `/health`) additionally receive restrictive CSP and no-store caching; future frontend routes can use their own CSP and cache policy.
TLS termination belongs at the deployment boundary; HSTS is not enabled for local HTTP.

Rate limiting allows 60 requests per minute per connection IP, including health
checks, preflight and invalid requests. Windows begin with an IP's first request;
a boundary can permit bursts across adjacent windows. Forwarded IP headers are
ignored, so deployments behind a proxy need an explicit trusted-proxy policy
before using forwarded client addresses. Limits are local to each process.
The mutex-protected map holds at most 10,000 IPs; when full, new IPs receive 429.
Expired entries are scanned at most once per minute when requests arrive; an individual expired window resets immediately on access. During idle periods, entries remain
bounded and are reclaimed on the next request; no cleanup goroutine is needed.
A 429 response includes `Retry-After` rounded up to the seconds remaining in that IP window. At map capacity, new IPs wait until the next scheduled cleanup.
These controls provide basic abuse protection, not comprehensive DDoS protection.


## Tests and coverage

See [TESTING.md](TESTING.md) for the issue #3 audit, reproducible commands, measured coverage and full-stack smoke procedure.
