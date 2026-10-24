# Build the SPA with locked dependencies.
FROM node:24-alpine AS frontend
WORKDIR /build/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build a standalone Linux binary with no runtime Go dependencies.
FROM golang:1.27.3-alpine AS backend
WORKDIR /build/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/calculator .

# Only the binary and compiled SPA are needed at runtime.
FROM scratch
WORKDIR /app
COPY --from=backend /out/calculator /app/calculator
COPY --from=frontend /build/frontend/dist /app/static
ENV ADDR=:8080 STATIC_DIR=/app/static
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/calculator"]
