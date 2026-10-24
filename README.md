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

## Single-container deployment

From the repository root:

```sh
docker build -t sezzle-calculator .
docker run --rm -p 8080:8080 sezzle-calculator
```

Open `http://localhost:8080`. Go serves the compiled React SPA at `/`, calculations
at `/api/v1/calculations`, and the health check at `/health`. Same-origin requests
need no CORS configuration. No Node process, nginx or Docker Compose is required.

The multi-stage build uses Node 24 for the frontend and Go 1.27.3 for the backend.
The scratch runtime contains only the static Go binary and compiled frontend,
running as UID/GID 65532. It has no shell or development dependencies. To use a
read-only filesystem, add `--read-only` to the run command; no server-side history
or writable runtime directories are required.

`STATIC_DIR` enables static serving and must contain a built `index.html`; an
invalid directory fails startup. Docker sets it to `/app/static`. When unset,
local Go development continues to serve only the API, alongside Vite. `ADDR`
defaults to `:8080`; changing its container port also requires updating the port
mapping. If local port 8080 is occupied, use `-p 8081:8080` and open port 8081.

Frontend responses allow same-origin scripts/API requests and inline styles for
MUI/Emotion. Scripts cannot execute inline. API responses retain their restrictive
CSP. Static responses use `Cache-Control: no-cache` for revalidation, and shared
security headers apply to both. Extensionless SPA routes return `index.html`;
missing assets, directories and reserved API routes never receive the SPA fallback.
The existing per-IP rate limit also applies to static requests and health checks.

Smoke check: open the app, calculate `10 / 2 = 5`, check history, then verify a
zero divisor produces a floating error. `/health` returns `{"status":"ok"}`.
Docker does not declare a shell-based HEALTHCHECK because the minimal runtime has
no shell or HTTP CLI; deployment tools can probe the HTTP health endpoint directly.

The Docker deployment requirements are tracked in GitHub issue #11; #6 is an
older CI pull request.

Deployment validation: the Docker image was built and run with a read-only filesystem. Production assets, SPA fallback, health, calculation success/error responses and browser history were verified; no CSP violations were observed during the browser check.
