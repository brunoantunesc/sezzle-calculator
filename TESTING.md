# Tests and coverage — issue #3

Reviewed and executed on 2026-10-24 (America/Sao_Paulo).

## Review and implementation plan

The implementation already had 22 frontend tests and table-driven Go tests for
arithmetic, HTTP validation, health, CORS, headers, rate limits, concurrency,
cleanup, retry timing and server settings. Existing frontend tests covered
editing, operators, unary requests, sequential calculations, results, history,
keyboard input, loading delays, duplicate prevention and error notices.

The audit found the following gaps. All planned additions are implemented:

| Issue requirement | Existing coverage | Additions in this change |
| --- | --- | --- |
| All six operations and numerical safety | Arithmetic and HTTP operation tables | Decimal rounding, irrational root, nonzero base with zero exponent, extreme positive/negative exponents; finite-result checks with relative/absolute tolerance |
| HTTP and health contracts | Status codes, validation, body boundaries, health JSON, methods | Exact envelope fields and error field counts; body-read failure hides transport details; a valid request succeeds after invalid requests |
| Security controls | CORS, headers, concurrent/independent IP limits, bounded map, cleanup, retry expiry | Fixed clocks replace remaining real-clock usage in security tests |
| Frontend HTTP integration | Hook mocks and separate service tests | App → real service → mocked fetch → display/history flow; all binary keyboard operators assert their HTTP operation and operands |
| API failure paths | Division by zero, network failure, server failure, malformed success | Fake-timer timeout/abort and timer disposal; invalid/unknown error envelopes, overlong/empty messages, rate-limit message |
| Stale responses | Pending requests ignore input and duplicates | Late success and failure after unmount cannot affect a new instance or record history |
| History | Save, reuse, clear, persistence, cap, corrupted/unavailable storage | Exact newest-first ordering; failed HTTP calculations do not add entries |
| Accessibility | Named controls, keyboard feedback, Escape/Enter, Drawer close methods | All buttons have names, result live region, ignored shortcuts in editable fields or with modifiers |
| Full stack | Previously checked during implementation | Repeatable smoke procedure below; real Go success/error flows verified again |

No production behavior changes were needed. These are regression tests for
existing behavior; no artificial failing assertions or implementation-only mocks
were added to inflate coverage. Existing TDD tests remain in place. The latest
interaction requirement takes precedence over the original issue: **all calculator
input, including AC, is ignored while a request is pending**.

## Reproduce automated checks

Frontend (from `frontend`):

```sh
npm ci
npm test
npm run test:coverage
npm run build
npm run lint
```

Vitest uses jsdom, React Testing Library and V8 coverage. All unit/component tests
mock HTTP boundaries and require no running service. Clock-dependent tests use
fake timers or injected clocks, without sleeps. Reports include all application
TypeScript/TSX sources, excluding tests, declarations and the `main.tsx` DOM-mount
entry point. CSS and layout are not measured by code coverage.

Generated outputs: `frontend/coverage/index.html`, `coverage/lcov.info`,
`coverage/coverage-summary.json` (relative to `frontend`).

Backend (PowerShell, from `backend`; quoted flags avoid PowerShell splitting
filenames at the period):

```powershell
go test ./...
go vet ./...
go test ./... "-coverprofile=coverage.out"
go tool cover "-func=coverage.out"
New-Item -ItemType Directory -Force coverage | Out-Null
go tool cover "-html=coverage.out" -o coverage/index.html
```

For bash, use `mkdir -p coverage` in place of `New-Item`.
Both test commands return non-zero exit codes on test failures. Coverage artifacts are ignored by Git and
are not committed. ESLint also ignores generated coverage scripts.

## Final measured coverage

| Backend scope | Statements covered | Statement coverage |
| --- | --- | --- |
| All packages, including process entry point | 179 / 189 | **94.7%** |
| Calculator | 22 / 22 | **100.0%** |
| API (handlers and security) | 148 / 153 | **96.7%** |
| Server configuration / entry point | 9 / 14 | **64.3%** |

`newServer` is 100% covered. `main` is not executed by unit tests because it starts
a blocking server and can terminate the process via `log.Fatal`. Remaining API
gaps include defensive calculator-error mappings already prevented by HTTP input
validation, a parser error branch, and a remote-address fallback. Tests do not
introduce new abstractions or inject impossible calculator failures just to cover
those lines.

| Frontend scope | Statements | Branches | Functions | Lines |
| --- | --- | --- | --- | --- |
| All application sources | **96.72%** (177/183) | **93.26%** (180/193) | **100%** (39/39) | **100%** (134/134) |
| Calculator hook | 96.25% | 91.95% | 100% | 100% |
| History hook | 93.75% | 92.3% | 100% | 100% |
| API service | 100% | 100% | 100% | 100% |

**38 frontend tests pass** across six files. All three Go package suites pass.
TypeScript/production build, ESLint and `go vet ./...` pass. Line coverage is less
informative for compact JSX, so statement and branch metrics are reported too.
No hard percentage gate was added: meaningful regressions take priority over an
arbitrary threshold, as requested by the issue.

## Lightweight full-stack smoke test

This is the explicitly manual integration check permitted by issue #3. It is
separate from deterministic unit tests and requires both local services.

1. In `backend`, run `go run .` (default port 8080).
2. In `frontend`, run `npm run dev -- --port 5174 --strictPort`.
3. Open `http://127.0.0.1:5174/`, then use AC to start a fresh calculation.
4. Enter `10`, division, `2`, equals. In browser developer tools, Network, verify
   `POST /api/v1/calculations` with `{"operation":"divide","operands":[10,2]}`.
5. Verify HTTP 200 with `{"result":5}` and a displayed result of `5`.
6. Open History; verify the newest entry is `10 ÷ 2`, result `5`. Reload and verify
   it persists. Reuse its result, then close history with Escape and verify input
   is retained. History can also close with its button or backdrop.
7. Enter division, `0`, equals. Verify HTTP 422 with code `DIVISION_BY_ZERO`, a
   floating `Cannot divide by zero` notice, and no additional history entry.
8. Dismiss the notice or allow it to expire. Do not clear pre-existing history
   unless using a disposable browser profile.

Executed against the running local frontend and Go backend on 2026-10-24:
`10 ÷ 2 = 5` displayed successfully and appeared at the top of history; `5 ÷ 0`
showed `Cannot divide by zero`. The mocked HTTP integration test automates the
same success/error contract without depending on external processes.

## Practical limits

- Go coverage generation needed an unsandboxed local run because Windows denied
  the sandbox's temporary coverage-file rename; this was an environment issue,
  not a test failure.
- Race detection previously required unavailable cgo/C compiler support on this
  machine. Concurrent limiter behavior is tested, but no passing race-detector
  result is claimed. Run `go test -race ./...` on a configured CI runner.
- Unit tests check timeout/header configuration, not real-time slow-client load
  or DDoS resistance. Browser responsiveness remains a visual QA concern, not
  something code coverage proves. The full-stack smoke test is documented/manual,
  without adding a browser automation framework or CI service orchestration.
