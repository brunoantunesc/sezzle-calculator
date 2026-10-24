
# Sezzle Calculator

A full-stack calculator built with React, TypeScript, and Go, supporting basic and advanced arithmetic operations through a REST API.

## 1. Setup Instructions

### Prerequisites

- Node.js 22+
- Go 1.27.3
- Docker

### Running the application

Build and run the entire application in a single container:

```bash
docker build -t sezzle-calculator .
docker run --rm -p 8080:8080 sezzle-calculator
```

Open `http://localhost:8080`.

### Testing and Coverage

Backend:

```bash
cd backend
go test ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Frontend:

```bash
cd frontend
npm ci
npm run test:coverage
```

Both applications are also validated automatically through GitHub Actions.

## 2. API Examples

The backend exposes a stateless REST API for arithmetic operations.

### Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/calculations` | Execute an arithmetic operation |
| GET | `/health` | Check service availability |

### Supported Operations

| Operation | Identifier | Operands |
|---|---|---|
| Addition | `add` | 2 |
| Subtraction | `subtract` | 2 |
| Multiplication | `multiply` | 2 |
| Division | `divide` | 2 |
| Exponentiation | `power` | 2 |
| Square root | `sqrt` | 1 |

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/calculations \
  -H "Content-Type: application/json" \
  -d '{"operation":"power","operands":[2,10]}'
```

Successful response (`200 OK`):

```json
{
  "result": 1024
}
```

Square root example:

```json
{
  "operation": "sqrt",
  "operands": [81]
}
```

Response:

```json
{
  "result": 9
}
```

### Error Handling

Invalid calculations return a consistent JSON error structure.

Example: division by zero (`422 Unprocessable Content`).

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Cannot divide by zero"
  }
}
```

The API validates input types, supported operations, operand counts, and numerical results. Invalid requests return appropriate HTTP status codes.

## 3. Design Decisions

### Go Backend

The backend uses `net/http` and `encoding/json`, avoiding unnecessary frameworks and external dependencies.

Arithmetic logic is separated from HTTP handling to keep the implementation testable and maintainable.

### React SPA

The frontend uses React, TypeScript, Vite, and Material UI.

A single-page application was chosen because the calculator requires only one interactive view, without routing or server-side rendering.

All mathematical operations are performed by the Go backend, while React manages user interactions and presentation.

### Stateless Architecture

The backend does not require a database or authentication.

Each calculation is independent, and the frontend stores the 10 most recent successful calculations in browser `localStorage`.

### Security and Validation

The API treats all incoming requests as untrusted, applying strict JSON validation, request size limits, numerical overflow checks, HTTP timeouts, CORS restrictions, and basic rate limiting.

The application uses `float64` arithmetic, which is appropriate for a general-purpose calculator but does not provide exact decimal precision for financial calculations.

### Deployment

A multi-stage Docker build compiles the React frontend and Go backend into a single runtime container.

The Go server serves both the REST API and the compiled React application under the same origin, simplifying deployment and avoiding cross-origin requests in production.
